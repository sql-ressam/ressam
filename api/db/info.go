package db

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a API) DBInfo(w http.ResponseWriter, r *http.Request) {
	info, err := a.exporter.GetDBInfo(r.Context())
	if err != nil {
		http.Error(w, "get db info: "+err.Error(), http.StatusInternalServerError)
	}

	res, err := json.Marshal(info)
	if err != nil {
		http.Error(w, "marshal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		log.Println("write:", err.Error())
	}
}
