package http

import (
	"encoding/json"
	"log"
)

type Response struct {
	Error   error `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type AddRequest struct {
	Record string `json:"record"`
}

func (r Response) ToJSON() []byte {
	d, err := json.Marshal(r)
	if err != nil {
		log.Printf("[ERR] ToJSON error: %s", err.Error())
	}
	return d
}