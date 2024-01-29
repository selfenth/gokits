package cmder

import (
	"context"

	"github.com/xdeeper/gokits/cmder/urfave"
)

func MustStart(root any) {
	err := urfave.Build(root).Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type Runable interface {
	Run(ctx context.Context) error
}
