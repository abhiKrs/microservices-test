package controller

import (
	"errors"
	"net/http"
	"strings"

	"web-api/app/source/models"
	// service "web-api/app/source/services"
	"web-api/app/utility/kafka"
	log "web-api/app/utility/logger"
	myRedis "web-api/app/utility/redis"
	"web-api/app/utility/respond"
	"web-api/app/utility/validate"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	// "github.com/google/uuid"
)

type SourceServiceInterface interface {
	GetSourcesByProfileId(profileId uuid.UUID) (*models.GetAllSourceResponse, error)
	CreateSource(pModel models.CreateSourceRequest, profileId uuid.UUID) (*models.CreateSourceResponse, error)
	GetSourceById(sourceId uuid.UUID, profileId uuid.UUID) (*models.GetSourceResponse, error)
	Delete(sourceId uuid.UUID, profileId uuid.UUID) (*models.DeleteSourceResponse, error)
	UpdateSource(profileId uuid.UUID, sourceId uuid.UUID, req *models.UpdateSourceRequest) (*models.GetSourceResponse, error)
}

type SourceController struct {
	sourceService SourceServiceInterface
	validator     *validator.Validate
}

// type SourceController struct {
// 	sourceService service.SourceService
// 	validator     *validator.Validate
// }

func NewSourceController(sourceService SourceServiceInterface, validator *validator.Validate) *SourceController {
	return &SourceController{
		sourceService: sourceService,
		validator:     validator,
	}
}

func (pc *SourceController) SourceSayHello(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello From Source!!!"))
}

// func (pc *SourceController) SourceOnboardSayHello(w http.ResponseWriter, r *http.Request) {

// 	w.Write([]byte("Hello From Source-Onboarding!!!"))
// }

func (pc *SourceController) GetAll(w http.ResponseWriter, r *http.Request) {
	// pc.sourceService.

	log.DebugLogger.Println("getting all sources")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
		return
	}

	// TODO use JWT token
	splitToken := strings.Split(reqToken, "bearer_")
	// log.Println(l)
	// x := strings.Join(l[1:], "")
	// splitToken := strings.Split(reqToken, "bearer_")
	// To refactor after jwt
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	if profileId != "" {

		pId, err := uuid.Parse(profileId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetAllSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.sourceService.GetSourcesByProfileId(pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetAllSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusOK, payload)
		return

	} else {
		response := models.GetAllSourceResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *SourceController) Create(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("Creating a Source")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)

	var reqBody models.CreateSourceRequest
	err := reqBody.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(pc.validator, reqBody)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	if profileId != "" {
		pId, err := uuid.Parse(profileId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.CreateSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.sourceService.CreateSource(reqBody, pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.CreateSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}

		// sourceTopic := "source_topic_" + payload.Data.SourceToken
		// log.DebugLogger.Println("created topic named: ", sourceTopic)

		// ---------------Temporary======================
		var sourceTopic string
		if payload.Data.ProfileId != uuid.Nil {
			sourceTopic = "user_topic_" + payload.Data.ProfileId.String()
			// ======================================
			log.DebugLogger.Println("creating topic named: ", sourceTopic)
			// Crete kafka topic for source
			go kafka.CreateTopic(sourceTopic)
		}

		// Create redis store item
		teamTopic := "team_topic_" + payload.Data.TeamId.String()
		err = myRedis.PushSourceToTeamSet(sourceTopic, teamTopic)
		if err != nil {
			log.DebugLogger.Println(err)
		}

		err = myRedis.AddSourceTeamPair(sourceTopic, teamTopic)
		if err != nil {
			log.DebugLogger.Println(err)
		}
		respond.Json(w, http.StatusCreated, payload)
		return

	} else {
		response := models.CreateSourceResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *SourceController) GetOne(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("Creating a Source")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	pId, err := uuid.Parse(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	log.InfoLogger.Println(profileId)
	if sourceID := chi.URLParam(r, "id"); sourceID != "" {
		sId, err := uuid.Parse(sourceID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.sourceService.GetSourceById(sId, pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *SourceController) Delete(w http.ResponseWriter, r *http.Request) {

	log.InfoLogger.Println("Deleting a Source")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	pId, err := uuid.Parse(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	log.InfoLogger.Println(profileId)
	if sourceID := chi.URLParam(r, "id"); sourceID != "" {
		sId, err := uuid.Parse(sourceID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.sourceService.Delete(sId, pId)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
	// log.Println("done")
}

func (pc *SourceController) UpdateSource(w http.ResponseWriter, r *http.Request) {
	log.InfoLogger.Println("Updating a Source")
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		log.ErrorLogger.Println("Missing Authorization Header")
		respond.Error(w, http.StatusProxyAuthRequired, errors.New("missing authorization header"))
	}

	// TODO use JWT token

	splitToken := strings.Split(reqToken, "bearer_")
	log.DebugLogger.Println(splitToken)
	profileId := strings.Join(splitToken[1:], "")
	log.InfoLogger.Println(profileId)
	pId, err := uuid.Parse(profileId)
	if err != nil {
		log.ErrorLogger.Println(err)
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}

	var reqBody models.UpdateSourceRequest
	err = reqBody.Bind(r.Body)
	if err != nil {
		log.ErrorLogger.Println(err)
		respond.Error(w, http.StatusBadRequest, err)
		return
	}

	errs := validate.Validate(pc.validator, reqBody)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	if sourceID := chi.URLParam(r, "id"); sourceID != "" {
		sId, err := uuid.Parse(sourceID)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		payload, err := pc.sourceService.UpdateSource(pId, sId, &reqBody)
		if err != nil {
			log.ErrorLogger.Println(err)
			response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
			respond.Json(w, http.StatusBadRequest, response)
			return
		}
		respond.Json(w, http.StatusAccepted, payload)
		return

	} else {
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{"Wrong or missing source Id"}}
		respond.Json(w, http.StatusBadRequest, response)
		return
	}
}
