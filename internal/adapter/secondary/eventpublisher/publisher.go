package eventpublisher

import (
	"context"
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher/outbox"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/domain"
	"f4allgo-restaurant/internal/core/port"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/rs/zerolog"
	tally "github.com/uber-go/tally/v4"
	"gorm.io/gorm"
)

type DomainEventOutboxPublisher struct {
	repository    outbox.OutboxRepository
	eventCounters map[string]tally.Counter
}

// Interface compliance verification.
var _ port.DomainEventPublisher = (*DomainEventOutboxPublisher)(nil)

func NewDomainEventOutboxPublisher(db *gorm.DB, ctxGetter *trmgorm.CtxGetter, logger zerolog.Logger, config *boot.Config, scope tally.Scope) *DomainEventOutboxPublisher {
	outboxRepository := outbox.NewOutboxPostgresRepository(db, ctxGetter, config, logger)
	if config.AppInitOutboxDispatcher {
		// Initializes the outbox dispatcher and forget about it (because it
		// runs in its own goroutine)
		dispatcher := outbox.NewOutboxDispatcher(outboxRepository, logger, config, scope)
		dispatcher.InitOutboxDispatcher()
	}

	var eventCounters map[string]tally.Counter
	if scope != nil {
		restaurantCreated := scope.Tagged(map[string]string{"event_type": "RestaurantCreated"}).Counter("outgoing_events")
		restaurantDeleted := scope.Tagged(map[string]string{"event_type": "RestaurantDeleted"}).Counter("outgoing_events")
		restaurantMenuUpdated := scope.Tagged(map[string]string{"event_type": "RestaurantMenuUpdated"}).Counter("outgoing_events")
		eventCounters = map[string]tally.Counter{
			"RestaurantCreated":     restaurantCreated,
			"RestaurantDeleted":     restaurantDeleted,
			"RestaurantMenuUpdated": restaurantMenuUpdated,
		}
	}

	return &DomainEventOutboxPublisher{
		repository:    outboxRepository,
		eventCounters: eventCounters,
	}
}

// Publish publishes a domain event to the outside world. In this particular
// implementation it just persist the event using an outbox repository that
// stores the event to the outbox table. Another process will be responsible
// of sending the event in a reliable to message broker way using the outbox
// pattern.
func (p *DomainEventOutboxPublisher) Publish(ctx context.Context, e domain.DomainEvent) error {
	result := p.repository.Save(ctx, e)
	if result == nil && p.eventCounters != nil {
		p.eventCounters[e.GetType()].Inc(1)
	}

	return result
}
