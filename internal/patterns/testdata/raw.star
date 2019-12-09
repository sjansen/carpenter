names = ("waldo|fred", "plugh")

def xyzzy(x):
    if x == "xyzzy":
        return "Z"
    return "X"

def thud(x):
    return ""

register_urls({
    "id": "root",
    "path": {
        "prefix": [],
        "suffix": "/?",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/": "/",
        "/Spoon!": None,
    },
}, {
    "id": "slash-required",
    "path": {
        "prefix": ["foo"],
        "suffix": "/",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo": None,
        "/foo/": "/foo/",
    },
}, {
    "id": "no-final-slash",
    "path": {
        "prefix": ["bar"],
        "suffix": "",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/bar": "/bar",
        "/bar/": None,
    },
}, {
    "id": "optional-slash",
    "path": {
        "prefix": ["baz"],
        "suffix": "/?",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/baz": "/baz",
        "/baz/": "/baz",
    },
}, {
    "id": "regex",
    "path": {
        "prefix": [("qux", "quux")],
        "suffix": "/",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/qux/": "/quux/",
    },
}, {
    "id": "query-no-final-slash",
    "path": {
        "prefix": ["search"],
        "suffix": "",
    },
    "query": {
        "dedup": "never",
        "params": {
            "q": "X",
            "utf8": None,
        },
    },
    "tests": {
        "/search?q=cats": "/search?q=X",
        "/search?utf8=\342\234\224": "/search",
        "/search?q=dogs&utf8=\342\234\224": "/search?q=X",
    },
}, {
    "id": "complex",
    "path": {
        "prefix": [
            "corge",
            ("grault", "garply"),
            names,
            (".+", xyzzy),
        ],
        "suffix": (".*", thud),
    },
    "query": {
        "dedup": "never",
        "params": {
            "n": lambda x: "even" if len(x) % 2 == 0 else "odd",
        },
    },
    "tests": {
        "/corge/grault/fred/42/": "/corge/garply/plugh/X",
        "/corge/grault/waldo/xyzzy": "/corge/garply/plugh/Z",
        "/corge/grault/fred/random/?n=left": "/corge/garply/plugh/X?n=even",
        "/corge/grault/fred/random?n=right": "/corge/garply/plugh/X?n=odd",
    },
})
