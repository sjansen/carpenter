add_url("root",
    path={
        "prefix": [],
        "suffix": "/?",
    },
    query={
        "dedup": "never",
        "params": {},
    },
    tests={
        "/": "/",
        "/rfc3092/": None,
    },
)

add_url("slash-required",
    path={
        "prefix": ["foo"],
        "suffix": "/",
    },
    query={
        "dedup": "first",
        "params": {},
    },
    tests={
        "/foo": None,
        "/foo/": "/foo/",
        "/bar": None, # replaced by no-final-slash
        "/baz/": None, # replaced by optional-slash
    },
)

add_url("no-final-slash",
    path={
        "prefix": ["bar"],
        "suffix": "",
    },
    query={
        "dedup": "last",
        "params": {},
    },
    tests={
        "/foo": None, # shadowed by slash-required
        "/foo/": None, # shadowed by slash-required
        "/bar": "/bar",
        "/bar/": None,
    },
)

add_url("optional-slash",
    path={
        "prefix": ["baz"],
        "suffix": "/?",
    },
    query={
        "dedup": "never",
        "params": {},
    },
    tests={
        "/foo": None, # shadowed by slash-required
        "/foo/": None, # shadowed by slash-required
        "/baz": "/baz",
        "/baz/": "/baz",
    },
)

add_url("regex",
    path={
        "prefix": [("qux", "quux")],
        "suffix": "/?",
    },
    query={
        "dedup": "never",
        "params": {},
    },
    tests={
        "/qux/": "/quux",
    },
)

add_url("prefix-regex",
    path={
        "prefix": [
          ("foo|bar", "baz"),
          ("qux|quux", lambda x: "n=%s" % len(x)),
         ],
        "suffix": "/",
    },
    query={
        "dedup": "first",
        "params": {
          "utf8": lambda x: chr(x),
        },
    },
    tests={
        "/foo/qux":          None,
        "/foo/qux/":         "/baz/n=3/",
        "/bar/qux/":         "/baz/n=3/",
        "/bar/quux/":        "/baz/n=4/",
    },
)

add_url("suffix-regex",
    path={
        "prefix": ["corge"],
        "suffix": (".+", lambda x: x.upper()),
    },
    query={
        "dedup": "last",
        "params": {
          "utf8": lambda x: x == "✔",
        },
    },
    tests={
        "/corge/": None,
        "/corge/grault":      "/corge/GRAULT",
        "/corge/garply":      "/corge/GARPLY",
        "/corge/waldo/":      "/corge/WALDO/",
    },
)
