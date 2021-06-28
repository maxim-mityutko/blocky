package cmd

import (
	"blocky/config"
	"blocky/log"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	configPath string
	cfg        config.Config
	apiHost    string
	apiPort    uint16
)

// NewRootCommand creates a new root cli command instance
func NewRootCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "blocky",
		Short: "blocky is a DNS proxy ",
		Long: `A fast and configurable DNS Proxy
and ad-blocker for local network.
		   
Complete documentation is available at https://github.com/0xERR0R/blocky`,
		Run: func(cmd *cobra.Command, args []string) {
			newServeCommand().Run(cmd, args)
		},
	}

	c.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "path to config file")
	c.PersistentFlags().StringVar(&apiHost, "apiHost", "localhost", "host of blocky (API)")
	c.PersistentFlags().Uint16Var(&apiPort, "apiPort", 4000, "port of blocky (API)")

	c.AddCommand(newRefreshCommand(),
		NewQueryCommand(),
		NewVersionCommand(),
		newServeCommand(),
		newBlockingCommand(),
		NewListsCommand())

	return c
}

func apiURL(path string) string {
	return fmt.Sprintf("http://%s:%d%s", apiHost, apiPort, path)
}

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cfg = config.NewConfig(configPath, false)
	log.ConfigureLogger(cfg.LogLevel, cfg.LogFormat, cfg.LogTimestamp)

	if cfg.HTTPPort != "" {
		split := strings.Split(cfg.HTTPPort, ":")
		port, err := strconv.Atoi(split[len(split)-1])

		if err == nil {
			apiPort = uint16(port)
		}
	}
}

// Execute starts the command
func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
