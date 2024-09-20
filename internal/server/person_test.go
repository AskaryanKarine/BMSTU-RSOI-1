package server

import (
	"encoding/json"
	"errors"
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/models"
	"github.com/AskaryanKarine/BMSTU-ds-1/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gojuno/minimock/v3"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var regularPerson = models.Person{
	ID:      1,
	Name:    "test",
	Age:     1,
	Address: "test",
	Work:    "test",
}

func TestServer_createPerson(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	type fields struct {
		echo *echo.Echo
		pr   personRepository
	}
	tests := []struct {
		name               string
		fields             fields
		body               *models.Person
		expectedHTTPStatus int
	}{
		{
			name: "http-201: create",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).CreatePersonMock.Return(regularPerson, nil),
			},
			body:               &regularPerson,
			expectedHTTPStatus: 201,
		},
		{
			name: "http-400: can not bind",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			body:               &models.Person{},
			expectedHTTPStatus: 400,
		},
		{
			name: "http-400: validation error",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			body:               &models.Person{Name: "test", Age: -10},
			expectedHTTPStatus: 400,
		},
		{
			name: "http-500: database error",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).CreatePersonMock.Return(models.Person{}, errors.New("database error")),
			},
			body:               &regularPerson,
			expectedHTTPStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				pr:   tt.fields.pr,
			}

			var err error
			jsonData := []byte("")
			if tt.body != nil {
				jsonData, err = json.Marshal(tt.body)
				if err != nil {
					t.Errorf("json marshal error")
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(string(jsonData)))
			req.Header.Set("Content-type", "application/json")
			rw := httptest.NewRecorder()
			c := s.echo.NewContext(req, rw)

			err = s.createPerson(c)
			if err != nil {
				t.Errorf("createPerson() error = %v", err)
			}

			code := rw.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("createPerson() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}

		})
	}
}

func TestServer_deletePersonByID(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	type fields struct {
		echo *echo.Echo
		pr   personRepository
	}

	tests := []struct {
		name               string
		fields             fields
		pathParams         string
		expectedHTTPStatus int
	}{
		{
			name: "http-204: deleted correctly",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).DeletePersonByIDMock.Return(nil),
			},
			pathParams:         "1",
			expectedHTTPStatus: 204,
		},
		{
			name: "http-400: can not parse path param",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			pathParams:         "qwerty",
			expectedHTTPStatus: 400,
		},
		{
			name: "http-500: database error",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).DeletePersonByIDMock.Return(errors.New("database error")),
			},
			pathParams:         "1",
			expectedHTTPStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				pr:   tt.fields.pr,
			}
			r := httptest.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()
			c := s.echo.NewContext(r, w)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParams)

			err := s.deletePersonByID(c)
			if err != nil {
				t.Errorf("deletePersonByID() error = %v", err)
			}

			code := w.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("deletePersonByID() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}
		})
	}
}

func TestServer_getPersonByID(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	type fields struct {
		echo *echo.Echo
		pr   personRepository
	}
	tests := []struct {
		name               string
		fields             fields
		pathParams         string
		expectedHTTPStatus int
		result             models.Person
	}{
		{
			name: "http-200: person found",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).GetPersonByIDMock.Return(regularPerson, nil),
			},
			pathParams:         "1",
			result:             regularPerson,
			expectedHTTPStatus: 200,
		},
		{
			name: "http-400: can not parse path param",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			pathParams:         "qwerty",
			expectedHTTPStatus: 400,
		},
		{
			name: "http-404: person not found",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).GetPersonByIDMock.Return(models.Person{}, gorm.ErrRecordNotFound),
			},
			pathParams:         "1",
			expectedHTTPStatus: 404,
		},
		{
			name: "http-500: database error",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).GetPersonByIDMock.Return(models.Person{}, errors.New("database error")),
			},
			pathParams:         "1",
			expectedHTTPStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				pr:   tt.fields.pr,
			}

			r := httptest.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()
			c := s.echo.NewContext(r, w)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParams)

			err := s.getPersonByID(c)
			if err != nil {
				t.Errorf("getPersonByID() error = %v", err)
			}

			code := w.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("getPersonByID() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}

			body, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Errorf("ReadAll error")
			}
			var res models.Person
			err = json.Unmarshal(body, &res)
			if err != nil {
				t.Errorf("json unmarshal error")
			}

			if tt.expectedHTTPStatus == http.StatusOK {
				if !reflect.DeepEqual(res, tt.result) {
					t.Errorf("getPersonByID() expected %v, but got %v", tt.result, res)
				}
			}
		})
	}
}

func TestServer_getPersons(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	type fields struct {
		echo *echo.Echo
		pr   personRepository
	}
	tests := []struct {
		name               string
		fields             fields
		expectedHTTPStatus int
		result             []models.Person
	}{
		{
			name: "http-200: persons found",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).GetAllPersonMock.Return([]models.Person{regularPerson}, nil),
			},
			expectedHTTPStatus: 200,
			result:             []models.Person{regularPerson},
		},
		{
			name: "http-500: database error",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).GetAllPersonMock.Return(nil, errors.New("database error")),
			},
			expectedHTTPStatus: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				pr:   tt.fields.pr,
			}

			r := httptest.NewRequest(http.MethodPost, "/test", nil)
			w := httptest.NewRecorder()
			c := s.echo.NewContext(r, w)

			err := s.getPersons(c)
			if err != nil {
				t.Errorf("getPersons() error = %v", err)
			}

			code := w.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("getPersons() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}

			if tt.expectedHTTPStatus == http.StatusOK {
				body, err := io.ReadAll(w.Result().Body)
				if err != nil {
					t.Errorf("ReadAll error")
				}
				var res []models.Person
				err = json.Unmarshal(body, &res)
				if err != nil {
					t.Errorf("json unmarshal error")
				}

				if !reflect.DeepEqual(res, tt.result) {
					t.Errorf("getPersons() expected %v, but got %v", tt.result, res)
				}
			}
		})
	}
}

func TestServer_updatePerson(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	type fields struct {
		echo *echo.Echo
		pr   personRepository
	}
	tests := []struct {
		name               string
		fields             fields
		pathParams         string
		expectedHTTPStatus int
		result             models.Person
		body               *models.Person
	}{
		{
			name: "http-200: success update",
			fields: fields{
				echo: e,
				pr: NewPersonRepositoryMock(mc).UpdatePersonByIDMock.Return(nil).
					GetPersonByIDMock.Return(regularPerson, nil),
			},
			pathParams:         "1",
			expectedHTTPStatus: 200,
			result:             regularPerson,
			body:               &regularPerson,
		},
		{
			name: "http-400: can not parse path param",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			pathParams:         "qwerty",
			expectedHTTPStatus: 400,
		},
		{
			name: "http-400: can not bind body",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			pathParams:         "1",
			expectedHTTPStatus: 400,
			body:               &models.Person{},
		},
		{
			name: "http-400: validation error",
			fields: fields{
				echo: e,
				pr:   nil,
			},
			pathParams:         "1",
			expectedHTTPStatus: 400,
			body:               &models.Person{Name: "test", Age: -10},
		},
		{
			name: "http-404: person not found",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).UpdatePersonByIDMock.Return(gorm.ErrRecordNotFound),
			},
			pathParams:         "1",
			expectedHTTPStatus: 404,
			body:               &regularPerson,
		},
		{
			name: "http-500: database error",
			fields: fields{
				echo: e,
				pr:   NewPersonRepositoryMock(mc).UpdatePersonByIDMock.Return(nil).GetPersonByIDMock.Return(models.Person{}, errors.New("database error")),
			},
			pathParams:         "1",
			expectedHTTPStatus: 500,
			body:               &regularPerson,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				pr:   tt.fields.pr,
			}

			var err error
			jsonData := []byte("")
			if tt.body != nil {
				jsonData, err = json.Marshal(tt.body)
				if err != nil {
					t.Errorf("json marshal error")
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(string(jsonData)))
			req.Header.Set("Content-type", "application/json")
			rw := httptest.NewRecorder()
			c := s.echo.NewContext(req, rw)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(tt.pathParams)

			err = s.updatePerson(c)
			if err != nil {
				t.Errorf("updatePerson() error = %v", err)
			}

			code := rw.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("updatePerson() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}

			body, err := io.ReadAll(rw.Result().Body)
			if err != nil {
				t.Errorf("ReadAll error")
			}
			var res models.Person
			err = json.Unmarshal(body, &res)
			if err != nil {
				t.Errorf("json unmarshal error")
			}

			if tt.expectedHTTPStatus == http.StatusOK {
				if !reflect.DeepEqual(res, tt.result) {
					t.Errorf("updatePerson() expected %v, but got %v", tt.result, res)
				}
			}
		})
	}
}
