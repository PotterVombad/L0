package main

import (
	"context"

	environment "github.com/Netflix/go-env"
	"github.com/PotterVombad/L0/internal/api"
	"github.com/PotterVombad/L0/internal/db"
	"github.com/PotterVombad/L0/internal/db/cashe"
	"github.com/PotterVombad/L0/internal/db/pgs"
	"github.com/PotterVombad/L0/internal/env"
	stan "github.com/PotterVombad/L0/internal/nats"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	if _, err := environment.UnmarshalFromEnviron(&env.Env); err != nil {
		panic(err)
	}

	url := env.Env.GetPostgresURL()

	pgsDb := pgs.MustNew(ctx, url)
	cashe := cashe.New()

	storage := db.MustNew(ctx, cashe, &pgsDb)

	// TODO env
	clusterId := "test-cluster" 
	clientId := "test-client"
	subject := "test-subject"
	natsUrl := "http://localhost:4222"

	stream := stan.MustNew(ctx, clusterId, clientId, subject, natsUrl, &storage)

	defer func(stan.Stream, pgs.PgsDB) {
		stream.Close()
		pgsDb.Close(ctx)
	}(stream, pgsDb)

	api := api.New(&storage)

	log.Info("start app")

	if err := api.Run(); err != nil {
		panic(err)
	}
}
