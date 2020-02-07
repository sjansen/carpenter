url(
    "root",
    path = {
        "prefix": [],
        "suffix": "/",
    },
    query = {
        "match": {
            "foo": lambda x: x,
        },
    },
    tests = {
        "/?foo=bar": "/?foo=bar",
    },
)

url(
    "basic-query",
    path = {
        "prefix": ["search"],
        "suffix": "/?",
    },
    query = {
        "dedup": "never",
        "match": {
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

url(
    "dedup-never",
    path = {
        "prefix": ["team", "membership"],
        "suffix": "/",
    },
    query = {
        "dedup": "never",
        "match": {
            "users[]": "X",
        },
    },
    tests = {
        "/team/membership/?users[]=alice": "/team/membership/?users%5B%5D=X",
        "/team/membership/?users[]=bob&users[]=eve": "/team/membership/?users%5B%5D=X&users%5B%5D=X",
    },
)

url(
    "dedup-first",
    path = {
        "prefix": ["dedup", "first"],
        "suffix": "/",
    },
    query = {
        "dedup": "first",
        "match": {
            "users[]": lambda x: x,
        },
    },
    tests = {
        "/dedup/first/?users[]=alice": "/dedup/first/?users%5B%5D=alice",
        "/dedup/first/?users[]=bob&users[]=eve": "/dedup/first/?users%5B%5D=bob",
    },
)

url(
    "dedup-last",
    path = {
        "prefix": ["dedup", "last"],
        "suffix": "/",
    },
    query = {
        "dedup": "last",
        "match": {
            "users[]": lambda x: x,
        },
    },
    tests = {
        "/dedup/last/?users[]=alice": "/dedup/last/?users%5B%5D=alice",
        "/dedup/last/?users[]=bob&users[]=eve": "/dedup/last/?users%5B%5D=eve",
    },
)

url(
    "extra-params",
    path = {
        "prefix": ["extra", "params"],
        "suffix": "/?",
    },
    query = {},
    tests = {
        "/extra/params?foo&bar=baz&qux=quux": "/extra/params?bar=baz&foo=&qux=quux",
    },
)
