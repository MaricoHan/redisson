package services

import (
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/dto"
)

type IRights interface {
	Register(params *dto.RegisterRequest) (*dto.RegisterResponse, error)
	EditRegister(params *dto.RegisterRequest) (*dto.RegisterResponse, error)
	QueryRegister(params *dto.QueryRegisterRequest) (*dto.QueryRegisterResponse, error)
}

type Rights struct {
	logger *log.Entry
}

func NewRights(logger *log.Logger) *Rights {
	return &Rights{
		logger: logger.WithField("service", "rights_jiangsu"),
	}
}

func (r Rights) Register(params *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	//logger := r.logger.WithField("params", params).WithField("func", "Register")
	return nil, nil
}

func (r Rights) EditRegister(params *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	//logger := r.logger.WithField("params", params).WithField("func", "EditRegister")
	return nil, nil
}

func (r Rights) QueryRegister(params *dto.QueryRegisterRequest) (*dto.QueryRegisterResponse, error) {
	//logger := r.logger.WithField("params", params).WithField("func", "QueryRegister")
	return nil, nil
}
