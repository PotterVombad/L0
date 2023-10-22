package nats

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/PotterVombad/L0/internal/db"
	"github.com/PotterVombad/L0/internal/models"
	stan "github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
)

type Stream struct {
	conn    stan.Conn
	sub     stan.Subscription
	storage db.Store
}

func (s Stream) handleEvent(
	ctx context.Context,
	bb []byte,
) error {
	var o models.Order
	if err := json.Unmarshal(bb, &o); err != nil {
		return fmt.Errorf("json unmarshal | %s", err)
	}

	if err := s.storage.Save(ctx, o); err != nil {
		return fmt.Errorf("storage save | %s", err)
	}

	return nil
}

func (s *Stream) Close() {
	if err := s.sub.Unsubscribe(); err != nil {
		log.Errorf("could not unsubscribe from subject: %v", err)
	}

	s.conn.Close()
}

func MustNew(
	ctx context.Context,
	clasterId, clientId, subject, natsUrl string,
	storage db.Store,
) Stream {
	conn, err := stan.Connect(clasterId, clientId, stan.NatsURL(natsUrl))
	if err != nil {
		panic(fmt.Sprintf("could not connect to NATS Streaming Server: %s", err))
	}

	stream := Stream{
		conn:    conn,
		storage: storage,
	}

	sub, err := conn.Subscribe(
		subject,
		func(m *stan.Msg) {
			if err := stream.handleEvent(ctx, m.Data); err != nil {
				log.Error(err)
			}
		},
	)
	if err != nil {
		// TODO: panic
		log.Panic("Could not subscribe to subject: %w", err)
	}

	stream.sub = sub

	return stream
}
