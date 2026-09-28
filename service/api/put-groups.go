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

func (rt *_router) addMember(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json") // Specifica che la risposta è JSon

	auth := r.Header.Get("Authorization")           // nell'header devo aggiungere "authorization"
	nickname := strings.TrimPrefix(auth, "Bearer ") // devo controllare che effettivamente questo utente è nel gruppo
	idUser, err := database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Auth not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	str := ps.ByName("groupId")
	id_group, err := strconv.Atoi(str)
	if err != nil {
		http.Error(w, "Invalid groupId parameter", http.StatusBadRequest)
		return
	}

	type Member struct {
		Username string `json:"user"`
	}
	var input Member
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil { // scrive nel JSON della risposta
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newMem, err := database.AppDatabase.IdFromNick(rt.db, input.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound) // 404
			return
		}
		http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
		return
	}

	// prima però dovremmo controllare che l'autore sia davvero nel gruppo
	check, err := database.AppDatabase.HaveRights(rt.db, idUser, id_group)
	if err != nil {
		http.Error(w, "Finding error", http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "Forbidden: You are not a group member", http.StatusForbidden)
		return
	}

	err = database.AppDatabase.AddMember(rt.db, id_group, newMem)
	if err != nil {
		if err.Error() == "member already exists" {
			w.WriteHeader(http.StatusConflict) // 409
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "L'utente è già presente nel gruppo.",
			}); err != nil {
				log.Println("Errore nell'encode JSON:", err)
			}
			return
		}
		log.Println("Errore sconosciuto AddMember:", err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	w.WriteHeader(http.StatusCreated) // 201

}

func (rt *_router) changeName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	idStr := ps.ByName("groupId")
	gid, err := strconv.Atoi(idStr) // prendo ID dall'URL
	if err != nil {
		http.Error(w, "ID not valid", http.StatusBadRequest)
		return
	}

	type Name struct {
		Nome string `json:"nome"`
	}

	var nameG Name
	if err := json.NewDecoder(r.Body).Decode(&nameG); err != nil { // lettura del body
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = database.AppDatabase.PutNameGroup(rt.db, nameG.Nome, gid)
	if err != nil {
		http.Error(w, "Error in database UPDATE: ", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) changeImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	idStr := ps.ByName("groupId")
	gid, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID not valid", http.StatusBadRequest)
		return
	}
	/* dato che l'immagine viene data tramite json, non posso passarla come []byte ma come stringa
	che verrà poi convertita in formato base 64 */
	type Photo struct {
		Photo string `json:"photo"`
	}
	var img Photo
	if err := json.NewDecoder(r.Body).Decode(&img); err != nil { // lettura del body
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// la foto avrà formato base64
	decodedPhoto, err := base64.StdEncoding.DecodeString(img.Photo)
	if err != nil {
		http.Error(w, "Incorrect format image", http.StatusBadRequest)
		return
	}

	err = database.AppDatabase.PutIconGroup(rt.db, decodedPhoto, gid)
	if err != nil {
		http.Error(w, "Error in database UPDATE", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
