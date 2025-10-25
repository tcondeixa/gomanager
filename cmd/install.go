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

var installOptions struct {
	name string
}

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install packages",
	Long:    `Install packages`,
	Example: fmt.Sprintf("  %s install github.com/tcondeixa/gomanager@latest", binaryName),
	RunE:    runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(
		&installOptions.name,
		"name",
		"n",
		"",
		"Force name of the binary (default to go install name)",
	)
}

func runInstall(_ *cobra.Command, args []string) error {
	if len(args) > 1 && installOptions.name != "" {
		return fmt.Errorf("cannot use --name when installing multiple packages")
	}

	db := storage.New[pkg.Package](rootOptions.storagePath)
	err := db.Start()
	if err != nil {
		return err
	}

	path, err := goBinPath()
	if err != nil {
		return fmt.Errorf("failed to determine go bin path: %v", err)
	}

	// install and save to storage
	for _, item := range args {
		pack, err := pkg.New(item)
		if err != nil {
			return fmt.Errorf("failed to create package from %s: %v", item, err)
		}

		if installOptions.name != "" {
			oldPath := filepath.Join(path, pack.Name)
			exists, err := fileExists(oldPath)
			if err != nil {
				return fmt.Errorf("failed to check if file exists: %v", err)
			}

			if exists {
				tmpPath := oldPath + ".gomanager"
				err = os.Rename(oldPath, tmpPath)
				if err != nil {
					return fmt.Errorf("failed to rename existing binary: %v", err)
				}

				defer func() {
					err = os.Rename(oldPath+".gomanager", oldPath)
					if err != nil {
						slog.Error("failed to restore original binary", "error", err)
					}
				}()
			}
		}

		output, err := pack.Install()
		if err != nil {
			return fmt.Errorf("failed to install package %s: %v", item, err)
		}

		fmt.Println(rootOptions.colorScheme.Text(output))

		if installOptions.name != "" {
			oldPath := filepath.Join(path, pack.Name)
			newPath := filepath.Join(path, installOptions.name)
			err = os.Rename(oldPath, newPath)
			if err != nil {
				return fmt.Errorf("failed to rename binary to %s: %v", installOptions.name, err)
			}

			pack.Name = installOptions.name
		}

		err = db.SaveItem(pack.ID(), *pack)
		if err != nil {
			return err
		}

		fmt.Println(rootOptions.colorScheme.Header("Installed package: " + pack.Name))
	}

	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
