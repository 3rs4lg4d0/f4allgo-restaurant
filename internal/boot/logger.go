package boot

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var loggerOnce sync.Once

var log zerolog.Logger

func GetLogger() zerolog.Logger {
	loggerOnce.Do(func() {
		logLevel := GetConfig().LogLevel
		var output io.Writer
		if GetConfig().LogBeautify {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
		} else {
			output = os.Stdout
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Caller().
			Logger()
	})

	return log
}
