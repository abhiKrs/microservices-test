package controller

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"web-api/app/config"
	"web-api/app/source/dao"
	"web-api/app/source/models"
	"web-api/app/source/models/mocks"
	"web-api/app/source/schema"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Fixture struct {
	Sources []schema.Source
}

func NewFixture() *Fixture {
	return &Fixture{}
}

func (f *Fixture) SetupSources(sources []schema.Source) {
	f.Sources = sources
}

func (f *Fixture) Teardown() {
	// Clean up the test data here
}

func (f *Fixture) GetSource(id uuid.UUID) (schema.Source, error) {
	for _, source := range f.Sources {
		if source.Base.ID == id {
			return source, nil
		}
	}
	return schema.Source{}, fmt.Errorf("source not found")
}

// // var SourceControl = NewSourceController(services.NewSourceService(*dao.New(&gorm.DB{}, config.CFG)), validator.New())
var SourceControl = NewSourceController(mocks.NewMockSourceService(*dao.New(&gorm.DB{}, config.CFG)), validator.New())

func TestGetAll(t *testing.T) {
	// Set up test data
	sources := []schema.Source{
		{
			ProfileId:  uuid.New(),
			TeamId:     uuid.New(),
			Name:       "First",
			SourceType: 1,
			Base: schema.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			ProfileId:  uuid.New(),
			TeamId:     uuid.New(),
			Name:       "Second",
			SourceType: 2,
			Base: schema.Base{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
	fixture := NewFixture()
	fixture.SetupSources(sources)
	defer fixture.Teardown()

	// SourceControl.sourceService.

	// Make request to API
	req, err := http.NewRequest("GET", "/sources", nil)
	if err != nil {
		t.Fatal(err)
	}

	profileId := uuid.New()
	req.Header.Set("authorization", fmt.Sprintf("bearer_%v", profileId.String()))
	rr := httptest.NewRecorder()
	// rr.WriteHeader()
	handler := http.HandlerFunc(SourceControl.GetAll)
	handler.ServeHTTP(rr, req)

	// 	// Check response
	// 	if status := rr.Code; status != http.StatusOK {
	// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	// 	}

	var response models.GetAllSourceResponse
	// log.Println(string(rr.Body.Bytes()))
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// log.Println(response)

	// data := response.Data

	// if len(data) != len(sources) {
	// 	t.Errorf("handler returned wrong number of sources: got %v want %v", len(data), len(sources))
	// }

	// for i, res := range data {
	// 	if res.SourceId != sources[i].Base.ID {
	// 		t.Errorf("handler returned wrong source ID: got %v want %v", res.SourceId, sources[i].Base.ID)
	// 	}
	// 	if res.Name != sources[i].Name {
	// 		t.Errorf("handler returned wrong source name: got %v want %v", res.Name, sources[i].Name)
	// 	}
	// }
}

// func TestCreateSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{Name: "Alice"}
// 	fixture := NewFixture()
// 	defer fixture.Teardown()
// func TestCreateSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{Name: "Alice"}
// 	fixture := NewFixture()
// 	defer fixture.Teardown()

// 	// Make request to API
// 	jsonData, err := json.Marshal(source)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/sources", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.Create)
// 	handler.ServeHTTP(rr, req)
// 	// Make request to API
// 	jsonData, err := json.Marshal(source)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/sources", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.Create)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 	}
// 	// Check response
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 	}

// 	var response schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var response schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if response.Name != source.Name {
// 		t.Errorf("handler returned wrong source name: got %v want %v", response.Name, source.Name)
// 	}
// }
// 	if response.Name != source.Name {
// 		t.Errorf("handler returned wrong source name: got %v want %v", response.Name, source.Name)
// 	}
// }

// func TestGetSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{
// 		ProfileId:  uuid.New(),
// 		TeamId:     uuid.New(),
// 		Name:       "First",
// 		SourceType: 1,
// 		Base: schema.Base{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()
// func TestGetSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{
// 		ProfileId:  uuid.New(),
// 		TeamId:     uuid.New(),
// 		Name:       "First",
// 		SourceType: 1,
// 		Base: schema.Base{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()

// 	// Make request to API
// 	req, err := http.NewRequest("GET", "/sources", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// Make request to API
// 	req, err := http.NewRequest("GET", "/sources", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.GetOne)
// 	handler.ServeHTTP(rr, req)
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.GetOne)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	var response schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var response schema.Source
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if response.Base.ID != source.Base.ID {
// 		t.Errorf("handler returned wrong source ID: got %v want %v", response.Base.ID, source.Base.ID)
// 	}
// 	if response.Name != source.Name {
// 		t.Errorf("handler returned wrong source name: got %v want %v", response.Name, source.Name)
// 	}
// }
// 	if response.Base.ID != source.Base.ID {
// 		t.Errorf("handler returned wrong source ID: got %v want %v", response.Base.ID, source.Base.ID)
// 	}
// 	if response.Name != source.Name {
// 		t.Errorf("handler returned wrong source name: got %v want %v", response.Name, source.Name)
// 	}
// }

// func TestUpdateSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{
// 		ProfileId:  uuid.New(),
// 		TeamId:     uuid.New(),
// 		Name:       "First",
// 		SourceType: 1,
// 		Base: schema.Base{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()
// func TestUpdateSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{
// 		ProfileId:  uuid.New(),
// 		TeamId:     uuid.New(),
// 		Name:       "First",
// 		SourceType: 1,
// 		Base: schema.Base{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()

// 	// Make request to API
// 	newName := "Alicia"
// 	jsonData, err := json.Marshal(schema.Source{Name: newName})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", fmt.Sprintf("/sources/%v", source.Base.ID), bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.UpdateSource)
// 	handler.ServeHTTP(rr, req)
// 	// Make request to API
// 	newName := "Alicia"
// 	jsonData, err := json.Marshal(schema.Source{Name: newName})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", fmt.Sprintf("/sources/%v", source.Base.ID), bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.UpdateSource)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check that source was updated in fixture
// 	updatedSource, err := fixture.GetSource(source.Base.ID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if updatedSource.Name != newName {
// 		t.Errorf("handler did not update source name: got %v want %v", updatedSource.Name, newName)
// 	}
// }
// 	// Check that source was updated in fixture
// 	updatedSource, err := fixture.GetSource(source.Base.ID)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if updatedSource.Name != newName {
// 		t.Errorf("handler did not update source name: got %v want %v", updatedSource.Name, newName)
// 	}
// }

// func TestDeleteSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{
// 		ProfileId:  uuid.New(),
// 		TeamId:     uuid.New(),
// 		Name:       "First",
// 		SourceType: 1,
// 		Base: schema.Base{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()
// func TestDeleteSource(t *testing.T) {
// 	// Set up test data
// 	source := schema.Source{
// 		ProfileId:  uuid.New(),
// 		TeamId:     uuid.New(),
// 		Name:       "First",
// 		SourceType: 1,
// 		Base: schema.Base{
// 			ID:        uuid.New(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		},
// 	}
// 	fixture := NewFixture()
// 	fixture.SetupSources([]schema.Source{source})
// 	defer fixture.Teardown()

// 	// Make request to API
// 	req, err := http.NewRequest("DELETE", fmt.Sprintf("/sources/%v", source.Base.ID), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.Delete)
// 	handler.ServeHTTP(rr, req)
// 	// Make request to API
// 	req, err := http.NewRequest("DELETE", fmt.Sprintf("/sources/%v", source.Base.ID), nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(SourceControl.Delete)
// 	handler.ServeHTTP(rr, req)

// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// 	// Check response
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check that source was deleted from fixture
// 	_, err = fixture.GetSource(source.Base.ID)
// 	if err == nil {
// 		t.Errorf("handler did not delete source")
// 	}
// }
// 	// Check that source was deleted from fixture
// 	_, err = fixture.GetSource(source.Base.ID)
// 	if err == nil {
// 		t.Errorf("handler did not delete source")
// 	}
// }
