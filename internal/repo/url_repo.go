package repo

import (
	"github.com/jmoiron/sqlx"
	"urls/pkg/database"
)

type Url struct {
	id    int    `db:"id"`
	long  string `db:"long"`
	short string `db:"short"`
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
	CreateUrl(hash, long string)
	GetLastId() int
}

type MysqlUrlRepo struct {
	conn *sqlx.DB
}

func NewMysqlUrlRepo() *MysqlUrlRepo {
	return &MysqlUrlRepo{
		conn: database.GetConnection(),
	}
}

func (r *MysqlUrlRepo) CreateUrl(hash, long string) {
	r.conn.MustExec("INSERT INTO urls (`hash`, `long`) VALUES (?, ?)", hash, long)
}

func (r *MysqlUrlRepo) GetByFull(url string) (u *Url) {
	err := r.conn.Select(u, "SELECT * FROM urls WHERE long = ?", url)
	if err != nil {
		return nil
	}

	return u
}

func (r *MysqlUrlRepo) GetByShort(url string) (u *Url) {
	err := r.conn.Select(u, "SELECT * FROM urls WHERE `short` = ?", url)
	if err != nil {
		return nil
	}

	return u
}

func (r *MysqlUrlRepo) GetLastId() int {
	var url Url
	err := r.conn.Select(&url, "SELECT id FROM urls ORDER BY id DESC LIMIT 1")
	if err != nil {
		return -1
	}

	return url.id
}
