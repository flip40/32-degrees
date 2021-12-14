package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/flip40/32-degrees/mysql"
)

const (
	AddDataFieldType   = "type"
	AddDataFieldSource = "source"
	AddDataFieldValue  = "value"

	DataTypeTemperature = "temperature"
	DataTypeHumidity    = "humidity"
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
	} else {
		switch item.Type {
		case DataTypeTemperature, DataTypeHumidity:
			// OK
		default:
			return nil, errors.New(fmt.Sprintf("type '%s' is invalid", item.Type))
		}
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

func (h *Handler) AddDataHandler(w http.ResponseWriter, r *http.Request) {
	fields, err := DataFieldsFromRequest(r)
	if err != nil {
		err := fmt.Sprintf("error: %s", err)
		fmt.Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// Debug logging to console
	fmt.Printf("%s: %s, %s: %s, %s: %s\n", AddDataFieldType, fields.Type, AddDataFieldSource, fields.Source, AddDataFieldValue, fields.Value)

	if _, err := h.MySQL.SaveData(&mysql.Data{
		Source: fields.Source,
		Type:   fields.Type,
		Value:  fields.Value,
	}); err != nil {
		fmt.Printf("failed to save data, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to save data, %s", err)))
	}

	// Response
	// w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func (h *Handler) GetDataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.MySQL.GetData()
	if err != nil {
		fmt.Printf("failed to get data, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get data, %s", err)))
	}

	response, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("failed to marshal data, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to marshal data, %s", err)))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
