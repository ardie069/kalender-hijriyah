package astronomy

import "math"

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) Scale(s float64) Vector3 {
	return Vector3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

func (v Vector3) Dot(other Vector3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vector3) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector3) Unit() Vector3 {
	n := v.Norm()
	if n == 0 {
		return Vector3{}
	}
	return Vector3{v.X / n, v.Y / n, v.Z / n}
}
