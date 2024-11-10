package log

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
)

func NewLogger() (*slog.Logger, error) {
	var slogLevel slog.Level
	if err := slogLevel.UnmarshalText([]byte(viper.GetString(ArgLevel))); err != nil {
		return nil, err
	}

	var logFormat Format
	if err := logFormat.UnmarshalText([]byte(viper.GetString(ArgFormat))); err != nil {
		return nil, err
	}

	switch logFormat {
	case FormatText:
		return slog.New(
			NewLogHandler(tint.NewHandler(os.Stdout, &tint.Options{
				Level: slogLevel,
			})),
		), nil
	case FormatJSON:
		return slog.New(
			NewLogHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slogLevel,
			})),
		), nil
	}

	return nil, fmt.Errorf("invalid log format: %s", logFormat)
}
