package cmd

import (
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/gomanager/internal/pkg"
	"github.com/tcondeixa/gomanager/internal/storage"
)

var unistallCmd = &cobra.Command{
	Use:               "uninstall",
	Short:             "Uninstall packages",
	Long:              `Uninstall packages`,
	Example:           fmt.Sprintf("  %s uninstall %s", binaryName, binaryName),
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: installedPackagesCompletion,
	RunE:              runUninstall,
}

func installedPackagesCompletion(
	_ *cobra.Command,
	_ []string,
	_ string,
) ([]cobra.Completion, cobra.ShellCompDirective) {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Load()
	if err != nil {
		return []cobra.Completion{}, cobra.ShellCompDirectiveError
	}

	allItems := db.GetAllItems()
	return slices.Collect(maps.Keys(allItems)), cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(unistallCmd)
}

func runUninstall(_ *cobra.Command, args []string) error {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Load()
	if err != nil {
		return err
	}

	path, err := goBinPath()
	if err != nil {
		return fmt.Errorf("failed to determine uninstall path: %w", err)
	}

	slog.Info("Current golang bin dir", "path", path)

	for _, name := range args {
		item, exists := db.GetItem(name)
		if !exists {
			return fmt.Errorf("package %s not found in storage", name)
		}

		binPath := filepath.Join(path, item.Name)
		err := os.Remove(binPath)
		if err != nil {
			return fmt.Errorf("failed to remove binary at %s: %w", binPath, err)
		}
		slog.Info("Removed binary", "path", binPath)

		db.DeleteItem(item.ID())
		err = db.Save()
		if err != nil {
			return err
		}

		fmt.Println("Uninstalled package: ", item.Name)
	}

	return nil
}
