package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"test-task/internal/controller"
	"test-task/internal/domain/enrichment"
)

type UseCaseHandler interface {
	Enrichment(fio enrichment.Fio) error
	DeleteUser(id int) error
	ModifyUser() error
}

type Handler struct {
	useCase UseCaseHandler
	log     *slog.Logger
}

func NewEnrichmentHandler(useCase UseCaseHandler, log *slog.Logger) Handler {
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

		var fio controller.Fio
		if err := json.NewDecoder(r.Body).Decode(&fio); err != nil {
			uh.log.Error("не удалось получить данные с клиента", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		err := uh.useCase.Enrichment(enrichment.Fio{
			Name:       fio.Name,
			Surname:    fio.Surname,
			Patronymic: fio.Patronymic,
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

func (uh Handler) DeleteUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var fio controller.DeleteFio
		if err := json.NewDecoder(r.Body).Decode(&fio); err != nil {
			uh.log.Error("не удалось получить данные с клиента", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}
		//err := uh.useCase.DeleteUser()
		//if err != nil {
		//	uh.log.Error("Ошибка", slog.String("err", err.Error()))
		//	http.Error(w, "", http.StatusInternalServerError)
		//
		//	return
		//}

		uh.log.Info("Получилось")

		w.WriteHeader(http.StatusOK)
	}
}
func (uh Handler) ModifyUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var fio controller.ModifyFio
		if err := json.NewDecoder(r.Body).Decode(&fio); err != nil {
			uh.log.Error("не удалось получить данные с клиента", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		err := uh.useCase.ModifyUser()
		if err != nil {
			uh.log.Error("Ошибка", slog.String("err", err.Error()))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		uh.log.Info("Получилось")

		w.WriteHeader(http.StatusOK)
	}
}
