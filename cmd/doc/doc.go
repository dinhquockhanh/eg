package main

import (
	"log"

	"github.com/spf13/cobra/doc"

	"github.com/dinhquockhanh/eg/cmd"
)

func main() {
	cmd.RootCmd().DisableAutoGenTag = true
	if err := doc.GenMarkdownTree(cmd.RootCmd(), "."); err != nil {
		log.Fatal(err)
	}
}
