package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasatext-1887184/service/database"
)

func (rt *_router) getMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	str := ps.ByName("groupId")
	id_group, err := strconv.Atoi(str)
	if err != nil {
		http.Error(w, "Invalid groupId parameter", http.StatusBadRequest)
		return
	}

	// FUNZIONE PER vedere memebri da dentro il gruppo
	newMem, err := database.AppDatabase.ShowMembers(rt.db, id_group)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No users in this group", http.StatusNotFound) // 404
			return
		}
		http.Error(w, "Error in get members", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK) // 200 OK

	if err := json.NewEncoder(w).Encode(newMem); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
