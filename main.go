package main

import "github.com/All-The-Vibes/ATV-StarterKit/cmd"

// Set via ldflags by goreleaser
var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.SetVersion(version, commit)
	cmd.Execute()
}
