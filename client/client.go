package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Client struct {
    ClientSecret string `json:"client_secret"`
    RedirectURI  string `json:"redirect_uri"`
}

func IsValid(clientID, clientSecret string) bool {
    file, err := os.Open(".keys/clients.json")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return false
    }
    defer file.Close()

    var clients map[string]Client
    if err := json.NewDecoder(file).Decode(&clients); err != nil {
        fmt.Println("Error decoding JSON:", err)
        return false
    }

    client, exists := clients[clientID]
    return exists && client.ClientSecret == clientSecret
}

func DecodeClientCredentials(base64Auth string) (string, string, error) {
    authData, err := base64.StdEncoding.DecodeString(base64Auth)
    if err != nil {
        return "", "", err
    }

    auth := strings.Split(string(authData), ":")
    return auth[0], auth[1], nil
}


func CreateClientsEnvFile() bool {
    clientsEnv := os.Getenv("CLIENTS")
    if clientsEnv == "" {
        return false
    }

    clients := make(map[string]Client)
    clientEntries := strings.Split(clientsEnv, "|")
    for i := 0; i < len(clientEntries); i += 3 {
        clients[clientEntries[i]] = Client{
            ClientSecret: clientEntries[i+1],
            RedirectURI:  clientEntries[i+2],
        }
    }

    file, err := os.Create(".keys/clients.json")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return false
    }
    defer file.Close()

    if err := json.NewEncoder(file).Encode(clients); err != nil {
        fmt.Println("Error encoding JSON:", err)
        return false
    }

    return true
}