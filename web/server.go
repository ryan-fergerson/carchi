package web

import (
	"carchi/internal"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

func StartServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/", RecentConversationViewHandler).Methods(http.MethodGet)
	r.HandleFunc("/c/{id}", ConversationViewHandler).Methods(http.MethodGet)
	r.HandleFunc("/search", SearchViewHandler).Methods(http.MethodGet)
	r.HandleFunc("/search", SearchViewHandler).Methods(http.MethodPost)

	http.Handle("/", r)

	return http.ListenAndServe(":8080", nil)
}

func SearchViewHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error processing search request", http.StatusBadRequest)
		return
	}

	query := r.FormValue("query")

	s, err := internal.NewConversationBrowserService()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error processing search request", http.StatusBadRequest)
		return
	}

	searchResults, err := s.SearchConversations(query)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error processing search request", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/search.html"))
	err = tmpl.Execute(w, searchResults)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error processing search request", http.StatusBadRequest)
		return
	}
}

func RecentConversationViewHandler(w http.ResponseWriter, r *http.Request) {
	s, err := internal.NewConversationBrowserService()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error viewing recent conversations", http.StatusBadRequest)
		return
	}

	recentConversations, err := s.GetRecentConversations()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error viewing recent conversations", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/recent.html"))
	err = tmpl.Execute(w, recentConversations)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error viewing recent conversations", http.StatusBadRequest)
		return
	}
}

func ConversationViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Missing conversation id", http.StatusBadRequest)
		return
	}

	s, err := internal.NewConversationBrowserService()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error viewing conversation: "+id, http.StatusBadRequest)
		return
	}

	messageParts, err := s.GetMessagePartsByConversationId(id)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error viewing conversation: "+id, http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/conversation.html"))
	err = tmpl.Execute(w, messageParts)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error viewing conversation: "+id, http.StatusBadRequest)
		return
	}
}
