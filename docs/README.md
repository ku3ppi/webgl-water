# WebGL Water Tutorial - Technical Documentation

This directory contains comprehensive technical documentation for the WebGL Water Tutorial project, including system architecture diagrams and component relationships.

## 📋 Overview

The WebGL Water Tutorial is a real-time water simulation application with the following key features:
- **Real-time water rendering** with reflection and refraction effects
- **Multi-pass rendering pipeline** using framebuffers
- **Go backend** serving assets and managing application state
- **WebGL frontend** with vanilla JavaScript for maximum performance
- **WebSocket communication** for real-time state synchronization

## 📁 Documentation Structure

The documentation is organized by file type for better navigation:

### 📊 PlantUML Diagrams (`puml/`)
- **System Architecture** (`architecture.puml` / `architecture-de.puml`)
- **Rendering Pipeline** (`rendering-pipeline.puml` / `rendering-pipeline-de.puml`)
- **Data Flow** (`data-flow.puml` / `data-flow-de.puml`)
- **Component Structure** (`components.puml` / `components-de.puml`)
- **Code Maps** - Detaillierte Code-Navigation:
  - `backend-codemap-de.puml` - Go Backend Struktur
  - `frontend-codemap-de.puml` - JavaScript Frontend
  - `shader-codemap-de.puml` - Shader Variablen & Konstanten
  - `dataflow-codemap-de.puml` - Nachrichten & Datentypen

### 🖼️ SVG Exports (`svg/`)
- Exportierte Diagramme als SVG Dateien
- Für bessere Integration in andere Dokumentation
- Deutsche und englische Versionen verfügbar

### 📄 PDF Exports (`pdf/`)
- High-quality PDF Versionen der Diagramme
- Ideal für Präsentationen und Ausdrucke
- Professionelle Darstellung aller Systemkomponenten

### 📝 Text Exports (`txt/`)
- PlantUML Text-Ausgaben für bessere Durchsuchbarkeit
- Backup-Format für alle Diagramminhalte

## 🏗️ Architektur Übersicht

### 1. System Architektur
**Zweck:** High-level Übersicht der Systemkomponenten und Beziehungen

Zeigt die Interaktion zwischen:
- Frontend Komponenten (WebGL, Canvas, JavaScript)
- Backend Komponenten (HTTP Server, Asset Manager, State Manager)
- Dateisystem Ressourcen (Shader, Texturen, Meshes)
- Laufzeit Datenfluss

### 2. Rendering Pipeline
**Zweck:** Zustandsdiagramm des WebGL Rendering Prozesses

Illustriert die komplette Rendering Pipeline:
- Initialisierung und Asset Loading
- Multi-Pass Rendering (Refraktion → Reflektion → Hauptszene)
- Echtzeit Animationsschleife
- Status Management und Updates

### 3. Datenfluss
**Zweck:** Aktivitätsdiagramm für Datenbewegung durch das System

Verfolgt Datenfluss von:
- Initialer Seitenladevorgang und WebGL Setup
- Benutzerinteraktionen und Status Updates
- Multi-Pass Rendering mit Framebuffern
- WebSocket Echtzeit Synchronisation

### 4. Komponenten Struktur
**Zweck:** Klassendiagramm der wichtigsten Systemkomponenten

Detailliert die Struktur von:
- Go Backend Klassen (Server, Assets, State, Camera)
- JavaScript Frontend Komponenten (WebGLWaterApp, ShaderProgram, etc.)
- Shader Pipeline Komponenten
- Math3D Utility Bibliothek

### 5. Code-Karten (Neu!)
**Zweck:** Detaillierte Navigation durch den Code

**Backend Code-Karte**: Alle Go Strukturen, Methoden, HTTP Handler
**Frontend Code-Karte**: JavaScript Klassen, WebGL Konstanten, Event Handler
**Shader Code-Karte**: GLSL Variablen, Uniforms, Konstanten mit Werten
**Datenfluss Code-Karte**: Message-Typen, HTTP Strukturen, WebSocket Nachrichten

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

## 📁 Projektstruktur

```
webgl-water/
├── docs/                           # Dokumentation (neu organisiert)
│   ├── README.md                   # Diese Datei
│   ├── puml/                       # PlantUML Quelldiagramme
│   │   ├── architecture.puml       # System Architektur (EN)
│   │   ├── architecture-de.puml    # System Architektur (DE)
│   │   ├── components.puml         # Komponenten (EN)
│   │   ├── components-de.puml      # Komponenten (DE)
│   │   ├── data-flow.puml          # Datenfluss (EN)
│   │   ├── data-flow-de.puml       # Datenfluss (DE)
│   │   ├── rendering-pipeline.puml # Rendering Pipeline (EN)
│   │   ├── rendering-pipeline-de.puml # Rendering Pipeline (DE)
│   │   ├── backend-codemap-de.puml # Go Backend Code-Karte
│   │   ├── frontend-codemap-de.puml # JS Frontend Code-Karte
│   │   ├── shader-codemap-de.puml  # Shader Code-Karte
│   │   └── dataflow-codemap-de.puml # Datenfluss Code-Karte
│   ├── svg/                        # SVG Exporte
│   │   ├── *.svg                   # Alle Diagramme als SVG
│   ├── pdf/                        # PDF Exporte
│   │   ├── *.pdf                   # Alle Diagramme als PDF
│   └── txt/                        # Text Exporte
│       └── *.txt                   # PlantUML Text Outputs
├── cmd/server/                     # Go Anwendungs-Einstiegspunkt
├── internal/                       # Go Backend Implementierung
│   ├── app/                        # HTTP Server und Routing
│   ├── assets/                     # Asset Management
│   ├── math3d/                     # 3D Math Utilities
│   └── state/                      # Anwendungsstatus
├── web/                           # Frontend Ressourcen
│   ├── static/webgl-water.js      # Haupt JavaScript Anwendung
│   └── shaders/                   # GLSL Shader Dateien
├── assets/                        # Laufzeit Assets
│   ├── *.png                      # Textur Dateien
│   └── meshes.json               # Mesh Daten
└── *.png                         # Original Textur Assets
```

## 🚀 Next Steps

1. **Fix Water Rendering**: Complete the water surface rendering pipeline
2. **Camera Controls**: Improve camera positioning and movement
3. **Performance**: Optimize rendering performance and memory usage
4. **Features**: Add additional water effects and scene elements
5. **Documentation**: Keep diagrams updated as system evolves

## 📖 Verwendung der Dokumentation

### PlantUML Diagramme ansehen:
1. **Online Viewer**: Besuche [plantuml.com](http://www.plantuml.com/plantuml/uml/) um Diagramme anzuzeigen
2. **IDE Integration**: Nutze PlantUML Plugins in VSCode, IntelliJ, oder anderen Editoren
3. **Lokale Installation**: Installiere PlantUML für lokales Rendering
4. **Fertige Exporte**: Nutze die bereits generierten SVG/PDF Dateien in `svg/` und `pdf/`

### Code Navigation:
1. **Code-Karten nutzen**: Öffne die `*-codemap-de.puml` Dateien für detaillierte Codeübersicht
2. **Variablennamen suchen**: Alle echten Variablen-/Funktionsnamen sind in den Code-Karten
3. **Debugging**: Nutze die Konstantenwerte und Strukturnamen für Breakpoints
4. **Feature Erweiterung**: Verstehe Datenflüsse um neue Features hinzuzufügen

### Export Formate:
- **SVG**: Für Web-Integration und skalierbare Grafiken
- **PDF**: Für Präsentationen und hochwertige Ausdrucke  
- **TXT**: Für Textsuche und Backup der Diagramminhalte

## 🤝 Contributing

When making changes to the system:
1. Update the relevant PlantUML diagrams
2. Keep this README in sync with actual implementation
3. Document any architectural decisions or major changes
4. Update performance characteristics if they change significantly

---

**Letzte Aktualisierung**: Dezember 2024  
**Status**: Aktive Entwicklung - Dokumentation komplett überarbeitet und strukturiert  
**Neu**: Code-Karten für detaillierte Navigation durch Go Backend, JS Frontend, Shader und Datenstrukturen