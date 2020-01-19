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
    "search",
    path = {
        "prefix": ["search"],
        "suffix": "/?",
    },
    query = {
        "params": {
            "q": "X",
        },
    },
    tests = {
        "/search?q=apples": "/search?q=X",
        "/search/?q=oranges": "/search?q=X",
    },
)
