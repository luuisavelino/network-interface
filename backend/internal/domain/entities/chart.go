package entities

type Chart map[string]*CoverageArea

type CoverageArea struct {
	X int
	Y int
	R float64
}
