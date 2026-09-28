package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasatext-1887184/service/database"
)

func (rt *_router) showUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userId := ps.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	nick, err := database.AppDatabase.GetUsername(rt.db, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	propic, err := database.AppDatabase.GetPropic(rt.db, id)
	if err != nil {
		http.Error(w, "Propic not found", http.StatusNotFound)
		return
	}
	propicBase64 := base64.StdEncoding.EncodeToString(propic) // un json non può avere tipi []byte

	// devo inviare risposta come JSON
	type Nickname struct {
		Username string `json:"username"`
		Propic   string `json:"propic"`
	}

	info := Nickname{Username: nick, Propic: propicBase64} // queste variabili sono quelle ricavate dalla query e poi messe nel json
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}
