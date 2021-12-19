package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

type PlotDataResponse struct {
	Sources map[string]*Values `json:"sources"`
}

type Values struct {
	Temperature PlotData `json:"temperature"`
	Humidity    PlotData `json:"humidity"`
}
type PlotData struct {
	Values []float64    `json:"values"`
	Times  []*time.Time `json:"times"`
}

// {												// PlotDataResponse
// 	"sources": {									// PlotDataResponse.Sources
// 		"device_1": {								// Values
// 			"temperature": [						// Values.Temperature / PlotData
// 				"values": [70],						// PlotData.Values
// 				"times": ["2021-12-18 21:08:00"]	// PlotData.Times
// 			],
// 			"humidity": [
// 				"values": [100],
// 				"times": ["2021-12-18 21:08:00"]
// 			]
// 		}
// 	}
// }

func (h *Handler) GetPlotDataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.MySQL.GetData()
	if err != nil {
		fmt.Printf("failed to get data, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to get data, %s", err)))
	}

	loc, _ := time.LoadLocation("America/Los_Angeles") // TODO: move this value to constants

	response := &PlotDataResponse{
		Sources: make(map[string]*Values),
	}

	previousVals := make(map[string]map[string]string)

	for _, row := range data {
		source, ok := response.Sources[row.Source]
		if !ok {
			source = &Values{}
		}

		// omit duplicate values in a row as a test of display
		if prevSource, ok := previousVals[row.Source]; !ok {
			previousVals[row.Source] = make(map[string]string)
		} else {
			if prevVal, ok := prevSource[row.Type]; ok && prevVal == row.Value {
				continue
			}
		}
		previousVals[row.Source][row.Type] = row.Value

		localCreatedAt := row.CreatedAt.In(loc)
		switch row.Type {
		case DataTypeTemperature:
			value, err := strconv.ParseFloat(row.Value, 64)
			if err != nil {
				fmt.Printf("error converting float data: %s\n", err)
				continue
			}

			source.Temperature.Values = append(source.Temperature.Values, value)
			source.Temperature.Times = append(source.Temperature.Times, &localCreatedAt)

		case DataTypeHumidity:
			value, err := strconv.ParseFloat(row.Value, 64)
			if err != nil {
				continue
			}

			source.Humidity.Values = append(source.Humidity.Values, value)
			source.Humidity.Times = append(source.Humidity.Times, &localCreatedAt)
		default:
			continue
		}

		response.Sources[row.Source] = source
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("failed to marshal data, %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to marshal data, %s", err)))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
