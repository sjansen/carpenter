resources = [
    "films",
    "people",
    "planets",
    "species",
    "starships",
    "vehicles",
]

def wookiee(x):
    if x == "wookiee":
        return "wookiee"
    return "INVALID"

def add_urls():
    add_url("root",
        path={
            "prefix": ["api"],
            "suffix": "/",
        },
        query={
            "dedup": "never",
            "params": {
                "format": wookiee,
            },
        },
        tests={
            "/": None,
            "/api": None,
            "/api/": "/api/",
        },
    )

    for x in resources:
        add_url("/%s/" % x,
            path={
                "prefix": ["api", x],
                "suffix": "/",
            },
            query={
                "dedup": "never",
                "params": {
                    "format": wookiee,
                    "search": "X",
                },
            },
            tests={
                "/api/%s/" % x: "/api/%s/" % x,
                "/api/%s/?search=resistance" % x: "/api/%s/?search=X" % x,
            },
        )

    for x in resources:
        add_url("/%s/:id/" % x,
            path={
                "prefix": ["api", x, ("[1-9][0-9]*", "ID")],
                "suffix": "/",
            },
            query={
                "dedup": "never",
                "params": {
                    "format": wookiee,
                },
            },
            tests={
                "/api/%s/1/" % x: "/api/%s/ID/" % x,
                "/api/%s/1/?format=csv" % x: "/api/%s/ID/?format=INVALID" % x,
                "/api/%s/1/?format=wookiee" % x: "/api/%s/ID/?format=wookiee" % x,
            },
        )

add_urls()
