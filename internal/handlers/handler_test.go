package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pwydra/shred/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExerciseDao is a mock implementation of the ExerciseDao interface
type MockExerciseDao struct {
	mock.Mock
}

func (m *MockExerciseDao) Create(exerciseRequest *model.ExerciseRequest) (*model.Exercise, error) {
	args := m.Called(exerciseRequest)
	return args.Get(0).(*model.Exercise), args.Error(1)
}

func (m *MockExerciseDao) Read(uuid uuid.UUID) (*model.Exercise, error) {
	args := m.Called(uuid)
	return args.Get(0).(*model.Exercise), args.Error(1)
}

func (m *MockExerciseDao) Update(exercise *model.Exercise) error {
	args := m.Called(exercise)
	return args.Error(0)
}

func (m *MockExerciseDao) Delete(uuid uuid.UUID) error {
	args := m.Called(uuid)
	return args.Error(0)
}

func newUpdateExerciseRequest() model.Exercise {
	return model.Exercise{
		ExerciseUuid: uuid.New(),
		ExerciseFields: model.ExerciseFields{
			Description:      "Updated Description",
			Instructions:     "Updated Instructions",
			Cues:             "Updated Cues",
			VideoUrl:         "http://example.com/updated.mp4",
			CategoryCode:     "updated-category",
			LicenseShortName: "CC-BY-NC",
			LicenseAuthor:    "Jane Doe",
		},
	}
}

func TestCreateExercise(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/exercises", handler.CreateExercise)

	exReq := model.ExerciseRequest{
		ExerciseFields: model.ExerciseFields{
			ExerciseName:     "Squat",
			Description:      "Lower Body",
			Instructions:     "Stand with feet shoulder-width apart",
			Cues:             "Keep chest up and back flat",
			VideoUrl:         "http://example.com/squat.mp4",
			CategoryCode:     "strength",
			LicenseShortName: "CC-BY",
			LicenseAuthor:    "John Doe",
		},
		CreatedBy: uuid.New(),
	}

	ex := model.Exercise{
		ExerciseUuid:   uuid.New(),
		ExerciseFields: exReq.ExerciseFields,
	}

	mockDao.On("Create", &exReq).Return(&ex, nil)

	body, _ := json.Marshal(exReq)
	req, _ := http.NewRequest(http.MethodPost, "/exercises", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response model.Exercise
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ex.CreatedBy, response.CreatedBy)
}

func TestGetExercise(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/exercises/:uuid", handler.GetExercise)

	exUuid := uuid.New()
	ex := model.Exercise{
		ExerciseUuid: exUuid,
		ExerciseFields: model.ExerciseFields{
			ExerciseName:     "Squat",
			Description:      "Lower Body",
			Instructions:     "Stand with feet shoulder-width apart",
			Cues:             "Keep chest up and back flat",
			VideoUrl:         "http://example.com/squat.mp4",
			CategoryCode:     "strength",
			LicenseShortName: "CC-BY",
			LicenseAuthor:    "John Doe",
		},
		AuditRecord: model.AuditRecord{
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockDao.On("Read", exUuid).Return(&ex, nil)

	req, _ := http.NewRequest(http.MethodGet, "/exercises/"+exUuid.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.Exercise
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ex.ExerciseUuid, response.ExerciseUuid)
}

func TestGetExercise_BadUuid(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/exercises/:uuid", handler.GetExercise)

	req, _ := http.NewRequest(http.MethodGet, "/exercises/badUuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, w.Body.String(), `{"error":"invalid UUID length: 7"}`)
}

func TestGetExercise_DbError(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/exercises/:uuid", handler.GetExercise)

	exUuid := uuid.New()
	mockDao.On("Select", exUuid).Return(assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/exercises/"+exUuid.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateExercise(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/exercises/:uuid", handler.UpdateExercise)

	ex := newUpdateExerciseRequest()

	mockDao.On("Update", &ex).Return(nil)

	body, _ := json.Marshal(ex)
	req, _ := http.NewRequest(http.MethodPut, "/exercises/"+ex.ExerciseUuid.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.Exercise
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ex.ExerciseUuid, response.ExerciseUuid)
}

func TestUpdateExercise_BadUuid(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/exercises/:uuid", handler.UpdateExercise)

	ex := newUpdateExerciseRequest()

	body, _ := json.Marshal(ex)
	req, _ := http.NewRequest(http.MethodPut, "/exercises/badUuid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":"invalid UUID length: 7"}`)
}

func TestUpdateExercise_UnmatchedUuid(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/exercises/:uuid", handler.UpdateExercise)

	ex := newUpdateExerciseRequest()

	body, _ := json.Marshal(ex)
	req, _ := http.NewRequest(http.MethodPut, "/exercises/"+uuid.New().String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":"UUID in path does not match UUID in request body"}`)
}

func TestUpdateExercise_BadRequestBody(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/exercises/:uuid", handler.UpdateExercise)

	body := []byte(`blah`) // Invalid JSON structure

	req, _ := http.NewRequest(http.MethodPut, "/exercises/"+uuid.New().String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), `{"error":"invalid character 'b' looking for beginning of value"}`)
}

func TestDeleteExercise(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/exercises/:uuid", handler.DeleteExercise)

	exUuid := uuid.New()

	mockDao.On("Delete", exUuid).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/exercises/"+exUuid.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteExercise_DbError(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/exercises/:uuid", handler.DeleteExercise)

	exUuid := uuid.New()

	mockDao.On("Delete", exUuid).Return(assert.AnError)

	req, _ := http.NewRequest(http.MethodDelete, "/exercises/"+exUuid.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteExercise_BadUuid(t *testing.T) {
	mockDao := new(MockExerciseDao)
	handler := NewHandler(mockDao)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/exercises/:uuid", handler.DeleteExercise)

	exUuid := uuid.New()

	mockDao.On("Delete", exUuid).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/exercises/badGuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
