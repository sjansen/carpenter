register_urls({
    "id": "bad-lambda",
    "path": {
        "prefix": [
            (".*", lambda x: 42),
        ],
        "suffix": "strip",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo": "/bar",
    },
})
