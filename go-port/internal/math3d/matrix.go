package math3d

import (
	"math"
)

// Mat4 represents a 4x4 matrix in column-major order (OpenGL style)
type Mat4 [16]float32

// NewMat4 creates a new 4x4 matrix from 16 values in row-major order
func NewMat4(
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float32,
) Mat4 {
	// Convert from row-major to column-major
	return Mat4{
		m00, m10, m20, m30, // first column
		m01, m11, m21, m31, // second column
		m02, m12, m22, m32, // third column
		m03, m13, m23, m33, // fourth column
	}
}

// Identity returns a 4x4 identity matrix
func Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Zero returns a 4x4 zero matrix
func Zero() Mat4 {
	return Mat4{}
}

// Get returns the value at row i, column j (0-indexed)
func (m Mat4) Get(row, col int) float32 {
	return m[col*4+row]
}

// Set sets the value at row i, column j (0-indexed)
func (m *Mat4) Set(row, col int, value float32) {
	m[col*4+row] = value
}

// Multiply multiplies this matrix by another matrix
func (m Mat4) Multiply(other Mat4) Mat4 {
	var result Mat4

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := float32(0)
			for k := 0; k < 4; k++ {
				sum += m.Get(i, k) * other.Get(k, j)
			}
			result.Set(i, j, sum)
		}
	}

	return result
}

// MultiplyVec4 multiplies this matrix by a Vec4
func (m Mat4) MultiplyVec4(v Vec4) Vec4 {
	return Vec4{
		X: m[0]*v.X + m[4]*v.Y + m[8]*v.Z + m[12]*v.W,
		Y: m[1]*v.X + m[5]*v.Y + m[9]*v.Z + m[13]*v.W,
		Z: m[2]*v.X + m[6]*v.Y + m[10]*v.Z + m[14]*v.W,
		W: m[3]*v.X + m[7]*v.Y + m[11]*v.Z + m[15]*v.W,
	}
}

// MultiplyVec3 multiplies this matrix by a Vec3 (treated as Vec4 with W=1)
func (m Mat4) MultiplyVec3(v Vec3) Vec3 {
	result := m.MultiplyVec4(v.Extend(1.0))
	return result.ToVec3Homogeneous()
}

// MultiplyVec3Point multiplies this matrix by a Vec3 point (W=1)
func (m Mat4) MultiplyVec3Point(v Vec3) Vec3 {
	return m.MultiplyVec3(v)
}

// MultiplyVec3Vector multiplies this matrix by a Vec3 vector (W=0)
func (m Mat4) MultiplyVec3Vector(v Vec3) Vec3 {
	result := m.MultiplyVec4(v.Extend(0.0))
	return result.ToVec3()
}

// Transpose returns the transpose of this matrix
func (m Mat4) Transpose() Mat4 {
	return Mat4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// Translation creates a translation matrix
func Translation(x, y, z float32) Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

// TranslationVec3 creates a translation matrix from a Vec3
func TranslationVec3(v Vec3) Mat4 {
	return Translation(v.X, v.Y, v.Z)
}

// Scale creates a scale matrix
func Scale(x, y, z float32) Mat4 {
	return Mat4{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	}
}

// ScaleUniform creates a uniform scale matrix
func ScaleUniform(s float32) Mat4 {
	return Scale(s, s, s)
}

// ScaleVec3 creates a scale matrix from a Vec3
func ScaleVec3(v Vec3) Mat4 {
	return Scale(v.X, v.Y, v.Z)
}

// RotationX creates a rotation matrix around the X axis (angle in radians)
func RotationX(angle float32) Mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	return Mat4{
		1, 0, 0, 0,
		0, c, s, 0,
		0, -s, c, 0,
		0, 0, 0, 1,
	}
}

// RotationY creates a rotation matrix around the Y axis (angle in radians)
func RotationY(angle float32) Mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	return Mat4{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

// RotationZ creates a rotation matrix around the Z axis (angle in radians)
func RotationZ(angle float32) Mat4 {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	return Mat4{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// LookAt creates a view matrix looking from eye to target with given up vector
func LookAt(eye, target, up Vec3) Mat4 {
	forward := target.Sub(eye).Normalize()
	right := forward.Cross(up).Normalize()
	realUp := right.Cross(forward)

	// Negate forward for right-handed coordinate system
	forward = forward.Scale(-1)

	return Mat4{
		right.X, realUp.X, forward.X, 0,
		right.Y, realUp.Y, forward.Y, 0,
		right.Z, realUp.Z, forward.Z, 0,
		-right.Dot(eye), -realUp.Dot(eye), -forward.Dot(eye), 1,
	}
}

// Perspective creates a perspective projection matrix
func Perspective(fovY, aspect, near, far float32) Mat4 {
	f := float32(1.0 / math.Tan(float64(fovY)*0.5))

	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / (near - far), -1,
		0, 0, (2 * far * near) / (near - far), 0,
	}
}

// Ortho creates an orthographic projection matrix
func Ortho(left, right, bottom, top, near, far float32) Mat4 {
	return Mat4{
		2 / (right - left), 0, 0, 0,
		0, 2 / (top - bottom), 0, 0,
		0, 0, -2 / (far - near), 0,
		-(right + left) / (right - left),
		-(top + bottom) / (top - bottom),
		-(far + near) / (far - near), 1,
	}
}

// Inverse calculates the inverse of this matrix (using Gaussian elimination)
func (m Mat4) Inverse() (Mat4, bool) {
	// Create augmented matrix [A|I]
	var aug [4][8]float32

	// Copy original matrix to left side
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			aug[i][j] = m.Get(i, j)
		}
	}

	// Add identity matrix to right side
	for i := 0; i < 4; i++ {
		aug[i][i+4] = 1
	}

	// Forward elimination
	for i := 0; i < 4; i++ {
		// Find pivot
		maxRow := i
		for k := i + 1; k < 4; k++ {
			if math.Abs(float64(aug[k][i])) > math.Abs(float64(aug[maxRow][i])) {
				maxRow = k
			}
		}

		// Swap rows
		aug[i], aug[maxRow] = aug[maxRow], aug[i]

		// Check for singular matrix
		if aug[i][i] == 0 {
			return Mat4{}, false
		}

		// Scale pivot row
		pivot := aug[i][i]
		for j := 0; j < 8; j++ {
			aug[i][j] /= pivot
		}

		// Eliminate column
		for k := 0; k < 4; k++ {
			if k != i {
				factor := aug[k][i]
				for j := 0; j < 8; j++ {
					aug[k][j] -= factor * aug[i][j]
				}
			}
		}
	}

	// Extract inverse matrix from right side
	var result Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result.Set(i, j, aug[i][j+4])
		}
	}

	return result, true
}

// Determinant calculates the determinant of this matrix
func (m Mat4) Determinant() float32 {
	return m[0]*(m[5]*(m[10]*m[15]-m[11]*m[14])-m[6]*(m[9]*m[15]-m[11]*m[13])+m[7]*(m[9]*m[14]-m[10]*m[13])) -
		m[1]*(m[4]*(m[10]*m[15]-m[11]*m[14])-m[6]*(m[8]*m[15]-m[11]*m[12])+m[7]*(m[8]*m[14]-m[10]*m[12])) +
		m[2]*(m[4]*(m[9]*m[15]-m[11]*m[13])-m[5]*(m[8]*m[15]-m[11]*m[12])+m[7]*(m[8]*m[13]-m[9]*m[12])) -
		m[3]*(m[4]*(m[9]*m[14]-m[10]*m[13])-m[5]*(m[8]*m[14]-m[10]*m[12])+m[6]*(m[8]*m[13]-m[9]*m[12]))
}

// GetTranslation extracts the translation component from this matrix
func (m Mat4) GetTranslation() Vec3 {
	return Vec3{X: m[12], Y: m[13], Z: m[14]}
}

// SetTranslation sets the translation component of this matrix
func (m *Mat4) SetTranslation(v Vec3) {
	m[12] = v.X
	m[13] = v.Y
	m[14] = v.Z
}

// ToSlice returns the matrix as a float32 slice (useful for OpenGL)
func (m Mat4) ToSlice() []float32 {
	return m[:]
}
