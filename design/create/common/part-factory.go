package common


type PartFactory interface {
	CreateWheel() Wheel
	CreateChar() Char
}


type BMWPartFactory struct {
}

func(BMWPartFactory) CreateWheel() Wheel {
	return &AWheel{BaseWheel{Name: "A"}}
}

func(BMWPartFactory) CreateChar() Char {
	return &AChar{BaseChar{Name: "A"}}
}


type BYDPartFactory struct {

}

func(BYDPartFactory) CreateWheel() Wheel {
	return &BWheel{BaseWheel{Name: "B"}}
}

func(BYDPartFactory) CreateChar() Char {
	return &BChar{BaseChar{Name: "B"}}
}



type BenChiPartFactory struct {
}

func(BenChiPartFactory) CreateWheel() Wheel {
	return &BWheel{BaseWheel{Name: "B"}}
}

func(BenChiPartFactory) CreateChar() Char {
	return &BChar{BaseChar{Name: "B"}}
}


