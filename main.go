package main

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
)

func loadConfig() *toml.Tree {
	file, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	config, err := toml.Load(string(file))
	if err != nil {
		panic(err)
	}

	return config
}

func main() {
	r := gin.Default()

	tmpl := template.Must(template.ParseGlob("templates/*"))
	cfg := loadConfig()

	r.GET("/", func(c *gin.Context) {
		tmpl.ExecuteTemplate(c.Writer, "index.html", cfg.ToMap())
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":    "pong",
			"serverTime": time.Now(),
		})
	})

	r.Run(":8080")
}
