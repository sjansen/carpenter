names = ('waldo|fred', 'plugh')

def fn(x):
    if x == 'xyzzy':
        return 'thud'
    return 'X'

register_urls({
    'id': 'views.root',
    'path': {
        'parts': [],
        'slash': 'always',
    },
    'query': {
        'dedup': 'never',
        'params': {},
    },
    'tests': {
        '/':      '/',
        '/Spoon!': None,
    },
}, {
    'id': 'views.always',
    'path': {
        'parts': ['foo'],
        'slash': 'always',
    },
    'query': {
        'dedup': 'never',
        'params': {},
    },
    'tests': {
        '/foo':   None,
        '/foo/': '/foo/',
    },
}, {
    'id': 'views.never',
    'path': {
        'parts': ['bar'],
        'slash': 'never',
    },
    'query': {
        'dedup': 'never',
        'params': {},
    },
    'tests': {
        '/bar': '/bar',
        '/bar/': None,
    },
}, {
    'id': 'views.strip',
    'path': {
        'parts': ['baz'],
        'slash': 'strip',
    },
    'query': {
        'dedup': 'never',
        'params': {},
    },
    'tests': {
        '/baz':  '/baz',
        '/baz/': '/baz',
    },
}, {
    'id': 'views.regex',
    'path': {
        'parts': [('qux', 'quux')],
        'slash': 'always',
    },
    'query': {
        'dedup': 'never',
        'params': {},
    },
    'tests': {
        '/qux/': '/quux/',
    },
}, {
    'id': 'query.never',
    'path': {
        'parts': ['search'],
        'slash': 'never',
    },
    'query': {
        'dedup': 'never',
        'params': {
            'q':   'X',
            'utf8': None,
        },
    },
    'tests': {
        '/search?q=cats':         '/search?q=X',
        '/search?utf8=✔':         '/search',
        '/search?q=dogs&utf8=✔':  '/search?q=X',
    },
}, {
    'id': 'views.multi',
    'path': {
        'parts': [
            'corge',
            ('grault', 'garply'),
            names,
            ('.+', fn),
        ],
        'slash': 'strip',
    },
    'query': {
        'dedup': 'never',
        'params': {
            'n': lambda x: "even" if len(x) % 2 == 0 else "odd",
        },
    },
    'tests': {
        '/corge/grault/waldo/xyzzy':         '/corge/garply/plugh/thud',
        '/corge/grault/fred/42/':            '/corge/garply/plugh/X',
        '/corge/grault/fred/random/?n=left': '/corge/garply/plugh/X?n=even',
			  '/corge/grault/fred/random?n=right': '/corge/garply/plugh/X?n=odd',
    },
})
