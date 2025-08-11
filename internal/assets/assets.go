package assets

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ku3ppi/webgl-water/internal/math3d"
)

// Assets manages all game assets (meshes, textures, etc.)
type Assets struct {
	meshes   map[string]*Mesh
	textures map[string]*Texture
	basePath string
}

// NewAssets creates a new asset manager
func NewAssets(basePath string) *Assets {
	return &Assets{
		meshes:   make(map[string]*Mesh),
		textures: make(map[string]*Texture),
		basePath: basePath,
	}
}

// Mesh represents a 3D mesh with vertices, normals, and indices
type Mesh struct {
	Name          string    `json:"name"`
	Vertices      []float32 `json:"vertices"`      // Position data (x, y, z, x, y, z, ...)
	Normals       []float32 `json:"normals"`       // Normal data (nx, ny, nz, nx, ny, nz, ...)
	TexCoords     []float32 `json:"texCoords"`     // Texture coordinates (u, v, u, v, ...)
	Indices       []uint16  `json:"indices"`       // Triangle indices
	VertexCount   int       `json:"vertexCount"`   // Number of vertices
	TriangleCount int       `json:"triangleCount"` // Number of triangles
}

// Texture represents texture metadata
type Texture struct {
	Name     string `json:"name"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format   string `json:"format"`
	FilePath string `json:"filePath"`
}

// MeshData represents the combined mesh data structure
type MeshData struct {
	Meshes []Mesh `json:"meshes"`
}

// LoadMeshes loads all meshes from the meshes data file
func (a *Assets) LoadMeshes() error {
	meshPath := filepath.Join(a.basePath, "../meshes.bytes")

	// Check if the binary file exists, if not try JSON
	jsonPath := filepath.Join(a.basePath, "meshes.json")

	var meshData MeshData
	var err error

	if _, statErr := os.Stat(meshPath); statErr == nil {
		// Load from binary file (if we implement binary format)
		meshData, err = a.loadMeshesFromBinary(meshPath)
	} else {
		// Load from JSON file
		meshData, err = a.loadMeshesFromJSON(jsonPath)
	}

	if err != nil {
		return fmt.Errorf("failed to load meshes: %w", err)
	}

	// Store meshes in the asset manager
	for _, mesh := range meshData.Meshes {
		meshCopy := mesh // Create a copy to avoid pointer issues
		a.meshes[mesh.Name] = &meshCopy
	}

	return nil
}

// loadMeshesFromJSON loads meshes from a JSON file
func (a *Assets) loadMeshesFromJSON(path string) (MeshData, error) {
	file, err := os.Open(path)
	if err != nil {
		return MeshData{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return MeshData{}, err
	}

	var meshData MeshData
	if err := json.Unmarshal(data, &meshData); err != nil {
		return MeshData{}, err
	}

	return meshData, nil
}

// loadMeshesFromBinary loads meshes from binary format (placeholder)
func (a *Assets) loadMeshesFromBinary(path string) (MeshData, error) {
	// TODO: Implement binary mesh loading
	// For now, return empty data and let it fall back to JSON
	return MeshData{}, fmt.Errorf("binary mesh loading not implemented")
}

// GetMesh returns a mesh by name
func (a *Assets) GetMesh(name string) (*Mesh, error) {
	mesh, exists := a.meshes[name]
	if !exists {
		return nil, fmt.Errorf("mesh '%s' not found", name)
	}
	return mesh, nil
}

// ListMeshes returns a list of all loaded mesh names
func (a *Assets) ListMeshes() []string {
	names := make([]string, 0, len(a.meshes))
	for name := range a.meshes {
		names = append(names, name)
	}
	return names
}

// CreateWaterMesh generates a simple water plane mesh
func (a *Assets) CreateWaterMesh(size float32, segments int) *Mesh {
	// Calculate vertex count
	vertexCount := (segments + 1) * (segments + 1)
	triangleCount := segments * segments * 2

	vertices := make([]float32, vertexCount*3)
	normals := make([]float32, vertexCount*3)
	texCoords := make([]float32, vertexCount*2)
	indices := make([]uint16, triangleCount*3)

	// Generate vertices, normals, and texture coordinates
	step := size / float32(segments)
	halfSize := size * 0.5

	for i := 0; i <= segments; i++ {
		for j := 0; j <= segments; j++ {
			index := i*(segments+1) + j
			vertIndex := index * 3
			texIndex := index * 2

			// Position (Y is always 0 for water plane)
			x := float32(j)*step - halfSize
			z := float32(i)*step - halfSize
			vertices[vertIndex] = x
			vertices[vertIndex+1] = 0.0 // Water is at Y = 0
			vertices[vertIndex+2] = z

			// Normal (always pointing up)
			normals[vertIndex] = 0.0
			normals[vertIndex+1] = 1.0
			normals[vertIndex+2] = 0.0

			// Texture coordinates
			texCoords[texIndex] = float32(j) / float32(segments)
			texCoords[texIndex+1] = float32(i) / float32(segments)
		}
	}

	// Generate indices for triangles
	indexCount := 0
	for i := 0; i < segments; i++ {
		for j := 0; j < segments; j++ {
			topLeft := uint16(i*(segments+1) + j)
			topRight := topLeft + 1
			bottomLeft := uint16((i+1)*(segments+1) + j)
			bottomRight := bottomLeft + 1

			// First triangle (top-left, bottom-left, top-right)
			indices[indexCount] = topLeft
			indices[indexCount+1] = bottomLeft
			indices[indexCount+2] = topRight
			indexCount += 3

			// Second triangle (top-right, bottom-left, bottom-right)
			indices[indexCount] = topRight
			indices[indexCount+1] = bottomLeft
			indices[indexCount+2] = bottomRight
			indexCount += 3
		}
	}

	mesh := &Mesh{
		Name:          "water_plane",
		Vertices:      vertices,
		Normals:       normals,
		TexCoords:     texCoords,
		Indices:       indices,
		VertexCount:   vertexCount,
		TriangleCount: triangleCount,
	}

	// Store the generated mesh
	a.meshes["water_plane"] = mesh

	return mesh
}

// CreateTerrainMesh generates a simple terrain mesh (for testing)
func (a *Assets) CreateTerrainMesh(size float32, segments int, heightScale float32) *Mesh {
	// Calculate vertex count
	vertexCount := (segments + 1) * (segments + 1)
	triangleCount := segments * segments * 2

	vertices := make([]float32, vertexCount*3)
	normals := make([]float32, vertexCount*3)
	texCoords := make([]float32, vertexCount*2)
	indices := make([]uint16, triangleCount*3)

	// Generate vertices and texture coordinates
	step := size / float32(segments)
	halfSize := size * 0.5

	for i := 0; i <= segments; i++ {
		for j := 0; j <= segments; j++ {
			index := i*(segments+1) + j
			vertIndex := index * 3
			texIndex := index * 2

			// Position with some simple height variation
			x := float32(j)*step - halfSize
			z := float32(i)*step - halfSize
			// Simple height function (could be replaced with noise)
			height := heightScale * (float32(i+j) / float32(segments*2))

			vertices[vertIndex] = x
			vertices[vertIndex+1] = -height // Negative so it's below water
			vertices[vertIndex+2] = z

			// Texture coordinates
			texCoords[texIndex] = float32(j) / float32(segments)
			texCoords[texIndex+1] = float32(i) / float32(segments)
		}
	}

	// Calculate normals (after vertices are set)
	a.calculateNormals(vertices, indices, normals, segments)

	// Generate indices for triangles
	indexCount := 0
	for i := 0; i < segments; i++ {
		for j := 0; j < segments; j++ {
			topLeft := uint16(i*(segments+1) + j)
			topRight := topLeft + 1
			bottomLeft := uint16((i+1)*(segments+1) + j)
			bottomRight := bottomLeft + 1

			// First triangle (top-left, bottom-left, top-right)
			indices[indexCount] = topLeft
			indices[indexCount+1] = bottomLeft
			indices[indexCount+2] = topRight
			indexCount += 3

			// Second triangle (top-right, bottom-left, bottom-right)
			indices[indexCount] = topRight
			indices[indexCount+1] = bottomLeft
			indices[indexCount+2] = bottomRight
			indexCount += 3
		}
	}

	mesh := &Mesh{
		Name:          "terrain",
		Vertices:      vertices,
		Normals:       normals,
		TexCoords:     texCoords,
		Indices:       indices,
		VertexCount:   vertexCount,
		TriangleCount: triangleCount,
	}

	// Store the generated mesh
	a.meshes["terrain"] = mesh

	return mesh
}

// calculateNormals calculates vertex normals for a mesh
func (a *Assets) calculateNormals(vertices []float32, indices []uint16, normals []float32, segments int) {
	// Initialize normals to zero
	for i := range normals {
		normals[i] = 0.0
	}

	// Calculate face normals and accumulate
	for i := 0; i < len(indices); i += 3 {
		i1, i2, i3 := int(indices[i]), int(indices[i+1]), int(indices[i+2])

		// Get vertices
		v1 := math3d.NewVec3(vertices[i1*3], vertices[i1*3+1], vertices[i1*3+2])
		v2 := math3d.NewVec3(vertices[i2*3], vertices[i2*3+1], vertices[i2*3+2])
		v3 := math3d.NewVec3(vertices[i3*3], vertices[i3*3+1], vertices[i3*3+2])

		// Calculate face normal
		edge1 := v2.Sub(v1)
		edge2 := v3.Sub(v1)
		normal := edge1.Cross(edge2).Normalize()

		// Accumulate normals for each vertex
		normals[i1*3] += normal.X
		normals[i1*3+1] += normal.Y
		normals[i1*3+2] += normal.Z

		normals[i2*3] += normal.X
		normals[i2*3+1] += normal.Y
		normals[i2*3+2] += normal.Z

		normals[i3*3] += normal.X
		normals[i3*3+1] += normal.Y
		normals[i3*3+2] += normal.Z
	}

	// Normalize all vertex normals
	for i := 0; i < len(normals); i += 3 {
		normal := math3d.NewVec3(normals[i], normals[i+1], normals[i+2]).Normalize()
		normals[i] = normal.X
		normals[i+1] = normal.Y
		normals[i+2] = normal.Z
	}
}

// RegisterTexture registers a texture with the asset manager
func (a *Assets) RegisterTexture(name, filePath string, width, height int, format string) {
	texture := &Texture{
		Name:     name,
		Width:    width,
		Height:   height,
		Format:   format,
		FilePath: filePath,
	}
	a.textures[name] = texture
}

// GetTexture returns a texture by name
func (a *Assets) GetTexture(name string) (*Texture, error) {
	texture, exists := a.textures[name]
	if !exists {
		return nil, fmt.Errorf("texture '%s' not found", name)
	}
	return texture, nil
}

// GetTextureFilePath returns the full file path for a texture
func (a *Assets) GetTextureFilePath(name string) (string, error) {
	texture, err := a.GetTexture(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(a.basePath, texture.FilePath), nil
}

// ListTextures returns a list of all registered texture names
func (a *Assets) ListTextures() []string {
	names := make([]string, 0, len(a.textures))
	for name := range a.textures {
		names = append(names, name)
	}
	return names
}

// Initialize sets up default assets
func (a *Assets) Initialize() error {
	// Create basic water and terrain meshes
	a.CreateWaterMesh(20.0, 64)        // 20x20 unit water plane with 64x64 segments
	a.CreateTerrainMesh(50.0, 32, 5.0) // 50x50 unit terrain with height variation

	// Register default textures (these should exist in the assets directory)
	a.RegisterTexture("dudvmap", "dudvmap.png", 512, 512, "rgba")
	a.RegisterTexture("normalmap", "normalmap.png", 512, 512, "rgba")
	a.RegisterTexture("stone", "stone-texture.png", 512, 512, "rgba")

	return nil
}
