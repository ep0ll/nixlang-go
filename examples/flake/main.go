// Example: parse a flake reference (no lock, no network, no CLI).
//
// Usage:
//
//	go run ./examples/flake 'github:NixOS/nixpkgs'
//	go run ./examples/flake 'path:./.'
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ep0ll/nixlang-go/fetchers"
	"github.com/ep0ll/nixlang-go/flake"
	"github.com/ep0ll/nixlang-go/util"
)

func main() {
	refStr := "github:NixOS/nixpkgs"
	if len(os.Args) > 1 {
		refStr = os.Args[1]
	}

	ctx := util.NewContext()
	defer ctx.Close()
	must(util.Init(ctx))

	fs, err := fetchers.New(ctx)
	must(err)
	defer fs.Close()

	settings, err := flake.NewSettings(ctx)
	must(err)
	defer settings.Close()

	flags, err := flake.NewParseFlags(ctx, settings)
	must(err)
	defer flags.Close()

	ref, fragment, err := flake.ParseReference(ctx, fs, settings, flags, refStr)
	must(err)
	defer ref.Close()

	fmt.Println("parsed flake reference:", refStr)
	if fragment != "" {
		fmt.Println("fragment:", fragment)
	}
	fmt.Println("(parse only — no lock, fetch, or nix CLI)")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
