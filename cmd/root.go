package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var rootOptions struct {
	logLevel string
}

var rootCmd = &cobra.Command{
	Use:   "goinstall",
	Args:  cobra.NoArgs,
	Short: "CLI to manage go install binaries",
	Long:  `CLI to manage go install binaries`,
}

func Execute(version string) {
	rootCmd.Version = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVar(
		&rootOptions.logLevel,
		"log-level",
		"error",
		"Choose log level [panic,fatal,error,warn,info,debug,trace]",
	)
}

func initLogging() {
	level, err := log.ParseLevel(rootOptions.logLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}

	if level == level == log.DebugLevel || level == log.TraceLevel {
		log.SetReportCaller(true)
	} else {
		log.SetReportCaller(false)
	}

	log.SetLevel(level)
}
