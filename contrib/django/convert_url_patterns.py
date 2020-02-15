import argparse
import csv
import re
import sys
import textwrap
from collections import defaultdict

from django.core.management import BaseCommand
from django.template import Context, Template

try:
    from django.urls.converters import get_converters

    converters = get_converters()
except:
    converters = {}


NAMED_REGEX_PART = re.compile(r"^\(\?P<(?P<name>[^>]+)>\(?(?P<regex>[^)]+)\)?\)$")
NAMED_TYPE_PART = re.compile(r"^<((?P<type>[^:>]+):)?(?P<name>[^:>]+)>$")
PLAIN_PART = re.compile(r"^[^.*?+^$|\\[\](){}]+$")


class Command(BaseCommand):
    def add_arguments(self, parser):
        parser.add_argument(
            "-i", "--input", nargs="?", type=argparse.FileType("r"),
        )
        parser.add_argument(
            "-o",
            "--output",
            nargs="?",
            default=sys.stdout,
            type=argparse.FileType("w"),
        )
        parser.add_argument("--self-test", action="store_true")
        parser.add_argument(
            "-t", "--test-values", nargs="?", type=argparse.FileType("r"),
        )
        parser.add_argument(
            "-u", "--unknown-regexes", nargs="?", type=argparse.FileType("w"),
        )

    def handle(self, *args, **options):
        output = options["output"]
        if options["self_test"]:
            self_test(output)
            return

        patterns = self.__load_patterns(options["input"])
        test_values = self.__load_test_values(options["test_values"])
        self.__render(output, patterns, test_values)
        if options["unknown_regexes"]:
            self.__dump_regexes(options["unknown_regexes"], patterns, test_values)

    def __dump_regexes(self, output, patterns, test_values):
        regexes = set()
        for p in patterns:
            for key in p.regexes:
                if key not in test_values:
                    regexes.update(p.regexes)
        w = csv.writer(output)
        w.writerow(["RegEx", "Name", "Example"])
        for row in sorted(regexes):
            w.writerow(row)

    def __load_patterns(self, input):
        if not input:
            return [Pattern(tc, tc) for tc in TEST_CASES.keys()]

        reader = csv.DictReader(input)
        return [
            Pattern(
                row["Handler"],
                row["Pattern"],
                test_cases={
                    row["Test Case"]: row["Expected"],
                }
            )
            if row["Test Case"] else
            Pattern(
                row["Handler"],
                row["Pattern"],
            )
            for row in reader
        ]

    def __load_test_values(self, input):
        test_values = defaultdict(lambda: set())
        if input:
            reader = csv.DictReader(input)
            for row in reader:
                regex = row["RegEx"]
                name = row.get("Name", "")
                value = row.get("Example")
                if regex and value:
                    test_values[(regex, name)].add(value)
        for k, v in test_values.items():
            test_values[k] = sorted(v)
        return test_values

    def __render(self, output, patterns, test_values):
        template = Template(URL_TEMPLATE)
        for p in patterns:
            context = Context({
                "pattern": p,
                "test_cases": create_test_cases(p, test_values),
            })
            output.write(template.render(context))


def create_test_cases(pattern, test_values):
    if not test_values:
        return {}

    expected = ""
    test_cases = [""]
    for part in pattern.prefix:
        if isinstance(part, PlainPart):
            expected = expected + "/" + part.value
            test_cases = [
                tc + "/" + part.value
                for tc in test_cases
            ]
        else:
            expected = expected + "/" + part.replacement
            values = test_values.get((part.regex, part.name))
            if not values:
                values = test_values.get((part.regex, ''))
            tmp = []
            if values:
                for v in values:
                    for tc in test_cases:
                        tmp.append(tc+"/"+v)
            test_cases = tmp

    if pattern.suffix == "/":
        expected += "/"
        test_cases = [tc + "/" for tc in test_cases]

    return {
        tc: expected
        for tc in test_cases
    }


def tokenize(pattern):
    pattern = pattern.lstrip("^").rstrip("/$")

    begin, brackets, parens, escaped = 0, 0, 0, False
    for i, c in enumerate(pattern):
        if c == "/" and (brackets + parens) < 1:
            if escaped:
                end = i - 1
            else:
                end = i
            yield pattern[begin:end]
            begin = i + 1
        if escaped:
            escaped = False
        else:
            if c == "\\":
                escaped = True
            elif c == "[":
                brackets += 1
            elif c == "]":
                brackets -= 1
            elif c == "(":
                parens += 1
            elif c == ")":
                parens -= 1

    if begin < len(pattern):
        yield pattern[begin:]


class Pattern(object):
    def __init__(self, handler, pattern, test_cases=None):
        self.handler = handler
        self.test_cases = test_cases if test_cases is not None else {}

        self.regexes = set()
        self.__parse(pattern)

    def __parse(self, pattern):
        self.prefix = []
        for token in tokenize(pattern):
            if self.__match_named_regex(token):
                continue
            elif self.__match_named_type(token):
                continue
            elif self.__match_plain(token):
                continue
            self.__add_regex(token, "")

        if pattern.endswith("/$") or len(self.prefix) < 1:
            self.suffix = "/"
        else:
            self.suffix = "/?"

    def __add_regex(self, regex, name):
        self.prefix.append(RegexPart(regex, name))
        self.regexes.add((regex, name))

    def __match_named_regex(self, token):
        m = NAMED_REGEX_PART.match(token)
        if not m:
            return False
        groups = m.groupdict()
        regex = groups["regex"]
        self.__add_regex(regex, groups["name"])
        return True

    def __match_named_type(self, token):
        m = NAMED_TYPE_PART.match(token)
        if not m:
            return False
        groups = m.groupdict()
        if groups.get("type"):
            regex = converters[groups["type"]].regex
        else:
            regex = "[^/]+"
        self.__add_regex(regex, groups["name"])
        return True

    def __match_plain(self, token):
        m = PLAIN_PART.match(token)
        if not m:
            return False
        self.prefix.append(PlainPart(token))
        return True


class PlainPart(object):
    def __init__(self, value):
        self.type = "plain"
        self.value = value

    def __eq__(self, other):
        if not type(self) == type(other):
            return False
        return self.value == other.value


class RegexPart(object):
    def __init__(self, regex, name):
        self.type = "regex"
        self.regex = regex
        self.name = name
        self.replacement = name.upper() if name else "TODO"

    def __eq__(self, other):
        if not type(self) == type(other):
            return False
        if not self.regex == other.regex:
            return False
        return self.replacement == other.replacement

    def __repr__(self):
        if '"' in self.regex:
            return 'r"""' + self.regex + '"""'
        else:
            return 'r"' + self.regex + '"'

def self_test(output):
    for tc, expected in EXPECTED_PATTERNS.items():
        pattern = Pattern("tc", tc)
        if expected != pattern.prefix:
            output.write("FAIL: {}\n".format(tc))
            output.write("  expected: {}\n".format(expected))
            output.write("    actual: {}\n".format(pattern.prefix))
            continue
        output.write("PASS: {}\n".format(tc))
        expected = EXPECTED_TEST_CASES[tc]
        actual = create_test_cases(pattern, TEST_VALUES)
        if expected != actual:
            output.write("FAIL: {}\n".format(tc))
            output.write("  expected: {}\n".format(expected))
            output.write("    actual: {}\n".format(actual))


EXPECTED_PATTERNS = {
    "": [],
    "articles/<int:year>/<int:month>/<slug:slug>/$": [
        PlainPart("articles"),
        RegexPart(r"[0-9]+", "year"),
        RegexPart(r"[0-9]+", "month"),
        RegexPart(r"[-a-zA-Z0-9_]+", "slug"),
    ],
    "^articles/(?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<slug>[\w-]+)/$": [
        PlainPart("articles"),
        RegexPart(r"[0-9]{4}", "year"),
        RegexPart(r"[0-9]{2}", "month"),
        RegexPart(r"[\w-]+", "slug"),
    ],
    "a|b|c": [RegexPart("a|b|c", "")],
    "^go/(?P<page>(a|b))": [PlainPart("go"), RegexPart(r"a|b", "page")],
    "groups/<gid>": [PlainPart("groups"), RegexPart(r"[^/]+", "gid")],
    "^users/(?P<uid>[^/]+)": [PlainPart("users"), RegexPart(r"[^/]+", "uid")],
}

EXPECTED_TEST_CASES = {
    "": {"/": "/"},
    "articles/<int:year>/<int:month>/<slug:slug>/$": {
        "/articles/2020/02/Slurms_MacKenzie/": "/articles/YEAR/MONTH/SLUG/",
    },
    "^articles/(?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<slug>[\w-]+)/$": {
        "/articles/1974/08/Philip_J_Fry/": "/articles/YEAR/MONTH/SLUG/",
    },
    "a|b|c": {
        "/a": "/TODO",
        "/b": "/TODO",
        "/c": "/TODO",
    },
    "^go/(?P<page>(a|b))": {
        "/go/a": "/go/PAGE",
        "/go/b": "/go/PAGE",
    },
    "groups/<gid>": {
        "/groups/wheel": "/groups/GID",
    },
    "^users/(?P<uid>[^/]+)": {
        "/users/sjansen": "/users/UID",
    },
}


TEST_VALUES = {
    (r"[0-9]+", "year"): ["2020"],
    (r"[0-9]+", "month"): ["02"],
    (r"[-a-zA-Z0-9_]+", "slug"): ["Slurms_MacKenzie"],
    (r"[0-9]{4}", "year"): ["1974"],
    (r"[0-9]{2}", "month"): ["08"],
    (r"[\w-]+", "slug"): ["Philip_J_Fry"],
    (r"a|b|c", ""): ["a", "b", "c"],
    (r"a|b", "page"): ["a", "b"],
    (r"[^/]+", "gid"): ["wheel"],
    (r"[^/]+", "uid"): ["sjansen"],
}


URL_TEMPLATE = textwrap.dedent(
    '''\
    url({% with p=pattern %}{% autoescape off %}
        "{{ p.handler }}",
        path = {
            "prefix": [{% for part in p.prefix %}{% if part.type == "plain" %}
                "{{ part.value }}",{% else %}
                ({{ part|stringformat:"r" }}, "{{ part.replacement }}"),{% endif %}{% endfor %}
            ],
            "suffix": "{{ p.suffix }}",
        },
        query = {
            "other": "X",
        },
        tests = {{% for test_case, expected in p.test_cases.items %}
            "{{ test_case }}": "{{ expected }}",{% endfor %}{% for test_case, expected in test_cases.items %}
            "{{ test_case }}": "{{ expected }}",{% endfor %}
        },
    {% endautoescape %}{% endwith %})

'''
)
