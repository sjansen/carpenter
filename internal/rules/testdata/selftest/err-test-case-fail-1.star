register_urls({
    "id": "unexpected-result",
    "path": {
        "parts": [],
        "slash": "strip",
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
