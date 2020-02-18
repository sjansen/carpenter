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


ANON_REGEX_PART  = re.compile(r"^(\(\?\!(?P<reject>[^)]+)\))?(?P<regex>.+)$")
NAMED_REGEX_PART = re.compile(r"^(\(\?\!(?P<reject>[^)]+)\))?\(\?P<(?P<name>[^>]+)>\(?(?P<regex>[^)]+)\)?\)$")
NAMED_TYPE_PART  = re.compile(r"^<((?P<type>[^:>]+):)?(?P<name>[^:>]+)>$")
PLAIN_PART       = re.compile(r"^[^.*?+^$|\\[\](){}]+$")
PLAIN_PART_GUESS = re.compile(r"^[.]?[a-zA-Z][-_a-zA-Z0-9]+(\\?[.][a-zA-Z]+)?$")


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
            test_cases = create_test_cases(p, test_values)
            for k, v in p.test_cases.items():
                test_cases[k] = v
            context = Context({
                "pattern": p,
                "test_cases": test_cases,
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
    elif isinstance(pattern.suffix, RegexPart):
        expected += "/SUFFIX"
        test_cases = [tc + "/" for tc in test_cases]

    return {
        tc: expected
        for tc in test_cases
    }


class Pattern(object):
    def __init__(self, handler, pattern, test_cases=None):
        self.handler = handler
        self.raw = pattern
        self.test_cases = test_cases if test_cases is not None else {}

        self.regexes = set()
        self.__parse(pattern)

    def __add_regex(self, regex, name, reject=""):
        if reject:
            self.prefix.append(RegexPart(regex, name, reject))
            self.regexes.add((reject, ""))
        else:
            self.prefix.append(RegexPart(regex, name))
        self.regexes.add((regex, name))

    def __parse(self, pattern):
        self.prefix = []
        for token in tokenize(pattern):
            if self.__match_named_regex(token):
                continue
            elif self.__match_named_type(token):
                continue
            elif self.__match_plain(token):
                continue
            m = ANON_REGEX_PART.match(token)
            groups = m.groupdict()
            self.__add_regex(groups["regex"], "", groups["reject"])

        if pattern.endswith("/$") or len(self.prefix) < 1:
            self.suffix = PlainPart("/")
        elif pattern.endswith("/"):
            self.suffix = RegexPart(".*", "SUFFIX")
        else:
            self.suffix = PlainPart("/?")

    def __match_named_regex(self, token):
        m = NAMED_REGEX_PART.match(token)
        if not m:
            return False
        groups = m.groupdict()
        self.__add_regex(groups["regex"], groups["name"], groups["reject"])
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
        if m:
            self.prefix.append(PlainPart(token))
            return True
        m = PLAIN_PART_GUESS.match(token)
        if m:
            token = token.replace(r"\.", ".")
            self.prefix.append(PlainPart(token))
            return True
        return False


class PlainPart(object):
    def __init__(self, value):
        self.type = "plain"
        self.value = value
        self.value_as_repr = repr(value)

    def __eq__(self, other):
        if type(self) == type(other):
            return self.value == other.value
        if isinstance(other, str):
            return self.value == other
        return False

    def __repr__(self):
        return "PlainPart(%r)" % self.value


class RegexPart(object):
    def __init__(self, regex, name, reject=""):
        self.type = "regex"
        self.name = name
        self.regex = regex
        self.regex_as_raw = as_raw(self.regex)
        self.reject = reject
        self.reject_as_raw = as_raw(self.reject)
        self.replacement = name.upper() if name else "TODO"
        self.replacement_as_repr = repr(self.replacement)

    def __eq__(self, other):
        if not type(self) == type(other):
            return False
        if not self.regex == other.regex:
            return False
        return self.replacement == other.replacement

    def __repr__(self):
        if self.reject:
            return "RegexPart(%r, %r, %r)" % (self.regex, self.replacement, self.reject)
        else:
            return "RegexPart(%r, %r)" % (self.regex, self.replacement)


def as_raw(value):
    if '"' in value:
        return 'r"""' + value + '"""'
    else:
        return 'r"' + value + '"'


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
    "^(?!users|groups)(?P<resource>[^/]+)/$": [RegexPart(r"[^/]+", "resource", r"users|groups")],
    "help/(?!search)(.*)": [PlainPart("help"), RegexPart(r".*", "", r"search")],
    "favicon.ico": [PlainPart("favicon.ico")],
    ".well-known/": [PlainPart(".well-known")],
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
    "favicon.ico": {
        "/favicon.ico": "/favicon.ico",
    },
    "^go/(?P<page>(a|b))": {
        "/go/a": "/go/PAGE",
        "/go/b": "/go/PAGE",
    },
    "groups/<gid>": {
        "/groups/wheel": "/groups/GID",
    },
    "help/(?!search)(.*)": {
        "/help/TODO": "/help/TODO",
    },
    "^users/(?P<uid>[^/]+)": {
        "/users/sjansen": "/users/UID",
    },
    "^(?!users|groups)(?P<resource>[^/]+)/$": {
        "/roles/": "/RESOURCE/",
    },
    ".well-known/": {
        "/.well-known/": "/.well-known/SUFFIX",
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
    (r"[^/]+", "resource"): ["roles"],
    (r".*", ""): ["TODO"],
}


URL_TEMPLATE = textwrap.dedent(
    '''\
    {% with p=pattern %}{% autoescape off %}# {{ p.raw }}
    url(
        "{{ p.handler }}",
        path = {
            "prefix": [{% for part in p.prefix %}{% if part.type == "plain" %}
                {{ part.value_as_repr }},{% else %}
                ({{ part.regex_as_raw }}, {{ part.replacement_as_repr }}{% if part.reject %}, {{ part.reject_as_raw }}{% endif %}),{% endif %}{% endfor %}
            ],{% with s=p.suffix %}
            "suffix": {% if s.type == "plain" %}{{ s.value_as_repr }}{% else %}({{ s.regex_as_raw }}, {{ s.replacement_as_repr }}){% endif %},
        },{% endwith %}
        query = {
            "other": "X",
        },
        tests = {{% for test_case, expected in test_cases.items %}
            "{{ test_case }}": "{{ expected }}",{% endfor %}
        },
    ){% endautoescape %}{% endwith %}

'''
)
