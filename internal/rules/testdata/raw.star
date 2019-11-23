names = ("waldo|fred", "plugh")

def fn(x):
    if x == "xyzzy":
        return "thud"
    return "X"

register_urls({
    "id": "root",
    "path": {
        "parts": [],
        "slash": "strip",
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
    "id": "always",
    "path": {
        "parts": ["foo"],
        "slash": "always",
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
    "id": "never",
    "path": {
        "parts": ["bar"],
        "slash": "never",
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
    "id": "strip",
    "path": {
        "parts": ["baz"],
        "slash": "strip",
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
        "parts": [("qux", "quux")],
        "slash": "always",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/qux/": "/quux/",
    },
}, {
    "id": "query.never",
    "path": {
        "parts": ["search"],
        "slash": "never",
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
    "id": "multi",
    "path": {
        "parts": [
            "corge",
            ("grault", "garply"),
            names,
            (".+", fn),
        ],
        "slash": "strip",
    },
    "query": {
        "dedup": "never",
        "params": {
            "n": lambda x: "even" if len(x) % 2 == 0 else "odd",
        },
    },
    "tests": {
        "/corge/grault/waldo/xyzzy": "/corge/garply/plugh/thud",
        "/corge/grault/fred/42/": "/corge/garply/plugh/X",
        "/corge/grault/fred/random/?n=left": "/corge/garply/plugh/X?n=even",
        "/corge/grault/fred/random?n=right": "/corge/garply/plugh/X?n=odd",
    },
})
