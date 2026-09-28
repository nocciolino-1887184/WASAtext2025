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

// tutto per avere una icona nel giusto formato
type Ico64 struct {
	ConvoId   int     `json:"id"`
	Ricevente *string `json:"ricevitore"`
	Nome      *string `json:"nome"`
	Categ     string  `json:"categoria"`
	Icona     string  `json:"icona"`
	Tipo      string  `json:"tipo_mess"`
	Mess      string  `json:"mess"`
	Lettura   bool    `json:"lettura"` // per poter evidenziare un nuovo mess
	Tempo     string  `json:"tempo"`   // semplice tipo per salvare le informazioni
}

type ReactionsJS struct {
	Emoji   string `json:"emoji"`
	User    int    `json:"user"`
	Message int    `json:"message"`
}

// per una migliore mappatura del JSON, riscriviamo i messaggi
type MessJSon struct {
	Id         int           `json:"id"`
	Sender     int           `json:"sender"`
	Content    *string       `json:"content"` // può non esserci
	Text       *string       `json:"text"`
	Media_type string        `json:"tipo_file"`
	Timestamp  string        `json:"tempo"`
	Status     *string       `json:"letto"`  // no oppure timestamp
	Reacts     []ReactionsJS `json:"reacts"` // qui va bene usarle
	Forwarded  bool          `json:"forwarded"`
	Reply      int           `json:"reply"`
}

func (rt *_router) showMyConvos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	userId := ps.ByName("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// voglio tutte le conversazioni di questo utente

	res, err := database.AppDatabase.GetAllConvos(rt.db, id)
	if err != nil {
		log.Printf("Errore da GetAllConvos: %v", err)
		http.Error(w, "Failed to get conversations", http.StatusInternalServerError)
		return
	} // sempre controllare eventuali errori

	// procedura per convertire l'icona restituita in base 64
	var responses []Ico64
	for _, c := range res {
		iconBase64 := ""
		if len(c.Icon) > 0 {
			iconBase64 = base64.StdEncoding.EncodeToString(c.Icon)
		}

		responses = append(responses, Ico64{
			ConvoId:   c.ConvoId,
			Ricevente: c.Ricevente,
			Nome:      c.Nome,
			Categ:     c.Categ,
			Icona:     iconBase64,
			Tipo:      c.Type,
			Mess:      c.Text,
			Lettura:   c.Lettura,
			Tempo:     c.Tempo, // ci vuole congruenza tra le due struct
		})
	}
	
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Bad request: show all convos", http.StatusBadRequest)
		return
	}
}

func (rt *_router) showAConv(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	// ID CONVERSAZIONE nel path
	chatId := ps.ByName("convoId")
	cid, err := strconv.Atoi(chatId)
	if err != nil {
		http.Error(w, "Invalid ID from parameters", http.StatusBadRequest)
		return
	}
	// nel body avrò l'ID utente
	auth := r.Header.Get("Authorization") // nell'header devo aggiungere "authorization"
	nickname := strings.TrimPrefix(auth, "Bearer ")

	userId, err := database.AppDatabase.IdFromNick(rt.db, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found in show 1 conv", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error id-from-nick", http.StatusInternalServerError)
		return
	}
	// se apro una chat, visualizzo messaggi
	err = database.AppDatabase.ReadMessage(rt.db, cid, userId)
	if err != nil {
		log.Println("Errore aggiornamento lettura:", err)
	}

	res, err := database.AppDatabase.GetMyConv(rt.db, userId, cid)
	if err != nil {
		http.Error(w, "Failed to get the conversation", http.StatusInternalServerError)
		return
	}

	var totalmes []MessJSon
	for _, m := range res {
		var base64Foto *string
		if m.Content != nil && len(*m.Content) > 0 {
			encoded := base64.StdEncoding.EncodeToString(*m.Content)
			base64Foto = &encoded
		} else {
			base64Foto = nil
		}
		// bisogna occuparsi anche delle reaction
		var reactions []ReactionsJS
		for _, r := range m.Reacts {
			reactions = append(reactions, ReactionsJS{
				Emoji:   r.Emoji,
				User:    r.User,
				Message: r.Message,
			})
		}

		totalmes = append(totalmes, MessJSon{
			Id:         m.MessId,
			Sender:     m.Sender,
			Content:    base64Foto,
			Text:       m.Text,
			Media_type: m.Media_type,
			Timestamp:  m.Timestamp,
			Status:     m.Status,
			Reacts:     reactions,
			Forwarded:  m.Forwarded,
			Reply:      m.Reply,
		})
	}
	if err := json.NewEncoder(w).Encode(totalmes); err != nil {
		http.Error(w, "Bad request: show my convo", http.StatusBadRequest)
		return
	}
}
