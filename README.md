# SMB Server and Client Test Project

This project provides a complete SMB server and client implementation for testing purposes.

## Project Structure

```
smb-test/
├── Dockerfile              # Docker image for Samba server
├── docker-compose.yml      # Docker Compose configuration
├── smb.conf               # Samba configuration file
├── startup.sh             # Container startup script
├── shared-data/           # Directory mounted as SMB share
│   └── hello.txt          # Sample file
├── main.go                # Go SMB client implementation
├── go.mod                 # Go module dependencies
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
- **Client connects to**: `127.0.0.1:4445`

## Prerequisites

- Docker and Docker Compose installed
- Go 1.x or later installed (for the client)

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

## Client Configuration

The client configuration is in main.go:12-18:

```go
serverAddress := "127.0.0.1:4445" // SMB server address and port
shareName := "public"             // Share name
username := "user"                // Username
password := "password"            // Password
fileNameToRead := "hello.txt"     // File to read
```

## Testing the Setup

1. Start the SMB server:
   ```bash
   docker-compose up -d
   ```

2. Wait a few seconds for the server to initialize

3. Run the client:
   ```bash
   go run main.go
   ```

4. Add more files to `shared-data/` directory to test reading different files

## Troubleshooting

### Connection Refused
- Ensure the Docker container is running: `docker ps`
- Check container logs: `docker-compose logs`
- Verify ports are not in use: `lsof -i :4445` or `lsof -i :1139`

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

Remove the container and network:
```bash
docker-compose down
```

Remove the Docker image:
```bash
docker rmi smb-test-samba
```