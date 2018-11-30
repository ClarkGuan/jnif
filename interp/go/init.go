package gob

import "github.com/ClarkGuan/jnif/interp"

func init() {
	interp.Register("lang:go", new(goTransform))
}
