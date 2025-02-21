package jwks_endpoint

import (
	"authz-server/key_generator"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
)

type JWK struct {
    Kty string `json:"kty"`
    E   string `json:"e"`
    N   string `json:"n"`
    Alg string `json:"alg"`
}

func HandleJWKSEndpoint(w http.ResponseWriter, r *http.Request) {
    keys, err := key_generator.GetKeys()
    if err != nil {
        http.Error(w, "Failed to read public key file", http.StatusInternalServerError)
        return
    }

    var jwks []JWK
    for _, key := range keys {
        jwk, err := getJwk(key.PublicKey)
        if err != nil {
            http.Error(w, "Failed to parse public key", http.StatusInternalServerError)
            return
        }
        jwks = append(jwks, jwk)
    }

    pubKeys := map[string][]JWK{"keys": jwks}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(pubKeys)
}


func getJwk(pemData []byte) (JWK, error) {
    block, _ := pem.Decode(pemData)
    if block == nil || block.Type != "PUBLIC KEY" {
        return JWK{}, errors.New("Failed to parse public key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return JWK{}, err
    }

    rsaPub, ok := pub.(*rsa.PublicKey)
    if !ok {
        return JWK{}, errors.New("Failed to parse public key")
    }

    n := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())
    e := base64.RawURLEncoding.EncodeToString([]byte{1, 0, 1}) // 65537 in big-endian

    return JWK{
        Kty: "RSA",
        E:   e,
        N:   n,
        Alg: "RS256",
    },nil
}