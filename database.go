package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const CreateAndConnectDbQuery = `
CREATE TABLE IF NOT EXISTS leagues (
	id TEXT PRIMARY KEY,
	name TEXT,
	secret_id TEXT
);
CREATE TABLE IF NOT EXISTS history (
	id INTEGER PRIMARY KEY,
	league_id TEXT,
	winner TEXT,
	loser TEXT,
	draw BOOLEAN,
	date_string TEXT
);
`

type DBLeague struct {
	ID       string
	Name     string
	SecretID string
}

type DBHistory struct {
	ID         int
	LeagueID   string
	Winner     string
	Loser      string
	Draw       bool
	DateString string
}

// Connect to the db, if it is empty then also create its tables
func DBCreateAndConnect(addr string) (*sql.DB, error) {
	// connect to sqlite db
	db, err := sql.Open("sqlite3", addr)
	if err != nil {
		return nil, err
	}

	// create tables if they don't exist
	_, err = db.Exec(CreateAndConnectDbQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DBLookupLeageByID(id string) (DBLeague, error) {
	tx, err := AppDB.Begin()
	if err != nil {
		return DBLeague{}, err
	}
	defer tx.Rollback()
	return dbLookupLeague(id, "", tx)
}

func DBLookupLeagueBySecret(secretId string) (DBLeague, error) {
	tx, err := AppDB.Begin()
	if err != nil {
		return DBLeague{}, err
	}
	defer tx.Rollback()
	return dbLookupLeague("", secretId, tx)
}

func DBLookupLeagueByIDWithHistory(id string) (DBLeague, []DBHistory, error) {
	tx, err := AppDB.Begin()
	if err != nil {
		return DBLeague{}, nil, err
	}
	defer tx.Rollback()
	league, err := dbLookupLeague(id, "", tx)
	if err != nil {
		return DBLeague{}, nil, err
	}
	history, err := dbLookupHistory(id, tx)
	if err != nil {
		return DBLeague{}, nil, err
	}
	return league, history, nil
}

func DBCreateLeage(name, id, secretId string) error {
	tx, err := AppDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = dbCreateLeage(name, id, secretId, tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func DBCreateHistory(secretID string, history []DBHistory) error {
	tx, err := AppDB.Begin()
	if err != nil {
		return err
	}
	league, err := dbLookupLeague("", secretID, tx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, h := range history {
		err = dbCreateHistory(league.ID, h.Winner, h.Loser, h.Draw, h.DateString, tx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

var ErrLeagueNotFound = fmt.Errorf("league not found")

func dbLookupLeague(id, secretId string, tx *sql.Tx) (DBLeague, error) {
	if id == "" && secretId == "" {
		return DBLeague{}, fmt.Errorf("must specify either id or secretId in lookup")
	} else if id != "" && secretId != "" {
		return DBLeague{}, fmt.Errorf("must specify only one of id or secretId in lookup")
	}

	var q, qid string
	if id != "" {
		q = `SELECT id, name, secret_id FROM leagues WHERE id = ?`
		qid = id
	} else {
		q = `SELECT id, name, secret_id FROM leagues WHERE secret_id = ?`
		qid = secretId
	}

	row := tx.QueryRow(q, qid)
	var dbLeague DBLeague
	err := row.Scan(&dbLeague.ID, &dbLeague.Name, &dbLeague.SecretID)
	if err == sql.ErrNoRows {
		return DBLeague{}, ErrLeagueNotFound
	} else if err != nil {
		return DBLeague{}, err
	}
	return dbLeague, nil
}

func dbLookupHistory(id string, tx *sql.Tx) ([]DBHistory, error) {
	q := `SELECT id, league_id, winner, loser, draw, date_string FROM history WHERE league_id = ?`
	rows, err := tx.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []DBHistory
	for rows.Next() {
		var h DBHistory
		err := rows.Scan(&h.ID, &h.LeagueID, &h.Winner, &h.Loser, &h.Draw, &h.DateString)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}

func dbCreateLeage(name, id, secretId string, tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO leagues (id, name, secret_id) VALUES (?, ?, ?)`, id, name, secretId)
	return err
}

func dbCreateHistory(leagueId, winner, loser string, draw bool, dateString string, tx *sql.Tx) error {
	q := `INSERT INTO history (league_id, winner, loser, draw, date_string) VALUES (?, ?, ?, ?, ?)`
	_, err := tx.Exec(q, leagueId, winner, loser, draw, dateString)
	return err
}
