package common


type Brand int

const (
	BrandBMW Brand = iota + 1
	BrandBYD
	BrandBenchi
)

func (b Brand) ToString() string {
	switch b {
	case BrandBenchi:
		return "BenChi"
	case BrandBMW:
		return "BMW"
	case BrandBYD:
		return "BYD"
	default:
		return "default"
	}
}
