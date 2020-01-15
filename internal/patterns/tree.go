package patterns

type slash int

const (
	maySlash slash = iota
	mustSlash
	neverSlash
)

type tree struct {
	id       string
	slash    slash
	params   params
	children []*child
}

type child struct {
	part part
	tree *tree
}

func (t *tree) addPattern(p *pattern, depth int) {
	if len(p.prefix) > depth {
		for _, part := range p.prefix {
			c := &child{
				part: part,
				tree: &tree{},
			}
			c.tree.addPattern(p, depth+1)
			t.children = append(t.children, c)
		}
		return
	}

	if p.suffix == nil {
		t.id = p.id
		t.slash = p.slash
		t.params = p.params
		return
	}

	c := &child{
		part: p.suffix,
		tree: &tree{
			id:     p.id,
			slash:  p.slash,
			params: p.params,
		},
	}
	t.children = append(t.children, c)
}