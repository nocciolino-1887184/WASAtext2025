package database

func (db *appdbimpl) PutNick(nick string, id int) error {

	_, err := db.c.Exec("UPDATE user SET username = ? WHERE id = ?", nick, id)
	return err
}

func (db *appdbimpl) PutPropic(id int, propic []byte) error {

	_, err := db.c.Exec("UPDATE user SET propic = ? WHERE id = ?", propic, id)
	return err
}
