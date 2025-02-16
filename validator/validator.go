package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt/v4"
)



var (
	keyCache = make(map[string]([]interface{}))
	cacheMutex = &sync.Mutex{}
)

func ValidateAccessToken(accessToken string) (bool, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(jwt.MapClaims);

	if !ok {
		return false, errors.New("invalid token claims")
	}

	issuer, err := validateIssuer(claims)

	if err != nil {
		return false, err
	}

	keys, err := getKeys(issuer)

	if err != nil {
		return false, err
	}

	return validateTokenWithKeys(accessToken, keys,issuer)
}



func validateTokenWithKeys(accessToken string, keys []interface{}, issuer string) (bool, error) {
	return validateTokenWithKeysRecursive(accessToken, keys, issuer, false)
}

func validateTokenWithKeysRecursive(accessToken string, keys []interface{}, issuer string, base bool) (bool, error) {
	for _, key := range keys {
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err == nil && token.Valid {
			return true, nil
		}
	}

	if base {
		return false, errors.New("invalid token")
	}
	
	// Otherwise, update the keys and try one more time
	keys,err := getUpdatedkeys(issuer)
	if err != nil {
		return false, err
	}

	return validateTokenWithKeysRecursive(accessToken, keys, issuer, true)
}


func validateIssuer(claims jwt.MapClaims) (string, error) {
	if issuer, ok := claims["iss"].(string); ok {
		return issuer, nil
	}
	return "", errors.New("issuer not found in token")
}

func getKeys(issuer string) (([]interface{}), error) {
	cacheMutex.Lock()
	if jwks, ok := keyCache[issuer]; ok {
		cacheMutex.Unlock()
		return jwks, nil
	}
	cacheMutex.Unlock()

	keys, err := getUpdatedkeys(issuer)
	if err != nil {
		return nil, err
	}
	cacheMutex.Lock()
	keyCache[issuer] = keys
	cacheMutex.Unlock()
	return keys, nil
}

func getUpdatedkeys(issuer string) ([]interface{}, error) {
	
	// Remove trailing slash from issuer
	if issuer[len(issuer)-1] == '/' {
		issuer = issuer[:len(issuer)-1]
	}

	url := fmt.Sprintf("%s/.well-known/jwks.json", issuer)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS from issuer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get JWKS from issuer: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response body: %w", err)
	}

	var jwks map[string]interface{}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWKS response body: %w", err)
	}

	// Find the JWKS URI in the response body
	jwksURI, ok := jwks["jwks_uri"].(string)
	if !ok {
		return nil, errors.New("jwks_uri not found in response body")
	}

	// Make a request to the JWKS URI to get the actual JWKS
	resp, err = http.Get(jwksURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS from JWKS URI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get JWKS from JWKS URI: status code %d", resp.StatusCode)
	}

	body=nil
	err=nil

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response body: %w", err)
	}
	// Return the keys from the JWKS URI
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWKS response body: %w", err)
	}

	keys, ok := jwks["keys"].([]interface{})
	if !ok {
		return nil, errors.New("keys not found in JWKS response body")
	}
	if keys == nil {
		return nil, errors.New("keys not found in JWKS response body")
	}

	return keys, nil
}
