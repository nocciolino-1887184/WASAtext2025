package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wasatext-1887184/service/database"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// informa il client che la risposta sarà json
	w.Header().Set("content-type", "application/json")
	// il sender viene preso dal token
	auth := r.Header.Get("Authorization") // nell'header devo aggiungere "authorization"
	nickname := strings.TrimPrefix(auth, "Bearer ")

	convoIdStr := ps.ByName("convoId")
	convoId, err := strconv.Atoi(convoIdStr)
	if err != nil {
		http.Error(w, "Invalid convoId parameter", http.StatusBadRequest)
		return
	}

	type Mess_req struct {
		// alcuni campi possono essere NULL *
		MessId      int     `json:"id"`
		ConvoId     int     `json:"convoid"`
		Content     *string `json:"content"` // nel frontend vogliamo una stringa
		Text        *string `json:"testo"`
		ContentType string  `json:"media_type"` //  "testo", "foto", "video", ecc.
		Reply       int     `json:"reply"`
		// tempo lettura e reaction sono create di default
		// NON salvo URL della foto perché dipende dalla mia memo locale
	}

	var msg Mess_req
	msg.ConvoId = convoId

	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Bad request: send message", http.StatusBadRequest)
		return
	}

	idUser, err := database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
		return
	}
	var contentBytes []byte // qui convertiamo in byte il contenuto, per darlo al database.
	if msg.Content != nil {
		decoded, err := base64.StdEncoding.DecodeString(*msg.Content)
		if err != nil {
			http.Error(w, "Errore nella decodifica base64", http.StatusBadRequest)
			return
		}
		contentBytes = decoded
	}

	idMess, err := database.AppDatabase.SendMessage(
		rt.db,
		idUser,
		msg.ConvoId,
		&contentBytes,
		msg.Text,
		msg.ContentType,
		msg.Reply,
	)

	if err != nil {
		http.Error(w, "Database INSERT error", http.StatusInternalServerError)
		return
	}
	msg.MessId = idMess
	log.Println("Messaggio inserito con ID:", idMess)

	risposta := map[string]int{"idMess": idMess}
	if err := json.NewEncoder(w).Encode(risposta); err != nil {
		http.Error(w, "Error encoding response idMess", http.StatusInternalServerError)
		return
	}

}

func (rt *_router) putRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	convoIdStr := ps.ByName("convoId")
	convoId, err := strconv.Atoi(convoIdStr)
	if err != nil {
		http.Error(w, "Invalid convoId parameter", http.StatusBadRequest)
		return
	}
	auth := r.Header.Get("Authorization") // nell'header devo aggiungere "authorization"
	nickname := strings.TrimPrefix(auth, "Bearer ")

	idUser, err := database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
		return
	}

	err = database.AppDatabase.ReadMessage(rt.db, convoId, idUser)
	if err != nil {
		http.Error(w, "Error in read a message: ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
