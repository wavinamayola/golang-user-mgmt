package main

import (
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	config "github.com/wavinamayola/user-management/internal/config"
	"github.com/wavinamayola/user-management/internal/routes"
	api "github.com/wavinamayola/user-management/internal/services"
	"github.com/wavinamayola/user-management/internal/storage"
	"github.com/wavinamayola/user-management/internal/utils"
)

//go:embed swagger-ui
var swaggerContent embed.FS
var cfg config.Config

func init() {
	var err error
	cfg, err = utils.LoadConfig("./env/")
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	log.Print("starting user management service")

	store, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize storage: %+v", err)
	}
	defer func() {
		log.Print("closing database conn")
		store.GetDB().Close()
	}()

	r := mux.NewRouter()
	fsys, _ := fs.Sub(swaggerContent, "swagger-ui")
	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(http.FS(fsys))))

	s := api.NewAPI(store)
	routes.SetupRoutes(r, s)
	r.Use(basicAuthMiddleware)

	httpListen(r)
}

func httpListen(r *mux.Router) {
	port := cfg.Server.Port
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("unable to listen on port %d", port)
	}

	go func() {
		log.Printf("listening on %v", l.Addr())
		if err := http.Serve(l, r); err != nil {
			log.Printf("failed to start http server: %+v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	l.Close()

	log.Print("stopping user management service")
}

func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "swagger") {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		payload, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != "admin" || pair[1] != "admin" {
			utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
