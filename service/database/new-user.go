package database

func (db *appdbimpl) NewUser(name string) (int, error) {
	// fai controllo se non esiste già questo username
	risp, err := db.c.Exec("INSERT INTO user (username) VALUES (?)", name) // l'id si incrementa da solo
	if err != nil {
		return 0, err
	}
	userId, err := risp.LastInsertId()
	return int(userId), err
}

func (db *appdbimpl) IsAlreadyIn(nick string) (bool, error) {
	// restituisce se l'utente è già nella tabella
	var r bool
	// exists si ferma al primo match
	query := `SELECT EXISTS (SELECT 1 FROM user WHERE username = ?)`
	err := db.c.QueryRow(query, nick).Scan(&r)
	// restituisce un errore solo se non ci sono righe
	if err != nil {
		return false, err
	}
	return r, nil
}
