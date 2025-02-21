package main

import (
	"authz-server/client"
	"authz-server/jwks_endpoint"
	"authz-server/key_generator"
	"authz-server/token_endpoint"
	"authz-server/well_known"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        h.handlePost(w, r)
    case http.MethodGet:
        h.handleGet(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *RootHandler) handlePost(w http.ResponseWriter, r *http.Request) {  
    switch r.URL.Path {
    case "/authorize":
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Authorization endpoint"}`))
        return;
    case "/token":
        token_endpoint.HandleTokenRequest(w, r);
        return;
    default:
        http.Error(w, "Endpoint not found", http.StatusNotFound)
    }
}

func (h *RootHandler) handleGet(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/.well-known/openid-configuration":
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        well_known.HandleWellKnownEndpoint(w,r)
    case "/jwks":
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        jwks_endpoint.HandleJWKSEndpoint(w,r)
    default:
        http.Error(w, "Endpoint not found", http.StatusNotFound)
    }
}

func main() {

    err := os.RemoveAll("./.keys")
    if err != nil {
        log.Printf("Error deleting .keys file: %v", err)
    } else {
        fmt.Println(".keys file deleted successfully")
    }


    fmt.Println("Generating new keys...")
    _,_,err = key_generator.GenerateNewKeys()
    if err != nil {
        log.Fatalf("Error generating new keys: %v", err)
    }

    fmt.Println("Creating clients...")
    client.CreateClientsEnvFile()

    port, err := strconv.Atoi(os.Getenv("PORT"))
    if err != nil {
        port = 8000
    }

    fmt.Printf("Starting http server on port %d...\n", port)
    http.Handle("/", &RootHandler{})
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}