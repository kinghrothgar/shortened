package data

import (
	"github.com/jackc/pgx"
)

var Database *pgx.ConnPool

func GetConnectionPool() (pool *pgx.ConnPool, err error) {
	config, err := pgx.ParseURI("postrgres://shortened:shortened@localhost/shortened?sslmode=disable")

	if err != nil {
		return pool, err
	}

	maxConnections := 50

	poolConfig := pgx.ConnPoolConfig{config, maxConnections, nil}

	pool, err = pgx.NewConnPool(poolConfig)

	if err != nil {
		return pool, err
	}

	return pool, err
}

func StoreUrl(url string) (id int64, err error) {
	sql := "INSERT INTO urls (url) VALUES ($1) RETURNING id"

	err = Database.QueryRow(sql, url).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
}

func GetUrl(id int64) (url string, err error) {
	sql := "SELECT url FROM urls WHERE id=$1"

	err = Database.QueryRow(sql, id).Scan(&url)

	if err != nil {
		return url, err
	}

	return url, err
}

func UpdateUrlUsage(urlId int64) (result bool, err error) {
	tx, err := Database.BeginIso("serializable")
	if err != nil {
		return false, err
	}

	selectSql := "SELECT id, visits FROM url_usages WHERE url_id=$1"

	var (
		id     int32
		visits int32
	)
	rows, err := tx.Query(selectSql, urlId)

	if err != nil {
		return false, err
	}

	for rows.Next() {
		err = rows.Scan(&id, &visits)

		if err != nil {
			return false, err
		}

		if id != 0 {
			break
		}
	}

	rows.Close()

	if id == 0 {
		insertSql := "INSERT INTO url_usages (url_id, visits) VALUES ($1, $2) RETURNING id"

		visits = 1
		err := tx.QueryRow(insertSql, urlId, visits).Scan(&id)

		if err != nil {
			return false, err
		}

		err = tx.Commit()

		if err != nil {
			return false, err
		}

		return true, err
	}

	visits = visits + 1

	updateSql := "UPDATE url_usages SET visits=$1 WHERE url_id=$2"
	_, err = tx.Exec(updateSql, visits, urlId)

	if err != nil {
		return false, err
	}

	err = tx.Commit()

	if err != nil {
		return false, err
	}

	return true, err
}
