package brew

import (
	"log/slog"
	"os"
)

func InitLogging(prefex string) {
	if prefex == "" {
		prefex = "brew"
	}
	// create a log file with the given prefix
	logFile, err := os.OpenFile(prefex+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "prefix" {
				a.Value = slog.StringValue(prefex)
			}
			return a
		},
	})))
}
