package database

import "log"

func (db *appdbimpl) IdFromNick(nickname string) (int, error) {
	var userID int
	err := db.c.QueryRow("SELECT id FROM user WHERE username = ?", nickname).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// funzione per controllare esiste una chat privata
func (db *appdbimpl) ExistChat(sender int, receiver int) (int, error) {
	var convo int
	err := db.c.QueryRow("SELECT convo_id FROM chat WHERE sender = ? AND receiver = ?", sender, receiver).Scan(&convo)
	if err != nil {
		return 0, err
	}
	return convo, err
}

func (db *appdbimpl) AllUsers() ([]string, error) {
	var utenti []string
	rows, err := db.c.Query("SELECT username FROM user")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		utenti = append(utenti, username)
	}
	return utenti, nil
}
