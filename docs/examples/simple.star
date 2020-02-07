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
