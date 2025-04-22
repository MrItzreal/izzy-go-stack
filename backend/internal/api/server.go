package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/your-username/your-repo/internal/config"
	"github.com/your-username/your-repo/internal/database"
	"github.com/your-username/your-repo/internal/handlers"
)

// Server holds the HTTP server and its dependencies
type Server struct {
	Router *chi.Mux
	Config *config.Config
	DB     *database.DB
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, db *database.DB) *Server {
	server := &Server{
		Router: chi.NewRouter(),
		Config: cfg,
		DB:     db,
	}

	// Set up middleware
	server.Router.Use(middleware.Logger)
	server.Router.Use(middleware.Recoverer)
	server.Router.Use(middleware.RealIP)
	server.Router.Use(middleware.RequestID)

	// Set up CORS
	server.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers
	productHandler := handlers.NewProductHandler(db)
	userHandler := handlers.NewUserHandler(db)
	orderHandler := handlers.NewOrderHandler(db, cfg)

	// Set up routes
	server.Router.Route("/api", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/products", productHandler.List)
			r.Get("/products/{id}", productHandler.Get)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth("api", map[string]string{
				"api": "secret",
			}))

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.List)
				r.Post("/", userHandler.Create)
				r.Get("/{id}", userHandler.Get)
				r.Put("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})

			// Order routes
			r.Route("/orders", func(r chi.Router) {
				r.Get("/", orderHandler.List)
				r.Post("/", orderHandler.Create)
				r.Get("/{id}", orderHandler.Get)
				r.Put("/{id}", orderHandler.Update)
				r.Delete("/{id}", orderHandler.Delete)
			})

			// Product management routes
			r.Route("/admin/products", func(r chi.Router) {
				r.Post("/", productHandler.Create)
				r.Put("/{id}", productHandler.Update)
				r.Delete("/{id}", productHandler.Delete)
			})
		})
	})

	return server
}
