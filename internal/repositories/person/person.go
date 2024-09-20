package person

import (
	"fmt"
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/models"
	"gorm.io/gorm"
)

const personTable = "persons"

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *storage {
	return &storage{db: db}
}

func (s *storage) GetAllPerson() ([]models.Person, error) {
	var persons []models.Person
	err := s.db.Table(personTable).Find(&persons).Error
	if err != nil {
		return nil, fmt.Errorf("error getting all persons: %w", err)
	}

	return persons, nil
}

func (s *storage) CreatePerson(person models.Person) (models.Person, error) {
	err := s.db.Table(personTable).Create(&person).Error
	if err != nil {
		return models.Person{}, fmt.Errorf("error creating person: %w", err)
	}
	return person, nil
}

func (s *storage) GetPersonByID(id int32) (models.Person, error) {
	var person models.Person
	err := s.db.Table(personTable).Where("id = ?", id).Take(&person).Error
	if err != nil {
		return models.Person{}, fmt.Errorf("error getting person by id: %w", err)
	}
	return person, nil
}

func (s *storage) DeletePersonByID(id int32) error {
	err := s.db.Table(personTable).Where("id = ?", id).Delete(&models.Person{}).Error
	if err != nil {
		return fmt.Errorf("error deleting person: %w", err)
	}
	return nil
}

func (s *storage) UpdatePersonByID(id int32, person models.Person) error {
	err := s.db.Table(personTable).Where("id = ?", id).Updates(&person).Error
	if err != nil {
		return fmt.Errorf("error updating person: %w", err)
	}
	return nil
}
