package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"jargon/pkg/api"
)

type Database struct {
	DB *sql.DB
}

func OpenDatabase(dbPath string) Database {
	log.Printf("Opening database at %v", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	return Database{DB: db}
}

func (db Database) CloseDatabase() {
	log.Println("Closing database")
	db.DB.Close()
}

func (db Database) initJargonTags() {
	query := `
	CREATE TABLE IF NOT EXISTS jargontags(
		tag_id UNSIGNED INT NO NULL,
		jargon_id UNSIGNED INT NO NULL,
		PRIMARY KEY(tag_id, jargon_id)
	)
	`

	log.Printf("Creating JargonTags table with query:\n%v\n", query)
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

}

func (db Database) initTags() {
	query := `
	CREATE TABLE IF NOT EXISTS tags(
		tag_id UNSIGNED INT NO NULL PRIMARY KEY,
		tag TEXT NOT NULL
	)
	`

	log.Printf("Creating Tag table with query:\n%v\n", query)
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Populate with tags
	log.Println("Populating tag table")
	query = `
	INSERT INTO Tags(tag_id, tag)
	VALUES(?,?)
	`
	stmt, err = db.DB.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for id, label := range api.TagLabels {

		fmt.Printf("add id:%v label:%v\n", id, label)
		_, err = stmt.Exec(id, label)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func (db Database) InitialiseTables() {

	stmt := `
	CREATE TABLE IF NOT EXISTS jargon(
		label_id UNSIGNED INT NOT NULL PRIMARY KEY,
		label TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT
	)
	`

	log.Printf("Creating table with statement:\n%v\n", stmt)
	_, err := db.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	db.initTags()
	db.initJargonTags()
}

func (db Database) AddJargon(j api.Jargon) error {
	ctx := context.Background()
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
	INSERT INTO jargon(label_id, label, title, description)
	VALUES(?, ?, ?, ?)
	`)

	fmt.Printf("Adding to Jargon to table: %v\n", j)
	result, err := stmt.Exec(j.LabelId, j.Label, j.Title, j.Description)
	if result == nil {
		fmt.Printf("label_id=%v already exists\n", j.LabelId)
		return nil
	}

	stmt, err = tx.Prepare(`
	INSERT INTO jargontags(tag_id, jargon_id)
	VALUES(?, ?)
	`)
	if err != nil {
		log.Fatal(err)
	}

	if j.Tags != nil {
		for _, tag := range j.Tags {
			_, err = stmt.Exec(tag, j.LabelId)
		}
		if err != nil {
			log.Fatal(err)
		}

	}

	tx.Commit()
	return nil
}

func (db Database) GetJargon(label_id uint32) (*api.Jargon, error) {
	stmt, _ := db.DB.Prepare(`
	SELECT label, title, description FROM jargon
	WHERE label_id == ?
	`)

	var (
		label       string
		title       string
		description string
	)

	if err := stmt.QueryRow(label_id).Scan(&label, &title, &description); err != nil {
		return nil, err
	}

	return &api.Jargon{label_id, label, title, description, nil}, nil
}

func (db Database) GetAllJargons() ([]api.Jargon, error) {
	query := `
	SELECT label_id, label, title, description FROM jargon	
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	jargons := make([]api.Jargon, 0)
	for rows.Next() {
		var (
			id          uint32
			label       string
			title       string
			description string
		)

		if err := rows.Scan(&id, &label, &title, &description); err != nil {
			return nil, err
		}

		jargon := api.MakeJargon(label, title, description, nil)
		if jargon.LabelId != id {
			log.Fatal("Label id issues")
		}

		jargons = append(jargons, jargon)

	}

	return jargons, nil
}