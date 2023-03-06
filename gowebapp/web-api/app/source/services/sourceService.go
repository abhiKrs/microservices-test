package services

import (
	"fmt"
	"strings"
	"web-api/app/source/models"

	// "web-api/app/source/schema"
	"web-api/app/source/dao"
	logs "web-api/app/utility/logger"
	"web-api/app/utility/respond"

	"github.com/google/uuid"
)

type SourceService struct {
	dao dao.SourceDataAccess
}

func NewSourceService(dao dao.SourceDataAccess) *SourceService {
	return &SourceService{
		dao: dao,
	}
}

func (ps *SourceService) GetSourcesByProfileId(profileId uuid.UUID) (*models.GetAllSourceResponse, error) {

	var err error

	// Get Source
	dbSources, err := ps.dao.SourceAccess.GetByProfileId(profileId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.GetAllSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	var dataList []models.SourceBodyResponse

	for _, dbSource := range *dbSources {
		data := models.SourceBodyResponse{
			SourceId:    dbSource.Base.ID,
			ProfileId:   dbSource.ProfileId,
			Name:        dbSource.Name,
			SourceType:  dbSource.SourceType,
			TeamId:      dbSource.TeamId,
			SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
			CreatedAt:   dbSource.Base.CreatedAt,
		}
		dataList = append(dataList, data)
	}

	bodyPayload := models.GetAllSourceResponse{
		IsSuccessful: true,
		Data:         dataList,
	}

	return &bodyPayload, nil

}

func (ps *SourceService) CreateSource(pModel models.CreateSourceRequest, profileId uuid.UUID) (*models.CreateSourceResponse, error) {

	var err error
	var pId uuid.UUID
	if pModel.TeamId != "" {
		pId = uuid.New()
	} else {
		pId = uuid.MustParse(pModel.TeamId)
	}

	// Get Source
	dbSource, err := ps.dao.SourceAccess.Create(profileId, pModel.Name, pModel.SourceType, pId)
	if err != nil {
		response := models.CreateSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	bodyPayload := models.SourceBodyResponse{
		SourceId:    dbSource.Base.ID,
		ProfileId:   dbSource.ProfileId,
		Name:        dbSource.Name,
		SourceType:  dbSource.SourceType,
		TeamId:      dbSource.TeamId,
		SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
		CreatedAt:   dbSource.Base.CreatedAt,
	}

	response := models.CreateSourceResponse{IsSuccessful: true, Data: bodyPayload}
	return &response, nil
}

func (ps *SourceService) GetSourceById(sourceId uuid.UUID, profileId uuid.UUID) (*models.GetSourceResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get Source
	dbSource, err := ps.dao.SourceAccess.GetById(sourceId)
	if err != nil {
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbSource.ProfileId != profileId {
		logs.ErrorLogger.Printf("Unathorised access to source: %v from different profile: %v /n", dbSource.Base.ID, profileId)
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	bodyPayload := models.SourceBodyResponse{
		SourceId:    dbSource.Base.ID,
		ProfileId:   dbSource.ProfileId,
		Name:        dbSource.Name,
		SourceType:  dbSource.SourceType,
		TeamId:      dbSource.TeamId,
		SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
		CreatedAt:   dbSource.Base.CreatedAt,
	}

	response := models.GetSourceResponse{IsSuccessful: true, Data: bodyPayload}
	return &response, nil

}

func (ps *SourceService) Delete(sourceId uuid.UUID, profileId uuid.UUID) (*models.DeleteSourceResponse, error) {

	// var dbEmail schema.Email
	var err error

	// Get Source
	dbSource, err := ps.dao.SourceAccess.GetById(sourceId)
	if err != nil {
		response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbSource.ProfileId != profileId {
		logs.ErrorLogger.Printf("Unathorised access to source: %v from different profile: %v /n", dbSource.Base.ID, profileId)
		response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	dbSource, err = ps.dao.SourceAccess.DeleteById(sourceId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.DeleteSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	response := models.DeleteSourceResponse{IsSuccessful: true}
	return &response, nil

}

func (ps *SourceService) UpdateSource(profileId uuid.UUID, sourceId uuid.UUID, req *models.UpdateSourceRequest) (*models.GetSourceResponse, error) {

	var err error

	// Get Source
	dbSource, err := ps.dao.SourceAccess.GetById(sourceId)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{err.Error()}}
		return &response, err
	}

	if dbSource.ProfileId != profileId {
		logs.ErrorLogger.Printf("Unathorised access to source: %v from different profile: %v /n", dbSource.Base.ID, profileId)
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	dbSource, err = ps.dao.SourceAccess.Update(dbSource, req.Name)
	if err != nil {
		logs.ErrorLogger.Println(err)
		response := models.GetSourceResponse{IsSuccessful: false, Message: []string{fmt.Sprintf("Unathorised access to source: %v from different profile", dbSource.Base.ID)}}
		return &response, respond.ErrBadRequest
	}

	bodyPayload := models.SourceBodyResponse{
		SourceId:    dbSource.Base.ID,
		ProfileId:   dbSource.ProfileId,
		Name:        dbSource.Name,
		SourceType:  dbSource.SourceType,
		TeamId:      dbSource.TeamId,
		SourceToken: strings.Join([]string{dbSource.TeamId.String(), dbSource.Base.ID.String()}, "_"),
		CreatedAt:   dbSource.Base.CreatedAt,
	}

	response := models.GetSourceResponse{IsSuccessful: true, Data: bodyPayload}
	return &response, nil
}
