package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// FetchCreateTrxTable creates a transaction table.
func FetchCreateTrxTable(db *sqlx.DB) error {
	sqlStr := `
		CREATE TABLE IF NOT EXISTS trx_trx (
			id SERIAL PRIMARY KEY,
			trx_type VARCHAR(64) NOT NULL,
			category VARCHAR(64) NOT NULL, 
			value DOUBLE PRECISION NOT NULL,
			trx_date DATE NOT NULL
		)
	`

	stmt, err := db.Preparex(sqlStr)
	if err != nil {
		errStr := fmt.Errorf("CreateTrxTable - Preparex, err: %s", err)
		return errStr
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		errStr := fmt.Errorf("CreateTrxTable - Exec, err: %s", err)
		return errStr
	}

	return nil
}

// FetchTrxSave saves transaction passed.
func FetchTrxSave(db *sqlx.DB, trx *Trx) error {
	sqlStr := `
		INSERT INTO trx_trx 
			(trx_type, category, value, trx_date)
		VALUES
			($1, $2, $3, $4)
	`

	stmt, err := db.Preparex(sqlStr)
	if err != nil {
		errStr := fmt.Errorf("FetchTrxSave - Preparex, err: %s", err)
		return errStr
	}
	defer stmt.Close()

	_, err = stmt.Exec(trx.TrxType, trx.Category, trx.Value, trx.TrxDate)
	if err != nil {
		errStr := fmt.Errorf("FetchTrxSave - Exec, err: %s", err)
		return errStr
	}

	return nil
}

// FetchTrxs returns all transactions.
func FetchTrxs(db *sqlx.DB) ([]*Trx, error) {
	sqlStr := `
		SELECT
			t.id,
			t.trx_type,
			t.category,
			t.value,
			t.trx_date
		FROM
			trx_trx t
	`

	stmt, err := db.Preparex(sqlStr)
	if err != nil {
		errStr := fmt.Errorf("FetchTrxs - Preparex, err: %s", err)
		return nil, errStr
	}
	defer stmt.Close()

	res, err := stmt.Queryx()
	if err != nil {
		errStr := fmt.Errorf("FetchTrxs - Queryx, err: %s", err)
		return nil, errStr
	}

	trxs := []*Trx{}
	for res.Next() {
		t := &Trx{}
		err := res.StructScan(t)
		if err != nil {
			errStr := fmt.Errorf("FetchTrxs - StructScan, err: %s", err)
			return nil, errStr
		}

		trxs = append(trxs, t)
	}

	return trxs, nil
}
