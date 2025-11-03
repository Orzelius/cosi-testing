package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Orzelius/cosi-testing/backend"
	mylog "github.com/Orzelius/cosi-testing/log"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var manifestsPath string
var backendName string
var logLevel string

func main() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&manifestsPath, "file", "", "path to the kubernetes manifests file to apply/diff")
	rootCmd.PersistentFlags().StringVar(&backendName, "backend", "", "logic to use to ('kubernetes' or 'ssa')")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level ('info' or 'debug')")

	cobra.MarkFlagRequired(rootCmd.PersistentFlags(), "file")
	cobra.MarkFlagRequired(rootCmd.PersistentFlags(), "backend")

	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(applyCmd)
}

var rootCmd = &cobra.Command{
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			return err
		}

		mylog.GetLogger().SetLevel(level)

		return nil
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "diff the local manifests against remote state",
	RunE: func(cmd *cobra.Command, args []string) error {
		be, err := getInitializedBackend()
		if err != nil {
			return nil
		}

		data, err := readManifestFile()
		if err != nil {
			return nil
		}

		return be.Diff(cmd.Context(), data)
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply the local manifests (with prune)",
	RunE: func(cmd *cobra.Command, args []string) error {
		be, err := getInitializedBackend()
		if err != nil {
			return nil
		}

		data, err := readManifestFile()
		if err != nil {
			return nil
		}

		return be.Apply(cmd.Context(), data)
	},
}

func getInitializedBackend() (backend.Backend, error) {
	var be backend.Backend
	switch backendName {
	case "ssa":
		be = &backend.FluxSSA{}
	case "kubernetes":
		be = &backend.Kubernetes{}
	default:
		return nil, fmt.Errorf("unknown backend: %q", backendName)
	}

	err := be.Init()
	return be, err
}

func readManifestFile() ([]byte, error) {
	return os.ReadFile(manifestsPath)
}
