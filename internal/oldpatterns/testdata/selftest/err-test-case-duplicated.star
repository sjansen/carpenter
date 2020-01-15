register_urls({
    "id": "original-test-case",
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
    "id": "duplicate-test-case",
    "path": {
        "prefix": ["foo"],
        "suffix": "",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo": "/foo",
        "/foo/": "",
    },
})