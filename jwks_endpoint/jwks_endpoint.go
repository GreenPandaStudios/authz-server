package jwks_endpoint

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"os"
)

type JWK struct {
    Kty string `json:"kty"`
    E   string `json:"e"`
    N   string `json:"n"`
    Alg string `json:"alg"`
}

func HandleJWKSEndpoint(w http.ResponseWriter, r *http.Request) {
    pemData, err := os.ReadFile(".keys/public_key.pem")
    if err != nil {
        http.Error(w, "Failed to read public key file", http.StatusInternalServerError)
        return
    }

    block, _ := pem.Decode(pemData)
    if block == nil || block.Type != "PUBLIC KEY" {
        http.Error(w, "Failed to decode PEM block", http.StatusInternalServerError)
        return
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        http.Error(w, "Failed to parse public key", http.StatusInternalServerError)
        return
    }

    rsaPub, ok := pub.(*rsa.PublicKey)
    if !ok {
        http.Error(w, "Not an RSA public key", http.StatusInternalServerError)
        return
    }

    n := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())
    e := base64.RawURLEncoding.EncodeToString([]byte{1, 0, 1}) // 65537 in big-endian

    jwk := JWK{
        Kty: "RSA",
        E:   e,
        N:   n,
        Alg: "RS256",
    }

    keys := map[string][]JWK{"keys": {jwk}}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(keys)
}