package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/copyleftdev/5l4pp3r/internal/config"
	"github.com/copyleftdev/5l4pp3r/internal/gatherer"
	"github.com/copyleftdev/5l4pp3r/internal/storage"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Setup logging
	level, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.Logging.Format == "json" {
		// default JSON format
	} else {
		// Use pretty console logging
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting 5l4pp3r...")

	// Setup storage
	strg, err := storage.NewStorage(cfg.Database.Type, cfg.Database.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize storage")
	}

	ctx := context.Background()
	if err := strg.InitSchema(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database schema")
	}
	defer strg.Close()

	// Gather info
	sysInfo, err := gatherer.GatherSystemInfo()
	if err != nil {
		log.Error().Err(err).Msg("Failed to gather system info")
		return
	}

	netIfaces, err := gatherer.GatherNetworkInfo()
	if err != nil {
		log.Error().Err(err).Msg("Failed to gather network info")
		return
	}

	configFiles, err := gatherer.GatherConfigFiles(cfg, cfg.Compression.Algorithm, cfg.Compression.Level)
	if err != nil {
		log.Error().Err(err).Msg("Failed to gather config files")
		return
	}

	// Store system info first, to get sysInfo.ID
	if err := strg.StoreSystemInfo(ctx, sysInfo); err != nil {
		log.Error().Err(err).Msg("Failed to store system info")
		return
	}

	// Assign system_id to all config files
	for _, cf := range configFiles {
		cf.SystemID = sysInfo.ID
	}

	// Store network interfaces
	for _, iface := range netIfaces {
		if err := strg.StoreNetworkInterface(ctx, iface); err != nil {
			log.Error().Err(err).Msgf("Failed to store network interface %s", iface.InterfaceName)
		}
	}

	// Store config files
	for _, cf := range configFiles {
		if err := strg.StoreConfigFile(ctx, cf); err != nil {
			log.Error().Err(err).Msgf("Failed to store config file %s", cf.FilePath)
		}
	}

	// Commit all changes
	if err := strg.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return
	}

	log.Info().Msg("Snapshot completed successfully.")
}
