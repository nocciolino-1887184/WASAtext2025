package api

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"wasatext-1887184/service/database"
)

func (rt *_router) delMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	strId := ps.ByName("groupId")
	groupId, err := strconv.Atoi(strId)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	userId := ps.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	//   una volta prese queste info, non mi resta che andare nella tabella member e cancellarmi
	del, err := database.AppDatabase.DeleteMember(rt.db, groupId, id)
	if err != nil {
		http.Error(w, "Error in DELETE member ", http.StatusInternalServerError)
		return
	}

	if !del {
		w.WriteHeader(http.StatusNotFound) //  404, errore
		return
	}

	w.WriteHeader(http.StatusNoContent) //  204, cancellato senza problemi
	log.Println("Member deleted")
	//  cancellato anche dal gruppo grazie a delete on cascade
}
