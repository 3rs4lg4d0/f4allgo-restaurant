package outbox

import (
	"f4allgo-restaurant/internal/boot"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog"
	tally "github.com/uber-go/tally/v4"
	"gorm.io/gorm"
)

const batchSize int = 300

const (
	successCounter string = "success"
	errorCounter   string = "error"
)

type OutboxDispatcher struct {
	repository     OutboxRepository
	logger         zerolog.Logger
	producer       *kafka.Producer
	reportCounters map[string]tally.Counter
}

func NewOutboxDispatcher(repository OutboxRepository, logger zerolog.Logger, config *boot.Config, scope tally.Scope) *OutboxDispatcher {
	producer, _ := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  config.KafkaBootstrapServers,
		"linger.ms":          500,
		"batch.size":         100 * 1024,
		"compression.type":   "lz4",
		"acks":               -1,
		"enable.idempotence": true,
	})

	var reportCounters map[string]tally.Counter
	if scope != nil {
		successes := scope.Tagged(map[string]string{"outcome": "success"}).Counter("outbox")
		errors := scope.Tagged(map[string]string{"outcome": "error"}).Counter("outbox")
		reportCounters = map[string]tally.Counter{
			"success": successes,
			"error":   errors,
		}
	}

	return &OutboxDispatcher{repository: repository, logger: logger, producer: producer, reportCounters: reportCounters}
}

// InitOutboxDispatcher initializes a background process (inside a go routine) that
// periodically polls the outbox table in order to send events to a message broker.
func (d *OutboxDispatcher) InitOutboxDispatcher() {
	d.logger.Debug().Msg("initializing the outbox dispatcher")
	go d.execute()
}

func (d *OutboxDispatcher) execute() {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		if acquired, err := d.acquireOutboxLock(); acquired {
			d.processOutbox()
			err := d.releaseOutboxLock()
			if err != nil {
				d.logger.Err(err).Msg("releasing the outbox lock")
			}
		} else if err != nil {
			d.logger.Debug().Msg("the lock is in use right now ¯\\_(ツ)_/¯")
		}
	}
}

func (d *OutboxDispatcher) acquireOutboxLock() (bool, error) {
	return d.repository.acquireLock()
}

func (d *OutboxDispatcher) releaseOutboxLock() error {
	return d.repository.releaseLock()
}

func (d *OutboxDispatcher) processOutbox() {
	var success []*Outbox
	var totalProcessed int
	var totalErr int
	var deliveryChan = make(chan kafka.Event, batchSize)
	var wg sync.WaitGroup

	d.logger.Debug().Msg("processing outbox messages")

	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					d.logger.Err(m.TopicPartition.Error).Msg("delivery problem")
					totalErr++
					d.reportCounters[errorCounter].Inc(1)
				} else {
					d.logger.Trace().Msgf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

					// Adding the successful outbox Id to the success slice for deletion.
					uuid, _ := m.Opaque.(uuid.UUID)
					success = append(success, &Outbox{Id: uuid})
					d.reportCounters[successCounter].Inc(1)
				}

				totalProcessed++
				wg.Done()

			default:
				d.logger.Debug().Msgf("Ignored event: %s", ev)
			}
		}
		d.logger.Debug().Msg("the goroutine for Kafka delivery reports has finished")
	}()

	err := d.repository.findInBatches(batchSize, func(batch *[]*Outbox, tx *gorm.DB) error {
		d.logger.Debug().Msgf("Sending %d messages to kafka", len(*batch))
		for _, o := range *batch {
			topic := buildOutboxTopicNamefromEventType(o.EventType)
			err := d.producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(o.AggregateId),
				Value:          o.Payload,
				Opaque:         o.Id,
			}, deliveryChan)

			if err != nil {
				d.logger.Err(err).Msg("when producing a message")
				// if any error happen sending the message we don't need to retry here,
				// the message will remain in the outbox table and will be sent in the
				// next outbox processing.
			} else {
				wg.Add(1)
			}
		}

		return nil
	})

	if err != nil {
		d.logger.Err(err).Msg("when trying to get outbox rows in batches")
	}

	// Wait until we get all the delivery reports from kafka client.
	wg.Wait()

	// We can safely close the channel because this is a dedicated channel only to
	// receive as many delivery reports as many messages are sent.
	close(deliveryChan)
	d.logger.Info().Msgf("%d messages where successfully delivered (with %d failed) from a total of %d processed from outbox", len(success), totalErr, totalProcessed)

	if len(success) > 0 {
		d.logger.Debug().Msgf("Deleting %d elements from outbox", len(success))
		err := d.repository.deleteInBatches(batchSize, success)
		if err != nil {
			d.logger.Err(err).Msg("when deleting sent outbox records in batches")
		}
	}
}

func buildOutboxTopicNamefromEventType(eventType string) string {
	return fmt.Sprintf("outbox-%s", strcase.ToKebab(eventType))
}
