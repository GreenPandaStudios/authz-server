package key_generator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateNewKeys() ([]byte, []byte, error) {
    // Generate RSA private key
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, err
    }

    // Serialize private key to PEM format
    privateKeyPEM := pem.EncodeToMemory(
        &pem.Block{
            Type:  "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
        },
    )

    // Generate RSA public key
    publicKey := &privateKey.PublicKey

    // Serialize public key to PEM format
    publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        return nil, nil, err
    }
    publicKeyPEM := pem.EncodeToMemory(
        &pem.Block{
            Type:  "PUBLIC KEY",
            Bytes: publicKeyBytes,
        },
    )

    // Ensure the directory exists
    err = os.MkdirAll(".keys", os.ModePerm)
    if err != nil {
        return nil, nil, err
    }

    // Save private key to a file
    err = os.WriteFile(".keys/private_key.pem", privateKeyPEM, 0600)
    if err != nil {
        return nil, nil, err
    }

    // Save public key to a file
    err = os.WriteFile(".keys/public_key.pem", publicKeyPEM, 0644)
    if err != nil {
        return nil, nil, err
    }

    return privateKeyPEM, publicKeyPEM, nil
}