package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"test-task/internal/controller"
	"test-task/internal/domain/user"
	"test-task/pkg/errs"
)

type UseCase interface {
	Enrichment(fio user.User) error
}

type Kafka struct {
	useCase UseCase
	log     *slog.Logger
}

func NewUserKafka(useCase UseCase, log *slog.Logger) Kafka {
	return Kafka{useCase: useCase, log: log}
}

func (k Kafka) Start() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "my-topic",
	})

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			k.log.Warn("Не удалось прочитать сообщение", slog.String("errs", err.Error()))
			continue
		}

		var fio controller.User
		err = json.Unmarshal(message.Value, &fio)
		if err != nil {
			k.log.Warn("Не удалось unmarshal данных", slog.String("errs", err.Error()))
			continue
		}

		err = k.useCase.Enrichment(user.User{
			Name:       fio.Name,
			Surname:    fio.Surname,
			Patronymic: fio.Patronymic,
		})
		if errors.Is(err, errs.FioFailedErr) {
			// Отправить ответное сообщение "FIOFAILED" в Kafka
			err = k.sendResponseMessage("FIOFAILED")
			if err != nil {
				k.log.Warn("Не удалось отправить ответное сообщение", slog.String("errs", err.Error()))
			}
			continue
		}
		if err != nil {
			k.log.Warn("Не удалось обогатить пользователя или добавить в базу")
			continue
		}

		k.log.Info("Успешно обогащён и добавлен в базу пользователь")

		reader.CommitMessages(context.Background(), message)
	}

}

func (k Kafka) sendResponseMessage(response string) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "response-topic", // Замените на соответствующую тему для ответных сообщений
	})

	defer writer.Close()

	message := kafka.Message{
		Value: []byte(response),
	}

	err := writer.WriteMessages(context.Background(), message)
	if err != nil {
		k.log.Error(err.Error())

		return err
	}

	return nil
}
