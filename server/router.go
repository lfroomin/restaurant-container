package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lfroomin/restaurant-container/controllers"
)

func NewRouter(env Env) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	restaurant := controllers.RestaurantController{
		Restaurant: env.Restaurant,
		Location:   env.Location,
	}

	router.POST("/", restaurant.Create)

	idGrp := router.Group("/:restaurantId")
	idGrp.GET("", restaurant.Read)
	idGrp.POST("", restaurant.Update)
	idGrp.DELETE("", restaurant.Delete)

	return router
}
