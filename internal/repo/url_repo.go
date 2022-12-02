package repo

import (
	"github.com/jmoiron/sqlx"
	"urls/pkg/database"
)

type Url struct {
	Id   uint8  `db:"id"`
	Long string `db:"long"`
	Hash string `db:"hash"`
}

func NewUrl(hash, long string) Url {
	return Url{
		Hash: hash,
		Long: long,
	}
}

type UrlRepo interface {
	GetByFull(url string) *Url
	GetByHash(url string) *Url
	CreateUrl(url Url) error
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

func (r *MysqlUrlRepo) CreateUrl(url Url) error {
	_, err := r.conn.NamedExec("INSERT INTO urls (`hash`, `long`) VALUES (:hash, :long)", url)

	return err
}

func (r *MysqlUrlRepo) GetByFull(url string) *Url {
	var u Url
	err := r.conn.Get(&u, "SELECT * FROM urls WHERE long = ? LIMIT 1", url)
	if err != nil {
		return nil
	}

	return &u
}

func (r *MysqlUrlRepo) GetByHash(hash string) *Url {
	var u Url
	err := r.conn.Get(&u, "SELECT * FROM urls WHERE `hash` = ? LIMIT 1", hash)
	if err != nil {
		return nil
	}

	return &u
}

func (r *MysqlUrlRepo) GetLastId() (result int) {
	if err := r.conn.QueryRow("SELECT MAX(id) from urls").Scan(&result); err == nil {
		return result
	}

	return 0
}
