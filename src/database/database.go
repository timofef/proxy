package database

import (
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type Database struct {
	dbConn *pgx.Conn
}

func InitConnection() (*Database, error) {
	conf := pgx.ConnConfig{
		User: "proxy",
		Password: "postgres",
		Database: "security",
	}

	conn, err := pgx.Connect(conf)
	if err!= nil {
		return nil, err
	}

	logrus.Info("Successfully connected to database")

	return &Database{dbConn: conn}, nil
}

func (db *Database) CloseConnection() {
	db.dbConn.Close()
}

func (db *Database) AddRequest(req Request) {
	_, err := db.dbConn.Exec("INSERT INTO requests (host, request) VALUES ($1, $2)", req.Host, req.Request)

	if err != nil {
		logrus.Error(err)
	}
}

func (db *Database) GetRequestList() ([]Request, error) {
	rows, err := db.dbConn.Query(
		"SELECT * FROM requests",
	)

	if err != nil {
		return nil, err
	}

	requests := make([]Request, 0)
	for rows.Next() {
		req := Request{}
		err := rows.Scan(&req.Id, &req.Host, &req.Request)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	rows.Close()

	return requests, nil
}

func (db *Database) GetRequest(id int) (Request, error) {
	row := db.dbConn.QueryRow(
		"SELECT * FROM requests WHERE id = $1", id)

	var request = Request{}

	err := row.Scan(&request.Id, &request.Host, &request.Request)
	if err != nil {
		logrus.Warn("Can't get single request from database")
		logrus.Error(err)
	}

	return request, err
}
