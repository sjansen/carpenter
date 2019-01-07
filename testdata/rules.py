def munge(x):
    if x == 'plug':
        return 'A'
    elif x == 'xyzzy':
        return 'B'
    return 'C'

metavariable = ('plug|xyzzy|thud', munge)

register_urls({
    id: "views.foo",
    parts: ['foo'],
    slash: 'always',
    tests: {
        '/foo/': '/foo/',
    },
}, {
    id: "views.bar",
    parts: ['bar', 'baz'],
    slash: 'never',
    tests: {
        '/bar/baz': '/bar/baz',
    },
}, {
    id: "views.qux",
    parts: ['qux', 'quux'],
    slash: 'strip',
    tests: {
        '/qux/quux':  '/qux/quux',
        '/qux/quux/': '/qux/quux',
    },
}, {
    id: "views.corge",
    parts: ['corge', ('g.....', 'X')],
    slash: 'strip',
    tests: {
        '/corge/grault/': '/corge/X',
        '/corge/garply/': '/corge/X',
    },
}, {
    id: "views.waldo",
    parts: [('waldo|fred', 'NAME'), metavariable],
    slash: 'always',
    tests: {
        '/waldo/plug/': '/NAME/A',
        '/fred/xyzzy/': '/NAME/B',
    },
})
