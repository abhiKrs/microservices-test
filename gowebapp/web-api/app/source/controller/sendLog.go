package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	// "web-api/app/source/models"
	"web-api/app/utility/kafka"
	log "web-api/app/utility/logger"
	"web-api/app/utility/respond"
	// "web-api/app/utility/validate"
	// "github.com/Shopify/sarama"
	// "github.com/go-chi/chi/v5"
	// "github.com/google/uuid"
)

// package kafkaproducer

type Message struct {
	Type int `json:"type"`
	// Data models.Data `json:"data"`
	Data interface{} `json:"data"`
}

func (pc *SourceController) SendLog(w http.ResponseWriter, r *http.Request) {

	// var reqBody models.SourceLog
	// err := reqBody.Bind(r.Body)
	// if err != nil {
	// 	log.ErrorLogger.Println(err)
	// 	respond.Error(w, http.StatusBadRequest, err)
	// 	return
	// }

	// log.DebugLogger.Println(reqBody)

	// errs := validate.Validate(pc.validator, reqBody)
	// if errs != nil {
	// 	respond.Errors(w, http.StatusBadRequest, errs)
	// 	return
	// }

	// if sourceId == "" {
	// 	response := models.SourceLogResponse{IsSuccessful: false, Message: []string{"Wrong or missing profile Id"}}
	// 	respond.Json(w, http.StatusBadRequest, response)
	// 	return
	// }

	var reqBody interface{}
	json.NewDecoder(r.Body).Decode(&reqBody)

	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(bearerToken, "bearer_")
	sourceToken := strings.Join(splitToken[1:], "")

	// source topic
	// sourceTopic := "source_topic_" + sourceToken
	// userId for now
	sourceTopic := "user_topic_" + sourceToken

	err := kafka.SendLog(sourceTopic, reqBody)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Json(w, http.StatusAccepted, map[string]bool{"isSuccessful": true})

}
