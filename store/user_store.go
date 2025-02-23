package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/mkhorosh/go-app/models"
)

type FileUserStore struct {
	filePath string
	mu       sync.Mutex
}

func NewFileUserStore(filePath string) *FileUserStore {
	return &FileUserStore{filePath: filePath}
}

func (s *FileUserStore) SaveUser(user *models.User) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	users = append(users, user)

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println("Writing to file:", s.filePath)

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func (s *FileUserStore) GetUsers() ([]*models.User, error) {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	file, err := os.ReadFile(s.filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var users []*models.User
	if len(file) > 0 {
		err = json.Unmarshal(file, &users)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}
