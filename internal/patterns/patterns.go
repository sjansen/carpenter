package patterns

type Patterns struct {
	tree
}

type pattern struct {
	id     string
	slash  slash
	prefix []part
	suffix *regexPart
	params params
	tests  map[string]string
}
