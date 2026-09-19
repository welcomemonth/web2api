package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 256<<20))
	dec.UseNumber()
	return dec.Decode(dst)
}
