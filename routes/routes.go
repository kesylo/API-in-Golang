package routes

import (
	"net/http"

	"../api"
	"../middleware"
	"../models"
	"../sessions"
	"../utils"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/logout", logoutGetHandler).Methods("GET")

	// facts teller
	r.HandleFunc("/facts", middleware.AuthRequired(factsGetHandler)).Methods("GET")
	r.HandleFunc("/facts", middleware.AuthRequired(factsPostHandler)).Methods("POST")
	//quotes avni
	r.HandleFunc("/quotes", middleware.AuthRequired(quotesGetHandler)).Methods("GET")
	r.HandleFunc("/quotes", middleware.AuthRequired(quotesPostHandler)).Methods("POST")
	// quotes avni v2
	r.HandleFunc("/quotesV2", middleware.AuthRequired(quotesV2GetHandler)).Methods("GET")
	r.HandleFunc("/quotesV2", middleware.AuthRequired(quotesPostHandler)).Methods("POST")
	//Calandar Antoine
	r.HandleFunc("/calandar", middleware.AuthRequired(calandarGetHandler)).Methods("GET")
	r.HandleFunc("/calandar", middleware.AuthRequired(calandarPostHandler)).Methods("POST")

	r.HandleFunc("/register", registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/{username}",
		middleware.AuthRequired(userGetHandler)).Methods("GET")
	return r
}

func factsPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func factsGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "facts.html", struct {
		Title string
	}{
		Title: api.Getfact(),
	})
}

// fonction Get et post Calendar Antoine
func calandarPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func quotesPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func quotesGetHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Quotes string
		Author string
	}{
		Quotes: api.GetQuotesB(1),
		Author: api.GetQuotesB(2),
	}
	utils.ExecuteTemplate(w, "quotes.html", data)
}
func quotesV2GetHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Quotes    string
		Author    string
		Imagelink string
		Title     string
		Message   string
	}{
		Quotes:    api.GetQuoteV2(2),
		Author:    api.GetQuoteV2(1),
		Imagelink: api.GetQuoteV2(3),
		Title:     api.GetQuoteV2(4),
		Message:   api.GetQuoteV2(5),
	}
	utils.ExecuteTemplate(w, "quotesV2.html", data)
	http.Redirect(w, r, "/", 304)
}

func calandarGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "calendar.html", struct {
		NextEventName        string
		NextEventdate        string
		NextEventlocation    string
		NextEventdescription string
		UpcommingEvents      string
	}{

		NextEventName:        api.GetEvent()[0][0],
		NextEventdate:        api.GetEvent()[0][1],
		NextEventlocation:    api.GetEvent()[0][2],
		NextEventdescription: api.GetEvent()[0][3],
		UpcommingEvents:      api.GetNextEvent(),
	})
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	updates, err := models.GetAllUpdates()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Title       string
		Updates     []*models.Update
		DisplayForm bool
	}{
		Title:       "All updates",
		Updates:     updates,
		DisplayForm: true,
	})
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	userId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/", 302)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserId := session.Values["user_id"]
	currentUserId, ok := untypedUserId.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	updates, err := models.GetUpdates(userId)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Title       string
		Updates     []*models.Update
		DisplayForm bool
	}{
		Title:       username,
		Updates:     updates,
		DisplayForm: currentUserId == userId,
	})
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "invalid login")
		default:
			utils.InternalServerError(w)
		}
		return
	}
	userId, err := user.GetId()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userId
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err == models.ErrUsernameTaken {
		utils.ExecuteTemplate(w, "register.html", "username taken")
		return
	} else if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/login", 302)
}
