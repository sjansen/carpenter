url(
    "first",
    path = {
        "prefix": ["foo"],
        "suffix": "/",
    },
    query = {},
    tests = {},
)

url(
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
