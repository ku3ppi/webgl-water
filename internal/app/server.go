package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/ku3ppi/webgl-water/internal/assets"
	"github.com/ku3ppi/webgl-water/internal/state"
)

// Server represents the main application server
type Server struct {
	router     *mux.Router
	assets     *assets.Assets
	appState   *state.State
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	staticPath string
	port       int
}

// NewServer creates a new server instance
func NewServer(assetsPath, staticPath string, port int) *Server {
	server := &Server{
		router:     mux.NewRouter(),
		assets:     assets.NewAssets(assetsPath),
		appState:   state.NewState(),
		staticPath: staticPath,
		port:       port,
		clients:    make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	// Static file serving
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticPath))))

	// Asset serving
	s.router.HandleFunc("/assets/{filename}", s.handleAssetFile).Methods("GET")

	// API endpoints
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/meshes", s.handleGetMeshes).Methods("GET")
	api.HandleFunc("/meshes/{name}", s.handleGetMesh).Methods("GET")
	api.HandleFunc("/textures", s.handleGetTextures).Methods("GET")
	api.HandleFunc("/state", s.handleGetState).Methods("GET")
	api.HandleFunc("/state/water", s.handleUpdateWater).Methods("POST")
	api.HandleFunc("/state/camera", s.handleUpdateCamera).Methods("POST")

	// WebSocket endpoint for real-time updates
	s.router.HandleFunc("/ws", s.handleWebSocket)

	// Shader serving
	s.router.HandleFunc("/shaders/{name}", s.handleShader).Methods("GET")

	// Main application route
	s.router.HandleFunc("/", s.handleIndex).Methods("GET")
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Initialize assets
	if err := s.assets.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize assets: %w", err)
	}

	log.Printf("Starting server on port %d", s.port)
	log.Printf("Static path: %s", s.staticPath)

	// Start state update ticker
	go s.startStateUpdates()

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}

// startStateUpdates starts a ticker to update application state
func (s *Server) startStateUpdates() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	lastTime := time.Now()

	for range ticker.C {
		now := time.Now()
		deltaTime := float32(now.Sub(lastTime).Milliseconds())
		lastTime = now

		// Update application state
		s.appState.Update(&state.AdvanceClockMessage{DeltaTime: deltaTime})

		// Broadcast state updates to connected WebSocket clients
		s.broadcastStateUpdate()
	}
}

// handleIndex serves the main application page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WebGL Water Tutorial - Go Port</title>
    <style>
        body { margin: 0; padding: 0; background: #000; overflow: hidden; }
        canvas { display: block; }
        #controls {
            position: absolute;
            top: 10px;
            right: 10px;
            background: rgba(0,0,0,0.8);
            color: white;
            padding: 15px;
            border-radius: 8px;
            font-family: Arial, sans-serif;
            font-size: 12px;
            min-width: 280px;
            max-width: 320px;
        }
        .control-group { margin-bottom: 12px; }
        label { display: inline-block; width: 120px; font-size: 11px; }
        input[type="range"] { width: 150px; }
        input[type="checkbox"] { margin-left: 8px; }
        h3 { margin-top: 0; margin-bottom: 15px; font-size: 14px; }
    </style>
</head>
<body>
    <canvas id="canvas" width="1200" height="800"></canvas>

    <div id="controls">
        <h3>Water Controls</h3>
        <div class="control-group">
            <label>Reflectivity:</label>
            <input type="range" id="reflectivity" min="0" max="1" step="0.01" value="0.6">
            <span id="reflectivity-value">0.6</span>
        </div>
        <div class="control-group">
            <label>Fresnel Strength:</label>
            <input type="range" id="fresnel" min="0" max="5" step="0.1" value="2.0">
            <span id="fresnel-value">2.0</span>
        </div>
        <div class="control-group">
            <label>Wave Speed:</label>
            <input type="range" id="wave-speed" min="0" max="0.1" step="0.001" value="0.03">
            <span id="wave-speed-value">0.03</span>
        </div>
        <div class="control-group">
            <label>Use Reflection:</label>
            <input type="checkbox" id="use-reflection" checked>
        </div>
        <div class="control-group">
            <label>Use Refraction:</label>
            <input type="checkbox" id="use-refraction" checked>
        </div>
        <div class="control-group">
            <label>Show Scenery:</label>
            <input type="checkbox" id="show-scenery" checked>
        </div>
    </div>

    <script src="/static/webgl-water.js"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// handleAssetFile serves asset files (textures, etc.)
func (s *Server) handleAssetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Try serving from current directory (where the original PNG files are)
	// Working directory is now webgl-water root
	rootPath := filepath.Join(".", filename)
	if _, err := os.Stat(rootPath); err == nil {
		w.Header().Set("Content-Type", getContentType(filename))
		http.ServeFile(w, r, rootPath)
		return
	}

	// Fall back to assets directory
	assetsPath := filepath.Join(".", "assets", filename)
	if _, err := os.Stat(assetsPath); err == nil {
		w.Header().Set("Content-Type", getContentType(filename))
		http.ServeFile(w, r, assetsPath)
		return
	}

	// File not found
	http.NotFound(w, r)
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

// handleGetMeshes returns a list of all available meshes
func (s *Server) handleGetMeshes(w http.ResponseWriter, r *http.Request) {
	meshNames := s.assets.ListMeshes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"meshes": meshNames,
	})
}

// handleGetMesh returns a specific mesh by name
func (s *Server) handleGetMesh(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meshName := vars["name"]

	mesh, err := s.assets.GetMesh(meshName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mesh)
}

// handleGetTextures returns a list of all available textures
func (s *Server) handleGetTextures(w http.ResponseWriter, r *http.Request) {
	textureNames := s.assets.ListTextures()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"textures": textureNames,
	})
}

// handleGetState returns the current application state
func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	camera := s.appState.GetCamera()
	water := s.appState.GetWater()

	response := map[string]interface{}{
		"clock":   s.appState.GetClock(),
		"scenery": s.appState.GetScenery(),
		"camera": map[string]interface{}{
			"position":   [3]float32{camera.GetPosition().X, camera.GetPosition().Y, camera.GetPosition().Z},
			"viewMatrix": camera.GetViewMatrix().ToSlice(),
		},
		"water": water,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WaterUpdateRequest represents a water property update request
type WaterUpdateRequest struct {
	Reflectivity    *float32 `json:"reflectivity,omitempty"`
	FresnelStrength *float32 `json:"fresnelStrength,omitempty"`
	WaveSpeed       *float32 `json:"waveSpeed,omitempty"`
	UseReflection   *bool    `json:"useReflection,omitempty"`
	UseRefraction   *bool    `json:"useRefraction,omitempty"`
}

// handleUpdateWater updates water properties
func (s *Server) handleUpdateWater(w http.ResponseWriter, r *http.Request) {
	var req WaterUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Apply updates
	if req.Reflectivity != nil {
		s.appState.Update(&state.SetReflectivityMessage{Value: *req.Reflectivity})
	}
	if req.FresnelStrength != nil {
		s.appState.Update(&state.SetFresnelMessage{Value: *req.FresnelStrength})
	}
	if req.WaveSpeed != nil {
		s.appState.Update(&state.SetWaveSpeedMessage{Value: *req.WaveSpeed})
	}
	if req.UseReflection != nil {
		s.appState.Update(&state.UseReflectionMessage{Value: *req.UseReflection})
	}
	if req.UseRefraction != nil {
		s.appState.Update(&state.UseRefractionMessage{Value: *req.UseRefraction})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// CameraUpdateRequest represents a camera update request
type CameraUpdateRequest struct {
	MouseDown *struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"mouseDown,omitempty"`
	MouseUp   *bool `json:"mouseUp,omitempty"`
	MouseMove *struct {
		X int32 `json:"x"`
		Y int32 `json:"y"`
	} `json:"mouseMove,omitempty"`
	Zoom *float32 `json:"zoom,omitempty"`
}

// handleUpdateCamera updates camera state
func (s *Server) handleUpdateCamera(w http.ResponseWriter, r *http.Request) {
	var req CameraUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Apply camera updates
	if req.MouseDown != nil {
		s.appState.Update(&state.MouseDownMessage{X: req.MouseDown.X, Y: req.MouseDown.Y})
	}
	if req.MouseUp != nil && *req.MouseUp {
		s.appState.Update(&state.MouseUpMessage{})
	}
	if req.MouseMove != nil {
		s.appState.Update(&state.MouseMoveMessage{X: req.MouseMove.X, Y: req.MouseMove.Y})
	}
	if req.Zoom != nil {
		s.appState.Update(&state.ZoomMessage{Delta: *req.Zoom})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// handleShader serves shader files
func (s *Server) handleShader(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderName := vars["name"]

	shaderPath := filepath.Join(s.staticPath, "..", "shaders", shaderName)

	w.Header().Set("Content-Type", "text/plain")
	http.ServeFile(w, r, shaderPath)
}

// handleWebSocket handles WebSocket connections for real-time updates
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	s.clients[conn] = true
	defer delete(s.clients, conn)

	log.Printf("WebSocket client connected")

	// Send initial state
	s.sendStateUpdate(conn)

	// Listen for client messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
		// For now, we just ignore client messages
		// In a more complete implementation, we could handle client-side state updates here
	}
}

// broadcastStateUpdate sends state updates to all connected WebSocket clients
func (s *Server) broadcastStateUpdate() {
	if len(s.clients) == 0 {
		return
	}

	for conn := range s.clients {
		if err := s.sendStateUpdate(conn); err != nil {
			log.Printf("Error sending state update: %v", err)
			delete(s.clients, conn)
			conn.Close()
		}
	}
}

// sendStateUpdate sends the current state to a specific WebSocket connection
func (s *Server) sendStateUpdate(conn *websocket.Conn) error {
	camera := s.appState.GetCamera()
	water := s.appState.GetWater()

	stateUpdate := map[string]interface{}{
		"type":    "state_update",
		"clock":   s.appState.GetClock(),
		"scenery": s.appState.GetScenery(),
		"camera": map[string]interface{}{
			"position":   [3]float32{camera.GetPosition().X, camera.GetPosition().Y, camera.GetPosition().Z},
			"viewMatrix": camera.GetViewMatrix().ToSlice(),
		},
		"water": water,
	}

	return conn.WriteJSON(stateUpdate)
}

// GetPort returns the server port
func (s *Server) GetPort() int {
	return s.port
}

// GetAssetsManager returns the assets manager
func (s *Server) GetAssetsManager() *assets.Assets {
	return s.assets
}

// GetAppState returns the application state
func (s *Server) GetAppState() *state.State {
	return s.appState
}
