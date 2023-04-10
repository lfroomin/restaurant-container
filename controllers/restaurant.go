package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lfroomin/restaurant-container/internal/model"
	"log"
	"net/http"
)

type RestaurantStorer interface {
	Save(restaurant model.Restaurant) error
	Get(restaurantId string) (model.Restaurant, bool, error)
	Update(restaurant model.Restaurant) error
	Delete(restaurantId string) error
}

type Geocoder interface {
	Geocode(address model.Address) (model.Location, string, error)
}

type RestaurantController struct {
	Restaurant RestaurantStorer
	Location   Geocoder
}

func (rc RestaurantController) Create(c *gin.Context) {
	var restaurant model.Restaurant
	err := c.ShouldBindJSON(&restaurant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "error binding request body"})
		return
	}

	id := uuid.NewString()
	restaurant.Id = &id
	log.Printf("RestaurantController.Create restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := rc.Location.Geocode(*restaurant.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
			return
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := rc.Restaurant.Save(restaurant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, restaurant)
}

func (rc RestaurantController) Read(c *gin.Context) {
	restaurantId := c.Param("restaurantId")

	log.Printf("RestaurantController.Read restaurantId: %s\n", restaurantId)

	// Validate input
	if restaurantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "restaurantId is empty"})
		return
	}

	restaurant, exists, err := rc.Restaurant.Get(restaurantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func (rc RestaurantController) Update(c *gin.Context) {
	restaurantId := c.Param("restaurantId")

	var restaurant model.Restaurant
	err := c.ShouldBindJSON(&restaurant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "error binding request body"})
		return
	}

	if restaurant.Id == nil || restaurantId != *restaurant.Id {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "restaurantId in URL path parameters and restaurant in body do not match"})
		return
	}

	log.Printf("RestaurantController.Update restaurantName: %s  restaurantId: %s\n", restaurant.Name, *restaurant.Id)

	// Get the geocode of the restaurant address
	if restaurant.Address != nil {
		location, timezoneName, err := rc.Location.Geocode(*restaurant.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
			return
		}

		restaurant.Address.Location = &location
		restaurant.Address.TimezoneName = &timezoneName
	}

	if err := rc.Restaurant.Update(restaurant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func (rc RestaurantController) Delete(c *gin.Context) {
	restaurantId := c.Param("restaurantId")

	// Validate input
	if restaurantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "restaurantId is empty"})
		return
	}

	log.Printf("RestaurantController.Delete restaurantId: %s\n", restaurantId)

	err := rc.Restaurant.Delete(restaurantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "")
}
