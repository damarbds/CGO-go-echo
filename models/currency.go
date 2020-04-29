package models

type Currency struct {
	Id 		int 		`json:"id"`
	Code 	string		`json:"code"`
	Name 	string		`json:"name"`
	Symbol string		`json:"symbol"`
}