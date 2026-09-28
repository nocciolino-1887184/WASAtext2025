package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wasatext-1887184/service/database"
)

type Reactions struct {
	Emoji  string `json:"emoji"`
	Autore string `json:"user"`
}

// un POST che aggiunge una reaction alla tabella e in una lista
func (rt *_router) addReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json") // Specifica che la risposta è JSon

	auth := r.Header.Get("Authorization")           // nell'header devo aggiungere "authorization"
	nickname := strings.TrimPrefix(auth, "Bearer ") // sarà poi nel database che ci sarà il controllo del ID!

	str := ps.ByName("messageId")
	id_mess, err := strconv.Atoi(str)
	if err != nil {
		http.Error(w, "Invalid messageId parameter", http.StatusBadRequest)
		return
	}

	type React struct {
		IdEmoji string `json:"idEmoji"`
	}
	var input React
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil { // scrive nel JSON della risposta
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	emoji := input.IdEmoji

	userId, err := database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found in add reaction", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
		return
	}

	err = database.AppDatabase.AddComment(rt.db, emoji, id_mess, userId)
	if err != nil {
		log.Printf("Errore in AddComment: emoji=%s, messId=%d, userId=%d, err=%v\n", emoji, id_mess, userId, err)
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}

	react := Reactions{
		Emoji:  emoji,
		Autore: nickname,
	}

	w.WriteHeader(http.StatusCreated) // per sapere in che stato ci troviamo es 201
	if err := json.NewEncoder(w).Encode(react); err != nil {
		http.Error(w, "Errore nella codifica della emoji", http.StatusInternalServerError)
		return
	}

}

// mi serve visualizzare i dati salvati
func (rt *_router) getReactions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	str := ps.ByName("messageId")
	id_mess, err := strconv.Atoi(str)
	if err != nil {
		http.Error(w, "Invalid messageId parameter", http.StatusBadRequest)
		return
	}

	res, err := database.AppDatabase.GetComment(rt.db, id_mess)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Bad request: show all reactions", http.StatusBadRequest)
		return
	}

}

func (rt *_router) delReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json") // Specifica che la risposta è JSon

	auth := r.Header.Get("Authorization")           // nell'header devo aggiungere "authorization"
	nickname := strings.TrimPrefix(auth, "Bearer ") // sarà poi nel database che ci sarà il controllo del ID!

	str := ps.ByName("messageId")
	mess, err := strconv.Atoi(str)
	if err != nil {
		http.Error(w, "Invalid messageId parameter", http.StatusBadRequest)
		return
	}
	reactId := ps.ByName("reactId")
	if reactId == "" {
		http.Error(w, "Missing reactId", http.StatusBadRequest)
		return
	}

	ok, err := database.AppDatabase.UnComment(rt.db, reactId, mess, nickname)
	if err != nil {
		http.Error(w, "Error in delete comment", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Reaction not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
