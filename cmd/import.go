package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/gomanager/internal/pkg"
	"github.com/tcondeixa/gomanager/internal/storage"
)

var importOptions struct {
	filePath string
}

var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "import installed packages from file",
	Long:    `import installed packages from file`,
	Example: fmt.Sprintf("  %s import -f /tmp/%s", binaryName, defaultExportFileName),
	Args:    cobra.NoArgs,
	RunE:    runimport,
}

func init() {
	rootCmd.AddCommand(importCmd)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	importCmd.Flags().StringVarP(
		&importOptions.filePath,
		"file",
		"f",
		filepath.Join(home, defaultExportFileName),
		"filepath to import list of installed packages",
	)
}

func runimport(_ *cobra.Command, _ []string) error {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Import(importOptions.filePath)
	if err != nil {
		return fmt.Errorf("failed to import storage: %w", err)
	}

	items := db.GetAllItems()
	for _, item := range items {
		slog.Info("Install package", "package", item.URI, "current_version", item.Version)
		output, err := item.Install()
		if err != nil {
			return fmt.Errorf("failed to install package %s: %v", item, err)
		}

		fmt.Println(rootOptions.colorScheme.Text(output))

		err = db.SaveItem(item.ID(), item)
		if err != nil {
			return fmt.Errorf("failed to save installed package %s: %v", item, err)
		}

		fmt.Println("Installed package: ", item.Name)
	}

	return nil
}
