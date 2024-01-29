package main

import (
	"github.com/xdeeper/gokits/cmder"
)

func main() {
	cmder.MustStart(Gokits{
		Description: "golang develop kits",
		Clone: Clone{
			Description: "clone golang project into $GOPATH/src",
		},
	})
}

type Gokits struct {
	Description string

	Clone Clone `cmd:"clone"`
}
