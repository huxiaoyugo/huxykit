package common


type Char interface {
	GetName()string
}

type BaseChar struct {
	Name string
}

func(b BaseChar) GetName()string {
	return b.Name
}

type AChar struct {
	BaseChar
}

type BChar struct {
	BaseChar
}




