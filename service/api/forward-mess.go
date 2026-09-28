package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"wasatext-1887184/service/database"
)

func (rt *_router) forwardMess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	auth := r.Header.Get("Authorization")
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

	// dato che posso inoltrare anche a conversazioni che non esistono, scriverò
	/*	convoIdStr := ps.ByName("convoId")
		convoId, err := strconv.Atoi(convoIdStr)
		if err != nil {
			http.Error(w, "Invalid convoId parameter", http.StatusBadRequest)
			return
		}  */

	type Inoltro struct {
		MessId       int    `json:"messId"`
		ConvoId      int    `json:"convoId"`
		Destinatario string `json:"dest"`
	}

	var msg Inoltro
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var tempConv int

	if msg.ConvoId > 0 {
		tempConv = msg.ConvoId
	} else if msg.Destinatario != "" {
		destId, err := database.AppDatabase.IdFromNick(rt.db, msg.Destinatario)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
			return
		}
// controllo se la chat a cui voglio inoltrare esiste o meno
		find, err := database.AppDatabase.ExistChat(rt.db, idUser, destId)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Database error checking chat", http.StatusInternalServerError)
				return
			}
			newChat, err := database.AppDatabase.PostChat(rt.db, idUser, msg.Destinatario)
			if err != nil {
				http.Error(w, "Errore creazione chat", http.StatusInternalServerError)
				return
			}
			tempConv = newChat
		} else {
			tempConv = find
		}
		// NON ESISTE: LA CREIAMO ORA

	} else {
		http.Error(w, "Specifica convoId o destNick", http.StatusBadRequest)
		return
	}

	newId, err := database.AppDatabase.ForwardMessage(rt.db, msg.MessId, tempConv, idUser)
	if err != nil {
		http.Error(w, "Error in forward message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"newMess":   newId,
		"convoId":   tempConv,
		"forwarded": true,
	}); err != nil {
		http.Error(w, "Bad request: forward message", http.StatusBadRequest)
		return
	}
}
