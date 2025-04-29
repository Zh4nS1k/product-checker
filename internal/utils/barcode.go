package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var countryPrefixes = map[string]string{
	"000-019": "USA", "020-029": "Restricted", "030-039": "USA",
	"040-049": "Restricted", "050-059": "USA", "060-139": "USA",
	"300-379": "France", "380": "Bulgaria", "383": "Slovenia",
	"385": "Croatia", "387": "Bosnia", "400-440": "Germany",
	"450-459": "Japan", "460-469": "Russia", "470": "Kyrgyzstan",
	"471": "Taiwan", "474": "Estonia", "475": "Latvia",
	"476": "Azerbaijan", "477": "Lithuania", "478": "Uzbekistan",
	"479": "Sri Lanka", "480": "Philippines", "481": "Belarus",
	"482": "Ukraine", "484": "Moldova", "485": "Armenia",
	"486": "Georgia", "487": "Kazakhstan", "489": "Hong Kong",
	"490-499": "Japan", "500-509": "UK", "520": "Greece",
	"528": "Lebanon", "529": "Cyprus", "530": "Albania",
	"531": "Macedonia", "535": "Malta", "539": "Ireland",
	"540-549": "Belgium", "560": "Portugal", "569": "Iceland",
	"570-579": "Denmark", "590": "Poland", "594": "Romania",
	"599": "Hungary", "600-601": "South Africa", "603": "Ghana",
	"608": "Bahrain", "609": "Mauritius", "611": "Morocco",
	"613": "Algeria", "616": "Kenya", "618": "Ivory Coast",
	"619": "Tunisia", "621": "Syria", "622": "Egypt",
	"624": "Libya", "625": "Jordan", "626": "Iran",
	"627": "Kuwait", "628": "Saudi Arabia", "629": "UAE",
	"640-649": "Finland", "690-699": "China", "700-709": "Norway",
	"729": "Israel", "730-739": "Sweden", "740": "Guatemala",
	"741": "El Salvador", "742": "Honduras", "743": "Nicaragua",
	"744": "Costa Rica", "745": "Panama", "746": "Dominican Republic",
	"750": "Mexico", "754-755": "Canada", "759": "Venezuela",
	"760-769": "Switzerland", "770": "Colombia", "773": "Uruguay",
	"775": "Peru", "777": "Bolivia", "779": "Argentina",
	"780": "Chile", "784": "Paraguay", "786": "Ecuador",
	"789-790": "Brazil", "800-839": "Italy", "840-849": "Spain",
	"850": "Cuba", "858": "Slovakia", "859": "Czech Republic",
	"860": "Serbia", "865": "Mongolia", "867": "North Korea",
	"868-869": "Turkey", "870-879": "Netherlands", "880": "South Korea",
	"884": "Cambodia", "885": "Thailand", "888": "Singapore",
	"890": "India", "893": "Vietnam", "896": "Pakistan",
	"899": "Indonesia", "900-919": "Austria", "930-939": "Australia",
	"940-949": "New Zealand", "950": "Global Office", "955": "Malaysia",
	"958": "Macau", "960-969": "Global Office", "977": "Serial Publications",
	"978-979": "Bookland", "980": "Refund Receipts", "981-984": "Coupons",
	"990-999": "Coupons",
}

func IsBarcodeValid(code string) bool {
	if len(code) != 8 && len(code) != 12 && len(code) != 13 && len(code) != 14 {
		return false
	}

	if _, err := strconv.Atoi(code); err != nil {
		return false
	}

	if len(code) == 13 {
		sum := 0
		for i, c := range code[:12] {
			digit, _ := strconv.Atoi(string(c))
			if i%2 == 0 {
				sum += digit * 1
			} else {
				sum += digit * 3
			}
		}
		checksum := (10 - (sum % 10)) % 10
		lastDigit, _ := strconv.Atoi(string(code[12]))
		return checksum == lastDigit
	}

	return true
}

func GetCountryFromBarcode(code string) string {
	if len(code) < 3 {
		return "Unknown"
	}

	prefix := code[:3]
	for rangeStr, country := range countryPrefixes {
		if strings.Contains(rangeStr, "-") {
			parts := strings.Split(rangeStr, "-")
			start := parts[0]
			end := parts[1]

			if prefix >= start && prefix <= end {
				return country
			}
		} else if prefix == rangeStr {
			return country
		}
	}

	return "Unknown"
}

func GetCountryFromIP(ip string) (string, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip + "?fields=country")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Country string `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Country, nil
}
