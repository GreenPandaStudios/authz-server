package token_endpoint

import (
	"authz-server/key_generator"
	"authz-server/well_known"
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

    keys,  err := key_generator.GetKeys()
    if err != nil {
        return "", "", err
    }

    //Use the first key pair to sign the tokens
    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keys[0].PrivateKey)
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