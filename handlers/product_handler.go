package handlers

import (
	"encoding/json"
	"net/http"
	"product-checker/utils"
)

func CheckProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Barcode string `json:"barcode"`
	}
	_ = json.NewDecoder(r.Body).Decode(&input)

	valid := utils.IsBarcodeValid(input.Barcode)

	response := struct {
		Barcode    string `json:"barcode"`
		IsOriginal bool   `json:"is_original"`
		Message    string `json:"message"`
	}{
		Barcode:    input.Barcode,
		IsOriginal: valid,
		Message:    "Product is authentic ✅",
	}

	if !valid {
		response.Message = "Warning: possibly counterfeit ⚠️"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
