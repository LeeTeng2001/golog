package common

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func ConfigureGlobalLog(output string) error {
	w := io.Discard
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("could not create debug output at %s :%w ", output, err)
		}
		w = f
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
	return nil
}
