import argparse
import json
import sys

from django.core.management import BaseCommand
from django.urls import resolve


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

    def handle(self, *args, **options):
        cases = json.load(options["input"])
        output = options["output"]

        failures = 0
        for path in sorted(cases.keys()):
            expected = cases[path]
            try:
                match = resolve(path)
                actual = lookup_str(match.func)
            except BaseException as e:
                actual = ''
            if expected == actual:
                output.write("PASS: {}\n".format(path))
            else:
                failures += 1
                output.write("FAIL: {}\n".format(path))
                output.write("  expected: {}\n".format(expected))
                output.write("    actual: {}\n".format(actual))

        if failures > 0:
            sys.exit(1)


def lookup_str(handler):
    module = handler.__module__
    name = handler.__qualname__
    return '.'.join([module, name])
