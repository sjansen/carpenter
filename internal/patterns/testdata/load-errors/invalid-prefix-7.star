url(
    "example",
    path = {
        "prefix": [
            (".*", lambda x: "<" + x + ">", lambda x: x == "hunter2"),
        ],
        "suffix": "/",
    },
    query = {},
    tests = {},
)
