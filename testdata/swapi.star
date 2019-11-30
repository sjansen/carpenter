resources = [
    "films",
    "people",
    "planets",
    "species",
    "starships",
    "vehicles",
]

def generate_urls():
    urls = [{
        "id": "root",
        "path": {
            "parts": ["api"],
            "slash": "always",
        },
        "query": {
            "dedup": "never",
            "params": {
                "format": wookiee,
            },
        },
        "tests": {
            "/": None,
            "/api": None,
            "/api/": "/api/",
        },
    }]

    for x in resources:
        urls.append({
            "id": "/%s/" % x,
            "path": {
                "parts": ["api", x],
                "slash": "always",
            },
            "query": {
                "dedup": "never",
                "params": {
                    "format": wookiee,
                    "search": "X",
                },
            },
            "tests": {
                "/api/%s/" % x: "/api/%s/" % x,
                "/api/%s/?search=resistance" % x: "/api/%s/?search=X" % x,
            },
        })

    for x in resources:
        urls.append({
            "id": "/%s/:id/" % x,
            "path": {
                "parts": ["api", x, ("[1-9][0-9]*", "ID")],
                "slash": "always",
            },
            "query": {
                "dedup": "never",
                "params": {
                    "format": wookiee,
                },
            },
            "tests": {
                "/api/%s/1/" % x: "/api/%s/ID/" % x,
                "/api/%s/1/?format=csv" % x: "/api/%s/ID/?format=INVALID" % x,
                "/api/%s/1/?format=wookiee" % x: "/api/%s/ID/?format=wookiee" % x,
            },
        })

    return urls

def wookiee(x):
    if x == "wookiee":
        return "wookiee"
    return "INVALID"

register_urls(
    *generate_urls()
)
