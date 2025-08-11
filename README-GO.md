# WebGL Water Tutorial - Go Port

A complete port of the [WebGL Water Tutorial](https://github.com/chinedufn/webgl-water-tutorial) from Rust/WASM to Go with JavaScript/WebGL frontend.

## Overview

This project demonstrates realistic water rendering using WebGL with reflections, refractions, and animated water waves. The original Rust/WASM implementation has been ported to a Go-based web server that serves a JavaScript WebGL frontend.

### Key Features

- **Realistic Water Rendering**: Reflection and refraction effects with dynamic wave animation
- **Interactive Camera**: Orbit camera with mouse controls and zoom
- **Real-time Controls**: Adjust water properties (reflectivity, fresnel strength, wave speed) in real-time
- **Framebuffer Rendering**: Separate reflection and refraction framebuffers for realistic water effects
- **Go Backend**: HTTP server with WebSocket support for real-time state synchronization
- **Modern Architecture**: Clean separation between backend state management and frontend rendering

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Backend    │    │   HTTP/WebSocket │    │  JS/WebGL       │
│                 │    │                 │    │  Frontend       │
│ • State Mgmt    │◄──►│ • REST API      │◄──►│ • WebGL Context │
│ • 3D Math       │    │ • Asset Serving │    │ • Shader Mgmt   │
│ • Asset Pipeline│    │ • Real-time Sync│    │ • Camera Ctrl   │
│ • Mesh Gen      │    │                 │    │ • UI Controls   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Backend (Go)

- **HTTP Server**: Serves static files and provides REST API endpoints
- **WebSocket Support**: Real-time state synchronization between server and client
- **3D Math Library**: Complete implementation of vectors, matrices, and quaternions
- **Asset Management**: Mesh generation and texture serving
- **State Management**: Camera controls, water properties, and application state

### Frontend (JavaScript/WebGL)

- **WebGL Rendering**: Full WebGL implementation with shader management
- **Camera Controls**: Interactive orbit camera with mouse/wheel controls
- **Framebuffer Management**: Reflection and refraction texture rendering
- **UI Controls**: Real-time adjustment of water properties
- **Asset Loading**: Dynamic loading of meshes and textures from Go backend

## Quick Start

### Prerequisites

- Go 1.21 or later
- Modern web browser with WebGL support
- Docker (optional)

### Local Development

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd webgl-water
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Run the server**:
   ```bash
   go run ./cmd/server
   ```

4. **Open your browser**:
   Navigate to `http://localhost:8080`

### Docker Development

1. **Build and run with Docker Compose**:
   ```bash
   docker-compose up --build
   ```

2. **Open your browser**:
   Navigate to `http://localhost:8080`

## Usage

### Controls

- **Mouse Drag**: Rotate camera around the water scene
- **Mouse Wheel**: Zoom in/out
- **UI Sliders**: Adjust water properties in real-time
  - Reflectivity: Controls how reflective the water surface appears
  - Fresnel Strength: Affects the viewing angle dependency of reflections
  - Wave Speed: Controls animation speed of water waves
  - Use Reflection/Refraction: Toggle reflection and refraction effects
  - Show Scenery: Toggle rendering of underwater terrain

### API Endpoints

The Go server provides several REST API endpoints:

- `GET /` - Main application page
- `GET /api/meshes` - List all available meshes
- `GET /api/meshes/{name}` - Get specific mesh data
- `GET /api/textures` - List all available textures
- `GET /api/state` - Get current application state
- `POST /api/state/water` - Update water properties
- `POST /api/state/camera` - Update camera state
- `GET /assets/{filename}` - Serve asset files
- `GET /shaders/{name}` - Serve shader files
- `WS /ws` - WebSocket endpoint for real-time updates

## Building

### Development

```bash
# Build for current platform
go build -o server ./cmd/server

# Run with custom options
go run ./cmd/server -port 3000 -assets ./assets -static ./web/static

# Environment variables
PORT=3000 ASSETS_PATH=./assets STATIC_PATH=./web/static go run ./cmd/server
```

### Docker

```bash
# Build Docker image
docker build -t webgl-water .

# Run with Docker
docker run -p 8080:8080 webgl-water

# Development environment
docker-compose up --build
```

## Project Structure

```
webgl-water/
├── docs/                    # 📚 Comprehensive documentation
│   ├── puml/               # PlantUML source diagrams (EN + DE)
│   ├── svg/                # SVG exports for web integration
│   ├── pdf/                # PDF exports for presentations
│   ├── txt/                # Text exports for searchability
│   └── README.md           # Detailed technical documentation
├── cmd/server/              # Main server executable
├── internal/
│   ├── app/                 # HTTP server and handlers
│   ├── assets/              # Asset management
│   ├── math3d/              # 3D mathematics library
│   └── state/               # Application state management
├── web/
│   ├── static/              # Static files (JS, CSS)
│   └── shaders/             # GLSL shader files
├── assets/                  # Runtime assets (textures, meshes)
├── *.png                   # Texture assets (dudvmap, normalmap, stone)
├── Dockerfile               # Production Docker image
├── Dockerfile.dev           # Development Docker image
├── docker-compose.yml       # Docker Compose configuration
└── README.md               # This file
```

## Technical Details

### 3D Math Library

The Go backend includes a complete 3D mathematics library with:

- **Vectors**: Vec2, Vec3, Vec4 with all standard operations
- **Matrices**: 4x4 matrices with perspective/orthographic projections
- **Quaternions**: Full quaternion support for rotations
- **Camera System**: Orbit camera with smooth controls

### WebGL Rendering Pipeline

1. **Refraction Pass**: Render underwater scene to framebuffer
2. **Reflection Pass**: Render above-water scene (mirrored) to framebuffer
3. **Water Rendering**: Combine reflection/refraction textures with water shaders
4. **Scene Rendering**: Render final scene with water and terrain
5. **Debug Views**: Small previews of framebuffer textures

### Shader System

The project includes several shader programs:

- **Water Shaders**: Advanced water rendering with reflection/refraction mixing
- **Mesh Shaders**: Standard mesh rendering with texture support
- **Textured Quad**: For framebuffer visualization and debug views

### Asset Pipeline

- **Mesh Generation**: Procedural generation of water plane and terrain
- **Texture Loading**: Dynamic texture loading from file system
- **Format Support**: JSON mesh format with binary support planned

## Configuration

The application can be configured through environment variables or command-line flags:

```bash
# Environment variables
export PORT=8080
export ASSETS_PATH=./assets
export STATIC_PATH=./web/static

# Command-line flags
./server -port 8080 -assets ./assets -static ./web/static
```

## Performance

### Optimization Features

- **VAO Extension**: Uses Vertex Array Objects when available
- **Framebuffer Optimization**: Efficient reflection/refraction rendering
- **State Caching**: Minimal state changes in WebGL
- **Asset Caching**: Browser caching for static assets

## Development

### Adding New Features

1. **Backend Changes**: Modify Go code in `internal/` directories
2. **Frontend Changes**: Update JavaScript code in `web/static/`
3. **Shaders**: Add/modify GLSL shaders in `web/shaders/`
4. **Assets**: Add new assets to `assets/` directory
5. **Documentation**: Update relevant diagrams in `docs/puml/`

For frontend changes, simply refresh the browser.

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Deployment

### Docker Production

```bash
# Build production image
docker build -t webgl-water .

# Run production container
docker run -p 8080:8080 webgl-water:latest
```

### Binary Distribution

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server-linux ./cmd/server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o server-windows.exe ./cmd/server
```

## Troubleshooting

### Common Issues

1. **WebGL Not Supported**: Ensure your browser supports WebGL
2. **Shader Compilation Errors**: Check browser console for shader errors
3. **Asset Loading Fails**: Verify asset paths and file permissions
4. **WebSocket Connection Issues**: Check firewall and proxy settings

### Debug Mode

Check the browser console for detailed error messages and performance metrics.

## Documentation

Comprehensive technical documentation is available in the `docs/` directory:

- **Architecture Diagrams**: System overview, rendering pipeline, data flow, components
- **Code Maps**: Detailed navigation through Go backend, JS frontend, shaders, and data structures
- **Multiple Formats**: PlantUML sources, SVG exports, PDF versions, and text outputs
- **English + German**: Complete documentation in both languages

See [docs/README.md](docs/README.md) for the complete technical documentation.

## Comparison with Original

| Feature | Rust/WASM | Go Port |
|---------|-----------|---------|
| Backend Language | Rust | Go |
| Frontend | Rust/WASM | JavaScript/WebGL |
| Architecture | Monolithic WASM | Client-Server |
| State Management | Local | Server-side + Sync |
| Asset Loading | Compile-time | Runtime |
| Hot Reload | ❌ | ✅ |
| Docker Support | ❌ | ✅ |
| Multi-platform | ✅ | ✅ |

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Update relevant documentation in `docs/` if needed
5. Run tests: `go test ./...`
6. Format code: `go fmt ./...`
7. Submit a pull request

## License

This project maintains the same license as the original:

- MIT License or Apache License 2.0

## Acknowledgments

- Original [WebGL Water Tutorial](https://github.com/chinedufn/webgl-water-tutorial) by [chinedufn](https://github.com/chinedufn)
- [ThinMatrix's OpenGL Water Tutorial](https://www.youtube.com/watch?v=HusvGeEDU_U&list=PLRIWtICgwaX23jiqVByUs0bqhnalNTNZh) for the water rendering techniques

## See Also

- [Original Rust/WASM Tutorial](http://chinedufn.com/3d-webgl-basic-water-tutorial/)
- [WebGL Fundamentals](https://webglfundamentals.org/)
- [Go WebGL Projects](https://github.com/topics/webgl+golang)