package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"wasatext-1887184/service/database"
)

func (rt *_router) changeNick(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr) // prendo ID dall'URL
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "ID non valido.",
		}); err != nil {
			log.Println("Errore nell'encode JSON:", err)
		}
		return
	}

	type Name struct { // la struct ha tanti campi quanti sono quelli del JSON
		Nickname string `json:"nickname"` // il nome è nel body
	}

	var nick Name
	if err := json.NewDecoder(r.Body).Decode(&nick); err != nil { // lettura del body
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "JSON non valido nel body della richiesta.",
		}); err != nil {
			log.Println("Errore nell'encode JSON:", err)
		}
		return
	}

	err = database.AppDatabase.PutNick(rt.db, nick.Nickname, id)
	if err != nil {
		w.WriteHeader(http.StatusConflict) // 409
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Nickname già in uso. Sceglierne un altro.",
		}); err != nil {
			log.Println("Errore nell'encode JSON:", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) changePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	idStr := ps.ByName("id")
	id, err := strconv.Atoi(idStr) // prendo ID dall'URL
	if err != nil {
		http.Error(w, "ID not valid", http.StatusBadRequest)
		return
	}

	type Photo struct {
		Propic string
	}
	var pfp Photo
	if err := json.NewDecoder(r.Body).Decode(&pfp); err != nil { // lettura del body
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	decodedPhoto, err := base64.StdEncoding.DecodeString(pfp.Propic)
	if err != nil {
		http.Error(w, "Incorrect format image", http.StatusBadRequest)
		return
	}

	err = database.AppDatabase.PutPropic(rt.db, id, decodedPhoto) // la foto avrà formato base64
	if err != nil {
		http.Error(w, "Error in database UPDATE", http.StatusInternalServerError)
		return
	}	
	w.WriteHeader(http.StatusOK)
}
