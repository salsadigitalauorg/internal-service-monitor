package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/cfg"
	"github.com/salsadigitalauorg/internal-services-monitor/internal/expectations"
	"gopkg.in/yaml.v3"
)

func SetupRoutes (r *gin.Engine, monitors []cfg.MonitorConfig) {
	for _, monitor := range monitors {

		r.GET("/monitor/"+monitor.Name, func(c *gin.Context) {
			expectationsMet := false
			var fails []string

			var expectation expectations.Expectation
			switch monitor.Type {
				case "http":
					expectation = &expectations.Http{}
					expectation.WithUrl(monitor.Url)
				case "tcp":
					expectation = &expectations.Tcp{}
					expectation.WithUrl(monitor.Url)
				case "stub":
					expectation = &expectations.Stub{}
					expectation.WithUrl(monitor.Url)

			}

			for _, expects := range monitor.Expects {
				ok, err := expectation.IsOK(expects)
				if ok {
					expectationsMet = true
				} else {
					if err != "" {
						fails = append(fails, err)
					} else {
						fails = append(fails, fmt.Sprintf(
							"Expected %s to be %s %s",
							expects.Field,
							expects.Op,
							expects.Value,
						))
					}
				}
			}

			if !expectationsMet {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Status check failed", "reasons": fails})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
	}

}

func main() {

	config := flag.String("config", "cfg.yml", "Path to configuration file")
	port := flag.String("port", "8080", "Port to start the application on")
	flag.Parse()

	if _, err := os.Stat(*config); os.IsNotExist(err) {
		fmt.Printf("Configuration file (%s) does not exist", *config)
		return
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
		    "message": "pong",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot connect to the database",
		})
	})

	var cfg cfg.Config
	f, err := ioutil.ReadFile(*config)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(f, &cfg)

	if err != nil {
		panic(err)
	}

	SetupRoutes(r, cfg.Monitors)

  r.Run(fmt.Sprintf(":%s", *port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
