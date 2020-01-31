url(
    "root",
    path = {
        "prefix": [],
        "suffix": "/",
    },
    query = {
        "dedup": "never",
        "params": {},
    },
    tests = {
        "/": "/",
        "/rfc3092/": None,
    },
)

url(
    "slash-required",
    path = {
        "prefix": ["foo"],
        "suffix": "/",
    },
    query = {
        "dedup": "first",
        "params": {},
    },
    tests = {
        "/foo": None,
        "/foo/": "/foo/",
        "/bar": None,  # replaced by no-final-slash
        "/baz/": None,  # replaced by optional-slash
    },
)

url(
    "no-final-slash",
    path = {
        "prefix": ["bar"],
        "suffix": "",
    },
    query = {
        "dedup": "last",
        "params": {},
    },
    tests = {
        "/foo": None,  # shadowed by slash-required
        "/foo/": None,  # shadowed by slash-required
        "/bar": "/bar",
        "/bar/": None,
    },
)

url(
    "optional-slash",
    path = {
        "prefix": ["baz"],
        "suffix": "/?",
    },
    query = {
        "dedup": "never",
        "params": {},
    },
    tests = {
        "/foo": None,  # shadowed by slash-required
        "/foo/": None,  # shadowed by slash-required
        "/baz": "/baz",
        "/baz/": "/baz",
    },
)

url(
    "regex",
    path = {
        "prefix": [("qux", "quux")],
        "suffix": "/?",
    },
    query = {
        "dedup": "never",
        "params": {},
    },
    tests = {
        "/qux/": "/quux",
    },
)

url(
    "goldilocks",
    path = {
        "prefix": [
            "corge",
            "grault",
            "garply",
        ],
        "suffix": "",
    },
    query = {},
    tests = {
        "/corge": None,
        "/corge/grault": None,
        "/corge/grault/garply": "/corge/grault/garply",
        "/corge/grault/garply/": None,
        "/corge/grault/garply/fred": None,
    },
)

url(
    "query",
    path = {
        "prefix": ["search"],
        "suffix": "/?",
    },
    query = {
        "dedup": "never",
        "params": {
            "q": "X",
            "utf8": None,
        },
    },
    tests = {
        "/search?utf8=✔": "/search",
        "/search/?q=cats": "/search?q=X",
        "/search/?q=dogs&utf8=✔": "/search?q=X",
    },
)
