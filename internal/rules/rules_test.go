package rules

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var script = `
register_urls({
    'id': 'views.foo',
    'parts': ['foo'],
    'slash': 'always',
}, {
    'id': 'views.bar',
    'parts': ['bar'],
    'slash': 'never',
}, {
    'id': 'views.baz',
    'parts': ['baz'],
    'slash': 'strip',
})
`
var expected = &Rules{Matchers: []Matcher{
	&PlainPart{
		id:    "views.foo",
		slash: "always",
		value: "foo",
	},
	&PlainPart{
		id:    "views.bar",
		slash: "never",
		value: "bar",
	},
	&PlainPart{
		id:    "views.baz",
		slash: "strip",
		value: "baz",
	},
}}

func TestLoad(t *testing.T) {
	require := require.New(t)

	r := strings.NewReader(script)
	actual, err := Load("<test script>", r)
	require.NoError(err)
	require.Equal(expected, actual)
}
