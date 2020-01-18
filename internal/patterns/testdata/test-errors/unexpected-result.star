add_url(
    "first",
    path = {
        "prefix": ["foo"],
        "suffix": "/",
    },
    query = {},
    tests = {},
)

add_url(
    "second",
    path = {
        "prefix": [(".+", "ANY")],
        "suffix": "/?",
    },
    query = {},
    tests = {
        "/foo": "/foo",
    },
)
