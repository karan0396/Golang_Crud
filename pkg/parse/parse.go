package parse

import (
	"api/pkg/validation"

	"encoding/json"
	"fmt"
	"net/http"

)
//decoding data to some variable and validating
func Parse(e http.ResponseWriter, r *http.Request, w interface{}) error {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&w); err != nil {
		return err
	}
	err :=validation.Validate(w)
	if err != nil {
		http.Error(e,
			fmt.Sprintf("validate Error %s", err),
			http.StatusBadRequest)
			return err
	}
	return nil
}