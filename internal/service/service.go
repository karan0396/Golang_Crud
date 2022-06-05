package service


import (
	"api/internal/model"
	"api/pkg/token"
	"api/pkg/logger"

	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo model.Repository
}

type NewService interface {
	LoginService(data model.Credential) (string, error)
	CreateService(data model.User) error
	GetService(p model.Params) ([]model.User, error)
	DeleteService(id int) error
	UpdateService(data model.User) error
}

func NewServ(repo model.Repository) NewService {
	return &service{repo}
}

func init() {
	logger.IntializeLogger()
	
}

type Signinresponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func (s service) LoginService(data model.Credential) (string, error) {

	password, err := s.repo.Getpasswordbyemail(data)
	if err != nil {
		logger.Logger.DPanic("Email not found", zap.Error(err))
		return "", err
	}
	err = CheckPasswordHash(data.Password, password)
	if err != nil {
		logger.Logger.DPanic("Password not match", zap.Error(err))
		return "", err
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	tknstr, error := token.CreateToken(data.Email, expirationTime)
	if error != nil {
		logger.Logger.DPanic("Token not created", zap.Error(err))
		return "", err
	}

	return tknstr, nil
}

func (s service) CreateService(data model.User) error {
	err := s.repo.Create(&data)
	if err != nil {
		logger.Logger.DPanic("data not added", zap.Error(err))
	}
	return err
}

func (s service) GetService(p model.Params) ([]model.User, error) {
	newuser, err := s.repo.Get(p)
	if err != nil {
		logger.Logger.DPanic("it does not get data", zap.Error(err))
	}
	return newuser, nil
}

func (s service) DeleteService(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		logger.Logger.DPanic("data not deleted", zap.Error(err))
	}
	return nil
}

func (s service) UpdateService(data model.User) error {
	err := s.repo.Update(data)
	if err != nil {
		logger.Logger.DPanic("data not updated", zap.Error(err))
	}
	return nil
}
