package models

type Animal struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Species     string `json:"species"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}
