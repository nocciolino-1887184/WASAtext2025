package database

import (
	"log"
)

/* prima creo una nuova convo di tipo chat
poi trovo l'ID del receiver che sarà anche il nome della chat
infine aggiungo una nuova chat alla tabella */

func (db *appdbimpl) PostChat(sender int, receiver string) (int, error) {
	risp, err := db.c.Exec(`INSERT INTO conversations(convoType) VALUES ('chat')`)
	if err != nil {
		log.Println("Errore INSERT chat:", err) // indica il tipo di errore
		return 0, err
	}
	idChat, err := risp.LastInsertId()
	if err != nil {
		return 0, err
	}

	var idRec int
	err = db.c.QueryRow(`SELECT id FROM user WHERE username = ?`, receiver).Scan(&idRec)
	if err != nil {
		return 0, err
	}
	if sender > idRec {
		sender, idRec = idRec, sender
	}

	// il problema sono i doppioni

	_, err = db.c.Exec(`INSERT INTO chat (convo_id, sender, receiver) VALUES (?,?,?) `, idChat, sender, idRec) // timestamp = now
	if err != nil {
		return 0, err
	} // l'id è quello della conversazione

	return int(idChat), nil // potrei usare int64 per numeri molto più grandi
}

func (db *appdbimpl) PostGroup(sender int, name string, members []int) (int, error) {
	risp, err := db.c.Exec(`INSERT INTO conversations(convoType) VALUES ('gruppo')`)
	if err != nil {
		return 0, err
	}
	idGroup, err := risp.LastInsertId()
	if err != nil {
		return 0, err
	}

	mem := len(members)
	_, err = db.c.Exec(`INSERT INTO groups (convo_id, name, members) VALUES (?,?,?) `, idGroup, name, mem) // timestamp = now
	if err != nil {
		return 0, err
	}

	for _, userID := range members {
		_, err := db.c.Exec(`INSERT INTO members (group_id, member) VALUES (?, ?)`, idGroup, userID)
		if err != nil {
			return 0, err
		}
	}
	return int(idGroup), nil
}
