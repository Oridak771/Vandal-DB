package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
Use:   "vandal",
Short: "Vandal-DB CLI",
Long:  `A command-line interface for managing Vandal-DB resources.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
