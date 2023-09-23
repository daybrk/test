package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"test-task/internal/controller"
	"test-task/internal/domain/user"
)

type UseCaseHandler interface {
	Enrichment(fio user.User) error
	DeleteUser(id int) error
	ModifyUser(enrichmentFio user.EnrichmentUser) error
}

type Handler struct {
	useCase UseCaseHandler
	log     *slog.Logger
}

func NewUserHandler(useCase UseCaseHandler, log *slog.Logger) Handler {
	return Handler{useCase: useCase, log: log}
}

func (uh Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/add-user", uh.AddUser())
	mux.HandleFunc("/delete-user", uh.DeleteUser())
	mux.HandleFunc("/modify-user", uh.ModifyUser())
	mux.HandleFunc("/get-users-with-filter", uh.AddUser())
}

func (uh Handler) AddUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var fio controller.User
		if err := json.NewDecoder(r.Body).Decode(&fio); err != nil {
			uh.log.Error("не удалось получить данные с клиента", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		uh.log.Info("Начато обогащение и добавление пользователя в базу", slog.Any("user", fio))

		err := uh.useCase.Enrichment(user.User{
			Name:       fio.Name,
			Surname:    fio.Surname,
			Patronymic: fio.Patronymic,
		})
		if err != nil {
			uh.log.Error("Ошибка", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		uh.log.Info("обогащение и добавление пользователя в базу закончилось")

		w.WriteHeader(http.StatusOK)
	}
}

func (uh Handler) DeleteUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			uh.log.Error("ошибка при конвертации id в int", slog.String("err", err.Error()))
		}

		if id == 0 {
			uh.log.Error("ошибка в отправляемых данных id = 0")
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		err = uh.useCase.DeleteUser(id)
		if err != nil {
			uh.log.Error("ошибка при удалении пользователя", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		uh.log.Info("Получилось")

		w.WriteHeader(http.StatusOK)
	}
}

// TODO
func (uh Handler) ModifyUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var fio controller.ModifyUser
		if err := json.NewDecoder(r.Body).Decode(&fio); err != nil {
			uh.log.Error("не удалось получить данные с клиента", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		if fio.Id == 0 {
			uh.log.Error("ошибка в отправляемых данных id = 0")
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		err := uh.useCase.ModifyUser(user.EnrichmentUser{
			Id:          fio.Id,
			Name:        fio.Name,
			Surname:     fio.Surname,
			Patronymic:  fio.Patronymic,
			Age:         fio.Age,
			Gender:      fio.Gender,
			Nationality: fio.Nationality,
		})
		if err != nil {
			uh.log.Error("Ошибка", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		uh.log.Info("Получилось")

		w.WriteHeader(http.StatusOK)
	}
}
