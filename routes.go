package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"runtime"
	"time"
)

func Home(w http.ResponseWriter, tmpl *template.Template, user_config map[string]interface{}) {
	tmpl.ExecuteTemplate(w, "index.html", user_config)
}

func GetServerInfo(w http.ResponseWriter, server_config map[string]interface{}, server_start_time time.Time) {
	z, _ := time.Now().Local().Zone()

	j, err := json.Marshal(ServerInfo{
		GoVersion:      runtime.Version(),
		ServerTimezone: z,
		ServerUptime:   time.Now().Sub(server_start_time).String(),

		ConfigInfo: server_config,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(j)
}
