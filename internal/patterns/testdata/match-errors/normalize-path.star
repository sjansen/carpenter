url(
    "example",
    path = {
        "prefix": [
            "foo",
            ("bar|baz", lambda x: int(x)),
        ],
        "suffix": "/?",
    },
    query = {},
    tests = {
        "/foo/bar": "/foo/baz",
    },
)
