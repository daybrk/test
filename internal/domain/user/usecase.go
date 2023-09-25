package user

import (
	"log/slog"
	"test-task/pkg/errs"
)

type Service interface {
	Validation(user *User, enrichmentUser *EnrichmentUser) bool
	EnrichmentAPI(user User) (EnrichmentUser, error)
	PutInDatabase(enrichmentUser EnrichmentUser) error
	CheckUserExist(id int) (bool, error)
	DeleteUser(id int) error
	ModifyUser(enrichmentUser EnrichmentUser) error
	GetFilteredUsers(filter Filter) ([]EnrichmentUser, error)
}

type userUseCase struct {
	service Service
	log     *slog.Logger
}

func NewUserUseCase(service Service, log *slog.Logger) *userUseCase {
	return &userUseCase{service: service, log: log}
}

func (uc userUseCase) Enrichment(user User) error {
	if uc.service.Validation(&user, nil) == false {
		uc.log.Info("Неправильные данные")

		return errs.FioFailedErr
	}

	enrichmentUser, err := uc.service.EnrichmentAPI(user)
	if err != nil {
		uc.log.Error("ошибка при обогащении пользователя", slog.String("errs", err.Error()))

		return err
	}
	uc.log.Info("Обогащение успешно прошло", slog.Any("user", enrichmentUser))

	err = uc.service.PutInDatabase(enrichmentUser)
	if err != nil {
		uc.log.Error("ошибка при добавлении пользователя", slog.String("errs", err.Error()))

		return err
	}
	uc.log.Info("Обогащение пользователя и добавление его в базу прошло успешно")

	return nil
}

func (uc userUseCase) DeleteUser(id int) error {
	if exist, err := uc.service.CheckUserExist(id); err == nil {
		if !exist {
			uc.log.Warn("пользователь не существует")

			return errs.UserNotFoundErr
		}
	} else {
		uc.log.Error("ошибка при проверке пользователя на существование", slog.String("errs", err.Error()))

		return err
	}

	err := uc.service.DeleteUser(id)
	if err != nil {
		uc.log.Error("не удалось удалить пользователя",
			slog.String("errs", err.Error()), slog.Int("user_id", id))

		return err
	}

	uc.log.Info("пользователь успешно удалён", slog.Int("user_id", id))

	return nil
}

func (uc userUseCase) ModifyUser(enrichmentUser EnrichmentUser) error {
	if exist, err := uc.service.CheckUserExist(enrichmentUser.Id); err == nil {
		if !exist {
			uc.log.Warn("пользователь не существует")

			return errs.UserNotFoundErr
		}
	} else {
		uc.log.Error("ошибка при проверке пользователя на существование", slog.String("errs", err.Error()))

		return err
	}

	if uc.service.Validation(nil, &enrichmentUser) == false {
		uc.log.Warn("Неправильные данные")

		return errs.FioFailedErr
	}

	err := uc.service.ModifyUser(enrichmentUser)
	if err != nil {
		uc.log.Error("не удалось изменить пользователя",
			slog.String("errs", err.Error()), slog.Int("user_id", enrichmentUser.Id))
		return err
	}

	uc.log.Info("пользователь успешно изменён", slog.Int("user_id", enrichmentUser.Id))

	return nil
}

func (uc userUseCase) Filter(filter Filter) ([]EnrichmentUser, error) {
	users, err := uc.service.GetFilteredUsers(filter)
	if err != nil {
		uc.log.Error("ошибка при попытке взять пользователей с использованием фильтров",
			slog.String("err", err.Error()))

		return nil, err
	}

	return users, nil
}
