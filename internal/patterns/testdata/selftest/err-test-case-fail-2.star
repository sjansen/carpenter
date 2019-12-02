register_urls({
    "id": "rfc3092",
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
    "id": "shadowed-pattern",
    "path": {
        "parts": ["foo"],
        "slash": "always",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo/": "/foo/",
    },
})

