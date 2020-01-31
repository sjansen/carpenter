url(
    "first",
    path = {
        "prefix": ["foo"],
        "suffix": "/",
    },
    query = {},
    tests = {
        "/foo": "",
        "/foo/": "/foo/",
        "/bar/": "",
    },
)

url(
    "second",
    path = {
        "prefix": [("b.+", "ANY")],
        "suffix": "/?",
    },
    query = {},
    tests = {
        "/foo/": "",
        "/bar/": "/ANY",
        "/baz/": "/ANY",
        "/qux/": "",
    },
)
