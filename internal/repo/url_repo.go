package repo

import (
	"github.com/jmoiron/sqlx"
	"net/url"
	"urls/pkg/database"
)

type Url struct {
	Id   uint8  `db:"id" json:"-"`
	Long string `db:"long" json:"long"`
	Hash string `db:"hash" json:"hash"`
}

func NewUrl(hash, long string) Url {
	return Url{
		Hash: hash,
		Long: url.QueryEscape(long),
	}
}

type UrlWriteRepo interface {
	CreateUrl(url Url) error
	BatchCreateUrl(urls []Url) error
}

type UrlReadRepo interface {
	GetByFull(url string) *Url
	GetByHash(url string) *Url
	GetLastId() int
}

type MysqlUrlReadRepo struct {
	conn *sqlx.DB
}

type MysqlUrlWriteRepo struct {
	conn *sqlx.DB
}

func NewMysqlUrlReadRepo() *MysqlUrlReadRepo {
	return &MysqlUrlReadRepo{
		conn: database.GetReadConnection(),
	}
}

func NewMysqlUrlWriteRepo() *MysqlUrlWriteRepo {
	return &MysqlUrlWriteRepo{
		conn: database.GetWriteConnection(),
	}
}

func (r *MysqlUrlWriteRepo) BatchCreateUrl(urls []Url) error {
	_, err := r.conn.NamedExec("INSERT INTO urls (`hash`, `long`) VALUES (:hash, :long)", urls)

	return err
}

func (r *MysqlUrlWriteRepo) CreateUrl(url Url) error {
	_, err := r.conn.NamedExec("INSERT INTO urls (`hash`, `long`) VALUES (:hash, :long)", url)

	return err
}

func (r *MysqlUrlReadRepo) GetByFull(url string) *Url {
	var u Url
	err := r.conn.Get(&u, "SELECT * FROM urls WHERE long = ? LIMIT 1", url)
	if err != nil {
		return nil
	}

	return &u
}

func (r *MysqlUrlReadRepo) GetByHash(hash string) *Url {
	var u Url
	err := r.conn.Get(&u, "SELECT * FROM urls WHERE `hash` = ? LIMIT 1", hash)
	if err != nil {
		return nil
	}

	return &u
}

func (r *MysqlUrlReadRepo) GetLastId() (result int) {
	if err := r.conn.QueryRow("SELECT MAX(id) from urls").Scan(&result); err == nil {
		return result
	}

	return 0
}
