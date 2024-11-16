package slogx

import "log/slog"

func Error(msg string, err error) {
	slog.Error(msg, "err", err)
}

func DebugAny(msg string, a any) {
	slog.Debug(msg, "obj", a)
}
