# WebGL Water Tutorial - Technical Documentation

This directory contains comprehensive technical documentation for the WebGL Water Tutorial project, including system architecture diagrams and component relationships.

## 📋 Overview

The WebGL Water Tutorial is a real-time water simulation application with the following key features:
- **Real-time water rendering** with reflection and refraction effects
- **Multi-pass rendering pipeline** using framebuffers
- **Go backend** serving assets and managing application state
- **WebGL frontend** with vanilla JavaScript for maximum performance
- **WebSocket communication** for real-time state synchronization

## 🏗️ Architecture Diagrams

### 1. System Architecture
**File:** `architecture.puml`  
**Purpose:** High-level overview of system components and their relationships

Shows the interaction between:
- Frontend components (WebGL, Canvas, JavaScript)
- Backend components (HTTP Server, Asset Manager, State Manager)
- File system resources (shaders, textures, meshes)
- Runtime data flow

### 2. Rendering Pipeline
**File:** `rendering-pipeline.puml`  
**Purpose:** State diagram of the WebGL rendering process

Illustrates the complete rendering pipeline:
- Initialization and asset loading
- Multi-pass rendering (refraction → reflection → main scene)
- Real-time animation loop
- State management and updates

### 3. Data Flow
**File:** `data-flow.puml`  
**Purpose:** Activity diagram showing data movement through the system

Tracks data flow from:
- Initial page load and WebGL setup
- User interactions and state updates
- Multi-pass rendering with framebuffers
- WebSocket real-time synchronization

### 4. Component Structure
**File:** `components.puml`  
**Purpose:** Class diagram of major system components

Details the structure of:
- Go backend classes (Server, Assets, State, Camera)
- JavaScript frontend components (WebGLWaterApp, ShaderProgram, etc.)
- Shader pipeline components
- Math3D utility library

## 🔧 Key Components

### Backend (Go)
- **Server**: HTTP server with REST API and WebSocket support
- **Assets**: Mesh and texture management with procedural generation
- **State**: Application state management with real-time updates
- **Camera**: 3D camera system with orbital controls
- **Math3D**: Vector and matrix math utilities

### Frontend (JavaScript/WebGL)
- **WebGLWaterApp**: Main application class managing the render loop
- **Shader Management**: Dynamic shader loading and compilation
- **Mesh Buffers**: GPU buffer management for 3D geometry
- **Framebuffers**: Render targets for multi-pass rendering
- **Texture Manager**: Texture loading and binding

### Rendering Pipeline
1. **Refraction Pass**: Render scene below water surface
2. **Reflection Pass**: Render mirrored scene above water
3. **Main Scene**: Composite water surface with textures
4. **Real-time Animation**: 60fps update loop with state sync

## 📊 Performance Characteristics

- **Frame Rate**: Target 60 FPS with smooth animation
- **Mesh Complexity**: 
  - Water plane: 4,225 vertices (65x65 grid)
  - Terrain: 1,089 vertices (33x33 grid)
- **Texture Resolution**: 512x512 for all textures
- **Framebuffer Size**: 
  - Reflection: 320x180
  - Refraction: 320x180

## 🔍 Debugging Information

### Current Status (Checkpoint)
- ✅ **Assets Loading**: All shaders, textures, and meshes load correctly
- ✅ **WebGL Context**: Successfully initialized with required extensions
- ✅ **Framebuffers**: Created with proper texture parameters
- ✅ **Shader Compilation**: All shaders compile without errors
- 🔄 **Rendering Issues**: Water appears as dark plate (vertex/fragment shader mismatch resolved)
- 🔄 **Progress**: Small reflective area visible (reflection pass working)

### Known Issues
1. **Vertex Shader**: Fixed position attribute from vec2 to vec3
2. **Texture Parameters**: Added CLAMP_TO_EDGE for framebuffer textures
3. **Canvas Dimensions**: Corrected JavaScript constants to match HTML canvas size
4. **Asset Pipeline**: All assets now served correctly from restructured paths

## 📁 File Structure

```
webgl-water/
├── docs/                           # This documentation
│   ├── README.md                   # This file
│   ├── architecture.puml           # System architecture
│   ├── rendering-pipeline.puml     # Rendering state diagram
│   ├── data-flow.puml              # Data flow activity diagram
│   └── components.puml             # Component class diagram
├── cmd/server/                     # Go application entry point
├── internal/                       # Go backend implementation
│   ├── app/                        # HTTP server and routing
│   ├── assets/                     # Asset management
│   ├── math3d/                     # 3D math utilities
│   └── state/                      # Application state
├── web/                           # Frontend resources
│   ├── static/webgl-water.js      # Main JavaScript application
│   └── shaders/                   # GLSL shader files
├── assets/                        # Runtime assets
│   ├── *.png                      # Texture files
│   └── meshes.json               # Mesh data
└── *.png                         # Original texture assets
```

## 🚀 Next Steps

1. **Fix Water Rendering**: Complete the water surface rendering pipeline
2. **Camera Controls**: Improve camera positioning and movement
3. **Performance**: Optimize rendering performance and memory usage
4. **Features**: Add additional water effects and scene elements
5. **Documentation**: Keep diagrams updated as system evolves

## 📖 How to Use This Documentation

1. **PlantUML Diagrams**: Use PlantUML to render the `.puml` files into images
2. **Online Viewer**: Visit [plantuml.com](http://www.plantuml.com/plantuml/uml/) to view diagrams
3. **IDE Integration**: Use PlantUML plugins in VSCode, IntelliJ, or other editors
4. **Export Formats**: Generate PNG, SVG, or PDF versions for documentation

## 🤝 Contributing

When making changes to the system:
1. Update the relevant PlantUML diagrams
2. Keep this README in sync with actual implementation
3. Document any architectural decisions or major changes
4. Update performance characteristics if they change significantly

---

**Last Updated**: August 2025  
**Status**: Active Development - Rendering Pipeline Debug Phase