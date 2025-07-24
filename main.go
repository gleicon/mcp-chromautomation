package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gleicon/mcp-chromautomation/internal/cli"
	"github.com/gleicon/mcp-chromautomation/internal/server"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "mcp-chromautomation",
		Short:   "Chrome automation MCP service with beautiful CLI",
		Long:    `A Model Context Protocol service for Chrome browser automation with local data storage and a modern CLI interface.`,
		Version: fmt.Sprintf("%s (%s)", version, commit),
	}

	// Server command
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start the MCP server",
		Long:  "Start the MCP server to provide Chrome automation tools to clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			srv := server.NewEnhanced()
			return srv.Start(ctx)
		},
	}

	// CLI command
	cliCmd := &cobra.Command{
		Use:   "ui",
		Short: "Launch interactive CLI",
		Long:  "Launch the interactive CLI interface for Chrome automation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.Start()
		},
	}

	// Add commands
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(cliCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}