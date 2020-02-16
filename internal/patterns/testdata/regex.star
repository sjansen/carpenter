url(
    "prefix-regex",
    path = {
        "prefix": [
            (r"foo|bar", "baz"),
            (r"qux|quux", lambda x: "n=%s" % len(x)),
        ],
        "suffix": "/",
    },
    query = {
        "dedup": "first",
        "match": {
            "utf8": lambda k, v: str(v == "✔"),
        },
    },
    tests = {
        "/foo/qux": None,
        "/foo/qux/": "/baz/n=3/",
        "/bar/qux/": "/baz/n=3/",
        "/bar/quux/": "/baz/n=4/",
        "/foo/quux/?utf8=✔": "/baz/n=4/?utf8=True",
    },
)

url(
    "suffix-regex",
    path = {
        "prefix": ["corge"],
        "suffix": (r".+", lambda x: x.upper()),
    },
    query = {
        "dedup": "last",
        "match": {
            "utf8": lambda k, v: str(v == "✔"),
        },
    },
    tests = {
        "/corge/": None,
        "/corge/grault": "/corge/GRAULT",
        "/corge/garply": "/corge/GARPLY",
        "/corge/waldo/": "/corge/WALDO/",
        "/corge/fred?utf8=✔": "/corge/FRED?utf8=True",
        "/corge/fred?utf8=!": "/corge/FRED?utf8=False",
    },
)

url(
    "reject-regex",
    path = {
        "prefix": [
            (r"[a-z]", "X", r"[aeiou]"),
        ],
        "suffix": "/?",
    },
    query = {},
    tests = {
        "/a": None,
        "/b": "/X",
        "/c/": "/X",
    },
)

url(
    "any-suffix",
    path = {
        "prefix": ["prefix"],
        "suffix": (r".*", "SUFFIX"),
    },
    query = {},
    tests = {
        "/prefix": None,
        "/prefix/": "/prefix/SUFFIX",
        "/prefix/suffix": "/prefix/SUFFIX",
    },
)
