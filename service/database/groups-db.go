package database

import (
	"errors"
	"log"
	"time"
)

type Group struct {
	Name     string
	Image    []byte
	LastMess string
	When     time.Time
}

func (db *appdbimpl) GetAllGroups(idUser int) ([]Group, error) {
	/* di ogni gruppo voglio Nome, icona, ultimo mess.
	quindi mi serve un join tra membri per dire che ne faccio parte,
	gruppo e messaggi per avere l'ultimo mandato */
	query := `SELECT g.name, g.icon, m.text, m.timestamp
			FROM groups g
			JOIN members mb ON g.convo_id = mb.group_id
			LEFT JOIN messages m ON m.id(
				SELECT id 
				FROM messages 
				WHERE convo_id = g.convo_id 
				ORDER BY timestamp DESC 
				LIMIT 1
			)
			WHERE mb.member = ?; ` // così prendiamo ultimo messaggio scritto

	risp, err := db.c.Query(query, idUser)
	if err != nil {
		return nil, err
	}
	// chiusura della connessione in modo sicuro 
	defer func() {
		if err := risp.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	// ora devo restituire una LISTA di gruppi, quindi
	var results []Group // qui raccoglierò tutte le soluzioni da restituire poi
	for risp.Next() {
		var g Group
		if err := risp.Scan(&g.Name, &g.Image, &g.LastMess, &g.When); err != nil {
			return nil, err
		}
		results = append(results, g)
	}

	if err := risp.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *appdbimpl) DeleteMember(idGroup int, id int) (bool, error) {
	query := `DELETE FROM members WHERE group_id = ? AND member = ?`
	result, err := db.c.Exec(query, idGroup, id)
	if err != nil {
		return false, err
	}
	// questo controllo servirà poi per restituire 204 oppure 404
	cambiamenti, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if cambiamenti == 0 {
		return false, nil
	}
	return true, nil
}

// fai check per vedere se l'user ha diritti e se c'è già un certo ID nella tabella members
func (db *appdbimpl) AddMember(id_g int, mem int) error {
	query := `INSERT OR IGNORE INTO members(group_id, member, joined) VALUES (?,?,CURRENT_TIMESTAMP)`
	res, err := db.c.Exec(query, id_g, mem)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Se le righe sono 0, significa che non è stato inserito nessun utente (c'era già o non esiste)
	if rows == 0 {
		return errors.New("member already exists")
	}
	return nil
}

// visualizza membri di un gruppo
func (db *appdbimpl) ShowMembers(gid int) ([]string, error) {
	var members []string

	query := `SELECT u.username 
              FROM user u
              JOIN members m ON u.id = m.member
              WHERE m.group_id = ?`

	rows, err := db.c.Query(query, gid)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		members = append(members, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func (db *appdbimpl) HaveRights(id int, gid int) (bool, error) {
	// restituisce se l'utente è nel gruppo o meno
	var r int
	query := `SELECT 1 FROM members WHERE member = ? AND group_id = ? LIMIT 1`
	risp := db.c.QueryRow(query, id, gid)
	err := risp.Scan(&r) // restituisce un errore solo se non ci sono righe
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *appdbimpl) PutNameGroup(name string, gid int) error {
	_, err := db.c.Exec("UPDATE groups SET name = ? WHERE convo_id = ?", name, gid)
	return err
}

func (db *appdbimpl) PutIconGroup(icon []byte, gid int) error {
	_, err := db.c.Exec("UPDATE groups SET icon = ? WHERE convo_id = ?", icon, gid)
	return err
}
