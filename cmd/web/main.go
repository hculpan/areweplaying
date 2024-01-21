package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hculpan/areweplaying/pkg/utils"
	"github.com/joho/godotenv"
)

var appKey []byte

func main() {
	var port string

	err := godotenv.Load()
	if err != nil {
		port = "8080"
	}

	port = os.Getenv("PORT")
	appKey = []byte(os.Getenv("APP_KEY"))
	if len(appKey) == 0 {
		log.Fatal("unable to load APP_KEY")
	}

	logFile := os.Getenv("WEB_LOG_FILE")
	if len(logFile) > 0 {
		utils.SetLogFile(logFile)
	}

	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	FileServer(r, "/static", http.Dir(filesDir))

	routes(r)

	log.Default().Printf("Starting up server on port :%s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
