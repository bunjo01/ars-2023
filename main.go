// Config Service API
//
//	   Title: Config Service API
//
//	   Schemes: http
//		  Version: 0.0.1
//		  BasePath: /
//
//		  Produces:
//			- application/json
//
// swagger:meta
package main

import (
	cdb "ars-2023/configdatabase"
	pm "ars-2023/prometheus"
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// Router instance
	router := mux.NewRouter()
	router.StrictSlash(true)
	// db api instance
	store, err := cdb.New()
	if err != nil {
		log.Fatal(err)
	}
	// server instance
	server := configServer{
		store: store,
		Keys:  make(map[string]string),
	}
	//Config operation handlers
	router.HandleFunc("/config/", pm.CountCreateConfig(server.createConfigHandler)).Methods("POST")
	router.HandleFunc("/config/all/", pm.CountGetAllConfig(server.getAllConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/all/", pm.CountGetConfigVersion(server.getConfigVersionsHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/all/", pm.CountDelConfigVersion(server.delConfigVersionsHandler)).Methods("DELETE")
	router.HandleFunc("/config/{id}/{version}/", pm.CountGetConfig(server.getConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}/{version}/", pm.CountDelConfig(server.delConfigHandler)).Methods("DELETE")
	// Group operation handlers
	router.HandleFunc("/group/", pm.CountCreateGroup(server.createGroupHandler)).Methods("POST")
	router.HandleFunc("/group/all/", pm.CountGetAllGroup(server.getAllGroupHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/all/", pm.CountGetGroupVersion(server.getGroupVersionsHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/all/", pm.CountDelGroupVersion(server.delGroupVersionsHandler)).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}/", pm.CountGetGroup(server.getGroupHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/", pm.CountDelGroup(server.delGroupHandler)).Methods("DELETE")
	// Group append handler
	router.HandleFunc("/group/append/", pm.CountAppendGroup(server.appendGroupHandler)).Methods("POST")
	// Label operations handlers
	router.HandleFunc("/group/{id}/{version}/{labels}/", pm.CountGetConfigByLabels(server.getConfigsByLabel)).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/{new}/{labels}/", pm.CountDelConfigByLabels(server.delConfigsByLabel)).Methods("DELETE")
	// Swagger handler
	router.HandleFunc("/swagger.yaml", server.swaggerHandler).Methods("GET")
	// Metrics handler
	router.Path("/metrics").Handler(pm.MetricsHandler())

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("...server starting...")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("...server shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("...server stopped...")

}
