import re

import django
from django.core.management import BaseCommand

try:
    from django.urls.converters import get_converters
except:

    def get_converters():
        raise NotImplementedError()


NAMED_REGEX_PART = re.compile(r"^\(\?P<(?P<name>[^>]+)>(?P<regex>.*)\)$")
NAMED_TYPE_PART = re.compile(r"^<((?P<type>[^:]+):)?(?P<name>[^:]+)>$")
PLAIN_PART = re.compile("[^|.*+?\\\[\](){}]+")


TEST_CASES = {
    "^$": [],
    "articles/<int:year>/<int:month>/<slug:slug>/": [
        "articles",
        (r"[0-9]+", "year"),
        (r"[0-9]+", "month"),
        (r"[-a-zA-Z0-9_]+", "slug"),
    ],
    "^articles/(?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<slug>[\w-]+)/$": [
        "articles",
        (r"[0-9]{4}", "year"),
        (r"[0-9]{2}", "month"),
        (r"[\w-]+", "slug"),
    ],
    "^users/(?P<id>[^/]+)": ["users", (r"[^/]+", "id")],
}


class Command(BaseCommand):
    def handle(self, *args, **options):
        self_test()


def parse(pattern):
    converters = get_converters()

    prefix = []
    for token in tokenize(pattern):
        if NAMED_REGEX_PART.match(token):
            groups = NAMED_REGEX_PART.match(token).groupdict()
            part = (groups["regex"], groups["name"])
        elif NAMED_TYPE_PART.match(token):
            groups = NAMED_TYPE_PART.match(token).groupdict()
            regex = converters[groups["type"]].regex
            part = (regex, groups["name"])
        elif PLAIN_PART.match(token):
            part = token
        else:
            part = (token, "TODO")
        prefix.append(part)

    return prefix


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


def self_test():
    for tc, expected in TEST_CASES.items():
        actual = parse(tc)
        if expected == actual:
            print("PASS: {}".format(tc))
        else:
            print("FAIL: {}".format(tc))
            print("  expected: {}".format(expected))
            print("    actual: {}".format(actual))
