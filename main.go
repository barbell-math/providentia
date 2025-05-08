package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/barbell-math/providentia/state"
)

func main() {
	globalCtxt, err := state.Parse(context.Background(), os.Args[1:])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println(state.FromContext(globalCtxt))

	s, ok := state.FromContext(globalCtxt)
	if !ok {
		fmt.Println("EHH?")
	}
	s.Log.Info("We are up and running!")
}
