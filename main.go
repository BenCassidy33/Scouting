package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"time"

	"github.com/pelletier/go-toml"
)

const (
	toml_type = iota
	json_type
)

type ServerInfo struct {
	GoVersion      string `json:"goVersion"`
	ServerTimezone string `json:"serverTimezone`
	ServerUptime   string `json:"serverUptime"`

	ConfigInfo map[string]interface{} `json:"configInfo"`
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
	entries, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}
	fmt.Println(entries)

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
				panic(fmt.Sprintf("Invalid config key: %s\n", k))
			} else if !validateType(v.(string), types) {
				panic(fmt.Sprintf("Invalid config type: %s, at key: %s\n", v, k))
			}
		}
	}
}

func loadConfig(name string, filetype int) map[string]interface{} {
	file, err := os.ReadFile(name)
	config := make(map[string]interface{})
	if err != nil {
		panic(err)
	}

	if filetype == toml_type {
		c, err := toml.Load(string(file))
		if err != nil {
			panic(err)
		}
		config = c.ToMap()
	} else if filetype == json_type {
		err = json.Unmarshal(file, &config)
		if err != nil {
			panic(err)
		}
	}

	return config
}

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	server_config := loadConfig("server.config.toml", toml_type)
	user_config := make(map[string]interface{})

	if server_config["Paramaters"].(map[string]interface{})["config_filetype"] == "json" {
		if _, err := os.Stat("config.json"); err != nil {
			panic("config.json not found")
		}
		user_config = loadConfig("config.json", json_type)
	} else if server_config["Paramaters"].(map[string]interface{})["config_filetype"] == "toml" {
		if _, err := os.Stat("config.toml"); err != nil {
			panic("config.toml not found")
		}
		user_config = loadConfig("config.toml", toml_type)
	} else {
		panic("Invalid Config Type")
	}

	// valiateServerConfig()
	valiateUserConfig(user_config, server_config)
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	router := http.NewServeMux()
	server_start_time := time.Now()
	fmt.Println("Server started at: ", server_start_time)

	router.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Home(w, tmpl, user_config)
	}))

	router.Handle("GET /api/serverInfo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GetServerInfo(w, server_config, server_start_time)
	}))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server.ListenAndServe()
}
