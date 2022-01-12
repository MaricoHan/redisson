package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitlab.bianjie.ai/irita-nftp/nftp-open-api/config"
	"gitlab.bianjie.ai/irita-nftp/nftp-open-api/internal/app"
)

var (
	defaultCLIHome = os.ExpandEnv("$HOME/.nftp")
	flagHome       = "home"
)

func main() {
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:   "nftp",
		Short: "nftp Daemon (server)",
	}
	rootCmd.AddCommand(StartCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed executing mkr command: %s, exiting...\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// StartCmd return the start command
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Example: "nftp start",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Load(cmd, flagHome); err != nil {
				return err
			}
			app.Start()
			return nil
		},
	}
	cmd.Flags().String(flagHome, defaultCLIHome, "nftp server config path")
	return cmd
}
