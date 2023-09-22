package enrichment

import (
	"fmt"
	"log/slog"
	"regexp"
	"test-task/internal/adapters/web"
)

type Storage interface {
	Insert() error
}

type Router interface {
	EnrichmentAge(name string) (web.Age, error)
	EnrichmentGender(name string) (web.Gender, error)
	EnrichmentNationality(name string) (web.Nationality, error)
}

type enrichmentService struct {
	storage Storage
	router  Router
	log     *slog.Logger
}

func NewEnrichmentService(storage Storage, router Router, log *slog.Logger) *enrichmentService {
	fmt.Println("NewEnrichmentService")
	return &enrichmentService{storage: storage, router: router, log: log}
}

func (e enrichmentService) Validation(fio Fio) bool {
	if fio.Name == "" || hasCyrillic(fio.Name) {
		return false
	}

	if fio.Surname == "" || hasCyrillic(fio.Surname) {
		return false
	}

	return true
}

func (e enrichmentService) EnrichmentAPI(fio Fio) (FioEnrichment, error) {
	var enrichmentFIO FioEnrichment
	enrichmentFIO.Name = fio.Name
	enrichmentFIO.Surname = fio.Surname
	enrichmentFIO.Patronymic = fio.Patronymic

	age, err := e.router.EnrichmentAge(fio.Name)
	if err != nil {
		return FioEnrichment{}, err
	}
	enrichmentFIO.Age = age.Age
	e.log.Info("обогащение возрастом",
		slog.Any("age", age.Age), slog.Any("после обогащения", enrichmentFIO))

	gender, err := e.router.EnrichmentGender(fio.Name)
	if err != nil {
		return FioEnrichment{}, err
	}
	enrichmentFIO.Gender = gender.Gender
	e.log.Info("обогащение полом",
		slog.Any("gender", gender.Gender), slog.Any("после обогащения", enrichmentFIO))

	nationality, err := e.router.EnrichmentNationality(fio.Name)
	if err != nil {
		return FioEnrichment{}, err
	}
	for _, value := range nationality.Country {
		enrichmentFIO.Nationality = append(enrichmentFIO.Nationality, value.CountryId)
	}
	e.log.Info("обогащение национальностью",
		slog.Any("nationality", nationality.Country), slog.Any("после обогащения", enrichmentFIO))

	return enrichmentFIO, nil
}

func (e enrichmentService) PutInDatabase(fio FioEnrichment) error {
	return nil
}

func hasCyrillic(input string) bool {
	re := regexp.MustCompile("[\u0400-\u04FF]+")

	return re.MatchString(input)
}
