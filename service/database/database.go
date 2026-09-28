/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):
// Eseguo le query

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	NewUser(name string) (int, error)
	PutNick(nick string, id int) error
	PutPropic(id int, propic []byte) error
	GetUsername(id int) (string, error)
	GetPropic(id int) ([]byte, error)

	PostChat(sender int, receiver string) (int, error)
	PostGroup(sender int, name string, members []int) (int, error)
	SendMessage(
		sender int,
		convo int,
		content *[]byte,
		text *string,
		mt string,
		reply int) (int, error)
	GetAllConvos(id int) ([]Conv, error)
	GetMyConv(id int, chat_id int) ([]Mess, error)
	GetAllGroups(idUser int) ([]Group, error)

	ForwardMessage(mess_id int, convo_id int, sender int) (int, error)
	DeleteMessage(convoId int, messId int) (bool, error)
	ReadMessage(convoId int, myId int) error
	AddComment(emoji string, messId int, id int) error
	GetComment(id_mess int) ([]Reaction, error)
	UnComment(emoji string, messId int, nick string) (bool, error)

	// DeleteGroup(idGroup int) (bool, error)
	AddMember(id_g int, mem int) error
	DeleteMember(idGroup int, id int) (bool, error)
	ShowMembers(gid int) ([]string, error)
	PutNameGroup(name string, gid int) error
	PutIconGroup(icon []byte, gid int) error

	// utili
	IdFromNick(nickname string) (int, error)
	AllUsers() ([]string, error)
	// PrintAllRows(name string) error
	ExistChat(sender int, receiver int) (int, error)
	// ResetConversations() error
	HaveRights(id int, gid int) (bool, error)
	IsAlreadyIn(nick string) (bool, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Necessary if we are using foreign keys
	_, errPragma := db.Exec("PRAGMA foreign_keys = ON")
	if errPragma != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPragma)
	} 
	
	/* _, err := db.Exec("DROP TABLE IF EXISTS user;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE user (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username VARCHAR(16) NOT NULL UNIQUE, propic BLOB);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	/* _, err = db.Exec("DROP TABLE IF EXISTS conversations;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='conversations';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE conversations (
					id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					convoType TEXT CHECK(convoType IN ('chat', 'gruppo')) NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	/* _, err = db.Exec("DROP TABLE IF EXISTS messages;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='messages';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) { // Posso salvare le foto direttamente nel DB
		sqlStmt := `CREATE TABLE messages
			(mess_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			sender INTEGER NOT NULL,
			convo_id INTEGER NOT NULL,
			content BLOB,
			text TEXT,
			media_type TEXT CHECK(media_type IN ('text', 'image')) NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			letto TIMESTAMP NULL DEFAULT NULL,
			forwarded BOOLEAN DEFAULT FALSE,
			reply INTEGER,
			
			FOREIGN KEY (sender) REFERENCES user(id),
			FOREIGN KEY (convo_id) REFERENCES conversations(id) ON DELETE CASCADE);` // controllo se un messaggio è in una CHAT o GRUPPO
		// come mettere chiave composta = primary key (x,y)
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	
	/*_, err = db.Exec("DROP TABLE IF EXISTS chat;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='chat';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// la chat ha chiave composta da 3
		sqlStmt := `CREATE TABLE chat
				(convo_id INTEGER NOT NULL PRIMARY KEY,
				sender INTEGER NOT NULL,
				receiver INTEGER NOT NULL,
				created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				
				UNIQUE(sender, receiver), 
				FOREIGN KEY (convo_id) REFERENCES conversations(id) ON DELETE CASCADE,
				FOREIGN KEY (sender) REFERENCES user(id),		
				FOREIGN KEY (receiver) REFERENCES user(id)
				);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	/*_, err = db.Exec("DROP TABLE IF EXISTS reactions;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='reactions';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) { // esistono a parte e si legano ad un messaggio
		sqlStmt := `CREATE TABLE reactions
			(id_emoji TEXT NOT NULL,	
			username INTEGER NOT NULL,
			message INTEGER NOT NULL, 

			PRIMARY KEY (username, message),
			FOREIGN KEY (username) REFERENCES user(id),
			FOREIGN KEY (message) REFERENCES messages(mess_id) ON DELETE CASCADE );`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	/*_, err = db.Exec("DROP TABLE IF EXISTS groups;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='groups';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE groups
			(convo_id INTEGER NOT NULL PRIMARY KEY,
			name VARCHAR(20) NOT NULL,
			icon BLOB,
			members INTEGER NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (convo_id) REFERENCES conversations(id) ON DELETE CASCADE
			);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	/*_, err = db.Exec("DROP TABLE IF EXISTS members;")
	if err != nil {
	    return nil, fmt.Errorf("error dropping table: %w", err)
	}  */
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='members';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) { // tabella che associa ogni gruppo ai suoi membri
		// possiamo avere una chiave più sicura unendo più campi
		sqlStmt := `CREATE TABLE members
			(group_id INTEGER NOT NULL,
			member INTEGER NOT NULL,
			joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_read TIMESTAMP NULL DEFAULT NULL,
			
			PRIMARY KEY(group_id,member),
			FOREIGN KEY (group_id) REFERENCES groups(convo_id) ON DELETE CASCADE,
			FOREIGN KEY (member) REFERENCES user(id) );` // quando ci si unisce al gruppo
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
