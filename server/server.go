package server

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/location"
	cfg "github.com/lfroomin/restaurant-container/config"
	"github.com/lfroomin/restaurant-container/controllers"
	"github.com/lfroomin/restaurant-container/internal/dynamo"
	"github.com/lfroomin/restaurant-container/internal/geocode"
	"log"
)

func Init(appCfg cfg.Config) {
	env := newEnv(appCfg)
	r := NewRouter(env)
	r.Run(appCfg.ServerAddress)
}

type Env struct {
	Restaurant controllers.RestaurantStorer
	Location   controllers.Geocoder
}

func newEnv(appCfg cfg.Config) Env {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("failed loading config: %w", err)
	}

	log.Printf("Config: RestaurantsTable: %s  LocationPlaceIndex: %s\n", appCfg.RestaurantsTable, appCfg.PlaceIndex)

	return Env{
		Restaurant: dynamo.RestaurantStorage{
			Client: dynamodb.NewFromConfig(awsCfg),
			Table:  appCfg.RestaurantsTable,
		},
		Location: geocode.LocationService{
			Client:     location.NewFromConfig(awsCfg),
			PlaceIndex: appCfg.PlaceIndex,
		},
	}
}
