url(
    "example",
    path = {
        "prefix": [],
        "suffix": "/",
    },
    query = {
        "match": {
            "x": lambda n: chr(n),
        },
    },
    tests = {
        "/?x=42": "/?x=*",
    },
)
