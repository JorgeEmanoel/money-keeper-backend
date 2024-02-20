package model

type Transaction struct {
	Id          int
	Name        string
	Description string
	Direction   string
	Value       float64
	Currency    string
	Period      string
	Status      string
}
