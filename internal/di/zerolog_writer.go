package di

import (
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/medincident/medincident-command-service/internal/config"
)

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
