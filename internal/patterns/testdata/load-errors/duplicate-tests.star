add_url(
    "slash-required",
    path = {
        "prefix": [
            ("foo|bar", "rfc3029"),
        ],
        "suffix": "/",
    },
    query = {},
    tests = {
        "/foo": None,
        "/foo/": "/rfc3092/",
        "/bar/": "/rfc3092/",
    },
)

add_url(
    "optional-slash",
    path = {
        "prefix": [
            ("[bB][aA][rRzZ]", lambda x: x.lower()),
        ],
        "suffix": "/?",
    },
    query = {},
    tests = {
        "/bar": "/bar",
        "/bar/": "/bar",
        "/Bar/": "/bar",
        "/baZ/": "/baz",
        "/BAZ/": "/baz",
    },
)
