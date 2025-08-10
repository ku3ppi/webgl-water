# WebGL Water Tutorial - Go Port

A modern implementation of a WebGL water simulation tutorial, ported from Rust/WASM to Go with a vanilla JavaScript frontend.

![WebGL Water Simulation](screenshot.png)

## Features

- **Real-time water simulation** with animated waves, reflection, and refraction
- **Interactive camera controls** with mouse orbit and zoom
- **Dynamic water properties** adjustable via UI controls
- **Modern Go backend** serving assets, shaders, and API endpoints
- **WebGL frontend** with vanilla JavaScript for maximum performance
- **Docker support** for easy deployment and development

## Quick Start

### Prerequisites

- Go 1.22 or later
- Modern web browser with WebGL support
- Docker (optional)

### Running Locally

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd webgl-water
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Start the server:**
   ```bash
   go run ./cmd/server
   ```

4. **Open your browser:**
   Navigate to `http://localhost:8080`

### Using Docker

1. **Build and run:**
   ```bash
   docker-compose up --build
   ```

2. **Access the application:**
   Open `http://localhost:8080`

## Project Structure

```
webgl-water/
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── app/             # HTTP server and routing
│   ├── assets/          # Asset management
│   └── state/           # Application state and 3D math
├── web/
│   ├── static/          # JavaScript frontend
│   └── shaders/         # GLSL shader files
├── assets/              # Texture files and meshes
├── *.png               # Texture assets
├── go.mod              # Go module definition
├── Dockerfile          # Production container
├── Dockerfile.dev      # Development container
└── docker-compose.yml  # Docker compose configuration
```

## Controls

### Camera
- **Mouse drag**: Orbit around the water
- **Mouse wheel**: Zoom in/out (now supports 3x more zoom range)

### Water Properties
- **Reflectivity**: Controls how reflective the water surface appears
- **Fresnel Strength**: Adjusts the Fresnel effect intensity
- **Wave Speed**: Controls animation speed of water waves
- **Use Reflection**: Toggle reflection rendering
- **Use Refraction**: Toggle refraction rendering
- **Show Scenery**: Toggle background scenery visibility

## Architecture

### Backend (Go)
- **HTTP Server**: Gorilla Mux router serving static files and API
- **Asset Management**: Loads and serves textures, meshes, and shaders
- **State Management**: Real-time application state with WebSocket updates
- **3D Math**: Vector and matrix operations for camera and transformations

### Frontend (JavaScript/WebGL)
- **WebGL Context**: Direct WebGL API usage for maximum performance
- **Shader Management**: Dynamic shader loading and compilation
- **Asset Loading**: Asynchronous loading of textures and meshes
- **Camera System**: Orbital camera with smooth controls
- **Water Rendering**: Multi-pass rendering with reflection/refraction

## Development

### Building
```bash
# Build for current platform
go build -o server ./cmd/server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o server-linux ./cmd/server
```

### Running with Custom Options
```bash
# Custom port and paths
go run ./cmd/server -port 3000 -assets ./assets -static ./web/static

# Environment variables
PORT=3000 ASSETS_PATH=./assets STATIC_PATH=./web/static go run ./cmd/server
```

### API Endpoints

- `GET /` - Main application page
- `GET /static/*` - Static files (JavaScript, etc.)
- `GET /assets/{filename}` - Asset files (textures, etc.)
- `GET /shaders/{name}` - GLSL shader files
- `GET /api/meshes` - List available meshes
- `GET /api/meshes/{name}` - Get specific mesh data
- `GET /api/state` - Current application state
- `POST /api/state/water` - Update water properties
- `POST /api/state/camera` - Update camera state
- `GET /ws` - WebSocket for real-time updates

## Recent Improvements

- ✅ **Larger Canvas**: Increased default size to 1200x800
- ✅ **Enhanced Zoom**: 3x more zoom-out capability (up to 150 units)
- ✅ **Better Controls**: Repositioned controls to avoid overlap, improved styling
- ✅ **Restructured Repository**: Moved to standard Go project layout
- ✅ **Asset Pipeline**: Fixed texture and shader serving

## Technical Details

### Water Rendering Pipeline
1. **Reflection Pass**: Render scene from below water surface
2. **Refraction Pass**: Render underwater scene
3. **Water Pass**: Render water surface with reflection/refraction textures
4. **Post-processing**: Combine passes with Fresnel effects

### Shader System
- `water-vertex.glsl` / `water-fragment.glsl`: Main water surface rendering
- `mesh-vertex.glsl` / `mesh-fragment.glsl`: Scene geometry rendering
- `textured-quad-vertex.glsl` / `textured-quad-fragment.glsl`: Post-processing

### Asset Loading
- **Textures**: PNG files for water effects and scene materials
- **Meshes**: JSON format with vertices, normals, and UV coordinates
- **Shaders**: GLSL files loaded dynamically

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is dual-licensed under:
- [Apache License 2.0](LICENSE-APACHE)
- [MIT License](LICENSE-MIT)

## Credits

Based on the WebGL Water Tutorial, ported from the original Rust/WASM implementation to Go with modern web technologies.