package token_endpoint

import (
	"authz-server/client"
	"encoding/json"
	"net/http"
	"strings"
)

func HandleTokenRequest(w http.ResponseWriter, r *http.Request) {
    var requestBody map[string]interface{}
    var requestHeaders = r.Header

    contentType := r.Header.Get("Content-Type")
    if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
        err := r.ParseForm()
        if err != nil {
            http.Error(w, `{"error": "invalid_request"}`, http.StatusBadRequest)
            return
        }
        requestBody = make(map[string]interface{})
        for key, values := range r.Form {
            if len(values) > 0 {
                requestBody[key] = values[0]
            }
        }
    } else {
        err := json.NewDecoder(r.Body).Decode(&requestBody)
        if err != nil {
            http.Error(w, `{"error": "invalid_request"}`, http.StatusBadRequest)
            return
        }
    }
  

    grantType, ok := requestBody["grant_type"].(string)
    if !ok {
        http.Error(w, `{"error": "unsupported_grant_type"}`, http.StatusBadRequest)
        return
    }

    switch grantType {
    case "authorization_code":
        handleAuthorizationCode(w, requestBody, requestHeaders)
    case "client_credentials":
        handleClientCredentials(w, requestHeaders)
    case "refresh_token":
        handleRefreshToken(w, requestBody, requestHeaders)
    default:
        http.Error(w, `{"error": "unsupported_grant_type"}`, http.StatusBadRequest)
    }
}

func handleAuthorizationCode(w http.ResponseWriter, _ map[string]interface{}, _ http.Header) {
    // TODO: validate the authorization code
    http.Error(w, `{"error": "unsupported_grant_type"}`, http.StatusBadRequest)
}

func handleClientCredentials(w http.ResponseWriter, requestHeaders http.Header) {

    authHeader := requestHeaders.Get("Authorization")
    if authHeader == "" {
        http.Error(w, `{"error": "invalid_client"}`, http.StatusBadRequest)
        return
    }
    

    base64Auth := strings.Split(authHeader, " ")[1]
    clientID, clientSecret, _ := client.DecodeClientCredentials(base64Auth)
    
    if client.IsValid(clientID, clientSecret) {
        token := getAccessToken(clientID, 3600)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(token)
        return
    }

    http.Error(w, `{"error": "invalid_client"}`, http.StatusBadRequest)
}

func handleRefreshToken(w http.ResponseWriter, _ map[string]interface{}, _ http.Header) {
    // TODO: validate the refresh token
    http.Error(w, `{"error": "unsupported_grant_type"}`, http.StatusBadRequest)
}

func getAccessToken(clientID string, expirationTime int) map[string]interface{} {

    accessToken, refreshToken, err:= issueNewAccessToken(clientID, expirationTime, nil)
    if err != nil {
        return nil;
    }

    return map[string]interface{}{
        "access_token":  accessToken,
        "token_type":    "bearer",
        "expires_in":    expirationTime,
        "refresh_token": refreshToken,
    }
}