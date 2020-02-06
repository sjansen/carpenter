from importlib import import_module

from django.conf import settings
from django.core.management import BaseCommand


class Command(BaseCommand):
    def handle(self, *args, **options):
        urlconf = import_module(settings.ROOT_URLCONF)
        patterns = extract_views_from_urlpatterns(urlconf.urlpatterns)
        max_length = 0
        for p in patterns:
            max_length = max(max_length, len(p[0]))
        for p in patterns:
            print("{pattern:{length}} | {lookup_str}".format(
                length=max_length,
                lookup_str=p[1],
                pattern=p[0],
            ))


def extract_pattern(p):
    if hasattr(p, 'pattern'):
        pattern = str(p.pattern)
    else:
        pattern = p.regex.pattern
    if pattern.startswith("^"):
        return pattern[1:]
    return pattern


def extract_views_from_urlpatterns(urlpatterns, base=''):
    """
    Return a list of views from a list of urlpatterns.
    Each object in the returned list is a two-tuple: (pattern, lookup_str)
    """
    views = []
    for p in urlpatterns:
        if hasattr(p, 'url_patterns'):
            try:
                patterns = p.url_patterns
            except ImportError:
                continue
            views.extend(extract_views_from_urlpatterns(
                patterns, base + extract_pattern(p),
            ))
        elif hasattr(p, 'lookup_str'):
            try:
                views.append((
                    base + extract_pattern(p), p.lookup_str,
                ))
            except ViewDoesNotExist:
                continue
        else:
            raise TypeError(_("%s does not appear to be a urlpattern object") % p)
    return views
