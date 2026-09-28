package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasatext-1887184/service/database"
)

func (rt *_router) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	type User struct {
		Nickname string `json:"nickname"`
		Id       int    `json:"id"`
	}
	var input User
	err := json.NewDecoder(r.Body).Decode(&input)	// vede se c'è input dal body
	if err != nil || input.Nickname == "" {
		http.Error(w, "Nickname mancante o non valido", http.StatusBadRequest)
		return
	}
	var userId int
	test, err := database.AppDatabase.IsAlreadyIn(rt.db, input.Nickname)
	if err != nil {
		http.Error(w, "Error in check existence", http.StatusInternalServerError)
		return
	}
	if !test {
		// puoi controllare anche se il nick già esiste (nel database)
		userId, err = database.AppDatabase.NewUser(rt.db, input.Nickname)
		if err != nil {
			http.Error(w, "Error in create user ", http.StatusInternalServerError)
			return
		}
	} else {
		// se esiste, recupera l'id dal DB
		userId, err = database.AppDatabase.IdFromNick(rt.db, input.Nickname)
		if err != nil {
			http.Error(w, "Error in retrive ID", http.StatusInternalServerError)
			return
		}
	}
	input.Id = userId
	if err := json.NewEncoder(w).Encode(input); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	/* restituiamo l'ID dell'utente che entra */

	// qui si controlla solo il body

}
