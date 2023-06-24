package model

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type GetItemRequest struct {
	Id string `json:"id"`
}

type ReloadRequest struct {
	Path string `json:"path"`
}

func DecodeGetItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := GetItemRequest{}
	id := strings.Split(r.URL.Path, "/")[2]
	request.Id = id
	return request, nil
}

func DecodeReloadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ReloadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
