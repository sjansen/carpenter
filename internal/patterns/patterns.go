package patterns

type Patterns struct {
	tests map[string]match
	tree  tree
}

type match struct {
	id  string
	url string
}

type pattern struct {
	id     string
	slash  slash
	prefix []part
	suffix *regexPart
	params params
	tests  map[string]string
}
