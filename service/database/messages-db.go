package database

import (
	"database/sql"
	"log"
)

// funzione che permette di salvare i messaggi mandati nel database, prendendo in input campi mess
func (db *appdbimpl) SendMessage(
	sender int,
	convo int,
	content *[]byte,
	text *string,
	mt string,
	reply int) (int, error) { // il read è già false e il tempo attuale
	// un messaggio all'inizio non ha emoji
	// il client fornisce dall'header il suo nickname e dal database noi prendiamo l'ID

	query := `
	INSERT INTO messages (sender, convo_id, content, text, media_type, letto, forwarded, reply)
	VALUES (?, ?, ?, ?, ?,NULL,false,?) `

	// se reply=0 significa che non sto rispondendo a nessuno
	var replyValue sql.NullInt64
	if reply != 0 {
		replyValue.Valid = true
		replyValue.Int64 = int64(reply)
	} else {
		replyValue.Valid = false
	}

	res, err := db.c.Exec(query, sender, convo, content, text, mt, replyValue) // EXEC usato se non vogliamo risultati per righe
	if err != nil {
		log.Println("INSERT ERROR:", err)

		return 0, err // il tempo viene salvato direttamente come local time
	}
	// eventuali errori sono gestiti dal handler
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *appdbimpl) ForwardMessage(mess_id int, convo_id int, sender int) (int, error) {
	// per inoltrare un messaggio servono tutte le info del messaggio originale + dove andrà e chi lo inoltra
	query := `SELECT content, text, media_type 
	          FROM messages
			  WHERE mess_id = ? `
	// uso variabili per il contenuto del messaggio
	var contenuto *[]byte
	var testo *string
	var mediaType string
	// non prenderò convoid, che verrà preso invece in input
	err := db.c.QueryRow(query, mess_id).Scan(&contenuto, &testo, &mediaType)
	if err != nil {
		return 0, err
	}

	insertQuery := `INSERT INTO messages (sender, convo_id, content, text, media_type, letto, forwarded)
	 VALUES (?, ?, ?, ?, ?, NULL, true)`

	res, err := db.c.Exec(insertQuery, sender, convo_id, contenuto, testo, mediaType)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *appdbimpl) DeleteMessage(convoId int, messId int) (bool, error) {
	query := `DELETE FROM messages WHERE mess_id = ? AND convo_id = ? `
	res, err := db.c.Exec(query, messId, convoId)
	if err != nil {
		return false, err
	}
	// questo controllo servirà poi per restituire 204 oppure 404
	cambiamenti, err := res.RowsAffected()
	if err != nil { // errore nell'operazione
		return false, err
	}

	if cambiamenti == 0 { // Nessun messaggio cancellato
		return false, nil
	}

	return true, nil // tutto ok
}

// qui intervengo nel database e modifico effettivamente le ultime visualizzazioni
func (db *appdbimpl) ReadMessage(convoId int, myId int) error {
	// se i messaggi sono in una chat, avrò un solo read, altrimenti in un gruppo controllo tutti i membri
	var conv string // obbligatorio per non avere errori di tipo
	query := `SELECT convoType FROM conversations WHERE id = ?`
	err := db.c.QueryRow(query, convoId).Scan(&conv)
	if err != nil {
		return err
	}
	if conv == "chat" {
		query := `UPDATE messages
				  SET letto = CURRENT_TIMESTAMP 
				  WHERE convo_id = ? 
				  AND sender != ? 
				  AND letto IS NULL`
		_, err := db.c.Exec(query, convoId, myId)
		if err != nil {
			return err
		}
	} else if conv == "gruppo" {

		query := `UPDATE members 
            	  SET last_read = CURRENT_TIMESTAMP 
				  WHERE group_id = ? 
				  AND member = ?`
		_, err := db.c.Exec(query, convoId, myId)
		if err != nil {
			return err
		}
		queryUpdateMessages := `
			UPDATE messages 
			SET letto = CURRENT_TIMESTAMP
			WHERE convo_id = ?        
			  AND letto IS NULL       
			  AND NOT EXISTS (          
				SELECT 1
				FROM members mem
				WHERE mem.group_id = messages.convo_id
				  AND mem.member != messages.sender       
				  AND (mem.last_read < messages.timestamp OR mem.last_read IS NULL)  )`
		// abbiamo escluso chi ha mandato il mess e controlliamo se qualcuno non lo ha ancora letto
		// ci accertiamo che non esista nessuno escluso dalla lettura (Not exist = true)

		_, err = db.c.Exec(queryUpdateMessages, convoId)
		if err != nil {
			return err
		}

	}
	return nil
}
