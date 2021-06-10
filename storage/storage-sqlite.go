package storage

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//var tableName = "elements"
var ErrNoSuchElement = errors.New("no element with such id")

type SQLiteStorage struct {
	db *sql.DB
}

func OpenSQLiteStorage(filename string) (Storage, error) {
	return openSQLiteStorage(filename)
}

func openSQLiteStorage(filename string) (Storage, error) {
	dir, _ := filepath.Split(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) && dir != "" {
		if err = os.MkdirAll(dir, 0755); err != nil { // 755 needs to be changed to 0700
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	var sqls SQLiteStorage
	sqls.db = db
	statement := `CREATE TABLE IF NOT EXISTS elements (id INTEGER PRIMARY KEY AUTOINCREMENT, payload TEXT, added DATETIME)`
	if _, err := sqls.db.Exec(statement); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // go-sqlite3 docs recommend this
	return &sqls, nil
}

func (sqls *SQLiteStorage) Get(id ...ID) ([]Element, error) {
	var rows *sql.Rows
	var err error
	switch {
	case len(id) > 0:
		rows, err = sqls.db.Query(`SELECT * FROM elements WHERE id = ?`, id[0])
	default:
		rows, err = sqls.db.Query(`SELECT * FROM elements`)
	}
	elements := []Element{}
	if err != nil {
		return elements, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Element
		if rows.Scan(&e.ID, &e.Payload, &e.Added); err == nil {
			elements = append(elements, e)
		} else {
			log.Println(err)
		}
	}
	if len(elements) == 0 {
		return elements, ErrNoSuchElement
	}
	return elements, nil
}

func (sqls *SQLiteStorage) Put(e Element) (ID, error) {
	if e.Payload == "" {
		return BadID, errors.New("can not add empty Element")
	}
	t := time.Now()
	id := BadID
	statement := `insert into elements values (?, ?, ?)`
	result, err := sqls.db.Exec(statement, nil, e.Payload, t)
	if err != nil {
		return id, err
	}
	n, err := result.LastInsertId()
	if err == nil {
		id = ID(n)
	}
	return id, err
}

func (sqls *SQLiteStorage) Upd(id ID, e Element) error {
	result, err := sqls.db.Exec(`update elements set payload = ? where id = ?`, e.Payload, id)
	if n, err := result.RowsAffected(); err != nil || n == 0 {
		return ErrNoSuchElement
	}
	return err
}

func (sqls *SQLiteStorage) Del(id ID) error {
	result, err := sqls.db.Exec(`delete from elements where id = ?`, id)
	if n, err := result.RowsAffected(); err != nil || n == 0 {
		return ErrNoSuchElement
	}
	return err
}

func (sqls *SQLiteStorage) Close() error {
	return sqls.db.Close()
}
