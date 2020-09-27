package exchange


import (
	"net/http"
	"encoding/json"	
	"errors"
	

)


func GetRates() (map[string]float32, error) {
	resp, err := http.Get("https://api.exchangeratesapi.io/latest?base=RUB")
	if err != nil {
		return nil,errors.New("Could not get rates")
	}
	defer resp.Body.Close()

	res := struct {
		Error string
		Rates map[string]float32
	}{

	}
	
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, errors.New("Could not parse rates")
	}

	if res.Error != "" { 
		return nil, errors.New("Could not get rates")
	}

	return res.Rates, nil
}
