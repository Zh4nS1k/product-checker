package models

type Product struct {
	ID         string `json:"id"`
	Barcode    string `json:"barcode"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	IsOriginal bool   `json:"is_original"`
}
