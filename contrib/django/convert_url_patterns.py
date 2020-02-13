import argparse
import csv
import re
import sys
import textwrap

from django.core.management import BaseCommand
from django.template import Context, Template

try:
    from django.urls.converters import get_converters

    converters = get_converters()
except:
    converters = {}


NAMED_REGEX_PART = re.compile(r"^\(\?P<(?P<name>[^>]+)>\(?(?P<regex>[^)]+)\)?\)$")
NAMED_TYPE_PART = re.compile(r"^<((?P<type>[^:>]+):)?(?P<name>[^:>]+)>$")
PLAIN_PART = re.compile("[^|.*+?\\\[\](){}]+")


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
        parser.add_argument("-t", "--self-test", action="store_true")
        parser.add_argument(
            "-u",
            "--unknown-regexes",
            nargs="?",
            type=argparse.FileType("w"),
        )

    def handle(self, *args, **options):
        output = options["output"]
        if options["self_test"]:
            self_test(output)
            return
        if options["input"]:
            reader = csv.DictReader(options["input"])
            patterns = [
                Pattern(
                    row["Handler"],
                    row["Pattern"],
                    test_case=row.get("Test Case"),
                    expected=row.get("Expected"),
                )
                for row in reader
            ]
        else:
            patterns = [Pattern(tc, tc) for tc in TEST_CASES.keys()]
        self.__render(output, patterns)
        if options["unknown_regexes"]:
            self.__dump_regexes(options["unknown_regexes"], patterns)

    def __dump_regexes(self, output, patterns):
        regexes = set()
        for p in patterns:
            regexes.update(p.regexes)
        w = csv.writer(output)
        w.writerow(["RegEx", "Name", "Example"])
        for row in sorted(regexes):
            w.writerow(row)

    def __render(self, output, patterns):
        template = Template(URL_TEMPLATE)
        for p in patterns:
            context = Context({"pattern": p})
            output.write(template.render(context))


def self_test(output):
    for tc, expected in TEST_CASES.items():
        actual = Pattern("tc", tc)
        if expected == actual.prefix:
            output.write("PASS: {}\n".format(tc))
        else:
            output.write("FAIL: {}\n".format(tc))
            output.write("  expected: {}\n".format(expected))
            output.write("    actual: {}\n".format(actual.prefix))


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
    def __init__(self, handler, pattern, test_case=None, expected=None):
        self.handler = handler
        self.test_case = test_case
        self.expected = expected

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
    def __init__(self, regex, replacement):
        self.type = "regex"
        self.regex = regex
        self.replacement = replacement

    def __eq__(self, other):
        if not type(self) == type(other):
            return False
        if not self.regex == other.regex:
            return False
        return self.replacement == other.replacement


TEST_CASES = {
    "^$": [],
    "articles/<int:year>/<int:month>/<slug:slug>/": [
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
    "^go/(?P<page>(a|b))": [PlainPart("go"), RegexPart(r"a|b", "page")],
    "groups/<gid>": [PlainPart("groups"), RegexPart(r"[^/]+", "gid")],
    "^users/(?P<uid>[^/]+)": [PlainPart("users"), RegexPart(r"[^/]+", "uid")],
}


URL_TEMPLATE = textwrap.dedent(
    '''\
    url({% with p=pattern %}{% autoescape off %}
        "{{ p.handler }}",
        path = {
            "prefix": [{% for part in p.prefix %}{% if part.type == "plain" %}
                "{{ part.value }}",{% else %}
                ({% if '"' in part.regex %}r"""{{ part.regex }}"""{% else %}r"{{ part.regex }}"{% endif %}, "{{ part.replacement|upper|default:"TODO" }}"),{% endif %}{% endfor %}
            ],
            "suffix": "{{ p.suffix }}",
        },
        query = {
            "other": "X",
        },
        tests = {{% if p.test_case %}
            "{{ p.test_case }}": "{{ p.expected }}",
        {% endif %}},
    {% endautoescape %}{% endwith %})

'''
)
