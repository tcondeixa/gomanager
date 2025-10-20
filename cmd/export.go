package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/gomanager/internal/pkg"
	"github.com/tcondeixa/gomanager/internal/storage"
)

const defaultExportFileName = binaryName + ".json"

var exportOptions struct {
	filePath string
}

var exportCmd = &cobra.Command{
	Use:     "export",
	Short:   "Export installed packages to file",
	Long:    `Export installed packages to file`,
	Example: fmt.Sprintf("  %s export -f /tmp/%s", binaryName, defaultExportFileName),
	Args:    cobra.NoArgs,
	RunE:    runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	exportCmd.Flags().StringVarP(
		&exportOptions.filePath,
		"file",
		"f",
		filepath.Join(home, defaultExportFileName),
		"filepath to export list of installed packages",
	)
}

func runExport(_ *cobra.Command, _ []string) error {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Load()
	if err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	err = db.Export(exportOptions.filePath)
	if err != nil {
		return fmt.Errorf("failed to export installed packages: %w", err)
	}

	fmt.Println("Installed packages exported to: ", exportOptions.filePath)

	return nil
}
