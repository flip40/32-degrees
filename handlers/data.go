package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	AddDataFieldType   = "type"
	AddDataFieldSource = "source"
	AddDataFieldValue  = "value"

	DataTypeTemperature = "tempurature"
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

// TODO: remove these once stored in DB
var temp, hum float64
var source string

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

	// TODO: remove these and save to DB instead
	switch fields.Type {
	case DataTypeTemperature:
		fmt.Println("saving temperature to globals")
		val, err := strconv.ParseFloat(fields.Value, 64)
		if err != nil {
			// TODO: log error
		}
		source = fields.Source
		temp = val
		fmt.Println("saved temperature to globals")
	case DataTypeHumidity:
		fmt.Println("saving humidity to globals")
		val, err := strconv.ParseFloat(fields.Value, 64)
		if err != nil {
			// TODO: log error
		}
		source = fields.Source
		hum = val
		fmt.Println("saved humidity to globals")
	default:
		fmt.Printf("no match for fields.Type: %s\n", fields.Type)
	}

	// Response
	// w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func (h *Handler) GetDataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("source: %s\ntempurature: %.1f\nhumidity: %.1f", source, temp, hum)))
}
