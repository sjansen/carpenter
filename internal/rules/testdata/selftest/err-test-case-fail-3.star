register_urls({
    "id": "bad-lambda",
    "path": {
        "parts": [
            (".*", lambda x: 42),
        ],
        "slash": "strip",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo": "/bar",
    },
})
