names = ('waldo|fred', 'plugh')

def fn(x):
    if x == 'xyzzy':
	return 'thud'
    return 'X'

register_urls({
    'id': 'views.always',
    'parts': ['foo'],
    'slash': 'always',
    'tests': {
        '/foo':   None,
        '/foo/': '/foo/',
    },
}, {
    'id': 'views.never',
    'parts': ['bar'],
    'slash': 'never',
    'tests': {
        '/bar': '/bar',
        '/bar/': None,
    },
}, {
    'id': 'views.strip',
    'parts': ['baz'],
    'slash': 'strip',
    'tests': {
        '/baz':  '/baz',
        '/baz/': '/baz',
    },
}, {
    'id': 'views.regex',
    'parts': [('qux', 'quux')],
    'slash': 'always',
    'tests': {
        '/qux/': '/quux/',
    },
}, {
    'id': 'views.multi',
    'parts': [
	'corge',
	('grault', 'garply'),
	names,
	('.+', fn),
    ],
    'slash': 'always',
    'tests': {
        '/corge/grault/waldo/xyzzy/': '/corge/garply/plugh/thud/',
        '/corge/grault/fred/42/':     '/corge/garply/plugh/X/',
    },
})
