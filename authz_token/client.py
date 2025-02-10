import os
import json
import base64

def is_valid_client(client_id: str, client_secret: str) -> bool:
    with open('.keys/clients.json', 'r') as f:
        clients = json.load(f)
    if clients.get(client_id):
        return clients[client_id]['client_secret'] == client_secret
    return False


def create_clients_env_file():
    clients_env = os.environ.get('CLIENTS')
    if not clients_env:
        return False
    clients = {}
    client_entries = clients_env.split('|')
    for i in range(0, len(client_entries), 3):
        clients[client_entries[i]] = {
            'client_secret': client_entries[i + 1],
            'redirect_uri': client_entries[i + 2]
        }
    
    with open('.keys/clients.json', 'w') as f:
        f.write(json.dumps(clients))