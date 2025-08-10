package math3d

import (
	"math"
)

// Quat represents a quaternion for 3D rotations
// Components are stored as X, Y, Z, W where W is the scalar part
type Quat struct {
	X, Y, Z, W float32
}

// NewQuat creates a new quaternion from components
func NewQuat(x, y, z, w float32) Quat {
	return Quat{X: x, Y: y, Z: z, W: w}
}

// Identity returns the identity quaternion (no rotation)
func QuatIdentity() Quat {
	return Quat{X: 0, Y: 0, Z: 0, W: 1}
}

// FromAxisAngle creates a quaternion from an axis and angle (in radians)
func QuatFromAxisAngle(axis Vec3, angle float32) Quat {
	halfAngle := angle * 0.5
	s := float32(math.Sin(float64(halfAngle)))
	c := float32(math.Cos(float64(halfAngle)))

	normalizedAxis := axis.Normalize()

	return Quat{
		X: normalizedAxis.X * s,
		Y: normalizedAxis.Y * s,
		Z: normalizedAxis.Z * s,
		W: c,
	}
}

// FromEuler creates a quaternion from Euler angles (in radians)
// Order: YXZ (yaw, pitch, roll)
func QuatFromEuler(yaw, pitch, roll float32) Quat {
	cy := float32(math.Cos(float64(yaw * 0.5)))
	sy := float32(math.Sin(float64(yaw * 0.5)))
	cp := float32(math.Cos(float64(pitch * 0.5)))
	sp := float32(math.Sin(float64(pitch * 0.5)))
	cr := float32(math.Cos(float64(roll * 0.5)))
	sr := float32(math.Sin(float64(roll * 0.5)))

	return Quat{
		X: sr*cp*cy - cr*sp*sy,
		Y: cr*sp*cy + sr*cp*sy,
		Z: cr*cp*sy - sr*sp*cy,
		W: cr*cp*cy + sr*sp*sy,
	}
}

// FromMat4 creates a quaternion from a rotation matrix
func QuatFromMat4(m Mat4) Quat {
	trace := m.Get(0, 0) + m.Get(1, 1) + m.Get(2, 2)

	if trace > 0 {
		s := float32(math.Sqrt(float64(trace + 1.0))) * 2 // s = 4 * qw
		return Quat{
			X: (m.Get(2, 1) - m.Get(1, 2)) / s,
			Y: (m.Get(0, 2) - m.Get(2, 0)) / s,
			Z: (m.Get(1, 0) - m.Get(0, 1)) / s,
			W: 0.25 * s,
		}
	} else if m.Get(0, 0) > m.Get(1, 1) && m.Get(0, 0) > m.Get(2, 2) {
		s := float32(math.Sqrt(float64(1.0+m.Get(0, 0)-m.Get(1, 1)-m.Get(2, 2)))) * 2 // s = 4 * qx
		return Quat{
			X: 0.25 * s,
			Y: (m.Get(0, 1) + m.Get(1, 0)) / s,
			Z: (m.Get(0, 2) + m.Get(2, 0)) / s,
			W: (m.Get(2, 1) - m.Get(1, 2)) / s,
		}
	} else if m.Get(1, 1) > m.Get(2, 2) {
		s := float32(math.Sqrt(float64(1.0+m.Get(1, 1)-m.Get(0, 0)-m.Get(2, 2)))) * 2 // s = 4 * qy
		return Quat{
			X: (m.Get(0, 1) + m.Get(1, 0)) / s,
			Y: 0.25 * s,
			Z: (m.Get(1, 2) + m.Get(2, 1)) / s,
			W: (m.Get(0, 2) - m.Get(2, 0)) / s,
		}
	} else {
		s := float32(math.Sqrt(float64(1.0+m.Get(2, 2)-m.Get(0, 0)-m.Get(1, 1)))) * 2 // s = 4 * qz
		return Quat{
			X: (m.Get(0, 2) + m.Get(2, 0)) / s,
			Y: (m.Get(1, 2) + m.Get(2, 1)) / s,
			Z: 0.25 * s,
			W: (m.Get(1, 0) - m.Get(0, 1)) / s,
		}
	}
}

// Add adds two quaternions
func (q Quat) Add(other Quat) Quat {
	return Quat{
		X: q.X + other.X,
		Y: q.Y + other.Y,
		Z: q.Z + other.Z,
		W: q.W + other.W,
	}
}

// Sub subtracts another quaternion from this one
func (q Quat) Sub(other Quat) Quat {
	return Quat{
		X: q.X - other.X,
		Y: q.Y - other.Y,
		Z: q.Z - other.Z,
		W: q.W - other.W,
	}
}

// Scale multiplies the quaternion by a scalar
func (q Quat) Scale(s float32) Quat {
	return Quat{
		X: q.X * s,
		Y: q.Y * s,
		Z: q.Z * s,
		W: q.W * s,
	}
}

// Multiply multiplies two quaternions (Hamilton product)
func (q Quat) Multiply(other Quat) Quat {
	return Quat{
		X: q.W*other.X + q.X*other.W + q.Y*other.Z - q.Z*other.Y,
		Y: q.W*other.Y - q.X*other.Z + q.Y*other.W + q.Z*other.X,
		Z: q.W*other.Z + q.X*other.Y - q.Y*other.X + q.Z*other.W,
		W: q.W*other.W - q.X*other.X - q.Y*other.Y - q.Z*other.Z,
	}
}

// Dot calculates the dot product of two quaternions
func (q Quat) Dot(other Quat) float32 {
	return q.X*other.X + q.Y*other.Y + q.Z*other.Z + q.W*other.W
}

// Length calculates the length (magnitude) of the quaternion
func (q Quat) Length() float32 {
	return float32(math.Sqrt(float64(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)))
}

// LengthSquared calculates the squared length of the quaternion
func (q Quat) LengthSquared() float32 {
	return q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
}

// Normalize returns a normalized quaternion
func (q Quat) Normalize() Quat {
	length := q.Length()
	if length == 0 {
		return QuatIdentity()
	}
	return q.Scale(1.0 / length)
}

// Conjugate returns the conjugate of the quaternion
func (q Quat) Conjugate() Quat {
	return Quat{X: -q.X, Y: -q.Y, Z: -q.Z, W: q.W}
}

// Inverse returns the inverse of the quaternion
func (q Quat) Inverse() Quat {
	lengthSq := q.LengthSquared()
	if lengthSq == 0 {
		return QuatIdentity()
	}
	return q.Conjugate().Scale(1.0 / lengthSq)
}

// RotateVec3 rotates a Vec3 by this quaternion
func (q Quat) RotateVec3(v Vec3) Vec3 {
	// Convert vector to quaternion
	vecQuat := Quat{X: v.X, Y: v.Y, Z: v.Z, W: 0}

	// Perform rotation: q * v * q^-1
	result := q.Multiply(vecQuat).Multiply(q.Conjugate())

	return Vec3{X: result.X, Y: result.Y, Z: result.Z}
}

// ToMat4 converts the quaternion to a 4x4 rotation matrix
func (q Quat) ToMat4() Mat4 {
	// Normalize the quaternion first
	norm := q.Normalize()

	x, y, z, w := norm.X, norm.Y, norm.Z, norm.W
	xx, yy, zz := x*x, y*y, z*z
	xy, xz, yz := x*y, x*z, y*z
	wx, wy, wz := w*x, w*y, w*z

	return Mat4{
		1 - 2*(yy+zz), 2*(xy+wz), 2*(xz-wy), 0,
		2*(xy-wz), 1 - 2*(xx+zz), 2*(yz+wx), 0,
		2*(xz+wy), 2*(yz-wx), 1 - 2*(xx+yy), 0,
		0, 0, 0, 1,
	}
}

// ToAxisAngle converts the quaternion to axis-angle representation
func (q Quat) ToAxisAngle() (Vec3, float32) {
	norm := q.Normalize()

	// Handle identity quaternion
	if norm.W >= 1.0 {
		return Vec3{1, 0, 0}, 0
	}

	angle := 2.0 * float32(math.Acos(float64(norm.W)))
	s := float32(math.Sqrt(float64(1 - norm.W*norm.W)))

	if s < 0.001 { // Avoid division by zero
		return Vec3{1, 0, 0}, 0
	}

	axis := Vec3{
		X: norm.X / s,
		Y: norm.Y / s,
		Z: norm.Z / s,
	}

	return axis, angle
}

// ToEuler converts the quaternion to Euler angles (in radians)
// Returns yaw, pitch, roll
func (q Quat) ToEuler() (float32, float32, float32) {
	norm := q.Normalize()
	x, y, z, w := norm.X, norm.Y, norm.Z, norm.W

	// Roll (x-axis rotation)
	sinr_cosp := 2 * (w*x + y*z)
	cosr_cosp := 1 - 2*(x*x+y*y)
	roll := float32(math.Atan2(float64(sinr_cosp), float64(cosr_cosp)))

	// Pitch (y-axis rotation)
	sinp := 2 * (w*y - z*x)
	var pitch float32
	if math.Abs(float64(sinp)) >= 1 {
		pitch = float32(math.Copysign(math.Pi/2, float64(sinp))) // Use 90 degrees if out of range
	} else {
		pitch = float32(math.Asin(float64(sinp)))
	}

	// Yaw (z-axis rotation)
	siny_cosp := 2 * (w*z + x*y)
	cosy_cosp := 1 - 2*(y*y+z*z)
	yaw := float32(math.Atan2(float64(siny_cosp), float64(cosy_cosp)))

	return yaw, pitch, roll
}

// Slerp performs spherical linear interpolation between two quaternions
func (q Quat) Slerp(other Quat, t float32) Quat {
	q1 := q.Normalize()
	q2 := other.Normalize()

	dot := q1.Dot(q2)

	// If the dot product is negative, the quaternions have opposite handed-ness
	// and slerp won't take the shorter path. Fix by reversing one quaternion.
	if dot < 0.0 {
		q2 = q2.Scale(-1)
		dot = -dot
	}

	// If the inputs are too close for comfort, linearly interpolate
	if dot > 0.9995 {
		result := q1.Add(q2.Sub(q1).Scale(t))
		return result.Normalize()
	}

	theta0 := float32(math.Acos(float64(dot)))
	theta := theta0 * t

	q2 = q2.Sub(q1.Scale(dot)).Normalize()

	cosTheta := float32(math.Cos(float64(theta)))
	sinTheta := float32(math.Sin(float64(theta)))

	return q1.Scale(cosTheta).Add(q2.Scale(sinTheta))
}

// Lerp performs linear interpolation between two quaternions
func (q Quat) Lerp(other Quat, t float32) Quat {
	result := q.Add(other.Sub(q).Scale(t))
	return result.Normalize()
}

// AngleTo calculates the angle between this quaternion and another (in radians)
func (q Quat) AngleTo(other Quat) float32 {
	dot := math.Abs(float64(q.Normalize().Dot(other.Normalize())))
	if dot >= 1.0 {
		return 0
	}
	return float32(math.Acos(dot)) * 2
}
