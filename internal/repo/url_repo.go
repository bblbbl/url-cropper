package repo

import (
	"database/sql"
	"urls/pkg/database"
)

type Url struct {
	long  string
	short string
}

func (u *Url) GetShort() string {
	return u.short
}

func (u *Url) GetLong() string {
	return u.long
}

type UrlRepo interface {
	GetByFull(url string) *Url
	GetByShort(url string) *Url
	CreateUrl(short, long string) error
	GetLastId() int
}

type MysqlUrlRepo struct {
	conn *sql.DB
}

func NewMysqlUrlRepo() *MysqlUrlRepo {
	return &MysqlUrlRepo{
		conn: database.GetConnection(),
	}
}

func (r *MysqlUrlRepo) CreateUrl(long, short string) error {
	_, err := r.conn.Exec("INSERT INTO urls (`long`, short) VALUES (?, ?)", long, short)

	return err
}

func (r *MysqlUrlRepo) GetByFull(url string) *Url {
	var u Url
	err := r.conn.
		QueryRow("SELECT short FROM urls WHERE long = ?", url).
		Scan(&u.short)

	if err != nil {
		return nil
	}

	u.long = url

	return &u
}

func (r *MysqlUrlRepo) GetByShort(url string) *Url {
	var u Url
	err := r.conn.
		QueryRow("SELECT long FROM urls WHERE short = ?", url).
		Scan(&u.long)

	if err != nil {
		return nil
	}

	u.short = url

	return &u
}

func (r *MysqlUrlRepo) GetLastId() int {
	var id int

	err := r.conn.
		QueryRow("SELECT id FROM urls ORDER BY id DESC LIMIT 1").
		Scan(&id)

	if err != nil {
		return -1
	}

	return id
}
