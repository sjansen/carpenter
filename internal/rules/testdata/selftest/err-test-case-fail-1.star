register_urls({
    "id": "first",
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
