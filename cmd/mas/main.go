// Package main provides the mas CLI for multi-agent-spec operations.
//
// Usage:
//
//	mas <command> [flags] [args]
//
// Commands:
//
//	render    Render TeamReport JSON to box or narrative format
//	version   Print version information
package main

import (
	"os"

	"github.com/plexusone/multi-agent-spec/cmd/mas/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
