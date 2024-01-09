package model

type Transaction struct {
	Id          int
	Name        string
	Description string
	Direction   string
	Value       int
	Currency    string
	Period      string
	Status      string
}
