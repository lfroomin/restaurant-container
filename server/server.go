package server

import (
	cfg "github.com/lfroomin/restaurant-container/config"
	"github.com/lfroomin/restaurant-container/controllers"
	"github.com/lfroomin/restaurant-container/internal/awsConfig"
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
	awsCfg, err := awsConfig.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Config: RestaurantsTable: %s  LocationPlaceIndex: %s\n", appCfg.RestaurantsTable, appCfg.PlaceIndex)

	return Env{
		Restaurant: dynamo.New(awsCfg, appCfg.RestaurantsTable),
		Location:   geocode.New(awsCfg, appCfg.PlaceIndex),
	}
}
