package handlers

import (
	"DB-forums/models"
	"encoding/json"
)

func toMessage(message string) []byte {
	jsonErr, err := json.Marshal(models.JsonError{Message: message})
	if err != nil {
		return []byte("")
	}
	return jsonErr
}
