register_urls({
    "id": "unexpected-result",
    "path": {
        "prefix": [],
        "suffix": "strip",
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
