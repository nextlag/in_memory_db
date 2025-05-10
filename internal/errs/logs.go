package errs

import "log/slog"

func ErrLog(err error) slog.Attr {
	return slog.String("error", err.Error())
}
