package gqlclient

// Anonymous ...
type Anonymous struct {
	cli *GQLClient
}

// Anonymous ...
func (g *GQLClient) Anonymous() *Anonymous {
	return &Anonymous{
		cli: g,
	}
}
