import argparse
import csv
import sys
from importlib import import_module

from django.conf import settings
from django.core.management import BaseCommand


class Command(BaseCommand):
    def add_arguments(self, parser):
        parser.add_argument(
            "-f", "--format", nargs="?", choices=["csv", "txt"], default="txt",
        )
        parser.add_argument(
            "-o",
            "--output",
            nargs="?",
            default=sys.stdout,
            type=argparse.FileType("w"),
        )
        parser.add_argument("-s", "--sort", nargs="?", choices=["handler", "pattern"])

    def handle(self, *args, **options):
        urlconf = import_module(settings.ROOT_URLCONF)
        patterns = extract_views_from_urlpatterns(urlconf.urlpatterns)

        if options["sort"] == "pattern":
            patterns.sort(key=lambda x: x[0])
        elif options["sort"] == "handler":
            patterns.sort(key=lambda x: x[1])

        if options["format"] == "csv":
            w = csv.writer(options["output"])
            w.writerow(["Pattern", "Handler", "Test Case", "Expected"])
            for p in patterns:
                w.writerow([p[0], p[1], "", ""])
        else:
            max_length = 0
            for p in patterns:
                max_length = max(max_length, len(p[0]))
            for p in patterns:
                options["output"].write(
                    "{pattern:{length}} | {handler}\n".format(
                        length=max_length, pattern=p[0], handler=p[1],
                    )
                )


def extract_pattern(base, p):
    if hasattr(p, "pattern"):
        pattern = str(p.pattern)
        if p.pattern._is_endpoint and not hasattr(p.pattern, "_regex"):
            pattern += "$"
    else:
        pattern = p.regex.pattern
    if base and pattern.startswith("^"):
        pattern = pattern[1:]
    return base + pattern


def extract_views_from_urlpatterns(urlpatterns, base=""):
    """
    Return a list of views from a list of urlpatterns.
    Each object in the returned list is a two-tuple: (pattern, handler)
    """
    views = []
    for p in urlpatterns:
        if hasattr(p, "url_patterns"):
            try:
                patterns = p.url_patterns
            except ImportError:
                continue
            views.extend(
                extract_views_from_urlpatterns(patterns, extract_pattern(base, p))
            )
        elif hasattr(p, "lookup_str"):
            try:
                views.append((extract_pattern(base, p), p.lookup_str))
            except ViewDoesNotExist:
                continue
        else:
            raise TypeError(_("%s does not appear to be a urlpattern object") % p)
    return views
