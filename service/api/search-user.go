package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasatext-1887184/service/database"
)

func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// un GET non ha un input, quindi potrei scrivere il parametro come query
	nickname := r.URL.Query().Get("nickname")
	if nickname == "" {
		http.Error(w, "Nickname mancante o non valido", http.StatusBadRequest)
		return
	}

	/* Prima controllo abbia scritto un nickname corretto
	poi controllo se esiste nel database come profilo
	se esiste, ne ricavo anche immagine profilo, mostrandolo in una finestra */

	var userId int
	test, err := database.AppDatabase.IsAlreadyIn(rt.db, nickname)
	if err != nil {
		http.Error(w, "Error in searcing NIckname", http.StatusInternalServerError)
		return
	}
	if !test {

		http.Error(w, "Error in check user ", http.StatusNotFound)
		return

	}
	// se esiste, recupera l'id dal DB
	userId, err = database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		http.Error(w, "Error in retrive ID", http.StatusInternalServerError)
		return
	}

	photoUser, err := database.AppDatabase.GetPropic(rt.db, userId)
	if err != nil {
		http.Error(w, "Error in retrive user's propic", http.StatusInternalServerError)
		return
	}
	propicBase64 := base64.StdEncoding.EncodeToString(photoUser)

	// Struct di output
	output := struct {
		Nickname string `json:"nickname"`
		Propic   string `json:"propic"` // base64 encoded image
	}{
		Nickname: nickname,
		Propic:   propicBase64,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		http.Error(w, "Errore nella risposta JSON", http.StatusInternalServerError)
	}
}

func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	utenti, err := database.AppDatabase.AllUsers(rt.db)
	if err != nil {
		http.Error(w, "Error in retrive all usernames", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(utenti); err != nil {
		http.Error(w, "Errore nella risposta JSON", http.StatusInternalServerError)
	}
}
