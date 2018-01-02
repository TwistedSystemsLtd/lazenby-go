package lazendata

type Lazenfile struct {
	Lazenkeys map[string]string
	Secrets   map[string]string
}

type RevealedSecret struct {
	Name string
	Value string
}