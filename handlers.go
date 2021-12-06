package main

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	AddDataFieldType   = "type"
	AddDataFieldSource = "source"
	AddDataFieldValue  = "value"
)

type DataFields struct {
	Type   string
	Source string
	Value  string
}

func DataFieldsFromRequest(r *http.Request) (*DataFields, error) {
	q := r.URL.Query()

	item := &DataFields{}
	item.Type = q.Get(AddDataFieldType)
	if item.Type == "" {
		return nil, errors.New(fmt.Sprintf("query field '%s' is required", AddDataFieldType))
	}

	item.Source = q.Get(AddDataFieldSource)
	if item.Source == "" {
		return nil, errors.New(fmt.Sprintf("query field '%s' is required", AddDataFieldSource))
	}

	item.Value = q.Get(AddDataFieldValue)
	if item.Value == "" {
		return nil, errors.New(fmt.Sprintf("query field '%s' is required", AddDataFieldValue))
	}

	return item, nil
}

type AddResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func AddDataHandler(w http.ResponseWriter, r *http.Request) {
	fields, err := DataFieldsFromRequest(r)
	if err != nil {
		err := fmt.Sprintf("error: %s", err)
		fmt.Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// Debug logging to console
	fmt.Printf("%s: %s, %s: %s, %s: %s\n", AddDataFieldType, fields.Type, AddDataFieldSource, fields.Source, AddDataFieldValue, fields.Value)
	// w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
