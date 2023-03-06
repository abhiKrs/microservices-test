package mocks

import (
	// "context"
	"log"
	"web-api/app/source/dao"
	"web-api/app/source/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSourceService struct {
	dao dao.SourceDataAccess
	mock.Mock
	getSourcesFn func(profileId uuid.UUID) (*models.GetAllSourceResponse, error)
	// createSourceFn func(pModel models.CreateSourceRequest, profileId uuid.UUID) (*models.CreateSourceResponse, error)
}

func NewMockSourceService(dao dao.SourceDataAccess) *MockSourceService {
	return &MockSourceService{
		dao: dao,
	}
}

func (m *MockSourceService) GetSourcesByProfileId(profileId uuid.UUID) (*models.GetAllSourceResponse, error) {
	log.Println("inside mock source service")

	if m != nil && m.getSourcesFn != nil {
		log.Print("return custom mock")
		return m.getSourcesFn(profileId)
	}
	return &models.GetAllSourceResponse{IsSuccessful: true, Data: []models.SourceBodyResponse{}}, nil
}

func (m *MockSourceService) CreateSource(pModel models.CreateSourceRequest, profileId uuid.UUID) (*models.CreateSourceResponse, error) {

	ret := m.Called(pModel, profileId)

	var r0 *models.CreateSourceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.CreateSourceResponse)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockSourceService) GetSourceById(sourceId uuid.UUID, profileId uuid.UUID) (*models.GetSourceResponse, error) {

	ret := m.Called(sourceId, profileId)

	var r0 *models.GetSourceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.GetSourceResponse)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockSourceService) Delete(sourceId uuid.UUID, profileId uuid.UUID) (*models.DeleteSourceResponse, error) {

	ret := m.Called(sourceId, profileId)

	var r0 *models.DeleteSourceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.DeleteSourceResponse)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

func (m *MockSourceService) UpdateSource(profileId uuid.UUID, sourceId uuid.UUID, req *models.UpdateSourceRequest) (*models.GetSourceResponse, error) {

	ret := m.Called(profileId, sourceId, req)

	var r0 *models.GetSourceResponse

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.GetSourceResponse)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
