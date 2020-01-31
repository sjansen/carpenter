url(
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

url(
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
