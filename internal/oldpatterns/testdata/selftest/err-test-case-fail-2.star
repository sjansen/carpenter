register_urls({
    "id": "rfc3092",
    "path": {
        "prefix": [
            ("foo|bar|baz", "rfc3092"),
        ],
        "suffix": "/?",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo": "/rfc3092",
        "/bar": "/rfc3092",
        "/baz": "/rfc3092",
    },
}, {
    "id": "shadowed-pattern",
    "path": {
        "prefix": ["foo"],
        "suffix": "/",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo/": "/foo/",
    },
})

