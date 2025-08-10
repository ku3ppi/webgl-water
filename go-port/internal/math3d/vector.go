package math3d

import (
	"math"
)

// Vec2 represents a 2D vector
type Vec2 struct {
	X, Y float32
}

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
}

// Vec4 represents a 4D vector
type Vec4 struct {
	X, Y, Z, W float32
}

// NewVec2 creates a new 2D vector
func NewVec2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

// NewVec3 creates a new 3D vector
func NewVec3(x, y, z float32) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

// NewVec4 creates a new 4D vector
func NewVec4(x, y, z, w float32) Vec4 {
	return Vec4{X: x, Y: y, Z: z, W: w}
}

// Vec2 methods
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vec2) Scale(s float32) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v Vec2) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Normalize() Vec2 {
	length := v.Length()
	if length == 0 {
		return Vec2{}
	}
	return Vec2{X: v.X / length, Y: v.Y / length}
}

// Vec3 methods
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

func (v Vec3) Scale(s float32) Vec3 {
	return Vec3{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vec3) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Normalize() Vec3 {
	length := v.Length()
	if length == 0 {
		return Vec3{}
	}
	return Vec3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

func (v Vec3) Distance(other Vec3) float32 {
	return v.Sub(other).Length()
}

// Vec4 methods
func (v Vec4) Add(other Vec4) Vec4 {
	return Vec4{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z, W: v.W + other.W}
}

func (v Vec4) Sub(other Vec4) Vec4 {
	return Vec4{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z, W: v.W - other.W}
}

func (v Vec4) Scale(s float32) Vec4 {
	return Vec4{X: v.X * s, Y: v.Y * s, Z: v.Z * s, W: v.W * s}
}

func (v Vec4) Dot(other Vec4) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

func (v Vec4) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)))
}

func (v Vec4) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
}

func (v Vec4) Normalize() Vec4 {
	length := v.Length()
	if length == 0 {
		return Vec4{}
	}
	return Vec4{X: v.X / length, Y: v.Y / length, Z: v.Z / length, W: v.W / length}
}

// ToVec3 converts Vec4 to Vec3 by dropping the W component
func (v Vec4) ToVec3() Vec3 {
	return Vec3{X: v.X, Y: v.Y, Z: v.Z}
}

// ToVec3Homogeneous converts Vec4 to Vec3 by dividing by W (homogeneous coordinates)
func (v Vec4) ToVec3Homogeneous() Vec3 {
	if v.W == 0 {
		return Vec3{X: v.X, Y: v.Y, Z: v.Z}
	}
	return Vec3{X: v.X / v.W, Y: v.Y / v.W, Z: v.Z / v.W}
}

// Extend Vec3 to Vec4 with W component
func (v Vec3) Extend(w float32) Vec4 {
	return Vec4{X: v.X, Y: v.Y, Z: v.Z, W: w}
}

// Common vector constants
var (
	Vec3Zero    = Vec3{0, 0, 0}
	Vec3One     = Vec3{1, 1, 1}
	Vec3Up      = Vec3{0, 1, 0}
	Vec3Forward = Vec3{0, 0, -1}
	Vec3Right   = Vec3{1, 0, 0}
)
