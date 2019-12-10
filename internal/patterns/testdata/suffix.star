def image(x):
    if len(x) < 1:
        return ""
    return x[len(x)-3:].upper()


register_urls({
    "id": "root",
    "path": {
        "prefix": [],
        "suffix": "/",
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/": "/",
        "/foo": None,
    },
}, {
    "id": "images",
    "path": {
        "prefix": ["images"],
        "suffix": ("(.gif|.jpg|.png)$", image)
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/images/": "/images",
        "/images/logo.png": "/images/PNG",
        "/images/photos/lake.jpg": "/images/JPG",
        "/images/emoji/smile.gif": "/images/GIF",
    },
}, {
    "id": "photos",
    "path": {
        "prefix": ["photos"],
        "suffix": (".[jJ][pP][eE]?[gG]$", "JPEG")
    },
    "query": {
        "dedup": "never",
        "params": {},
    },
    "tests": {
        "/photos/": "/photos",
        "/photos/sjansen/original.jpg": "/photos/JPEG",
        "/photos/sjansen/thumbnail.jpg": "/photos/JPEG",
    },
})
