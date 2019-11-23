register_urls({
    "id": "first",
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
    "id": "conflict",
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
