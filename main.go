package main

import (
	"fmt"
	"html/template"
	"os"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml"
)

type ServerInfo struct {
	version         string
	server_timezone string
	server_uptime   string
}

type Types struct {
	types []reflect.Type
}

func validateType(t string, types []string) bool {
	for _, v := range types {
		if v == t {
			return true
		}
	}

	return false
}

func valiateUserConfig(cfg map[string]interface{}, server_config map[string]interface{}) {
	server_params := server_config["Paramaters"].(map[string]interface{})
	types_unchecked := server_config["Input_Types"].(map[string]interface{})["types"]
	typesSlice, ok := types_unchecked.([]interface{})

	if !ok {
		panic("Invalid Types: Expected 'types' to be a list of strings")
	}

	var types []string
	for _, v := range typesSlice {
		str, ok := v.(string)

		if !ok {
			panic("Invalid Types: not a string")
		}

		types = append(types, str)
	}

	for k, v := range cfg {
		if reflect.TypeOf(v) != reflect.TypeOf(map[string]interface{}{}) {
			if !server_params["allow_bare_inputs"].(bool) {
				os.Stderr.WriteString(fmt.Sprintf("Invalid config key: %s\n", k))
				os.Exit(1)
			} else if !validateType(v.(string), types) {
				os.Stderr.WriteString(fmt.Sprintf("Invalid config type: %s, at key: %s\n", v, k))
				os.Exit(1)
			}
		}
	}
}

func loadConfig(name string) *toml.Tree {
	file, err := os.ReadFile(name)
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
	user_config := loadConfig("config.toml").ToMap()
	server_config := loadConfig("server.config.toml").ToMap()

	// valiateServerConfig()
	valiateUserConfig(user_config, server_config)

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	r.GET("/", func(c *gin.Context) {
		tmpl.ExecuteTemplate(c.Writer, "index.html", user_config)
	})

	r.GET("/server", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	r.Run(":8080")
}
