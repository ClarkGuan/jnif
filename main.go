package main

import (
	"log"
	"os"

	"github.com/ClarkGuan/jnif/interp"
	_ "github.com/ClarkGuan/jnif/interp/go"
)

func main() {
	if err := interp.Print(os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}
