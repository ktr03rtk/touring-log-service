package service

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"

	"github.com/pkg/errors"
)

type UserService interface {
	IsExists(model.Email) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{userRepository: ur}
}

func (s *userService) IsExists(email model.Email) (bool, error) {
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return false, errors.Wrapf(err, "failed to find user, email: %s", email)
	} else if u == nil {
		return false, nil
	}

	return true, nil
}
