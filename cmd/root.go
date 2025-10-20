package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	binaryName            = "gomanager"
	gomanagerConfigDirEnv = "gomanager_CONFIG_DIR"
	gomanagerDir          = "gomanager"
	gomanagerStorage      = "storage.json"
)

var logLevels = []string{"error", "warn", "info", "debug"}

var rootOptions struct {
	logLevel    string
	configDir   string
	storagePath string
}

var rootCmd = &cobra.Command{
	Use:   binaryName,
	Args:  cobra.NoArgs,
	Short: "CLI to manage go binaries",
	Long:  `CLI to manage go binaries`,
}

func Execute(version string) {
	rootCmd.Version = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initLogging, getConfigDir)

	rootCmd.PersistentFlags().StringVarP(
		&rootOptions.logLevel,
		"log",
		"l",
		"error",
		"Choose log level: "+strings.Join(logLevels, ", "),
	)

	cobra.CheckErr(rootCmd.RegisterFlagCompletionFunc(
		"log",
		cobra.FixedCompletions(logLevels, cobra.ShellCompDirectiveDefault),
	))
}

func initLogging() {
	level := slog.LevelInfo
	switch rootOptions.logLevel {
	case "error":
		level = slog.LevelError.Level()
	case "warn":
		level = slog.LevelWarn.Level()
	case "info":
		level = slog.LevelInfo.Level()
	case "debug":
		level = slog.LevelDebug.Level()
	default:
		slog.Error("invalid log level", "level", rootOptions.logLevel)
		os.Exit(1)
	}

	slog.SetLogLoggerLevel(level)
}

func getConfigDir() {
	rootOptions.configDir = os.Getenv(gomanagerConfigDirEnv)
	if rootOptions.configDir == "" {
		userDir, err := os.UserConfigDir()
		if err != nil {
			slog.Error("failed to get user config dir", "error", err)
			os.Exit(1)
		}

		rootOptions.configDir = filepath.Join(userDir, gomanagerDir)
	}

	// ensure the config dir exists
	slog.Info("using config dir", "path", rootOptions.configDir)
	err := os.MkdirAll(rootOptions.configDir, 0o755)
	if err != nil {
		slog.Error("failed to create config dir", "error", err, "path", rootOptions.configDir)
		os.Exit(1)
	}

	rootOptions.storagePath = filepath.Join(rootOptions.configDir, gomanagerStorage)
	slog.Info("using storage file", "path", rootOptions.storagePath)
}

func goBinPath() (string, error) {
	gobin := os.Getenv("GOBIN")
	if gobin != "" {
		return gobin, nil
	}

	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return fmt.Sprintf("%s/bin", gopath), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/go/bin", home), nil
}
