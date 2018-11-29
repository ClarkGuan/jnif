package gob

import "github.com/ClarkGuan/jnif/interp"

func init() {
	interp.Register("go", new(goTransform))
}
