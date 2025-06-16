package main

import (
	"embed"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Burtcam/encounter-builder-backend/config"
	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/Burtcam/encounter-builder-backend/utils"
)

// struct encounter {
// 	difficulty string
// 	pSize      int
// 	level      int
// }

//go:embed static
var staticFS embed.FS

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	todos   = []Todo{}
	nextID  = 1
	todosMu sync.Mutex
)

func main() {
	cfg := config.Load()
	logger.Log.Info("Backend Initializing",
		slog.String("version", "1.0.0"),
		slog.String("env", "development"),
	)

	//setup the sync cron for the db.
	go utils.ManageDBSync(*cfg)
	// //TODO Remove this else everytime the ap starts it'll rebuild the db.
	// err := utils.KickOffSync(*cfg)
	// if err != nil {
	// 	logger.Log.Error(err.Error())
	// }
	// setup the UI
	// serve /api/todos
	http.HandleFunc("/api/todos", todosHandler)

	// serve all static files under / (index.html, alpine.js via CDN in HTML)
	fs := http.FileServer(http.FS(staticFS))
	http.Handle("/", fs)

	log.Println("listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todosMu.Lock()
	defer todosMu.Unlock()

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(todos)
	case http.MethodPost:
		var t struct{ Text string }
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "bad payload", http.StatusBadRequest)
			return
		}
		newTodo := Todo{ID: nextID, Text: t.Text}
		nextID++
		todos = append(todos, newTodo)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
