package api

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"wasatext-1887184/service/database"
)

func (rt *_router) delMess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	IdStr := ps.ByName("convoId")
	convoId, err := strconv.Atoi(IdStr)
	if err != nil {
		http.Error(w, "Invalid convoId parameter", http.StatusBadRequest)
		return
	}
	IdStr = ps.ByName("messageId")
	messId, err := strconv.Atoi(IdStr)
	if err != nil {
		http.Error(w, "Invalid messageId parameter", http.StatusBadRequest)
		return
	}

	del, err := database.AppDatabase.DeleteMessage(rt.db, convoId, messId)
	if err != nil {
		http.Error(w, "Error in DELETE ", http.StatusInternalServerError)
		return
	}

	if !del {
		w.WriteHeader(http.StatusNotFound) // 404, errore
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204, cancellato senza problemi
	log.Println("Message deleted")

}
