# SMB Server and Client Test Project with Tailscale Proxy

This project provides a complete SMB server and client implementation with Tailscale proxy support for secure remote access.

## Project Structure

```
smb-test-with-tailscale/
├── Dockerfile              # Docker image for Samba server
├── docker-compose.yml      # Docker Compose configuration
├── smb.conf               # Samba configuration file
├── startup.sh             # Container startup script
├── shared-data/           # Directory mounted as SMB share
│   └── hello.txt          # Sample file
├── main.go                # Go SMB client implementation
├── go.mod                 # Go module dependencies
├── tailsacle-proxy/       # Tailscale proxy configuration
│   ├── config.yaml        # Port mapping configuration
│   └── README.md          # Proxy setup instructions
└── README.md              # This file
```

## SMB Server Configuration

The Samba server is configured with:
- **Share name**: `public`
- **Username**: `user`
- **Password**: `password`
- **Share path**: `/shared/public` (mounted from `./shared-data`)
- **Ports**:
  - Host port `1139` → Container port `139` (NetBIOS)
  - Host port `4445` → Container port `445` (SMB)

## Tailscale Proxy Configuration

This project includes a Tailscale proxy setup that allows secure remote access to SMB servers over Tailscale network.

### Architecture

```
SMB Client (main.go)
    ↓ (localhost:4789)
Tailscale Proxy Container
    ↓ (Tailscale network)
100.105.49.5:4445 (Remote SMB Server)
```

### Configuration

The proxy configuration is in `tailsacle-proxy/config.yaml`:

```yaml
portMappings:
  4789: "100.105.49.5:4445"
```

- **Local port**: `4789` - The port the client connects to
- **Remote target**: `100.105.49.5:4445` - The SMB server on Tailscale network

### Running the Tailscale Proxy

1. Ensure you have a Tailscale auth key
2. Navigate to the `tailsacle-proxy` directory
3. Run the proxy container:

```bash
docker run -d \
  -e TS_AUTH_KEY='<your_tailscale_auth_key>' \
  -v ./config.yaml:/config.yaml \
  -v tailscale-local:/.local \
  -v tailscale-run:/var/run/tailscale \
  -p 4789:4789 \
  --name tailscale-proxy \
  vajiraprabuddhaka/tailscale-proxy:latest
```

### Client Configuration with Tailscale Proxy

When using the Tailscale proxy, the client connects to `127.0.0.1:4789` (see main.go:14):

```go
serverAddress := "127.0.0.1:4789" // Connect via Tailscale proxy
```

## Prerequisites

- Docker and Docker Compose installed
- Go 1.x or later installed (for the client)
- For Tailscale proxy usage:
  - Tailscale account with an auth key
  - Remote SMB server accessible via Tailscale network
  - Tailscale proxy Docker image: `vajiraprabuddhaka/tailscale-proxy:latest`

## Starting the SMB Server

### Option 1: Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

To view logs:
```bash
docker-compose logs -f
```

To stop the server:
```bash
docker-compose down
```

### Option 2: Using Docker directly

Build the image:
```bash
docker build -t smb-server .
```

Run the container:
```bash
docker run -d \
  -p 1139:139 \
  -p 4445:445 \
  -v $(pwd)/shared-data:/shared/public \
  --name smb-server \
  smb-server
```

## Running the SMB Client

The Go client connects to the SMB server and reads files from the share.

### Install Dependencies

```bash
go mod download
```

### Run the Client

Make sure the SMB server is running, then:

```bash
go run main.go
```

Expected output:
```
Successfully connected to SMB server at 127.0.0.1.
Successfully mounted the 'public' share.

Files in share:
- hello.txt

Content of 'hello.txt':
Hello from SMB Server!
...
```

## Usage Scenarios

### Scenario 1: Local Testing (Direct Connection)

For local testing without Tailscale:

1. Update main.go to connect directly:
   ```go
   serverAddress := "127.0.0.1:4445"
   ```

2. Start the SMB server:
   ```bash
   docker-compose up -d
   ```

3. Run the client:
   ```bash
   go run main.go
   ```

### Scenario 2: Remote Access (via Tailscale Proxy)

For connecting to a remote SMB server over Tailscale:

1. Update `tailsacle-proxy/config.yaml` with your remote SMB server's Tailscale IP:
   ```yaml
   portMappings:
     4789: "<tailscale-ip>:4445"
   ```

2. Start the Tailscale proxy (see "Running the Tailscale Proxy" section above)

3. Ensure main.go is configured to use the proxy:
   ```go
   serverAddress := "127.0.0.1:4789"
   ```

4. Run the client:
   ```bash
   go run main.go
   ```

## Testing the Setup

### With Tailscale Proxy (Remote Access)

1. Start the Tailscale proxy:
   ```bash
   cd tailsacle-proxy
   docker run -d \
     -e TS_AUTH_KEY='<your_tailscale_auth_key>' \
     -v ./config.yaml:/config.yaml \
     -v tailscale-local:/.local \
     -v tailscale-run:/var/run/tailscale \
     -p 4789:4789 \
     --name tailscale-proxy \
     vajiraprabuddhaka/tailscale-proxy:latest
   ```

2. Verify the proxy is running:
   ```bash
   docker logs tailscale-proxy
   ```

3. Run the client:
   ```bash
   go run main.go
   ```

### Without Tailscale (Local Testing)

1. Start the SMB server:
   ```bash
   docker-compose up -d
   ```

2. Update main.go to use direct connection (`127.0.0.1:4445`)

3. Run the client:
   ```bash
   go run main.go
   ```

## Troubleshooting

### Connection Refused
- Ensure the Docker container is running: `docker ps`
- Check container logs: `docker-compose logs`
- Verify ports are not in use: `lsof -i :4445` or `lsof -i :1139`
- If using Tailscale proxy, ensure port 4789 is not in use: `lsof -i :4789`

### Tailscale Proxy Issues

#### Proxy not connecting to Tailscale
- Check the Tailscale proxy logs: `docker logs tailscale-proxy`
- Verify the TS_AUTH_KEY is valid and not expired
- Ensure the Tailscale service is running: `docker exec tailscale-proxy tailscale status`

#### Cannot reach remote SMB server
- Verify the remote Tailscale IP is correct in `config.yaml`
- Check if the remote server is accessible: `docker exec tailscale-proxy ping <tailscale-ip>`
- Ensure the remote SMB server is running and accepting connections on port 4445
- Verify firewall rules on the remote server allow Tailscale connections

#### Port mapping not working
- Check the config.yaml format is correct
- Restart the Tailscale proxy container: `docker restart tailscale-proxy`
- Verify the proxy is listening on port 4789: `netstat -an | grep 4789`

### Authentication Failed
- Verify username and password match in both server and client
- Check Samba logs in the container: `docker exec smb-server cat /var/log/samba/log.smbd`

### File Not Found
- Ensure the file exists in the `shared-data/` directory
- Check file permissions
- Update `fileNameToRead` in main.go to match the actual filename

## Modifying the Shared Files

Any files placed in the `shared-data/` directory will be accessible via the SMB share. The directory is mounted as a Docker volume.

## Security Note

This setup uses simple authentication and is intended for development and testing only. For production use, implement proper security measures including:
- Strong passwords
- Encrypted connections
- Firewall rules
- User access controls

## Cleaning Up

### Remove SMB Server

Remove the container and network:
```bash
docker-compose down
```

Remove the Docker image:
```bash
docker rmi smb-test-samba
```

### Remove Tailscale Proxy

Stop and remove the Tailscale proxy container:
```bash
docker stop tailscale-proxy
docker rm tailscale-proxy
```

Remove Tailscale volumes (optional):
```bash
docker volume rm tailscale-local tailscale-run
```

### Complete Cleanup

To remove everything (SMB server, Tailscale proxy, and all volumes):
```bash
# Stop and remove all containers
docker-compose down
docker stop tailscale-proxy
docker rm tailscale-proxy

# Remove images
docker rmi smb-test-samba
docker rmi vajiraprabuddhaka/tailscale-proxy:latest

# Remove volumes
docker volume rm tailscale-local tailscale-run
```