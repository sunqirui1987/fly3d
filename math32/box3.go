package math32

type Box3 struct {
	Minimum *Vector3
	Maximum *Vector3
}

func NewBox3(min, max *Vector3) *Box3 {

	b := new(Box3)
	b.Minimum = min.Clone()
	b.Maximum = max.Clone()
	return b
}
