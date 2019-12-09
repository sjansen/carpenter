register_urls({
    "id": "unexpected-result",
    "path": {
        "prefix": [],
        "suffix": "/?",
    },
    "query": {
        "dedup": "never",
        "params": {
            "foo": "",
        },
    },
    "tests": {
        "/?foo": "/?foo",
    },
})
