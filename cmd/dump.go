package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/goinstall/internal/pkg"
	"github.com/tcondeixa/goinstall/internal/storage"
)

const defaultDumpFileName = "goinstall.json"

var dumpOptions struct {
	filePath string
}

var dumpCmd = &cobra.Command{
	Use:     "dump",
	Short:   "Dump installed packages to file",
	Long:    `Dump installed packages to file`,
	Example: fmt.Sprintf("  %s dump -f /tmp/goinstall.json", binaryName),
	Args:    cobra.NoArgs,
	RunE:    runDump,
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	dumpCmd.Flags().StringVarP(
		&dumpOptions.filePath,
		"file",
		"f",
		filepath.Join(home, defaultDumpFileName),
		"filepath to dump installed packages",
	)
}

func runDump(_ *cobra.Command, _ []string) error {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Load()
	if err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	err = db.Dump(dumpOptions.filePath)
	if err != nil {
		return fmt.Errorf("failed to dump installed packages: %w", err)
	}

	fmt.Println("Installed packages dumped to: ", dumpOptions.filePath)

	return nil
}
