package main

import (
	"context"
	"fmt"
	"os"

	"github.com/senforsce/tndr/storybook/example"
)

func main() {
	s := example.Storybook()
	if err := s.ListenAndServeWithContext(context.Background()); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
