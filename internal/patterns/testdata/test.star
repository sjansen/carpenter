add_url(
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

add_url(
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
