package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/korableg/V8I.Manager/internal/config"
	"github.com/korableg/V8I.Manager/internal/globals"
	"github.com/korableg/V8I.Manager/internal/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	_cfg = "cfg"
)

var Version = "dev"

var (
	_cfgFile string
	cmd      *cobra.Command
	cfg      config.Config
)

func _rootCmd(cmd *cobra.Command, args []string) error {

	err := initConfig()
	if err != nil {
		return err
	}

	err = initConfigOs()
	if err != nil {
		return err
	}

	if cfg.IsWindowsService() {
		return runService()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewWorker(cfg.Lsts(), cfg.V8is())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		w.StartWatchingContext(ctx)
	}()

	select {
	case <-ctx.Done():
	case <-sig:
	}

	return nil

}

func main() {

	cmd := &cobra.Command{
		Use:     filepath.Base(os.Args[0]),
		Short:   globals.Description,
		Version: Version,
		RunE:    _rootCmd,
	}

	addCommands(cmd)

	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.PersistentFlags().StringVarP(&_cfgFile, _cfg, "c", "", "file with the application's settings")

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}

func initConfig() (err error) {

	cfg, err = config.New(_cfgFile)

	if err != nil {
		return err
	}

	return nil

}
