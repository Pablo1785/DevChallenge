package handlers

import (
	"DevChallenge/model"
	"DevChallenge/service"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type messageServiceHandler struct {
	messageService service.MessageService
}

func NewMessageServiceHandler(messageService service.MessageService) messageServiceHandler {
	return messageServiceHandler{messageService: messageService}
}

func (msh *messageServiceHandler) SendMessage() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var message model.Message

		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()

		if err := jsonDecoder.Decode(&message); err != nil {
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

		messageRecipients, err := msh.messageService.Send(ctx, &message)
		if err != nil {
			logrus.Error(err.Error())
			errorResponse := ErrorResponse{
				Timestamp: time.Now(),
				Status:    http.StatusInternalServerError,
				Error:     err.Error(),
				Message:   "Internal Server Error",
				Path:      r.RequestURI,
			}
			w.WriteHeader(errorResponse.Status)
			json.NewEncoder(w).Encode(&errorResponse)
			return
		}

		err = json.NewEncoder(w).Encode(messageRecipients)
		if err != nil {
			logrus.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	})
}
