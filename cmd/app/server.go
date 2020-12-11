package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	//"time"

	"github.com/Habibullo-1999/crud/cmd/app/middleware"
	"github.com/Habibullo-1999/crud/pkg/customers"
	"github.com/Habibullo-1999/crud/pkg/security"
	"github.com/gorilla/mux"

)

type Server struct {
	mux          *mux.Router
	customersSvc *customers.Service
	securitySvc  *security.Service
}

func NewServer(mux *mux.Router, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

func (s *Server) Init() {
	// s.mux.Handle("/customers", middleware.Logger(http.HandlerFunc(s.handleGetCustomerAll)))
	// s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)

	
	s.mux.HandleFunc("/customers", s.handleGetCustomerAll).Methods(GET)
	s.mux.HandleFunc("/customers", s.handleGetCustomerSave).Methods(POST)
	// chMd := middleware.CheckHeader("Content-Type", "application/json")
	// s.mux.Handle("/customers",chMd(http.HandlerFunc(s.handleGetCustomerSave))).Methods(POST)

	s.mux.HandleFunc("/customers/active", s.handleGetCustomerAllActive).Methods(GET)

	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleGetCustomerBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleGetCustomerUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerRemoveByID).Methods(DELETE)
	s.mux.Use(middleware.Basic(s.securitySvc.Auth))

}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	//idParam := r.URL.Query().Get("id")
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.ByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}

func (s *Server) handleGetCustomerAll(w http.ResponseWriter, r *http.Request) {

	item, err := s.customersSvc.All(r.Context())

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerAllActive(w http.ResponseWriter, r *http.Request) {

	item, err := s.customersSvc.AllActive(r.Context())

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerSave(w http.ResponseWriter, r *http.Request) {
	// idParam := r.URL.Query().Get("id")
	// name := r.URL.Query().Get("name")
	// phone := r.URL.Query().Get("phone")

	// id, err := strconv.ParseInt(idParam, 10, 64)

	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	// if name == "" && phone == "" {
	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }
	// itemR := &customers.Customer{
	// 	ID:    id,
	// 	Name:  name,
	// 	Phone: phone,
	// 	// Active:  bool(true),
	// 	// Created: time.Now(),
	// }

	// item, err := s.customersSvc.Save(r.Context(), itemR)

	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// 	return
	// }
	//получаем данные из параметра запроса
	// idP := r.FormValue("id")
	// name := r.FormValue("name")
	// phone := r.FormValue("phone")

	// id, err := strconv.ParseInt(idP, 10, 64)

	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	// if name == "" && phone == "" {
	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	// 	return
	// }

	// itemR := &customers.Customer{
	// 	ID:    id,
	// 	Name:  name,
	// 	Phone: phone,
	// 	/* Active:true,
	// 	Created:time.Now() */
	// }
	var item *customers.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err = s.customersSvc.Save(r.Context(), item)

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerRemoveByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.RemoveById(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.BlockByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
func (s *Server) handleGetCustomerUnBlockByID(w http.ResponseWriter, r *http.Request) {
	idP := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		return
	}

	item, err := s.customersSvc.UnBlockByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return
	}

}
