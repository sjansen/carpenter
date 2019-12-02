register_urls({
    "id": "original-test-case",
    "path": {
        "parts": [
            ("foo|bar|baz", "rfc3092"),
        ],
        "slash": "strip",
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
        "parts": ["foo"],
        "slash": "never",
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
