package api

/* creiamo un'unica conversazione, distinguendo però nel database in base alle tipologie
+ chiaro
+ ordinato
+ flessibile */

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strings"
	"wasatext-1887184/service/database"
)

func (rt *_router) addConvo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	auth := r.Header.Get("Authorization")
	nickname := strings.TrimPrefix(auth, "Bearer ")
	idUser, err := database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found in post conv", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
		return
	}

	// tutto questo per prendere l'ID da un nickname
	/*	gli input saranno
		- tipo convo
		- tutti i dati della relativa conv */
	type Convo struct {
		Creator int       `json:"creator"`
		Member  *[]string `json:"members"` // posso decidere alcuni partecipanti
		Name    string    `json:"name"`    // se privata = nick altro utente
		Tipo    string    `json:"tipo"`
	}
	var conversation Convo
	conversation.Creator = idUser
	if err := json.NewDecoder(r.Body).Decode(&conversation); err != nil {
		http.Error(w, "Bad request: create convo", http.StatusBadRequest)
		return
	}

	var memId []int // qui converto tutti i nick dei membri in ID

	if conversation.Member != nil { // potrei star creando una chat
		for _, nick := range *conversation.Member {
			id, err := database.AppDatabase.IdFromNick(rt.db, nick)
			if err != nil {
				// prendo i dati dall'input e restituisco un JSON esplicativo
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)

				if err := json.NewEncoder(w).Encode(map[string]string{
					"error": "L'utente '" + nick + "' non esiste o non è valido.",
				}); err != nil {
					log.Println("Errore nell'encode JSON:", err)
				}
				return
			}
			memId = append(memId, id)
		}
	}
	// tutte le righe qua sotto servono a togliere la duplicità del creatore tra i membri
	unici := make(map[int]struct{})
	for _, id := range memId {
		if id != conversation.Creator {
			unici[id] = struct{}{}
		}
	}
	// con struct{}{} creo un set leggero con il creatore una volta sola
	unici[conversation.Creator] = struct{}{}

	var unique []int
	for id := range unici {
		unique = append(unique, id)
	}

	var convoId int // controllo sia di tipo chat
	if strings.EqualFold(conversation.Tipo, "chat") {
		idChat, err := database.AppDatabase.PostChat(rt.db, conversation.Creator, conversation.Name) // time=ora
		if err != nil {
			log.Printf("Errore nella creazione della chat: %v", err)
			// l'utente non esiste se non riesco a creare le chat
			w.Header().Set("Content-Type", "application/json") // userò un json per dare la risposta
			w.WriteHeader(http.StatusNotFound)                 // 404

			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "L'utente richiesto non esiste.",
			}); err != nil {
				log.Println("Errore nell'encode JSON:", err)

			}
			return
		}
		convoId = idChat
		log.Printf("Chat di %s creata: %d", nickname, idChat)

	} else if strings.EqualFold(conversation.Tipo, "gruppo") {
		idGroup, err := database.AppDatabase.PostGroup(rt.db, conversation.Creator, conversation.Name, unique)
		if err != nil {
			log.Printf("Errore nella creazione del gruppo: %v", err)
			http.Error(w, "Error in create group", http.StatusInternalServerError)
			return
		}
		convoId = idGroup
		log.Printf("Gruppo di %s creato: %d", nickname, idGroup) // così restituisco l'ID, con più dati meglio encode
	} else {
		http.Error(w, "Conversation type doesn't exist ", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	resp := map[string]int{"id": convoId}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
	/* err = database.AppDatabase.PrintAllRows(rt.db, "conversations")
	if err != nil {
		log.Println("Error in database READ:", err)
	}
	*/
}
