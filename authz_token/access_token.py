import jwt
import datetime
from well_known import get_issuer

def issue_new_access_token(client_id: str, expires_in: int = 3600, scopes: dict = None) -> tuple[str,str]:
    # TODO create a new access token
    iss = get_issuer()
    now = datetime.datetime.now(tz=datetime.timezone.utc)
    access_token_payload = {
        'exp': now + datetime.timedelta(seconds=expires_in),
        'aud': client_id,  # replace 'your_audience' with the actual audience
        'iss':iss,
        'sub': client_id
    }
    if scopes:
        access_token_payload['scopes'] = scopes
    refresh_token_payload = {
        'exp': now + datetime.timedelta(days=30),
        'aud': client_id,
        'iss': iss,
        'sub': client_id,
        'token_type': 'refresh_token'
    }

    with open('.keys/private_key.pem', 'r') as key_file:
        private_key = key_file.read()

    access_token = jwt.encode(access_token_payload, private_key, algorithm='RS256')
    refresh_token = jwt.encode(refresh_token_payload, private_key, algorithm='RS256')
    return access_token, refresh_token


def valdiate_token(token: str, is_refresh_token: bool = False) -> dict:
    # TODO validate the access token
    with open('.keys/public_key.pem', 'r') as key_file:
        public_key = key_file.read()
  
    try:
        payload = jwt.decode(token, public_key, algorithms=['RS256'])
        if is_refresh_token:
            if payload.get('token_type') != 'refresh_token':
                return None
        else:
            if payload.get('token_type') == 'refresh_token':
                return None
            
        return payload
    except jwt.ExpiredSignatureError:
        return None
    except jwt.InvalidTokenError:
        return None
    except jwt.InvalidSignatureError:
        return None
    except jwt.InvalidIssuerError:
        return None
    except jwt.InvalidAudienceError:
        return None
    