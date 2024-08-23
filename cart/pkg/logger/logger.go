package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func SetDefaultLogger(lvl string) {
	if lvl == "debug" {
		w := os.Stdout
		logger := slog.New(
			tint.NewHandler(w, &tint.Options{
				NoColor: !isatty.IsTerminal(w.Fd()),
			}),
		)
		slog.SetDefault(logger)
		return
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}
