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
        "/bar": "/bar",
        "/bar/": None,
    },
)
