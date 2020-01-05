package main

import "time"

// Trx represents transaction row.
type Trx struct {
	ID       int64     `db:"id" json:"id"`
	TrxType  string    `db:"trx_type" json:"trx_type"`
	Category string    `db:"category" json:"category"`
	Value    float64   `db:"value" json:"value"`
	TrxDate  time.Time `db:"trx_date" json:"trx_date"`
}
