// Package bootstrap provides construction helpers used by the three
// binaries (command-server, query-server, gateway-server) during
// startup. Each function here is a straight Go constructor — no DI
// framework, no magic. Callers (main.go) wire things together
// explicitly and drive shutdown via `defer`.
package bootstrap

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/oops"
	oopszerolog "github.com/samber/oops/loggers/zerolog"

	"github.com/medincident/medincident-command-service/internal/config"
)

// Error codes emitted by BuildZerolog.
const (
	ErrCodeZerologInvalidLevel       = "invalid_level"
	ErrCodeZerologInvalidOutputLevel = "invalid_output_level"
	ErrCodeZerologDuplicateWriter    = "duplicate_writer"
	ErrCodeZerologDuplicateFile      = "duplicate_file"
	ErrCodeZerologMissingFilePath    = "missing_file_path"
	ErrCodeZerologUnknownOutputType  = "unknown_output_type"
)

// zerologGlobalsOnce gates the one-time mutation of process-wide
// zerolog state (TimeFieldFormat, ErrorMarshalFunc, ErrorStackMarshaler,
// TimestampFunc). Subsequent calls reuse whatever the first caller set.
var zerologGlobalsOnce sync.Once

// BuildZerolog returns a configured *zerolog.Logger and a cleanup
// function. The cleanup closes file writers opened for ZerologOutputTypeFile
// outputs; for console-only configs it's a no-op returning nil.
//
// The caller is responsible for invoking cleanup, typically via
// `defer cleanup()` immediately after the error check.
func BuildZerolog(cfg *config.ZerologConfig) (*zerolog.Logger, func() error, error) {
	eb := oops.In("bootstrap.zerolog")

	globalLevel, err := zerolog.ParseLevel(string(cfg.Level))
	if err != nil {
		return nil, nil, eb.Code(ErrCodeZerologInvalidLevel).Errorf("invalid log level %q", cfg.Level)
	}

	var (
		outputs   []zerolog.LevelWriter
		closers   []io.Closer
		writers   []*os.File
		filePaths []string
	)

	for idx := range cfg.Outputs {
		out := &cfg.Outputs[idx]

		if out.Level != "" {
			if _, err := zerolog.ParseLevel(string(out.Level)); err != nil {
				return nil, nil, eb.Code(ErrCodeZerologInvalidOutputLevel).With("output_index", idx).
					Errorf("invalid output level %q", out.Level)
			}
		}

		switch out.Type {
		case config.ZerologOutputTypeConsole:
			w := consoleTarget(out.Target)
			if slices.Contains(writers, w) {
				return nil, nil, eb.Code(ErrCodeZerologDuplicateWriter).With("output_index", idx).
					Errorf("duplicate console output for %s", out.Target)
			}
			writers = append(writers, w)
			outputs = append(outputs, buildOutputWriter(w, out, cfg.TimeFormat))

		case config.ZerologOutputTypeFile:
			if out.Path == "" {
				return nil, nil, eb.Code(ErrCodeZerologMissingFilePath).With("output_index", idx).
					Errorf("file output requires a non-empty path")
			}
			abs, err := filepath.Abs(out.Path)
			if err != nil {
				return nil, nil, eb.With("output_index", idx).Wrap(err)
			}
			if slices.Contains(filePaths, abs) {
				return nil, nil, eb.Code(ErrCodeZerologDuplicateFile).With("output_index", idx).
					Errorf("duplicate file output %q", abs)
			}
			filePaths = append(filePaths, abs)

			f, err := os.OpenFile(out.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
			if err != nil {
				return nil, nil, eb.With("output_index", idx).Wrap(err)
			}
			closers = append(closers, f)
			outputs = append(outputs, buildOutputWriter(f, out, cfg.TimeFormat))

		default:
			return nil, nil, eb.Code(ErrCodeZerologUnknownOutputType).With("output_index", idx).
				Errorf("unknown output type %q", out.Type)
		}
	}

	var w io.Writer
	switch len(outputs) {
	case 0:
		w = io.Discard
	case 1:
		w = outputs[0]
	default:
		ws := make([]io.Writer, len(outputs))
		for i, lw := range outputs {
			ws[i] = lw
		}
		w = zerolog.MultiLevelWriter(ws...)
	}

	zerologGlobalsOnce.Do(func() {
		zerolog.TimeFieldFormat = cfg.TimeFormat
		zerolog.ErrorMarshalFunc = oopszerolog.OopsMarshalFunc
		zerolog.ErrorStackMarshaler = oopszerolog.OopsStackMarshaller
		if cfg.TimeUTC {
			zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
		}
	})

	logger := zerolog.New(w).Level(globalLevel)
	ctx := logger.With()

	if cfg.Timestamp {
		ctx = ctx.Timestamp()
	}

	switch {
	case cfg.Caller && cfg.CallerSkip > 0:
		ctx = ctx.CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + cfg.CallerSkip)
	case cfg.Caller:
		ctx = ctx.Caller()
	}

	for key, value := range cfg.Fields {
		ctx = ctx.Str(key, value)
	}

	cleanup := func() error { return nil }
	if len(closers) > 0 {
		cleanup = func() error {
			var errs []error
			for i := len(closers) - 1; i >= 0; i-- {
				if err := closers[i].Close(); err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}
	}

	result := ctx.Logger()
	return &result, cleanup, nil
}

// consoleTarget maps a ZerologConsoleTarget to the corresponding *os.File.
func consoleTarget(target config.ZerologConsoleTarget) *os.File {
	if target == config.ZerologConsoleTargetStdout {
		return os.Stdout
	}
	return os.Stderr
}

// buildOutputWriter wraps w in a zerolog.LevelWriter, applying pretty
// formatting and per-output level filtering as configured.
func buildOutputWriter(w io.Writer, out *config.ZerologOutputConfig, globalTimeFormat string) zerolog.LevelWriter {
	tf := globalTimeFormat
	if out.TimeFormat != "" {
		tf = out.TimeFormat
	}

	var base zerolog.LevelWriter
	if out.Pretty {
		cw := zerolog.ConsoleWriter{
			Out:        w,
			NoColor:    out.NoColor,
			TimeFormat: tf,
		}
		if len(out.PartsOrder) > 0 {
			cw.PartsOrder = out.PartsOrder
		}
		if len(out.PartsExclude) > 0 {
			cw.PartsExclude = out.PartsExclude
		}
		base = zerolog.LevelWriterAdapter{Writer: cw}
	} else {
		base = zerolog.LevelWriterAdapter{Writer: w}
	}

	if out.Level != "" {
		level, err := zerolog.ParseLevel(string(out.Level))
		if err == nil && level > zerolog.TraceLevel {
			return &zerolog.FilteredLevelWriter{Writer: base, Level: level}
		}
	}

	return base
}
