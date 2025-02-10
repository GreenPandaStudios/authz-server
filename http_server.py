from http.server import SimpleHTTPRequestHandler, HTTPServer
from authz_token.token_endpoint import handle_token_request
from authz_token.client import create_clients_env_file
from key_generator import generate_new_keys
from well_known.well_known_endpoint import handle_well_known_endpoint
from jwks_endpoint.jwks_endpoint import handle_jwks_endpoint
import json
import os

class RootHandler(SimpleHTTPRequestHandler):
    def do_POST(self):
        try: 
            print(f'POST request,\nPath: {self.path}\nHeaders:\n{self.headers}\n')
            post_data = self.process_request_body()
            print(f'Body:\n{post_data}\n')
            if post_data==None:
                return  
            if self.path == '/authorize':            
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(b'{"message": "Authorization endpoint"}')
            elif self.path == '/token':
                [status,body] = handle_token_request(post_data, dict(self.headers))
                self.send_response(status)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(body).encode('utf-8'))
            else:
                self.send_response(404)
                self.send_header('Content-type', 'text/html')
                self.end_headers()
                self.wfile.write(b'Endpoint not found')
        except Exception as e:
            print(f'Error: {e}')
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(b'{"error": "Internal server error"}')

    def do_GET(self):
        try:
            if self.path == '/.well-known/openid-configuration':
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(handle_well_known_endpoint()).encode('utf-8'))
            elif self.path == '/jwks':
                self.send_response(200)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(handle_jwks_endpoint()).encode('utf-8'))
        except Exception as e:
            print(f'Error: {e}')
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(b'{"error": "Internal server error"}')

    def process_request_body(self)-> dict:
        content_length = int(self.headers['Content-Length']) if 'Content-Length' in self.headers else 0
        post_data = self.rfile.read(content_length)
        content_type = self.headers.get('Content-Type')

        if content_type == 'application/x-www-form-urlencoded':
            post_data = post_data.decode('utf-8')
            post_data = dict(x.split('=') for x in post_data.split('&'))
            return post_data
        elif content_type == 'application/json':
            post_data = json.loads(post_data.decode('utf-8'))
            return post_data
        else:
            self.send_response(400)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(b'{"error": "Unsupported Content-Type"}')
            return None


def run(server_class=HTTPServer, handler_class=RootHandler, port=int(os.getenv('PORT', 8000))):

    print('Generating new keys...')
    generate_new_keys()

    print('Creating clients...')
    create_clients_env_file()

    server_address = ('', port)
    httpd = server_class(server_address, handler_class)
    print(f'Starting http server on port {port}...')
    httpd.serve_forever()

if __name__ == "__main__":
    run()