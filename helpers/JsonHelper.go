package helpeers

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data interface{}, status int, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		println(err)
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
