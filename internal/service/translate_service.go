package service

import (
	"rndmBackend/internal/gateway"
	//"rndmBackend/internal/domain"
)

type TranslateService struct{
	repo *gateway.TranslateGateway
}

func NewTranslateService(repo *gateway.TranslateGateway) *TranslateService {
	return &TranslateService{repo}
}

func (s *TranslateService) GetTranslate(text []string, lang string) ([]string, error) {
	s.repo.RequestTranslate(text, lang)
	return []string{}, nil
}