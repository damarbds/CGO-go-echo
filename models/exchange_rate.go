package models

type ExChangeRate struct {
	Id 			int 	`json:"id"`
	Date 		string	`json:"date"`
	From		string	`json:"from"`
	To 			string 	`json:"to"`
	Rates		float64 	`json:"rates"`
}
