# This is a containerized functioning Authorization Server Written in Python

## Client Configuration and Management

### Adding Clients

Clients are provided as | delimited string in the environment variable "CLIENTS".

The structure is expected to follow this format:

```text
client_id|client_secret|redirect_uri|clientid...
```

[See the Dockerfile for more details](./Dockerfile)

### Environment Variables

- `PORT`: The port on which the server will run. Default is `8080`.
- `DOMAIN`: The domain name for the server. Default is `localhost`.
- `PROTOCOL`: The protocol to be used by the server. Default is `http`.
