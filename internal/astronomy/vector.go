package astronomy

import (
	"math"
)

// Vector3 nampung koordinat X, Y, Z (biasanya dalam satuan KM atau AU)
type Vector3 struct {
	X, Y, Z float64
}

// Sub: Pengurangan Vektor (v1 - v2)
func (v1 Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

// Dot: Dot Product
func (v1 Vector3) Dot(v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Norm: Panjang/Magnitudo Vektor
func (v1 Vector3) Norm() float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)
}
