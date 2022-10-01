package handlers

import (
	"DevChallenge/model"
	"DevChallenge/service"
	"encoding/json"
	"net/http"
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		addedPerson, err := psh.personService.Add(ctx, &person)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(addedPerson)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	})
}
