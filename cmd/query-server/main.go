// Command query-server is the read-side gRPC server for the medincident
// platform. This file is a placeholder added in Plan 1 of the monorepo
// merge: it boots zerolog, prints a startup banner, waits for SIGINT/
// SIGTERM, and exits cleanly. Real readers, handlers, and the Zitadel
// NATS consumer are wired in Plan 3.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("query-server placeholder starting")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info().Msg("query-server placeholder shutting down")
}
