package outbox

import (
	"context"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/domain"
	"fmt"
	"strconv"
	"time"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const lockMaxDuration time.Duration = 30
const restaurantAggregateType string = "Restaurant"

// OutboxRepository manages outbox persistent operations on domain events before they
// are published to a message broker. It is part of the outbox pattern implementation.
type OutboxRepository interface {

	// Save persists an outbox event in the external storage. This operation is used
	// by the event publisher inside a business transaction where the event is raised.
	Save(ctx context.Context, e domain.DomainEvent) error

	// acquireLock gets a lock on the outbox table.
	acquireLock() (bool, error)

	// releaseLock releases a lock on the outbox table.
	releaseLock() error

	// findInBatches restrieves all the registered events in the outbox table to
	// be processed in batches.
	findInBatches(batchSize int, fc func(*[]*Outbox, *gorm.DB) error) error

	// deleteInBatches deletes the provided records from the outbox table in batches.
	deleteInBatches(batchSize int, records []*Outbox) error
}

type OutboxPostgresRepository struct {
	mapper    Mapper
	db        *gorm.DB
	ctxGetter *trmgorm.CtxGetter
	logger    zerolog.Logger
}

// Interface compliance verification.
var _ OutboxRepository = (*OutboxPostgresRepository)(nil)

func NewOutboxPostgresRepository(db *gorm.DB, ctxGetter *trmgorm.CtxGetter, logger zerolog.Logger) *OutboxPostgresRepository {
	return &OutboxPostgresRepository{mapper: DefaultMapper{}, db: db, ctxGetter: ctxGetter, logger: logger}
}

func (r *OutboxPostgresRepository) Save(ctx context.Context, event domain.DomainEvent) error {
	var outboxRow *Outbox
	var avroRecord any

	switch e := event.(type) {
	case *domain.RestaurantCreated:
		outboxRow = &Outbox{
			Id:            uuid.New(),
			AggregateType: restaurantAggregateType,
			AggregateId:   strconv.FormatUint(e.Restaurant.Id, 10),
			EventType:     e.GetType(),
		}
		avroRecord = r.mapper.fromRestaurantCreated(e)
	case *domain.RestaurantDeleted:
		outboxRow = &Outbox{
			Id:            uuid.New(),
			AggregateType: restaurantAggregateType,
			AggregateId:   strconv.FormatUint(e.RestaurantId, 10),
			EventType:     e.GetType(),
		}
		avroRecord = r.mapper.fromRestaurantDeleted(e)
	case *domain.RestaurantMenuUpdated:
		outboxRow = &Outbox{
			Id:            uuid.New(),
			AggregateType: restaurantAggregateType,
			AggregateId:   strconv.FormatUint(e.RestaurantId, 10),
			EventType:     e.GetType(),
		}
		avroRecord = r.mapper.fromRestaurantMenuUpdated(e)
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(boot.GetConfig().KafkaSchemaRegistry))
	if err != nil {
		return fmt.Errorf("creating the schema registry client: %w", err)
	}

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		return fmt.Errorf("creating the avro serializer: %w", err)
	}

	avroBytes, err := ser.Serialize(buildOutboxTopicNamefromEventType(outboxRow.EventType), avroRecord)
	if err != nil {
		return fmt.Errorf("serializing to avro: %w", err)
	}

	outboxRow.Payload = avroBytes

	return r.ctxGetter.DefaultTrOrDB(ctx, r.db).Clauses(clause.OnConflict{DoNothing: true}).Create(outboxRow).Error
}

func (r *OutboxPostgresRepository) findInBatches(batchSize int, fc func(*[]*Outbox, *gorm.DB) error) error {
	var entries []*Outbox
	result := r.db.FindInBatches(&entries, batchSize, func(tx *gorm.DB, batch int) error {
		return fc(&entries, tx)
	})

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (r *OutboxPostgresRepository) deleteInBatches(batchSize int, records []*Outbox) error {
	totalRecords := len(records)
	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}
		batch := records[i:end]
		if err := r.db.Delete(&batch).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *OutboxPostgresRepository) acquireLock() (bool, error) {
	r.logger.Debug().Msg("Trying to get the lock on outbox table")

	var lock OutboxLock
	if err := r.db.Where("id = ?", 1).First(&lock).Error; err != nil {
		r.logger.Err(err).Msg("while getting the lock row from outbox_lock table")
		return false, err
	}

	if lock.Locked && lock.LockedUntil.After(time.Now()) {
		return false, nil
	} else {
		now := time.Now()
		lock.Locked = true
		lock.LockedAt = &now
		until := now.Add(lockMaxDuration * time.Second)
		lock.LockedUntil = &until

		if err := r.db.Save(&lock).Error; err != nil {
			r.logger.Err(err).Msg("while saving the locked lock row into outbox_lock table")
			return false, err
		}

		return true, nil
	}
}

func (r *OutboxPostgresRepository) releaseLock() error {
	r.logger.Debug().Msg("Trying to release the lock on outbox table")

	var lock OutboxLock
	if err := r.db.Where("id = ?", 1).First(&lock).Error; err != nil {
		r.logger.Err(err).Msg("while getting the lock row from outbox_lock table")
		return err
	}

	if !lock.Locked {
		r.logger.Warn().Msg("The lock is already free")
		return nil
	} else {
		lock.Locked = false
		lock.LockedAt = nil
		lock.LockedUntil = nil

		if err := r.db.Save(&lock).Error; err != nil {
			r.logger.Err(err).Msg("while saving the released lock row into outbox_lock table")
			return err
		}

		return nil
	}
}
