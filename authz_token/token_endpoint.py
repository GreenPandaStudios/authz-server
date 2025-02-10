from .access_token import issue_new_access_token
from .client import is_valid_client

def handle_token_request(request_body: dict, request_headers: dict) -> tuple[int, dict]:
    # TODO determine the grant type
    grant_type = request_body.get('grant_type')
    if grant_type == 'authorization_code':
        return _handle_authorization_code(request_body,request_headers)
    elif grant_type == "client_credentials":
        return _handle_client_credentials(request_body,request_headers)
    elif grant_type == 'refresh_token':
        return _handle_refresh_token(request_body,request_headers)
    else:
        return [400,
            {
                "error": "unsupported_grant_type",
            }
        ]

def _handle_authorization_code(request_body: dict, request_headers: dict) -> tuple[int, dict]:
    # TODO validate the authorization code
    return [400,
            {
                "error": "unsupported_grant_type",
            }]
    
def _handle_client_credentials(request_body: dict, request_headers: dict) -> tuple[int, dict]:
    # TODO validate the client credentials
    client_id = request_body.get('client_id')
    client_secret = request_headers.get('Authorization').split(' ')[1]
    if is_valid_client(client_id, client_secret):
        return [200, get_access_token(client_id)]
    return [400, {
        "error": "invalid_client"
    }]

def _handle_refresh_token(request_body: dict, request_headers: dict) -> tuple[int, dict]:
    # TODO validate the refresh token
    return [400, {}]




def get_access_token(client_id: str, expiration_time=3600) -> dict:
    [access_token, refresh_token] = issue_new_access_token(client_id, expiration_time)
    return  {
            "access_token": access_token,
            "token_type": "bearer",
            "expires_in": expiration_time,
            "refresh_token": refresh_token
        }