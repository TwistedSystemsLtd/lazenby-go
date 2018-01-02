package toml

type LazenToml struct {
	lazenkeys map[string]string
	secrets [][]byte
}

type Keypair struct {
	publickey []byte
	lazenkey []byte
}

type RevealedSecret struct {
	name  string
	value string
}
