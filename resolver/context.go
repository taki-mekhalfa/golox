package resolver

type functionCtx int

const (
	none functionCtx = iota
	function
	initializer
)
