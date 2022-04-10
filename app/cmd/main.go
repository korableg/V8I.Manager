package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/korableg/V8I.Manager/app/internal/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	configFlag      string = "config"
	configFlagShort string = "c"
)

var Version = "dev"

func main() {
	cmd := &cobra.Command{
		Use:     filepath.Base(os.Args[0]),
		Version: Version,
		RunE:    rootCmd,
	}

	cmd.PersistentFlags().StringP(configFlag, configFlagShort, "./config.yaml", "path to config")

	if err := cmd.Execute(); err != nil {
		logrus.Fatalf("application fatal error: %s", err.Error())
	}
}

func rootCmd(cmd *cobra.Command, args []string) error {
	cfgPath := cmd.Flag(configFlag).Value.String()

	en, err := engine.NewEngine(cfgPath)
	if err != nil {
		return fmt.Errorf("init engine: %w", err)
	}

	errChan := make(chan error, 1)

	go func() {
		if err = en.Start(); err != nil {
			errChan <- err
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	select {
	case <-sig:

		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
		defer cancel()

		if err = en.Shutdown(ctx); err != nil {
			return err
		}

		logrus.Info("application gracefully shutdown")

		return nil
	case err, ok := <-errChan:
		if !ok {
			return errors.New("error chan was closed")
		}
		return err
	}

	return nil
}
