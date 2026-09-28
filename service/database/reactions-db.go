package database

import "log"

type Reaction struct {
	Emoji string `json:"emoji"`
	User  string `json:"user"`
}

func (db *appdbimpl) AddComment(emoji string, messId int, id int) error {

	query := `INSERT INTO reactions (id_emoji, message, username)
			VALUES ( ?, ?, ? );`
	_, err := db.c.Exec(query, emoji, messId, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetComment(messId int) ([]Reaction, error) {
	query := `SELECT id_emoji, user.username
			FROM reactions
			JOIN user ON user.id = reactions.username
			WHERE message = ?`
	rows, err := db.c.Query(query, messId)
	if err != nil {
		return nil, err
	}
	// voglio direttamente il nome utente
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var reactions []Reaction

	for rows.Next() {
		var emoji string
		var user string
		if err := rows.Scan(&emoji, &user); err != nil {
			return nil, err
		}
		reactions = append(reactions, Reaction{
			Emoji: emoji,
			User:  user,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

func (db *appdbimpl) UnComment(emoji string, messId int, nick string) (bool, error) {
	risp, err := db.c.Exec(`DELETE FROM reactions
		WHERE id_emoji = ? AND message = ? AND username = (SELECT id FROM user WHERE username = ?);`, emoji, messId, nick)
	if err != nil {
		return false, err
	}
	// questo controllo servirà poi per restituire 204 oppure 404
	cambiamenti, err := risp.RowsAffected()
	if err != nil {
		return false, err
	}

	if cambiamenti == 0 {
		return false, nil
	}

	return true, nil // tutto ok

}
