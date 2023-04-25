package server

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lfroomin/restaurant-container/controllers"
	"io"
	"log"
	"strings"
)

func NewRouter(env Env) *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	router.Use(cors.New(config))

	router.Use(logRequest)

	restaurant := controllers.Restaurant{
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

func logRequest(c *gin.Context) {
	byteBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n  Request: [%s] %s%s (%s)\n", c.Request.Method, c.Request.Host, c.Request.URL, c.Request.Proto))
	sb.WriteString(fmt.Sprintf("  Header: %+v\n", c.Request.Header))
	if len(byteBody) > 0 {
		sb.WriteString(fmt.Sprintf("  Body: %s\n", string(byteBody)))
	}
	log.Print(sb.String())

	c.Next()
}
