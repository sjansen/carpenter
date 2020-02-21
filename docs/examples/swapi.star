resources = [
    "films",
    "people",
    "planets",
    "species",
    "starships",
    "vehicles",
]

def wookiee(k, v):
    if v == "wookiee":
        return "wookiee"
    return "INVALID"

def urls():
    url(
        "root",
        path = {
            "prefix": ["api"],
            "suffix": "/",
        },
        query = {
            "dedup": "never",
            "match": {
                "format": wookiee,
            },
        },
        tests = {
            "/": None,
            "/api": None,
            "/api/": "/api/",
        },
    )

    for resource in resources:
        url(
            "/%s/" % resource,
            path = {
                "prefix": ["api", resource],
                "suffix": "/",
            },
            query = {
                "dedup": "never",
                "match": {
                    "format": wookiee,
                    "search": "X",
                },
            },
            tests = {
                "/api/%s/" % resource: "/api/%s/" % resource,
                "/api/%s/?search=resistance" % resource: "/api/%s/?search=X" % resource,
            },
        )

    for resource in resources:
        url(
            "/%s/:id/" % resource,
            path = {
                "prefix": ["api", resource, ("[1-9][0-9]*", "ID")],
                "suffix": "/",
            },
            query = {
                "dedup": "never",
                "match": {
                    "format": wookiee,
                },
            },
            tests = {
                "/api/%s/1/" % resource: "/api/%s/ID/" % resource,
                "/api/%s/1/?format=csv" % resource: "/api/%s/ID/?format=INVALID" % resource,
                "/api/%s/1/?format=wookiee" % resource: "/api/%s/ID/?format=wookiee" % resource,
            },
        )

urls()
