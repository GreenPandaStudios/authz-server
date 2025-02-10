import os

def handle_well_known_endpoint():
    return {
        "issuer": get_issuer(),
        "authorization_endpoint": get_authorization_endpoint(),
        "token_endpoint": get_token_endpoint(),
        "jwks_uri": get_jwks_uri(),
        "scopes_supported": get_scopes_supported(),
        "response_types_supported": get_response_types_supported(),
        "grant_types_supported": get_grant_types_supported(),
        "token_endpoint_auth_methods_supported": get_token_endpoint_auth_methods_supported(),
        "code_challenge_methods_supported": get_code_challenge_methods_supported(),
        "token_endpoint_auth_signing_alg_values_supported": get_token_endpoint_auth_signing_alg_values_supported(),
        "id_token_signing_alg_values_supported": get_id_token_signing_alg_values_supported(),
        "request_object_signing_alg_values_supported": get_request_object_signing_alg_values_supported(),
        "response_modes_supported": get_response_modes_supported(),
        "subject_types_supported": get_subject_types_supported(),
        "userinfo_signing_alg_values_supported": get_userinfo_signing_alg_values_supported()
    }
def get_issuer():
    return __get_base_url()
def get_authorization_endpoint():
    return f"{__get_base_url()}/authorize"
def get_token_endpoint():
    return f"{__get_base_url()}/token"
def get_jwks_uri():
    return f"{__get_base_url()}/jwks"
def get_scopes_supported():
    return ["read", "write", "delete"]
def get_response_types_supported():
    return ["code", "token"]
def get_grant_types_supported():
    return ["authorization_code", "client_credentials", "refresh_token"]
def get_token_endpoint_auth_methods_supported():
    return ["client_secret_basic", "client_secret_post"]
def get_code_challenge_methods_supported():
    return ["plain", "S256"]
def get_token_endpoint_auth_signing_alg_values_supported():
    return ["RS256"]
def get_id_token_signing_alg_values_supported():
    return ["RS256"]
def get_request_object_signing_alg_values_supported():
    return ["RS256"]
def get_response_modes_supported():
    return ["query", "fragment", "form_post"]
def get_subject_types_supported():
    return ["public"]
def get_userinfo_signing_alg_values_supported():
    return ["RS256"]

def __get_base_url():
    domain = os.getenv("DOMAIN", "localhost")
    port = os.getenv("PORT", "8000")
    protocol = os.getenv("PROTOCOL", "https")
    
    if (protocol == "https" and port == "443") or (protocol == "http" and port == "80"):
        return f"{protocol}://{domain}"
    return f"{protocol}://{domain}:{port}"