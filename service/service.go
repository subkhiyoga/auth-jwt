package service

import (
	"github.com/subkhiyoga/auth-jwt/model"
	"github.com/subkhiyoga/auth-jwt/repo"
)

type LoginService interface {
	Login(username string, password string) (*model.Credentials, error)
}

type loginService struct {
	loginRepo repo.LoginRepo
}

func (s *loginService) Login(username string, password string) (*model.Credentials, error) {
	mahasiswa, err := s.loginRepo.GetByUnameAndPassword(username, password)
	if err != nil {
		return nil, err
	}
	if mahasiswa == nil {
		return nil, nil
	}

	return mahasiswa, nil
}

func NewLoginService(loginRepo repo.LoginRepo) LoginService {
	return &loginService{
		loginRepo: loginRepo,
	}
}
