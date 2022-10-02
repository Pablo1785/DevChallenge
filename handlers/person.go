package handlers

import (
	"DevChallenge/model"
	"DevChallenge/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type personServiceHandler struct {
	personService service.PersonService
}

func NewPersonServiceHandler(personService service.PersonService) personServiceHandler {
	return personServiceHandler{personService: personService}
}

func (psh *personServiceHandler) AddPerson() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var person model.Person

		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()

		if err := jsonDecoder.Decode(&person); err != nil {
			logrus.Info(err.Error())
			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "Bad Request Format",
				Path:      r.RequestURI,
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		addedPerson, err := psh.personService.Add(ctx, &person)
		if err != nil {
			logrus.Error(err.Error())
			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "Internal Server Error",
				Path:      r.RequestURI,
			}
			if strings.Contains(err.Error(), "already exists") {
				errorResponse.Status = http.StatusConflict
				errorResponse.Message = "Entity you were trying to create already exists"
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		err = json.NewEncoder(w).Encode(addedPerson)
		if err != nil {
			logrus.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	})
}

func (psh *personServiceHandler) CreateOrUpdateTrustConnections() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// URL params
		personId := mux.Vars(r)["personId"]

		var trustConnections model.TrustConnections

		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()

		if err := jsonDecoder.Decode(&trustConnections); err != nil {
			logrus.Info(err.Error())
			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Status:    http.StatusBadRequest,
				Error:     err.Error(),
				Message:   "Bad Request Format",
				Path:      r.RequestURI,
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		if err := trustConnections.Validate(); err != nil {
			logrus.Error(err.Error())
			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Status:    http.StatusUnprocessableEntity,
				Error:     "Invalid Data",
				Message:   err.Error(),
				Path:      r.RequestURI,
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		err := psh.personService.CreateOrUpdateTrustConnections(ctx, personId, trustConnections)
		if err != nil {
			logrus.Error(err.Error())
			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "Internal Server Error",
				Path:      r.RequestURI,
			}
			if strings.Contains(err.Error(), "is missing") {
				errorResponse.Status = http.StatusNotFound
				errorResponse.Message = fmt.Sprintf("Person or one or more people specified as trust connections do not exist")
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	})
}
