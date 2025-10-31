package cmd

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/spf13/cobra"
	"github.com/tcondeixa/gomanager/internal/color"
)

//go:embed images/gopher.jpeg
var image []byte

const (
	binaryName              = "gomanager"
	configDirEnv            = "GOMANAGER_CONFIG_DIR"
	configDir               = "gomanager"
	storageFile             = "storage.json"
	colorSchemeEnv          = "GOMANAGER_COLOR_SCHEME"
	defaultTextKey          = "tx"
	defaultTextKeyWithSep   = defaultTextKey + ":"
	defaultHeaderKey        = "hd"
	defaultHeaderKeyWithSep = defaultHeaderKey + ":"
	defaultErKey            = "er"
	defaultErrKeyWithSep    = defaultErKey + ":"
	defaultTextColor        = "#f5e0dc"
	defaultHeaderColor      = "#cba6f7"
	defaultErrColor         = "#f38ba8"
)

var (
	logLevels           = []string{"error", "warn", "info", "debug"}
	defaultColorOptions = []string{
		defaultTextKeyWithSep + defaultTextColor,
		defaultHeaderKeyWithSep + defaultHeaderColor,
		defaultErrKeyWithSep + defaultErrColor,
	}
	defaultColorScheme = strings.Join(defaultColorOptions, ",")
)

var rootOptions struct {
	logLevel    string
	configDir   string
	storagePath string
	noColor     bool
	colorScheme color.Scheme
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
	cobra.OnInitialize(initLogging, getConfigDir, colorScheme)

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

	rootCmd.PersistentFlags().BoolVar(
		&rootOptions.noColor,
		"no-color",
		false,
		"output with colors",
	)

	printLogo()
}

func checkTerminalImageSupport() string {
	// Check for iTerm2 support
	term := os.Getenv("TERM_PROGRAM")
	if term == "iTerm.app" || term == "WezTerm" {
		return "iterm2"
	}

	// Check for Kitty support
	if os.Getenv("KITTY_WINDOW_ID") != "" ||
		os.Getenv("KITTY_PID") != "" ||
		strings.Contains(os.Getenv("TERM"), "kitty") {
		return "kitty"
	}

	return ""
}

func printLogo() {
	term_img_protocol := checkTerminalImageSupport()
	encodedString := base64.StdEncoding.EncodeToString(image)
	switch term_img_protocol {
	case "iterm2":
		fmt.Printf("\033]1337;File=inline=2;width=10%%;preserveAspectRatio=1:%s\a\n", encodedString)
	case "kitty":
		fmt.Printf("\033_Ga=T,f=100,s=100,v=100,m=0;%s\033\\", encodedString)
	}
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
	rootOptions.configDir = os.Getenv(configDirEnv)
	if rootOptions.configDir == "" {
		userDir, err := os.UserConfigDir()
		if err != nil {
			slog.Error("failed to get user config dir", "error", err)
			os.Exit(1)
		}

		rootOptions.configDir = filepath.Join(userDir, configDir)
	}

	// ensure the config dir exists
	slog.Info("using config dir", "path", rootOptions.configDir)
	err := os.MkdirAll(rootOptions.configDir, 0o755)
	if err != nil {
		slog.Error("failed to create config dir", "error", err, "path", rootOptions.configDir)
		os.Exit(1)
	}

	rootOptions.storagePath = filepath.Join(rootOptions.configDir, storageFile)
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

func colorScheme() {
	colorSchemeText := os.Getenv(colorSchemeEnv)
	if colorSchemeText == "" {
		colorSchemeText = defaultColorScheme
	}

	slog.Info("Color option", "no-color", rootOptions.noColor)
	textColor := defaultTextColor
	headerColor := defaultHeaderColor
	errColor := defaultErrColor
	if rootOptions.noColor {
		rootOptions.colorScheme = *color.NewScheme(rootOptions.noColor, textColor, headerColor, errColor)
	}

	// format expected "tx:#f5e0dc,hd:#cba6f7,er:#f38ba8"
	colorsIter := strings.SplitSeq(strings.TrimSpace(colorSchemeText), ",")
	for c := range colorsIter {
		switch {
		case strings.HasPrefix(c, defaultTextKeyWithSep):
			textColor = strings.TrimPrefix(c, defaultTextKeyWithSep)
			_, err := color.HexToRGB(textColor)
			if err != nil {
				textColor = defaultTextColor
			}
		case strings.HasPrefix(c, defaultHeaderKeyWithSep):
			headerColor = strings.TrimPrefix(c, defaultHeaderKeyWithSep)
			_, err := color.HexToRGB(headerColor)
			if err != nil {
				headerColor = defaultHeaderColor
			}
		case strings.HasPrefix(c, defaultErrKeyWithSep):
			errColor = strings.TrimPrefix(c, defaultErrKeyWithSep)
			_, err := color.HexToRGB(errColor)
			if err != nil {
				errColor = defaultErrColor
			}
		}
	}

	if !rootOptions.noColor {
		slog.Info("color scheme", "text", textColor, "header", headerColor, "error", errColor)
	}
	rootOptions.colorScheme = *color.NewScheme(rootOptions.noColor, textColor, headerColor, errColor)
}
