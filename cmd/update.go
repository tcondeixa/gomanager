package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/gomanager/internal/pkg"
	"github.com/tcondeixa/gomanager/internal/storage"
)

var updateOptions struct {
	name           string
	forceNonLatest bool
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update packages",
	Long:    `Update packages`,
	Example: fmt.Sprintf("  %s update --name %s", binaryName, binaryName),
	RunE:    runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(
		&updateOptions.name,
		"name",
		"n",
		"",
		"name to be updated (default all packages)",
	)
	cobra.CheckErr(updateCmd.RegisterFlagCompletionFunc(
		"name",
		func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			return installedPackagesCompletion(cmd, args, toComplete)
		},
	))

	updateCmd.Flags().BoolVarP(
		&updateOptions.forceNonLatest,
		"force",
		"f",
		false,
		"force also non-latest versions",
	)
}

func runUpdate(_ *cobra.Command, _ []string) error {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Start()
	if err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	if updateOptions.name != "" {
		item, found := db.GetItem(updateOptions.name)
		if !found {
			return fmt.Errorf("package %s not found in storage", updateOptions.name)
		}

		item.UpdateVersion("latest")
		output, err := item.Install()
		if err != nil {
			return fmt.Errorf("failed to install package %s: %v", item, err)
		}

		fmt.Println(rootOptions.colorScheme.Text(output))

		item.UpdateVersion("latest")
		err = db.SaveItem(item.ID(), item)
		if err != nil {
			return fmt.Errorf("failed to save updated package %s: %v", item, err)
		}

		fmt.Println("Package updated successfully")
	}

	items := db.GetAllItems()
	for _, item := range items {
		if item.Version == "latest" || updateOptions.forceNonLatest {
			slog.Info("Updating package", "package", item.URI, "current_version", item.Version)
			item.UpdateVersion("latest")
			output, err := item.Install()
			if err != nil {
				return fmt.Errorf("failed to install package %s: %v", item, err)
			}

			fmt.Println(rootOptions.colorScheme.Text(output))

			err = db.SaveItem(item.ID(), item)
			if err != nil {
				return fmt.Errorf("failed to save updated package %s: %v", item, err)
			}

			fmt.Println(rootOptions.colorScheme.Text("Package " + item.Name + " updated successfully"))
		}
	}

	return nil
}
