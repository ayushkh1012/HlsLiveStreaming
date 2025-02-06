# Live HLS Server with Ad Insertion

A Go-based HTTP Live Streaming (HLS) server that supports live streaming with dynamic ad insertion. The server provides seamless integration of pre-roll and mid-roll advertisements in a live streaming context.

## Features

- Live HLS streaming with dynamic manifest updates
- Mid-roll ad insertion
- SCTE-35 markers for ad boundaries
- Docker containerization
- Kubernetes deployment via Helm
- Multiple quality variants support
- Cross-platform compatibility

## Directory Structure Details

### `/pkg/manifest`
Contains the core manifest handling logic for HLS streaming and ad insertion.

### `/utils`
Utility functions and configuration handling.

### `/config`
Configuration files for the server and deployment.

### `/media` and `/ads`
Directories containing the media segments and advertisement content.

### `/helm`
Helm chart for Kubernetes deployment.

### `/scripts`
Deployment and utility scripts.

## Quick Start

### Local Development

1. Clone the repository:

```bash
git clone https://github.com/yourusername/livehls.git
cd livehls
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

### Docker Deployment

1. Build and run using Docker Compose:
```bash
docker-compose up --build
```

2. Access the server at `http://localhost:8080`

### Kubernetes Deployment (Kind)

1. Deploy to local Kind cluster:
```bash
./scripts/deploy-kind.sh
```

2. Access the server at `http://localhost:30080`

## Configuration

### Server Configuration (config/config.yaml)
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  window_size: 5

paths:
  media: "./media"
  ads: "./ads"
  manifests: "./manifests"

ads:
  - id: "adv1"
    duration: 30
```

### Manifest Configuration

- Window Size: 5 segments
- Segment Duration: 10 seconds
- Ad Duration: 30 seconds (3 segments)
- Ad Insertion: Every 10 segments
- Pre-roll Ad: Enabled

## API Endpoints

- `/`: HTML interface with available endpoints
- `/media/`: Media segments directory
- `/ads/`: Advertisement segments directory
- `/manifests/`: HLS manifest files

## Development

### Building from Source

```bash
go build -o livehls main.go
```

### Running Tests

```bash
go test ./...
```

## Docker Support

Build the image:
```bash
docker build -t livehls:latest .
```

Run the container:
```bash
docker run -p 8080:8080 livehls:latest
```

## Kubernetes Deployment

### Using Helm

1. Install the chart:
```bash
helm upgrade --install livehls ./helm/livehls \
    --namespace livehls \
    --create-namespace
```

2. Uninstall:
```bash
helm uninstall livehls -n livehls
```

## Implementation Details

### Live Streaming
- Uses sliding window approach
- Maintains live streaming specifications

### Ad Insertion
- SCTE-35 markers for ad boundaries
- Smooth transitions with discontinuity markers

### Container Support
- Multi-stage Docker builds
- Volume mounting for media files
- Proper cache control
- Cross-platform compatibility

### Kubernetes Features
- Helm chart for deployment
- ConfigMap for configuration
- Volume management
- Health checks
- Service exposure