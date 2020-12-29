package main

import (
	"database/sql"
)

// IXC contém informações sobre uma instalação do IXC.
type IXC struct {
	db              *sql.DB
	stmtAddrByLogin *sql.Stmt
}

// OpenDB abre uma conexão com o MySQL do IXC.
func (ixc *IXC) OpenDB(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	err = db.Ping()
	if err != nil {
		return err
	}

	ixc.stmtAddrByLogin, err = db.Prepare(`
		SELECT framedipaddress
		FROM radacct
		WHERE username = ?
		AND acctstoptime IS NULL
		ORDER BY acctupdatetime DESC
		LIMIT 1`)
	if err != nil {
		return err
	}

	if ixc.db != nil {
		ixc.db.Close()
	}
	ixc.db = db
	db = nil

	return nil
}
