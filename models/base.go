package models

type Count struct {
	Count int `json:"count"`
}
type ResponseDelete struct {
	Id 			string	`json:"id"`
	Message		string	`json:"message"`
} 
type MetaPagination struct {
	Page          int `json:"page"`
	Total         int `json:"total_pages"`
	TotalRecords  int `json:"total_records"`
	Prev          int `json:"prev"`
	Next          int `json:"next"`
	RecordPerPage int `json:"record_per_page"`
}
