add_url("root",
    path={
        "prefix": [],
        "suffix": "/?",
    },
    query={
        "params": {
            "foo": lambda x: x,
        },
    },
    tests={
        "/?foo=bar": "/?foo=bar",
    },
)

add_url("basic-query",
    path={
        "prefix": ["search"],
        "suffix": "/?",
    },
    query={
        "dedup": "never",
        "params": {
          "q": "X",
          "utf8": None,
        },
    },
    tests={
        "/search?utf8=✔":         "/search",
        "/search/?q=cats":        "/search?q=X",
        "/search/?q=dogs&utf8=✔": "/search?q=X",
    },
)

add_url("dedup-never",
    path={
        "prefix": ["team", "membership"],
        "suffix": "/",
    },
    query={
        "dedup": "never",
        "params": {
          "users[]": "X"
        },
    },
    tests={
        "/team/membership/?users[]=alice":
            "/team/membership/?users%5B%5D=X",
        "/team/membership/?users[]=bob&users[]=eve":
            "/team/membership/?users%5B%5D=X&users%5B%5D=X",
    },
)

add_url("dedup-first",
    path={
        "prefix": ["dedup", "first"],
        "suffix": "/",
    },
    query={
        "dedup": "first",
        "params": {
          "users[]": lambda x: x,
        },
    },
    tests={
        "/dedup/first/?users[]=alice":
            "/dedup/first/?users%5B%5D=alice",
        "/dedup/first/?users[]=bob&users[]=eve":
            "/dedup/first/?users%5B%5D=bob",
    },
)

add_url("dedup-last",
    path={
        "prefix": ["dedup", "last"],
        "suffix": "/",
    },
    query={
        "dedup": "last",
        "params": {
          "users[]": lambda x: x,
        },
    },
    tests={
        "/dedup/last/?users[]=alice":
            "/dedup/last/?users%5B%5D=alice",
        "/dedup/last/?users[]=bob&users[]=eve":
            "/dedup/last/?users%5B%5D=eve",
    },
)

add_url("extra-params",
    path={
        "prefix": ["extra", "params"],
        "suffix": "/?",
    },
    query={},
    tests={
        "/extra/params?foo&bar=baz&qux=quux": "/extra/params?bar=baz&foo=&qux=quux",
    },
)
