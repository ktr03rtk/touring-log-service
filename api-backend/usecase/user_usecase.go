package usecase

import (
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"

	"github.com/ktr03rtk/touring-log-service/api-backend/domain/model"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/service"

	"github.com/pkg/errors"
)

type UserUsecase interface {
	SignUp(email, password, unit string) error
}

type userUsecase struct {
	userRepository repository.UserRepository
	userService    service.UserService
}

func NewUserUsecase(ur repository.UserRepository, us service.UserService) UserUsecase {
	return &userUsecase{
		userRepository: ur,
		userService:    us,
	}
}

func (u *userUsecase) SignUp(email, password, unit string) error {
	// TODO: check unit uniqueness
	ok, err := u.userService.IsExists(model.Email(email))
	if ok {
		return errors.Errorf("already registered email. email: %s", email)
	} else if err != nil {
		return err
	}

	id := model.CreateUUID()

	user, err := model.NewUser(model.UserID(id), model.Email(email), password, unit)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	if err := u.userRepository.Create(user); err != nil {
		return errors.Wrap(err, "failed to store user")
	}

	return nil
}
