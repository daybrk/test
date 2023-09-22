package enrichment

import (
	"errors"
	"log/slog"
)

type Service interface {
	Validation(fio Fio) bool
	EnrichmentAPI(fio Fio) (FioEnrichment, error)
	PutInDatabase(fio FioEnrichment) error
}

type enrichmentUseCase struct {
	service Service
	log     *slog.Logger
}

func NewEnrichmentUseCase(service Service, log *slog.Logger) *enrichmentUseCase {
	return &enrichmentUseCase{service: service, log: log}
}

func (uc enrichmentUseCase) Enrichment(fio Fio) error {
	check := uc.service.Validation(fio)
	if !check {
		uc.log.Info("Неправильные данные")
		return errors.New("required field is empty")
	}

	enrichment, err := uc.service.EnrichmentAPI(fio)
	if err != nil {
		return err
	}
	uc.log.Info("Обогащение успешно прошло", slog.Any("user", enrichment))

	err = uc.service.PutInDatabase(enrichment)
	if err != nil {
		return err
	}
	uc.log.Info("Добавление обогащённых данных прошло успешно")

	return nil
}

func (uc enrichmentUseCase) DeleteUser(id int) error {
	panic("implement me")
}

func (uc enrichmentUseCase) ModifyUser() error {
	panic("implement me")
}
