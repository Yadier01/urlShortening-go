package internal

import (
	"database/sql"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Router *http.ServeMux
	Conn   http.Server
	Adress string
	DB     *sql.DB
}

func NewServer(db *sql.DB) (*Server, error) {
	router := http.NewServeMux()
	server := &Server{
		Router: router,
		Conn: http.Server{
			Handler: router,
			Addr:    ":8082",
		},
		DB: db,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	server.Router.HandleFunc("GET /shorten", index)
	server.Router.HandleFunc("/{url}", server.findShortenUrl)
	server.Router.HandleFunc("POST /shorten", server.form)
}

type pageData struct {
	Url string
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t := template.Must(template.New("").ParseGlob("internal/templates/*"))
	t.ExecuteTemplate(w, tmpl, data)
}
func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", pageData{})
}

func (server *Server) form(w http.ResponseWriter, r *http.Request) {
	str := r.FormValue("url")

	text := addCharToText()
	shortenURl := r.Host + "/" + string(text)

	//check if url ArreadyExist
	normalURL, shortUrl, err := server.findNormalUrl(str)
	if err != nil {
		log.Fatal(err)
	}
	if normalURL == str {
		renderTemplate(w, "url-mssg", pageData{
			Url: shortUrl,
		})
		return
	}
	//check if shortenUrl Already exist
	_, shortUrl, err = server.findShortURL(shortenURl)
	if err != nil {
		log.Fatal(err)
	}

	//fix this later
	if strings.HasSuffix(shortUrl, "/"+string(text)) {
		text = addCharToText()
		return
	}

	server.addUrlDB(str, shortenURl)
}

// util. change location later
func addCharToText() []byte {
	rand.NewSource(time.Now().UnixNano())
	shortCode := make([]byte, 4)
	for i := range shortCode {
		randomIndex := rand.Intn(len(alphabet))
		shortCode[i] = alphabet[randomIndex]
	}

	return shortCode
}

func (server *Server) findShortenUrl(w http.ResponseWriter, r *http.Request) {
	url := r.Host + r.URL.String()
	normalUrl, _, err := server.findShortURL(url)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, normalUrl, http.StatusPermanentRedirect)
}
