package internal

import (
	"log"
)

type result struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func newResult(id int64, name string) *result {
	return &result{
		Id:   id,
		Name: name,
	}
}
func (server *Server) addUrlDB(urlName, shortenUrl string) {
	_, err := server.DB.Exec(`INSERT INTO foo (name, shortUrl) VALUES(?, ?);`, urlName, shortenUrl)
	if err != nil {
		log.Fatal(err)
	}
}
func (server *Server) findShortURL(url string) (normalUrl string, shortUrl string, err error) {
	sqlSt := `
	select * from foo
	WHERE shortUrl = ?;
	`
	rows, err := server.DB.Query(sqlSt, url)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlSt)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var normalUrl string
		var shortUrl string
		err := rows.Scan(&id, &normalUrl, &shortUrl)
		if err != nil {
			log.Fatal(err)
		}
		return normalUrl, shortUrl, nil
	}

	return "", "", err
}

func (server *Server) findNormalUrl(url string) (normalUrl string, shortUrl string, err error) {
	sqlSt := `
	select * from foo
	WHERE name= ?;
	`
	rows, err := server.DB.Query(sqlSt, url)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlSt)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var normalUrl string
		var shortUrl string
		err := rows.Scan(&id, &normalUrl, &shortUrl)
		if err != nil {
			log.Fatal(err)
		}
		return normalUrl, shortUrl, nil
	}

	return "", "", err
}
