package user

import (
	"errors"
	"log/slog"
)

type Service interface {
	Validation(user *User, enrichmentUser *EnrichmentUser) bool
	EnrichmentAPI(user User) (EnrichmentUser, error)
	PutInDatabase(enrichmentUser EnrichmentUser) error
	CheckUserExist(id int) (bool, error)
	DeleteUser(id int) error
	ModifyUser(enrichmentUser EnrichmentUser) error
}

type enrichmentUseCase struct {
	service Service
	log     *slog.Logger
}

func NewUserUseCase(service Service, log *slog.Logger) *enrichmentUseCase {
	return &enrichmentUseCase{service: service, log: log}
}

func (uc enrichmentUseCase) Enrichment(user User) error {
	if uc.service.Validation(&user, nil) == false {
		uc.log.Info("Неправильные данные")
		//TODO: ошибка нужна
		return errors.New("required field is empty")
	}

	enrichmentUser, err := uc.service.EnrichmentAPI(user)
	if err != nil {
		uc.log.Error("ошибка при обогащении пользователя", slog.String("err", err.Error()))

		return err
	}
	uc.log.Info("Обогащение успешно прошло", slog.Any("user", enrichmentUser))

	err = uc.service.PutInDatabase(enrichmentUser)
	if err != nil {
		uc.log.Error("ошибка при добавлении пользователя", slog.String("err", err.Error()))

		return err
	}
	uc.log.Info("Обогащение пользователя и добавление его в базу прошло успешно")

	return nil
}

func (uc enrichmentUseCase) DeleteUser(id int) error {
	if exist, err := uc.service.CheckUserExist(id); err == nil {
		if !exist {
			uc.log.Warn("пользователь не существует")
			//TODO: вернуть ошибку отсутсвия юзера
			return nil
		}
	} else {
		uc.log.Error("ошибка при проверке пользователя на существование", slog.String("err", err.Error()))

		return err
	}

	err := uc.service.DeleteUser(id)
	if err != nil {
		uc.log.Error("не удалось удалить пользователя",
			slog.String("err", err.Error()), slog.Int("user_id", id))

		return err
	}

	uc.log.Info("пользователь успешно удалён", slog.Int("user_id", id))

	return nil
}

func (uc enrichmentUseCase) ModifyUser(enrichmentUser EnrichmentUser) error {
	if exist, err := uc.service.CheckUserExist(enrichmentUser.Id); err == nil {
		if !exist {
			uc.log.Warn("пользователь не существует")
			//TODO: вернуть ошибку отсутсвия юзера
			return nil
		}
	} else {
		uc.log.Error("ошибка при проверке пользователя на существование", slog.String("err", err.Error()))

		return err
	}

	if uc.service.Validation(nil, &enrichmentUser) == false {
		uc.log.Info("Неправильные данные")
		//TODO: ошибка нужна
		return errors.New("required field is empty")
	}

	err := uc.service.ModifyUser(enrichmentUser)
	if err != nil {
		uc.log.Error("не удалось изменить пользователя",
			slog.String("err", err.Error()), slog.Int("user_id", enrichmentUser.Id))
		return err
	}

	uc.log.Info("пользователь успешно изменён", slog.Int("user_id", enrichmentUser.Id))

	return nil
}
