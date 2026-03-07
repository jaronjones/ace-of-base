package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/eljefe/islandtime/internal/config"
	"github.com/eljefe/islandtime/internal/server"

	// Register all widget types via init()
	_ "github.com/eljefe/islandtime/internal/widgets/bookmark"
	_ "github.com/eljefe/islandtime/internal/widgets/clock"
	_ "github.com/eljefe/islandtime/internal/widgets/todo"
	_ "github.com/eljefe/islandtime/internal/widgets/weather"
)

func main() {
	configPath := flag.String("config", "config/dashboard.json", "path to dashboard config file")
	devMode := flag.Bool("dev", false, "dev mode: load templates from disk on each request")
	logLevel := flag.String("log", "info", "log level: debug, info, warn, error")
	flag.Parse()

	// Configure structured logging
	var level slog.Level
	if err := level.UnmarshalText([]byte(*logLevel)); err != nil {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "path", *configPath, "err", err)
		os.Exit(1)
	}

	slog.Info("loaded config", "title", cfg.Title, "widgets", len(cfg.Widgets))

	// Run server (blocks)
	if err := server.Run(cfg, "templates", *devMode); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
