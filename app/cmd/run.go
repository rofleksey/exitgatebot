package cmd

import (
	"context"
	"exitgatebot/app/client/openai"
	"exitgatebot/app/client/steam"
	"exitgatebot/app/client/twitch_api"
	"exitgatebot/app/client/twitch_irc"
	"exitgatebot/app/config"
	"exitgatebot/app/service/checker"
	"exitgatebot/app/util/mylog"
	"log/slog"
	"os"
	"os/signal"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

var configPath string

var Run = &cobra.Command{
	Use:   "run",
	Short: "Run notifier",
	Run:   runNotifier,
}

func init() {
	Run.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to config yaml file (required)")
}

func runNotifier(_ *cobra.Command, _ []string) {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di := do.New()
	do.ProvideValue(di, appCtx)

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("Failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1) //nolint:gocritic
		return
	}
	do.ProvideValue(di, cfg)

	if err = mylog.Init(cfg); err != nil {
		slog.Error("Failed to init logging",
			slog.Any("error", err),
		)
		os.Exit(1)
		return
	}
	slog.InfoContext(appCtx, "Starting service...",
		slog.Bool("telegram", true),
	)

	do.Provide(di, twitch_api.NewClient)
	do.Provide(di, twitch_irc.NewClient)
	do.Provide(di, steam.NewClient)
	do.Provide(di, openai.NewClient)
	do.Provide(di, checker.New)

	go do.MustInvoke[*twitch_api.Client](di).RunRefreshLoop(appCtx)
	go do.MustInvoke[*twitch_irc.Client](di).RunRefreshLoop(appCtx)
	go do.MustInvoke[*twitch_irc.Client](di).RunConnectLoop(appCtx)
	go do.MustInvoke[*checker.Service](di).RunCheckLoop(appCtx)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("Shutting down server...")

		cancel()
	}()

	<-appCtx.Done()

	slog.Info("Waiting for services to finish...")
	_ = di.Shutdown()
}
