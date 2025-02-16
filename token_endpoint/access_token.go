package token_endpoint

import (
	"authz-server/well_known"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)


func issueNewAccessToken(clientID string, expiresIn int, scopes map[string]interface{}) (string, string, error) {
    if expiresIn == 0 {
        expiresIn = 3600
    }
    iss := well_known.GetIssuer()
    now := time.Now().UTC()

    accessTokenClaims := jwt.MapClaims{
        "exp": now.Add(time.Duration(expiresIn) * time.Second).Unix(),
        "aud": clientID,
        "iss": iss,
        "sub": clientID,
    }
    if scopes != nil {
        accessTokenClaims["scopes"] = scopes
    }

    refreshTokenClaims := jwt.MapClaims{
        "exp":        now.Add(30 * 24 * time.Hour).Unix(),
        "aud":        clientID,
        "iss":        iss,
        "sub":        clientID,
        "token_type": "refresh_token",
    }

    privateKeyData, err := os.ReadFile(".keys/private_key.pem")
    if err != nil {
        return "", "", err
    }

    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
    if err != nil {
        return "", "", err
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
    accessTokenString, err := accessToken.SignedString(privateKey)
    if err != nil {
        return "", "", err
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
    refreshTokenString, err := refreshToken.SignedString(privateKey)
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}