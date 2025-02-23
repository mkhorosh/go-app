package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mkhorosh/go-app/models"
)

type AnimalStore struct {
	filePath string
}

func NewAnimalStore(filePath string) *AnimalStore {
	return &AnimalStore{filePath: filePath}
}

func (s *AnimalStore) GetAnimals() ([]models.Animal, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var animals []models.Animal
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&animals); err != nil {
		return nil, err
	}

	return animals, nil
}

func (s *AnimalStore) SaveAnimals(animals []models.Animal) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(animals)
}

func (s *AnimalStore) AddAnimal(animal *models.Animal) error {
	animals, err := s.GetAnimals()
	if err != nil {
		return err
	}

	animals = append(animals, *animal)
	return s.SaveAnimals(animals)
}

func (s *AnimalStore) UpdateAnimal(id string, updatedAnimal *models.Animal) error {
	animals, err := s.GetAnimals()
	if err != nil {
		return err
	}

	for i, animal := range animals {
		if animal.ID == id {
			animals[i] = *updatedAnimal
			return s.SaveAnimals(animals)
		}
	}
	return fmt.Errorf("animal with ID %s not found", id)
}

func (s *AnimalStore) DeleteAnimal(id string) error {
	animals, err := s.GetAnimals()
	if err != nil {
		return err
	}

	for i, animal := range animals {
		if animal.ID == id {
			animals = append(animals[:i], animals[i+1:]...)
			return s.SaveAnimals(animals)
		}
	}
	return fmt.Errorf("animal with ID %s not found", id)
}
