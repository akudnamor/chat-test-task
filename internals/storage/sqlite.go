package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		log.Println(op, "failed to open", err)
		return nil, err
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS messages(
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"from" TEXT NOT NULL,
		"text" TEXT,
		"createdAt" DATE);		
	`)
	if err != nil {
		log.Println(op, "failed to prepare query", err)
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(op, "failed to exec stmt", err)
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (st *Storage) AddMessage(msg Message) (int64, error) {
	const op = "sqlite.AddMessage"

	stmt, err := st.db.Prepare(`INSERT INTO messages ("from", "text", "createdAt") VALUES (?, ?, ?)`)
	if err != nil {
		log.Println(op, "failed to prepare query", err)
		return -1, err
	}

	res, err := stmt.Exec(msg.From, msg.Text, time.Now())
	if err != nil {
		log.Println(op, "failed to exec stmt", err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(op, "failed to get last insert id", err)
		return -1, err
	}

	return id, nil
}

func (st *Storage) GetMessages(limit, skip int) ([]Message, error) {
	const op = "sqlite.GetMessages"

	stmt, err := st.db.Prepare(`SELECT * from messages ORDER BY "createdAt" DESC LIMIT (?) OFFSET (?)`)
	if err != nil {
		log.Println(op, "failed to prepare query", err)
		return nil, err
	}

	row, err := stmt.Query(limit, skip)
	if err != nil {
		log.Println(op, "failed to get rows", err)
		return nil, err
	}
	var messages []Message
	for row.Next() {
		var message Message
		err = row.Scan(&message.ID, &message.From, &message.Text, &message.CreatedAt)
		if err != nil {
			log.Println(op, "failed to scan into struct", err)
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
