package types

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetTicketIDRequest(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	uuidStr, ok := vars["booking_id"]
	if !ok || len(uuidStr) == 0 {
		return uuid.Nil, BadUUID
	}
	id, err := uuid.Parse(uuidStr)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
