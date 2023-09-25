package graphQL

import (
	"context"
	"log/slog"
	"test-task/internal/domain/user"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type UseCase interface {
	Enrichment(fio user.User) error
	DeleteUser(id int) error
	ModifyUser(enrichmentFio user.EnrichmentUser) error
	Filter(filter user.Filter) ([]user.EnrichmentUser, error)
}

type Resolver struct {
	useCase UseCase
	log     *slog.Logger
}

func NewUserResolver(useCase UseCase, log *slog.Logger) *Resolver {
	return &Resolver{useCase: useCase, log: log}
}

func (r *mutationResolver) AddUser(ctx context.Context, input User) (*Result, error) {
	r.log.Info("Начато обогащение и добавление пользователя в базу", slog.Any("user", input))

	err := r.useCase.Enrichment(user.User{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	})
	if err != nil {
		r.log.Error("Ошибка", slog.String("errs", err.Error()))

		return nil, nil
	}

	return &Result{
		Success: true,
		Message: nil,
		Error:   nil,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input DeleteUser) (*Result, error) {
	if input.ID == 0 {
		r.log.Error("ошибка в отправляемых данных id = 0")

		return nil, nil
	}

	err := r.useCase.DeleteUser(input.ID)
	if err != nil {
		r.log.Error("ошибка при удалении пользователя", slog.String("errs", err.Error()))

		return nil, err
	}

	return &Result{
		Success: true,
		Message: nil,
		Error:   nil,
	}, nil
}

// ModifyUser is the resolver for the modifyUser field.
func (r *mutationResolver) ModifyUser(ctx context.Context, input ModifyUser) (*Result, error) {
	if input.ID == 0 {
		r.log.Error("ошибка в отправляемых данных id = 0")

		return nil, nil
	}

	err := r.useCase.ModifyUser(user.EnrichmentUser{
		Id:          input.ID,
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         input.Age,
		Gender:      input.Gender,
		Nationality: input.Nationality,
	})
	if err != nil {
		r.log.Error("Ошибка", slog.String("errs", err.Error()))

		return nil, err
	}

	r.log.Info("Получилось")

	return &Result{
		Success: true,
		Message: nil,
		Error:   nil,
	}, nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, input Filter) ([]*FilteredUsers, error) {
	result, err := r.useCase.Filter(user.Filter{
		Name:        input.Name,
		Surname:     input.Surname,
		Patronymic:  input.Patronymic,
		Age:         input.Age,
		Gender:      input.Gender,
		Nationality: input.Nationality,
	})
	if err != nil {
		r.log.Error("Ошибка", slog.String("errs", err.Error()))

		return nil, err
	}

	var users []*FilteredUsers
	for _, value := range result {
		users = append(users, &FilteredUsers{
			Name:        value.Name,
			Surname:     value.Surname,
			Patronymic:  value.Patronymic,
			Age:         value.Age,
			Gender:      value.Gender,
			Nationality: value.Nationality,
		})
	}

	r.log.Info("Получилось")

	return users, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
