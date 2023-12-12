package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Short: "simulation service",
}

func main() {
	rootCmd.Execute()
}
