import base64
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.backends import default_backend
def handle_jwks_endpoint():
    with open('.keys/public_key.pem', 'rb') as key_file:
        pem_data = key_file.read()
        public_key_obj = serialization.load_pem_public_key(pem_data, backend=default_backend())
        numbers = public_key_obj.public_numbers()
        n = numbers.n
        n_bytes = n.to_bytes((n.bit_length() + 7) // 8, 'big')
        public_key = base64.urlsafe_b64encode(n_bytes).decode('utf-8').rstrip('=')
    
    return {
        "keys": [
            {
                "kty": "RSA",
                "e": "AQAB",
                "n": public_key,
                "alg": "RS256"
            }
        ]
    }