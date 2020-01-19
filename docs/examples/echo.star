add_url(
    "root",
    path = {
        "prefix": [],
        "suffix": "/",
    },
    query = {},
    tests = {
        "/": "/",
    },
)

add_url(
    "echo",
    path = {
        "prefix": [],
        "suffix": (".*", lambda x: x),
    },
    query = {},
    tests = {
        "/.well-known/apple-app-site-association": "/.well-known/apple-app-site-association",
    },
)
