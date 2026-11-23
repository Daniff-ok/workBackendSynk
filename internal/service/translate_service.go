package service

import (
	"rndmBackend/internal/repository"
	//"rndmBackend/internal/domain"
)

type TranslateService struct{
	repo *repository.TranslateRepo
}

func NewTranslateService(repo *repository.TranslateRepo) *TranslateService {
	return &TranslateService{repo}
}

func (s *TranslateService) GetTranslate(text []string, lang string) ([]string, error) {
	s.repo.RequestTranslate(text, lang)
	return []string{}, nil
}