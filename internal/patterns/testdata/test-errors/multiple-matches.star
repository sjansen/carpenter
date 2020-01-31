url(
    "first",
    path = {
        "prefix": ["foo"],
        "suffix": "/",
    },
    query = {},
    tests = {
        "/foo/": "/foo/",
    },
)

url(
    "second",
    path = {
        "prefix": [(".+", "ANY")],
        "suffix": "/",
    },
    query = {},
    tests = {},
)
