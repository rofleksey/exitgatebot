package main

import (
	_ "embed"
	"exitgatebot/app/cmd"
	"exitgatebot/app/util"
	"exitgatebot/app/util/mylog"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.szostok.io/version/extension"
)

func main() {
	mylog.Preinit()

	fmt.Fprintln(os.Stderr, util.Banner)

	rootCmd := &cobra.Command{Use: "exitgatebot"}
	rootCmd.AddCommand(cmd.Run)
	rootCmd.AddCommand(extension.NewVersionCobraCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
		return
	}
}
