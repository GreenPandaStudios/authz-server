package key_generator

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

   return saveNewKeys(privateKeyPEM, publicKeyPEM);
}



func saveNewKeys(privateKeyPEM, publicKeyPEM []byte) ([]byte, []byte, error) {
    sa := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    if sa == "" {
        return saveNewKeysToFile(privateKeyPEM, publicKeyPEM)
    } else {
        return saveNewKeysToFirestore(privateKeyPEM, publicKeyPEM, sa)
    }
}


func saveNewKeysToFile(privateKeyPEM, publicKeyPEM []byte) ([]byte, []byte, error) {
    // Ensure the directory exists
    err := os.MkdirAll(".keys", os.ModePerm)
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


type Key struct {
    PrivateKey []byte `firestore:"private_key"`
    PublicKey  []byte `firestore:"public_key"`
    CreatedAt  map[string]interface{} `firestore:"created_at"`
}


type KeyPair struct {
    PrivateKey []byte
    PublicKey  []byte
}


// saveNewKeysToFirestore saves the provided private and public keys to Firestore.
// It first checks if there are existing keys in the "current_keys" document. If found,
// it moves the existing keys to the "last_key" document. Then, it saves the new keys
// to the "current_keys" document.
//
// Parameters:
// - privateKeyPEM: The private key in PEM format.
// - publicKeyPEM: The public key in PEM format.
// - creds: The Firestore credentials in JSON format.
//
// Returns:
// - privateKeyPEM: The saved private key in PEM format.
// - publicKeyPEM: The saved public key in PEM format.
// - error: An error if any occurred during the process, otherwise nil.
func saveNewKeysToFirestore(privateKeyPEM, publicKeyPEM []byte, creds string) ([]byte, []byte, error) {
    ctx := context.Background()
    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        return nil, nil, errors.New("GOOGLE_CLOUD_PROJECT environment variable is not set")
    }
    // Create a new Firestore client with the provided credentials
    client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(creds))
    if err != nil {
        return nil, nil, err
    }

    // Get the current keys document from Firestore
    doc, err := client.Collection("keys").Doc("current_keys").Get(ctx)
    if err != nil && status.Code(err) != codes.NotFound {
        defer client.Close()
        return nil, nil, err
    }

    // If the current keys document exists, move it to the "last_key" document
    if doc.Exists() {
        var currentKey Key
        err = doc.DataTo(&currentKey)
        if err != nil {
            return nil, nil, err
        }

        _, err = client.Collection("keys").Doc("last_key").Set(ctx, currentKey)
        if err != nil {
            defer client.Close()
            return nil, nil, err
        }
    }

    // Create a new key object with the provided private and public keys
    newKey := Key{
        PrivateKey: privateKeyPEM,
        PublicKey:  publicKeyPEM,
        CreatedAt:  map[string]interface{}{"timestamp": firestore.ServerTimestamp},
    }

    // Save the new keys to the "current_keys" document
    _, err = client.Collection("keys").Doc("current_keys").Set(ctx, newKey)
    if err != nil {
        defer client.Close()
        return nil, nil, err
    }

    // Close the Firestore client
    defer client.Close()

    return privateKeyPEM, publicKeyPEM, nil
}


// GetKeys retrieves a list of KeyPair objects. It first checks if the
// environment variable "GOOGLE_APPLICATION_CREDENTIALS" is set. If it is,
// the function fetches the keys from Firestore using the provided service
// account credentials. If the environment variable is not set, it falls
// back to retrieving the keys from a local file.
//
// Returns:
//   - []KeyPair: A slice of KeyPair objects.
//   - error: An error if the keys could not be retrieved.
func GetKeys() ([]KeyPair, error) {
   sa := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    if sa == "" {
        return getKeysFromFile()
    } else {
        return getKeysFromFirestore(sa)
    }
}


func getKeysFromFile() ([]KeyPair, error) {
    privateKeyPEM, err := os.ReadFile(".keys/private_key.pem")
    if err != nil {
        return nil, err
    }

    publicKeyPEM, err := os.ReadFile(".keys/public_key.pem")
    if err != nil {
        return nil, err
    }

    keyPair := KeyPair{
        PrivateKey: privateKeyPEM,
        PublicKey:  publicKeyPEM,
    }

    return []KeyPair{keyPair}, nil
}


func getKeysFromFirestore(creds string) ([]KeyPair, error) {
    ctx := context.Background()

    keys := make([]KeyPair, 0)

    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        return nil, errors.New("GOOGLE_CLOUD_PROJECT environment variable is not set")
    }
    client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(creds)))
    if err != nil {
        return nil, err
    }

    doc, err := client.Collection("keys").Doc("current_keys").Get(ctx)
    if err != nil {
        defer client.Close()
        return nil, err
    }

    keys = append(keys, KeyPair{
        PrivateKey: doc.Data()["private_key"].([]byte),
        PublicKey:  doc.Data()["public_key"].([]byte),
    })

    doc, err = client.Collection("keys").Doc("last_key").Get(ctx)
    if err != nil && status.Code(err) != codes.NotFound {
        defer client.Close()
        return nil, err
    }
    defer client.Close()

    if err == nil {
        keys = append(keys, KeyPair{
            PrivateKey: doc.Data()["private_key"].([]byte),
            PublicKey:  doc.Data()["public_key"].([]byte),
        })
    }

    return keys, nil
}