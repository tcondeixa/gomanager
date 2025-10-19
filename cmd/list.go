package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/goinstall/internal/pkg"
	"github.com/tcondeixa/goinstall/internal/storage"
)

var availableOutputs = []string{"text", "json"}

var listOptions struct {
	outputFormat string
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List installed packages",
	Long:    `List installed packages`,
	Example: fmt.Sprintf("  %s list -o json", binaryName),
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(
		&listOptions.outputFormat,
		"output",
		"o",
		"text",
		"Output format: "+strings.Join(availableOutputs, ", "),
	)
	cobra.CheckErr(listCmd.RegisterFlagCompletionFunc(
		"output",
		cobra.FixedCompletions(availableOutputs, cobra.ShellCompDirectiveDefault),
	))
}

func runList(_ *cobra.Command, _ []string) error {
	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Load()
	if err != nil {
		return fmt.Errorf("failed to load storage: %w", err)
	}

	items := db.GetAllItems()
	switch listOptions.outputFormat {
	case "json":
		return printAsJSON(items)
	case "text":
		return printAsText(items)
	}

	return nil
}

func printAsJSON(items map[string]pkg.Package) error {
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode items to json: %w", err)
	}

	fmt.Println(string(bytes))

	return nil
}

func printAsText(items map[string]pkg.Package) error {
	if len(items) == 0 {
		fmt.Println("No installed packages found.")
		return nil
	}

	fmt.Println("Installed Packages:")
	fmt.Println("-------------------")
	for _, item := range items {
		fmt.Println(item.String())
		fmt.Println("-------------------")
	}

	return nil
}
