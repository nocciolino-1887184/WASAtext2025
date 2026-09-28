package database

import (
	"database/sql"
	"log"
	"strings"
)

type Conv struct {
	ConvoId   int     `json:"id"`
	Ricevente *string `json:"ricevente"` // per chat private
	Nome      *string `json:"nome"`      // per chat di gruppo
	Categ     string  `json:"categoria"`
	Icon      []byte  `json:"icona"`
	Type      string  `json:"tipo_mess"`
	Text      string  `json:"text"`
	Lettura   bool    `json:"lettura"`
	Tempo     string  `json:"tempo"`
}

type Reactions struct {
	Emoji   string
	User    int
	Message int
}

type Mess struct {
	MessId     int         `json:"id"`
	Sender     int         `json:"sender"`
	Content    *[]byte     `json:"content"`
	Text       *string     `json:"text"`
	Media_type string      `json:"tipo_file"`
	Timestamp  string      `json:"tempo"`
	Status     *string     `json:"letto"` // quando un membro lo ha letto
	Reacts     []Reactions `json:"reacts"`
	Forwarded  bool        `json:"forwarded"`
	Reply      int         `json:"reply"`
	// il messaggio non si salva le emoji, c'è a parte una tabella
}

func (db *appdbimpl) GetAllConvos(id int) ([]Conv, error) {

	query := // ci sono modifiche dovute al fatto che esistono chat senza messaggi
		`SELECT c.id, u.username, NULL AS nome, c.convoType, u.propic,
				COALESCE(m.media_type, 'none'), COALESCE(m.text, '') AS text, COALESCE(m.timestamp,ch.created) AS timestamp,
				COALESCE( (m.sender != ? AND m.letto IS NULL),0 ) AS lettura
		FROM conversations c
		JOIN chat ch ON (ch.convo_id = c.id)
		JOIN user u ON u.id = 
			CASE 
				WHEN ch.sender = ? THEN ch.receiver
				ELSE ch.sender
			END
		LEFT JOIN messages m ON m.mess_id = (
			SELECT mess_id
			FROM messages 
			WHERE convo_id = c.id 
			ORDER BY timestamp DESC, mess_id DESC
			LIMIT 1
		)
		WHERE (ch.sender = ? OR ch.receiver = ?) AND c.convoType = 'chat'
		
	UNION

		SELECT c.id, NULL AS ricevente, g.name, c.convoType, g.icon,
			COALESCE(m.media_type, 'none'),COALESCE(m.text, '') AS text, COALESCE(m.timestamp,g.created) AS timestamp,
			COALESCE( (m.sender != ? AND (mb.last_read IS NULL OR m.timestamp > mb.last_read)),0 ) AS lettura
		FROM conversations c
		JOIN members mb ON mb.group_id = c.id
		JOIN groups g ON g.convo_id = c.id
		LEFT JOIN messages m ON m.mess_id= (
			SELECT mess_id FROM messages 
			WHERE convo_id = c.id 
			ORDER BY timestamp DESC, mess_id DESC
			LIMIT 1
		)
		WHERE mb.member = ? AND c.convoType = 'gruppo'
		
		ORDER BY timestamp DESC`

// le chat qui sono mostrate in ordine cronologico. COALESCE evita errori nel caso di campi null (fornisce alternativa) 
		// case when mostra a chi stai scrivendo
	row, err := db.c.Query(query, id, id, id, id, id, id) // corretto per poter vedere la stessa chat da tutti i membri in forntend
	if err != nil {
		return nil, err
	}
	// per una chiusura più sicura
	defer func() {
		if err := row.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var results []Conv // qui raccoglierò tutte le soluzioni da restituire poi
	for row.Next() {
		var ct Conv
		if err := row.Scan(&ct.ConvoId, &ct.Ricevente, &ct.Nome, &ct.Categ,
				&ct.Icon, &ct.Type, &ct.Text, &ct.Tempo, &ct.Lettura); err != nil { return nil, err	}
		results = append(results, ct)
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *appdbimpl) GetMyConv(id int, convo_id int) ([]Mess, error) {
	// devo mostrare tutti i messaggi di quella chat/gruppo

	// devo capire di che tipo è la chat, quindi
	var convType string

	err := db.c.QueryRow(`SELECT convoType FROM conversations WHERE id=?`, convo_id).Scan(&convType)
	if err != nil {
		return nil, err
	}
	var last_user = ""
	if convType == "gruppo" {
		var minimo sql.NullString
		err := db.c.QueryRow(`
            SELECT MIN(last_read) 
            FROM members 
            WHERE group_id = ? AND member != ?`, convo_id, id).Scan(&minimo)

		if err == nil && minimo.Valid {
			last_user = minimo.String
		}
	}

	query :=
		`SELECT m.mess_id, m.sender, m.content, m.text, m.media_type, m.timestamp, m.letto, m.forwarded, COALESCE(m.reply, 0) as reply
	FROM messages m
	WHERE m.convo_id = ?
	ORDER BY m.timestamp ASC `

	rows, err := db.c.Query(query, convo_id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var output []Mess

	for rows.Next() {
		var msg Mess
		var lettoDB *string
		if err := rows.Scan(&msg.MessId,
			&msg.Sender,
			&msg.Content,
			&msg.Text,
			&msg.Media_type,
			&msg.Timestamp,
			&lettoDB,
			&msg.Forwarded,
			&msg.Reply); err != nil {
			log.Println("ERRORE SCAN messaggio:", err)
			return nil, err // gestione errore
		}
		if convType == "chat" {
			// Per le chat singole, ci fidiamo del campo 'letto' del messaggio
			msg.Status = lettoDB
		} else { // GRUPPI
			if last_user != "" {
				// puliamo il timestamp
				giustoTmp := strings.Replace(msg.Timestamp, "T", " ", 1)
				giustoTmp = strings.Replace(giustoTmp, "Z", "", 1)

				if len(giustoTmp) > 19 {
					giustoTmp = giustoTmp[:19]
				}		// la formattazione del tempo è lunga 19 CHAR
		// se qualcuno ha aperto la chat DOPO la ricezione dell'ultimo messaggio, =letto
				if giustoTmp <= last_user {
					letta := last_user
					msg.Status = &letta
				} else {
					msg.Status = nil
				}
			} else {
				// se last_user è vuoto, nessuno ha letto niente
				msg.Status = nil
			}
		}
		output = append(output, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}
