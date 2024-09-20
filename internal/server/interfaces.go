package server

import "github.com/AskaryanKarine/BMSTU-ds-1/internal/models"

//go:generate minimock -o mocks_storage.go -g
type personRepository interface {
	GetAllPerson() ([]models.Person, error)
	CreatePerson(person models.Person) (models.Person, error)
	GetPersonByID(id int32) (models.Person, error)
	DeletePersonByID(id int32) error
	UpdatePersonByID(id int32, person models.Person) error
}
