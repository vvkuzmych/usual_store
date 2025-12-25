package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "usual-store-cli",
	Short: "Usual Store Admin CLI",
	Long: `A command-line interface for managing the Usual Store application.

Provides tools for:
  - User management (create, list, reset password)
  - Database operations (migrate, seed, backup)
  - Testing utilities (generate test data)
  - Health checks and monitoring
  - Stripe integration testing`,
	Version: version,
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(healthCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
