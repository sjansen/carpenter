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
    'tests': {
        '/qux/': '/quux/',
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
        'slash': 'always',
    },
    'tests': {
        '/corge/grault/waldo/xyzzy/': '/corge/garply/plugh/thud/',
        '/corge/grault/fred/42/':     '/corge/garply/plugh/X/',
    },
})
