package api

import (
	"aviasearch/engine"
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetTicketsEndpoint(w http.ResponseWriter, r *http.Request) {
	response := &Response{}

	sort := r.URL.Query().Get("sort")
	orderBy := r.URL.Query().Get("order")
	variants, err := engine.GetVariants(sort, orderBy)
	w.Header().Set("Content-Type", "application/json")
	meta := Meta{
		Code:    http.StatusOK,
		Message: "OK",
	}
	response.Meta = meta
	if err != nil {
		resp, _ := json.Marshal(response)
		_, err := w.Write([]byte(resp))
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
	response.Data = variants
	resp, _ := json.Marshal(response)
	_, err = w.Write([]byte(resp))
	if err != nil {
		log.Println(err.Error())
	}
}
func GetTicketEndpoint(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	orderBy := r.URL.Query().Get("order")

	variants, err := engine.GetVariants(sort, orderBy)
	if err != nil {
		log.Println(err.Error())
		return
	}
	response := &Response{
		Data: variants[0],
	}
	response.Meta.Code = http.StatusOK
	response.Meta.Message = "Ok"
	resp, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(resp))
	if err != nil {
		log.Println(err.Error())
	}
}
