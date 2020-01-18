add_url("example",
    path={
        "prefix": [],
        "suffix": "/?",
    },
    query={
        "params": {
            "x": lambda n: chr(n),
        },
    },
    tests={
        "/?x=42": "/?x=*",
    },
)
