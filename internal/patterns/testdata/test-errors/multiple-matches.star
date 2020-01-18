add_url(
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

add_url(
    "second",
    path = {
        "prefix": [(".+", "ANY")],
        "suffix": "/",
    },
    query = {},
    tests = {},
)
