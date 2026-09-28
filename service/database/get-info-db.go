package database

func (db *appdbimpl) GetUsername(id int) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM user WHERE id= ?", id).Scan(&username) // scan permette di restituire il nick
	return username, err
}

func (db *appdbimpl) GetPropic(id int) ([]byte, error) {
	var propic []byte
	err := db.c.QueryRow("SELECT propic FROM user WHERE id= ?", id).Scan(&propic) // scan salva il valore nella var
	return propic, err
}
