register_urls({
    "id": "bad-lambda",
    "path": {
        "prefix": [
            (".*", lambda x: 42),
        ],
        "suffix": "/?",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/foo": "/bar",
    },
})
