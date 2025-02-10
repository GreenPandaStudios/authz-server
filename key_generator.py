from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
import os

def generate_new_keys():
    # Generate RSA private key
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend()
    )

    # Serialize private key
    private_key_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.NoEncryption()
    )

    # Generate RSA public key
    public_key = private_key.public_key()

    # Serialize public key
    public_key_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )

    # Ensure the directory exists
    os.makedirs('.keys', exist_ok=True)

    # Save private key to a file
    with open('.keys/private_key.pem', 'wb') as f:
        f.write(private_key_pem)

    # Save public key to a file
    with open('.keys/public_key.pem', 'wb') as f:
        f.write(public_key_pem)
    return private_key_pem, public_key_pem