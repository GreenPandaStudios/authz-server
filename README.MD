# This is a containerized Authorization Server Written in Go

## Client Configuration and Management

### Adding Clients

Clients are provided as | delimited string in the environment variable "CLIENTS".

The structure is expected to follow this format:

```text
client_id|client_secret|redirect_uri|clientid...
```

[See the Dockerfile for more details](./Dockerfile)

### Environment Variables

- `EXTERNAL_URL`: The external url of the application
- `PORT`: The port to expose in the image

# Authentication - Client Credentials Grant

The **Client Credentials Grant** is used when a client (such as a backend service) needs to authenticate itself and obtain an access token without user involvement.

## Request Example

Use the following `curl` command to request an access token:

```sh
curl -X POST http://localhost:8080/token \
     -H "Authorization: Basic $(echo -n 'client1:MyClientSecret1' | base64)" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     --data "grant_type=client_credentials"
```
