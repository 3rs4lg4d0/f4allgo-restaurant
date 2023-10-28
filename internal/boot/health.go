package boot

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/3rs4lg4d0/go-kafka-checker"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/InVisionApp/go-health/v2/handlers"
)

var healthOnce sync.Once

var handler http.HandlerFunc

func GetHealthHandler(db *sql.DB) http.HandlerFunc {
	healthOnce.Do(func() {
		l := GetLogger()
		h := health.New()
		h.DisableLogging()

		// Create a kafka check skipping the first three consumer timeouts if any.
		kafkaCheck, err := kafka.NewKafka(kafka.KafkaConfig{
			BootstrapServers:     GetConfig().KafkaBootstrapServers,
			SkipConsumerTimeouts: 3,
		})
		if err != nil {
			l.Err(err).Msg("unable to create a kafka check")
		} else {
			err := h.AddCheck(&health.Config{
				Name:     "check-kafka",
				Checker:  kafkaCheck,
				Interval: 5 * time.Second,
				Fatal:    false,
			})
			if err != nil {
				l.Err(err).Msg("kafka check was not added")
			}
		}

		// Create a simple SQL check.
		sqlCheck, err := checkers.NewSQL(&checkers.SQLConfig{
			Pinger: db,
		})
		if err != nil {
			l.Err(err).Msg("unable to create a database check")
		} else {
			err := h.AddCheck(&health.Config{
				Name:     "check-db",
				Checker:  sqlCheck,
				Interval: 3 * time.Second,
				Fatal:    true,
			})
			if err != nil {
				l.Err(err).Msg("sql check was not added")
			}
		}

		if err := h.Start(); err != nil {
			l.Err(err).Msg("unable to start healthcheck")
		}

		handler = handlers.NewJSONHandlerFunc(h, nil)
	})

	return handler
}
