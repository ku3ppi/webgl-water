# WebGL Water Tutorial Go Port - Deployment Guide

This document provides comprehensive instructions for deploying and running the WebGL Water Tutorial Go Port in various environments.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Local Development](#local-development)
3. [Docker Deployment](#docker-deployment)
4. [Production Deployment](#production-deployment)
5. [Environment Configuration](#environment-configuration)
6. [Troubleshooting](#troubleshooting)
7. [Performance Tuning](#performance-tuning)

## Quick Start

The fastest way to get the application running:

### Prerequisites

- Go 1.21 or later
- Modern web browser with WebGL support
- Git (for cloning the repository)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd webgl-water/go-port
```

### 2. Run with Make (Recommended)

```bash
# Setup development environment and run
make setup
make run
```

### 3. Manual Setup

```bash
# Download dependencies
go mod tidy

# Copy environment configuration
cp .env.example .env

# Build and run
go build -o build/server ./cmd/server
./build/server
```

### 4. Access the Application

Open your browser and navigate to:
```
http://localhost:8080
```

## Local Development

### Development Environment Setup

1. **Install Development Tools** (optional but recommended):
   ```bash
   # Hot reload tool
   go install github.com/cosmtrek/air@latest
   
   # Code formatting
   go install golang.org/x/tools/cmd/goimports@latest
   
   # Linting
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
   ```

2. **Setup with Make**:
   ```bash
   make setup    # Creates directories, copies assets, sets up .env
   make dev      # Runs with hot reload using Air
   ```

3. **Manual Development Setup**:
   ```bash
   # Create necessary directories
   mkdir -p assets web/static web/shaders build
   
   # Copy assets from parent directory
   cp ../dudvmap.png assets/
   cp ../normalmap.png assets/
   cp ../stone-texture.png assets/
   
   # Create .env file
   cp .env.example .env
   
   # Run with hot reload
   air
   ```

### Development Workflow

1. **Code Changes**: Modify Go code in `internal/` or `cmd/` directories
2. **Frontend Changes**: Update JavaScript in `web/static/` or shaders in `web/shaders/`
3. **Testing**: Use `make test` to run the test suite
4. **Formatting**: Use `make format` to format code
5. **Linting**: Use `make lint` to check code quality

### File Watching

The development server automatically reloads when Go files change (using Air). For frontend changes, simply refresh the browser.

## Docker Deployment

### Development with Docker

1. **Using Docker Compose** (Recommended):
   ```bash
   docker-compose --profile dev up --build
   ```

2. **Using Development Dockerfile**:
   ```bash
   docker build -f Dockerfile.dev -t webgl-water-dev .
   docker run -p 8080:8080 -v $(pwd):/app webgl-water-dev
   ```

### Production with Docker

1. **Build Production Image**:
   ```bash
   docker build -t webgl-water-go:latest .
   ```

2. **Run Production Container**:
   ```bash
   docker run -d \
     --name webgl-water \
     -p 8080:8080 \
     --restart unless-stopped \
     webgl-water-go:latest
   ```

3. **Using Docker Compose**:
   ```bash
   docker-compose up -d --build
   ```

### Docker Environment Variables

```bash
# Override default settings
docker run -p 8080:8080 \
  -e PORT=3000 \
  -e ASSETS_PATH=/custom/assets \
  -e STATIC_PATH=/custom/static \
  webgl-water-go:latest
```

## Production Deployment

### Binary Distribution

1. **Build for Production**:
   ```bash
   make build-all  # Creates binaries for multiple platforms
   ```

2. **Create Release Package**:
   ```bash
   make release    # Creates compressed release packages
   ```

3. **Deploy Binary**:
   ```bash
   # Copy binary and assets to production server
   scp -r build/webgl-water-server-linux-amd64 user@server:/opt/webgl-water/
   scp -r web/ user@server:/opt/webgl-water/
   scp -r assets/ user@server:/opt/webgl-water/
   ```

### Systemd Service (Linux)

Create `/etc/systemd/system/webgl-water.service`:

```ini
[Unit]
Description=WebGL Water Tutorial Go Port
After=network.target

[Service]
Type=simple
User=webgl
WorkingDirectory=/opt/webgl-water
ExecStart=/opt/webgl-water/webgl-water-server
Restart=always
RestartSec=3
Environment=PORT=8080
Environment=ASSETS_PATH=/opt/webgl-water/assets
Environment=STATIC_PATH=/opt/webgl-water/web/static

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable webgl-water
sudo systemctl start webgl-water
sudo systemctl status webgl-water
```

### Nginx Reverse Proxy

Create `/etc/nginx/sites-available/webgl-water`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket support
    location /ws {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Static asset caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        proxy_pass http://localhost:8080;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/webgl-water /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL with Let's Encrypt

```bash
sudo certbot --nginx -d your-domain.com
```

## Environment Configuration

### Configuration Methods

1. **Environment Variables**: Set in shell or systemd service
2. **`.env` File**: Place in working directory
3. **Command Line Flags**: Override other settings

### Priority Order

1. Command line flags (highest priority)
2. Environment variables
3. `.env` file
4. Default values (lowest priority)

### Key Configuration Options

| Variable | Flag | Default | Description |
|----------|------|---------|-------------|
| `PORT` | `-port` | 8080 | HTTP server port |
| `ASSETS_PATH` | `-assets` | ./assets | Path to asset files |
| `STATIC_PATH` | `-static` | ./web/static | Path to static files |
| `DEBUG` | N/A | false | Enable debug logging |
| `LOG_LEVEL` | N/A | info | Logging level |

### Production Configuration

```bash
# Production environment variables
export PORT=80
export GO_ENV=production
export DEBUG=false
export LOG_LEVEL=warn
export ASSETS_PATH=/var/www/webgl-water/assets
export STATIC_PATH=/var/www/webgl-water/static
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use

**Error**: `listen tcp :8080: bind: address already in use`

**Solutions**:
```bash
# Find process using the port
lsof -i :8080
netstat -tulpn | grep 8080

# Kill the process or use different port
export PORT=3000
./server -port 3000
```

#### 2. Assets Not Found

**Error**: `404 Not Found` for textures or meshes

**Solutions**:
```bash
# Check asset paths
ls -la assets/
ls -la web/static/

# Copy assets manually
cp ../dudvmap.png assets/
cp ../normalmap.png assets/
cp ../stone-texture.png assets/

# Verify permissions
chmod -R 644 assets/
chmod -R 644 web/
```

#### 3. WebGL Not Supported

**Error**: Browser console shows "WebGL not supported"

**Solutions**:
- Use a modern browser (Chrome, Firefox, Safari, Edge)
- Enable hardware acceleration in browser settings
- Update graphics drivers
- Try different browser

#### 4. Shader Compilation Errors

**Error**: Browser console shows shader compilation failures

**Solutions**:
```bash
# Check shader files exist and are readable
ls -la web/shaders/
cat web/shaders/water-vertex.glsl

# Verify shader syntax
# Check browser console for specific error messages
```

#### 5. WebSocket Connection Failed

**Error**: "WebSocket disconnected" in browser console

**Solutions**:
- Check firewall settings
- Verify proxy configuration supports WebSocket
- Try disabling browser extensions
- Check server logs for connection errors

### Debug Mode

Enable debug logging:

```bash
# Environment variable
export DEBUG=true

# Command line
./server -debug

# Check logs
tail -f /var/log/webgl-water.log
```

### Performance Issues

1. **High Memory Usage**:
   ```bash
   # Monitor memory
   top -p $(pidof webgl-water-server)
   
   # Check for memory leaks
   go tool pprof http://localhost:8080/debug/pprof/heap
   ```

2. **High CPU Usage**:
   ```bash
   # Profile CPU usage
   go tool pprof http://localhost:8080/debug/pprof/profile
   ```

3. **Slow Rendering**:
   - Check browser performance tools
   - Verify WebGL extensions are available
   - Reduce framebuffer sizes in configuration

## Performance Tuning

### Server Optimization

1. **Build Optimization**:
   ```bash
   # Optimized build
   CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/server
   
   # Further compression
   upx --best --lzma server
   ```

2. **Memory Settings**:
   ```bash
   # Garbage collection tuning
   export GOGC=100
   export GOMEMLIMIT=512MiB
   ```

3. **Connection Limits**:
   ```bash
   # System limits
   ulimit -n 65535
   
   # In systemd service
   LimitNOFILE=65535
   ```

### Client Optimization

1. **Framebuffer Sizes**: Reduce in `.env`:
   ```
   REFLECTION_TEXTURE_WIDTH=160
   REFLECTION_TEXTURE_HEIGHT=90
   REFRACTION_TEXTURE_WIDTH=640
   REFRACTION_TEXTURE_HEIGHT=360
   ```

2. **Mesh Complexity**: Reduce segments:
   ```
   WATER_MESH_SEGMENTS=32
   TERRAIN_MESH_SEGMENTS=16
   ```

3. **Browser Settings**:
   - Enable hardware acceleration
   - Close unnecessary tabs
   - Update browser and drivers

### Monitoring

1. **Health Checks**:
   ```bash
   # HTTP health check
   curl -f http://localhost:8080/ || exit 1
   
   # With Docker
   HEALTHCHECK --interval=30s --timeout=3s \
     CMD wget --spider http://localhost:8080/ || exit 1
   ```

2. **Metrics Collection**:
   ```bash
   # Basic monitoring
   while true; do
     echo "$(date): $(curl -s http://localhost:8080/api/state | jq -r '.clock')"
     sleep 10
   done
   ```

3. **Log Analysis**:
   ```bash
   # Error monitoring
   tail -f /var/log/webgl-water.log | grep -i error
   
   # Performance monitoring
   tail -f /var/log/webgl-water.log | grep -i "response_time"
   ```

## Security Considerations

### Production Security

1. **Run as Non-root User**:
   ```bash
   # Create dedicated user
   sudo useradd -r -s /bin/false webgl
   sudo chown -R webgl:webgl /opt/webgl-water
   ```

2. **Firewall Configuration**:
   ```bash
   # UFW example
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw deny 8080/tcp  # Block direct access
   ```

3. **Environment Variables**:
   ```bash
   # Don't expose sensitive information
   unset DEBUG
   unset GO_ENV
   ```

4. **File Permissions**:
   ```bash
   # Secure file permissions
   chmod 755 /opt/webgl-water/server
   chmod -R 644 /opt/webgl-water/assets
   chmod -R 644 /opt/webgl-water/web
   ```

### Development Security

1. **CORS Configuration**: Already configured for development
2. **Debug Information**: Only enable in development
3. **Asset Access**: Ensure assets don't contain sensitive data

This completes the deployment guide. The application should now be ready for deployment in any environment from development to production.