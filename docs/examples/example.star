def rename(path):
    return "renamed-%s" % path

set_rename_filter(rename)

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
        "suffix": (r".*", lambda x: x, r"^search/?$"),
    },
    query = {},
    tests = {
        "/.well-known/apple-app-site-association": "/.well-known/apple-app-site-association",
    },
)

url(
    "search",
    path = {
        "prefix": ["search"],
        "suffix": "/?",
    },
    query = {
        "match": {
            "q": "X",
        },
    },
    tests = {
        "/search?q=apples": "/search?q=X",
        "/search/?q=oranges": "/search?q=X",
    },
)
