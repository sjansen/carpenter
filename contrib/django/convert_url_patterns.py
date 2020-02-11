import re

import django
from django.core.management import BaseCommand

try:
    from django.urls.converters import get_converters

    converters = get_converters()
except:
    converters = {}


NAMED_REGEX_PART = re.compile(r"^\(\?P<(?P<name>[^>]+)>(?P<regex>.*)\)$")
NAMED_TYPE_PART = re.compile(r"^<((?P<type>[^:]+):)?(?P<name>[^:]+)>$")
PLAIN_PART = re.compile("[^|.*+?\\\[\](){}]+")


class Command(BaseCommand):
    def handle(self, *args, **options):
        self_test()


def self_test():
    for tc, expected in TEST_CASES.items():
        actual = Pattern("tc", tc)
        if expected == actual.prefix:
            print("PASS: {}".format(tc))
        else:
            print("FAIL: {}".format(tc))
            print("  expected: {}".format(expected))
            print("    actual: {}".format(actual.prefix))


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
    def __init__(self, handler, pattern):
        self.handler = handler
        self.prefix = []
        self.__parse(pattern)

    def __parse(self, pattern):
        for token in tokenize(pattern):
            if self.__match_named_regex(token):
                continue
            elif self.__match_named_type(token):
                continue
            elif self.__match_plain(token):
                continue
            prefix.append(RegexPart(token, "TODO"))

    def __match_named_regex(self, token):
        m = NAMED_REGEX_PART.match(token)
        if not m:
            return False
        groups = m.groupdict()
        self.prefix.append(RegexPart(groups["regex"], groups["name"]))
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
        self.prefix.append(RegexPart(regex, groups["name"]))
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
    "groups/<gid>": [PlainPart("groups"), RegexPart(r"[^/]+", "gid")],
    "^users/(?P<uid>[^/]+)": [PlainPart("users"), RegexPart(r"[^/]+", "uid")],
}
