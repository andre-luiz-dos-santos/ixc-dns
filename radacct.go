package main

import (
	"database/sql"
	"errors"
)

// GetAddrByLogin retorna o IP de uma das conexões de login.
func (ixc *IXC) GetAddrByLogin(login string) (string, error) {
	var addr string
	err := ixc.stmtAddrByLogin.QueryRow(login).Scan(&addr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return addr, nil
}
