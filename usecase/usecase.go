package usecase

import (
	"github.com/subkhiyoga/auth-jwt/model"
	"github.com/subkhiyoga/auth-jwt/repository"
)

type LoginUsecase interface {
	Login(username string, password string) (*model.Credentials, error)
}

type loginUsecase struct {
	loginRepository repository.LoginRepository
}

func (u *loginUsecase) Login(username string, password string) (*model.Credentials, error) {
	mahasiswa, err := u.loginRepository.GetByUnameAndPassword(username, password)
	if err != nil {
		return nil, err
	}
	if mahasiswa == nil {
		return nil, nil
	}

	return mahasiswa, nil
}

func NewLoginUsecase(loginRepository repository.LoginRepository) LoginUsecase {
	return &loginUsecase{
		loginRepository: loginRepository,
	}
}
