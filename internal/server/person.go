package server

import (
	"errors"
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/models"
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (s *Server) getPersons(c echo.Context) error {
	persons, err := s.pr.GetAllPerson()
	if err != nil {
		log.Errorf("database error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "internal server error",
		})
	}
	return c.JSON(http.StatusOK, persons)
}

func (s *Server) createPerson(c echo.Context) error {
	var req models.Person
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("can not bind request: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "bad json request",
		})
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": err.Error(),
		})
	}

	req, err = s.pr.CreatePerson(req)
	if err != nil {
		log.Errorf("database error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"errors": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{})
}

func (s *Server) getPersonByID(c echo.Context) error {
	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		log.Errorf("can not parse id %v, err: %v", rawId, err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "bad request",
		})
	}

	person, err := s.pr.GetPersonByID(int32(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("person not found, id=%v", id)
			return c.JSON(http.StatusNotFound, echo.Map{})
		}
		log.Errorf("databese error %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{})
	}

	return c.JSON(http.StatusOK, person)
}

func (s *Server) updatePerson(c echo.Context) error {
	var req models.Person
	err := c.Bind(&req)
	if err != nil {
		log.Errorf("can not bind request: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": "bad json request",
		})
	}

	if err = c.Validate(req); err != nil {
		log.Errorf("validation error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"errors": err.Error(),
		})
	}

	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		log.Errorf("can not parse id %v, err: %v", rawId, err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "bad request",
		})
	}

	req.ID = int32(id)
	err = s.pr.UpdatePersonByID(req.ID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("person not found, id=%v", id)
			return c.JSON(http.StatusNotFound, echo.Map{})
		}
		log.Errorf("databese error %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{})
	}

	person, err := s.pr.GetPersonByID(int32(id))
	if err != nil {
		log.Errorf("databese error %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{})
	}

	return c.JSON(http.StatusOK, person)
}

func (s *Server) deletePersonByID(c echo.Context) error {
	rawId := c.Param("id")
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		log.Errorf("can not parse id %v, err: %v", rawId, err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "bad request",
		})
	}

	err = s.pr.DeletePersonByID(int32(id))
	if err != nil {
		log.Errorf("databese error %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{})
	}
	return c.JSON(http.StatusNoContent, echo.Map{})
}
