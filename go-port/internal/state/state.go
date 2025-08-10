package state

import (
	"math"
	"sync"
	"time"

	"github.com/webgl-water-go/go-port/internal/math3d"
)

// State represents the complete application state
type State struct {
	mu       sync.RWMutex
	clock    float32
	camera   *Camera
	mouse    *Mouse
	water    *Water
	scenery  bool
	lastTime time.Time
}

// NewState creates a new application state
func NewState() *State {
	return &State{
		clock:    0.0,
		camera:   NewCamera(),
		mouse:    NewMouse(),
		water:    NewWater(),
		scenery:  true,
		lastTime: time.Now(),
	}
}

// GetClock returns the current clock time in milliseconds
func (s *State) GetClock() float32 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clock
}

// GetCamera returns a copy of the camera state
func (s *State) GetCamera() Camera {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.camera
}

// GetWater returns a copy of the water state
func (s *State) GetWater() Water {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.water
}

// GetScenery returns whether scenery should be shown
func (s *State) GetScenery() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scenery
}

// Update processes a state message
func (s *State) Update(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch m := msg.(type) {
	case *AdvanceClockMessage:
		s.clock += m.DeltaTime
	case *MouseDownMessage:
		s.mouse.SetPressed(true)
		s.mouse.SetPos(m.X, m.Y)
	case *MouseUpMessage:
		s.mouse.SetPressed(false)
	case *MouseMoveMessage:
		if !s.mouse.GetPressed() {
			return
		}
		oldX, oldY := s.mouse.GetPos()
		xDelta := float32(oldX - m.X)
		yDelta := float32(m.Y - oldY)

		s.camera.OrbitLeftRight(xDelta / 50.0)
		s.camera.OrbitUpDown(yDelta / 50.0)
		s.mouse.SetPos(m.X, m.Y)
	case *ZoomMessage:
		s.camera.Zoom(m.Delta)
	case *SetReflectivityMessage:
		s.water.Reflectivity = m.Value
	case *SetFresnelMessage:
		s.water.FresnelStrength = m.Value
	case *SetWaveSpeedMessage:
		s.water.WaveSpeed = m.Value
	case *UseReflectionMessage:
		s.water.UseReflection = m.Value
	case *UseRefractionMessage:
		s.water.UseRefraction = m.Value
	case *ShowSceneryMessage:
		s.scenery = m.Value
	}
}

// Camera represents the camera state
type Camera struct {
	position    math3d.Vec3
	target      math3d.Vec3
	up          math3d.Vec3
	distance    float32
	yaw         float32
	pitch       float32
	minDistance float32
	maxDistance float32
	minPitch    float32
	maxPitch    float32
}

// NewCamera creates a new camera with default settings
func NewCamera() *Camera {
	return &Camera{
		position:    math3d.NewVec3(0, 5, 10),
		target:      math3d.NewVec3(0, 0, 0),
		up:          math3d.Vec3Up,
		distance:    15.0,
		yaw:         0.0,
		pitch:       0.3,
		minDistance: 5.0,
		maxDistance: 50.0,
		minPitch:    -1.5,
		maxPitch:    1.5,
	}
}

// GetViewMatrix returns the view matrix for this camera
func (c *Camera) GetViewMatrix() math3d.Mat4 {
	c.updatePosition()
	return math3d.LookAt(c.position, c.target, c.up)
}

// GetPosition returns the camera position
func (c *Camera) GetPosition() math3d.Vec3 {
	c.updatePosition()
	return c.position
}

// OrbitLeftRight rotates the camera left/right around the target
func (c *Camera) OrbitLeftRight(delta float32) {
	c.yaw += delta
}

// OrbitUpDown rotates the camera up/down around the target
func (c *Camera) OrbitUpDown(delta float32) {
	c.pitch += delta
	if c.pitch < c.minPitch {
		c.pitch = c.minPitch
	}
	if c.pitch > c.maxPitch {
		c.pitch = c.maxPitch
	}
}

// Zoom changes the camera distance from the target
func (c *Camera) Zoom(delta float32) {
	c.distance += delta
	if c.distance < c.minDistance {
		c.distance = c.minDistance
	}
	if c.distance > c.maxDistance {
		c.distance = c.maxDistance
	}
}

// updatePosition updates the camera position based on yaw, pitch, and distance
func (c *Camera) updatePosition() {
	x := c.distance * float32(math.Cos(float64(c.pitch))) * float32(math.Sin(float64(c.yaw)))
	y := c.distance * float32(math.Sin(float64(c.pitch)))
	z := c.distance * float32(math.Cos(float64(c.pitch))) * float32(math.Cos(float64(c.yaw)))

	c.position = c.target.Add(math3d.NewVec3(x, y, z))
}

// Mouse represents mouse input state
type Mouse struct {
	x       int32
	y       int32
	pressed bool
}

// NewMouse creates a new mouse state
func NewMouse() *Mouse {
	return &Mouse{
		x:       0,
		y:       0,
		pressed: false,
	}
}

// SetPos sets the mouse position
func (m *Mouse) SetPos(x, y int32) {
	m.x = x
	m.y = y
}

// GetPos returns the current mouse position
func (m *Mouse) GetPos() (int32, int32) {
	return m.x, m.y
}

// SetPressed sets whether the mouse is pressed
func (m *Mouse) SetPressed(pressed bool) {
	m.pressed = pressed
}

// GetPressed returns whether the mouse is currently pressed
func (m *Mouse) GetPressed() bool {
	return m.pressed
}

// Water represents water rendering properties
type Water struct {
	Reflectivity    float32
	FresnelStrength float32
	WaveSpeed       float32
	UseReflection   bool
	UseRefraction   bool
}

// NewWater creates new water state with default properties
func NewWater() *Water {
	return &Water{
		Reflectivity:    0.6,
		FresnelStrength: 2.0,
		WaveSpeed:       0.03,
		UseReflection:   true,
		UseRefraction:   true,
	}
}

// GetDudvOffset calculates the current dudv texture offset based on time and wave speed
func (w *Water) GetDudvOffset(clockTime float32) float32 {
	return (clockTime / 1000.0) * w.WaveSpeed
}

// Message represents a state update message
type Message interface {
	message()
}

// AdvanceClockMessage updates the application clock
type AdvanceClockMessage struct {
	DeltaTime float32
}

func (*AdvanceClockMessage) message() {}

// MouseDownMessage represents a mouse press event
type MouseDownMessage struct {
	X, Y int32
}

func (*MouseDownMessage) message() {}

// MouseUpMessage represents a mouse release event
type MouseUpMessage struct{}

func (*MouseUpMessage) message() {}

// MouseMoveMessage represents a mouse move event
type MouseMoveMessage struct {
	X, Y int32
}

func (*MouseMoveMessage) message() {}

// ZoomMessage represents a zoom event
type ZoomMessage struct {
	Delta float32
}

func (*ZoomMessage) message() {}

// SetReflectivityMessage sets water reflectivity
type SetReflectivityMessage struct {
	Value float32
}

func (*SetReflectivityMessage) message() {}

// SetFresnelMessage sets water fresnel strength
type SetFresnelMessage struct {
	Value float32
}

func (*SetFresnelMessage) message() {}

// SetWaveSpeedMessage sets water wave animation speed
type SetWaveSpeedMessage struct {
	Value float32
}

func (*SetWaveSpeedMessage) message() {}

// UseReflectionMessage toggles water reflection
type UseReflectionMessage struct {
	Value bool
}

func (*UseReflectionMessage) message() {}

// UseRefractionMessage toggles water refraction
type UseRefractionMessage struct {
	Value bool
}

func (*UseRefractionMessage) message() {}

// ShowSceneryMessage toggles scenery rendering
type ShowSceneryMessage struct {
	Value bool
}

func (*ShowSceneryMessage) message() {}
