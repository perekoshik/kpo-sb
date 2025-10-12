package main

import (
	"context"
	"fmt"

	"github.com/perekoshik/zoo/internal/infra/cli"
	"github.com/perekoshik/zoo/internal/infra/di"
)

func main() {
	c, err := di.Build()
	if err != nil {
		panic(err)
	}
	err = c.Invoke(func(r *cli.Router) error {
		return r.Run(context.Background())
	})
	if err != nil {
		fmt.Println("fatal:", err)
	}
}
