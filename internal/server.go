package internal

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
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
	Url template.HTML
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
	userGivenURl := r.FormValue("url")

	text := addCharToText()
	shortenURl := r.Host + "/" + string(text)

	server.checkIfUrlExist(w, userGivenURl)

	dbShortUrl := server.addUrlDB(userGivenURl, shortenURl)
	htmlResponse := makeResponse(dbShortUrl)

	renderTemplate(w, "url-mssg", pageData{
		Url: template.HTML(htmlResponse),
	})
	return
}
func (server *Server) checkIfUrlExist(w http.ResponseWriter, userGivenUrl string) {
	//check if url AlreadyExist
	normalURL, shortUrl, err := server.findNormalUrl(userGivenUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	htmlResponse := makeResponse(shortUrl)

	if normalURL == userGivenUrl {
		w.WriteHeader(http.StatusFound)
		renderTemplate(w, "url-mssg", pageData{
			Url: template.HTML(htmlResponse),
		})
		return
	}
}
func makeResponse(url string) string {
	return fmt.Sprintf(`
    <input readonly id="myInput" value="%s">
    <button onclick="getCopy()">Copy</button>
`, url)
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
