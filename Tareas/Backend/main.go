package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Estudiante struct {
	Nombre string `json:"nombre"`
	Carnet string `json:"carnet"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		estudiante := Estudiante{
			Nombre: "Pablo Javier Batz Contreras",
			Carnet: "201902698",
		}

		w.Header().Set("Content-Type", "application/json")

		// Configurar CORS
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},              // Permitir solicitudes desde cualquier origen
			AllowedMethods:   []string{"GET", "OPTIONS"}, // Permitir los métodos GET y OPTIONS
			AllowedHeaders:   []string{"Content-Type"},   // Permitir el encabezado Content-Type
			AllowCredentials: true,                       // Permitir credenciales si es necesario
		})

		// Manejar la ruta con CORS habilitado
		handlerWithCORS := c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(estudiante)
		}))

		handlerWithCORS.ServeHTTP(w, r)
	})

	port := 8080
	fmt.Printf("Servidor escuchando en el puerto %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
