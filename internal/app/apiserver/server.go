package apiserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/assanoff/http-crud-server/internal/app/model"
	"github.com/assanoff/http-crud-server/internal/app/store"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type server struct {
	router   *mux.Router
	logger   *logrus.Logger
	store    store.Store
	endpoint string
}

func newServer(store store.Store, endpoint string) *server {
	s := &server{
		router:   mux.NewRouter(),
		logger:   logrus.New(),
		store:    store,
		endpoint: endpoint,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

	s.router.Use(s.logRequest)
	s.router.HandleFunc("/", HandlerHello)
	s.router.HandleFunc(s.endpoint+"/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc(s.endpoint+"/users/{id:[0-9]+}", s.handleUsersByID()).Methods("GET")
	s.router.HandleFunc(s.endpoint+"/users/", s.handleUsersByField()).Queries("field", "{field}", "val", "{value}").Methods("GET")
	s.router.HandleFunc(s.endpoint+"/users", s.handleUsers()).Methods("GET")
	s.router.HandleFunc(s.endpoint+"/users/{id:[0-9]+}", s.handleUpdateUsersByID()).Methods("PUT")
	s.router.HandleFunc(s.endpoint+"/users/{id:[0-9]+}", s.handleDeleteUsersByID()).Methods("DELETE")

}

// HandlerHello ...
func HandlerHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("crud-server"))
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Name:  req.Name,
			Email: req.Email,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleUsersByID() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().GetUserByID(id)

		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}
func (s *server) handleDeleteUsersByID() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.User().DeleteUserByID(id)

		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUpdateUsersByID() http.HandlerFunc {

	type request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Name:  req.Name,
			Email: req.Email,
		}

		u, err = s.store.User().UpdateUserByID(id, u)

		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}
func (s *server) handleUsersByField() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		fieldName := r.FormValue("field")
		value := r.FormValue("val")

		u, err := s.store.User().GetUserByField(fieldName, value)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleUsers() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.store.User().GetUsers()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
