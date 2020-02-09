url(
    "example",
    path = {
        "prefix": [],
        "suffix": "/",
    },
    query = {
        "match": {
            "x": lambda k, v: chr(v),
        },
    },
    tests = {
        "/?x=42": "/?x=*",
    },
)
