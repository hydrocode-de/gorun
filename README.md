# gorun

A Go HTTP server that provides a simple interface to run tool-spec containers. It allows you to execute containerized tools through a REST API, with built-in authentication and file management.

## Features

- 🚀 Run containerized tools through a simple HTTP API
- 🔒 Built-in authentication system
- 📁 File upload and download capabilities
- 🏗️ Support for custom mount points
- 📊 Job status tracking
- 🎨 Modern web interface for job management
- 🔄 Persistent storage for jobs and credentials

## Prerequisites

### Development
- Go 1.23 or later
- Node.js and npm (for frontend development)
- Make
- Docker (for testing)

### Deployment
- Docker
- Access to Docker daemon (via socket or TCP)
- Persistent storage for data

## Quick Start

### Using Docker (Recommended)

1. Build the Docker image:
   ```bash
   docker build -t gorun .
   ```

2. Run the container:
   ```bash
   docker run -d \
     -p 8080:8080 \
     -v /path/to/local/data:/data/gorun \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -e GORUN_PORT=8080 \
     -e GORUN_SECRET=your-secret-here \
     -e GORUN_MOUNT_PATH=/data/gorun/mounts \
     -e GORUN_DB=/data/gorun/gorun.db \
     --name gorun \
     gorun
   ```

### Environment Variables

- `GORUN_PORT` (Optional, default: 8080)
  - Port for the web interface
- `GORUN_SECRET` (Required)
  - Secret key for authentication
- `GORUN_MOUNT_PATH` (Optional)
  - Directory for container mounts
- `GORUN_DB` (Optional)
  - Path to the SQLite database
- `GORUN_PATH` (Optional)
  - Base directory for all gorun data

### Local Development

1. Install dependencies:
   ```bash
   go mod download
   cd manager && npm install
   ```

2. Build the frontend:
   ```bash
   make frontend-build
   ```

3. Run the development server:
   ```bash
   make dev
   ```

## API Documentation

The API documentation is available at `/api/docs` when running the server.

### Example API Usage

```bash
# Create a new job
curl -X POST http://localhost:8080/api/jobs \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "ubuntu",
    "command": ["echo", "Hello, World!"]
  }'

# Get job status
curl http://localhost:8080/api/jobs/{job_id} \
  -H "Authorization: Bearer your-token"
```

## Security Considerations

- Always use a strong `GORUN_SECRET`
- Keep your Docker socket secure
- Regularly backup your data directory
- Use HTTPS in production


