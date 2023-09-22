package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"test-task/internal/controller"
	"test-task/internal/domain/enrichment"
)

type UseCaseKafka interface {
	Enrichment(fio enrichment.Fio) error
}

type Kafka struct {
	useCase UseCaseKafka
	log     *slog.Logger
}

func NewEnrichmentKafka(useCase UseCaseKafka, log *slog.Logger) Kafka {
	fmt.Println("NewEnrichmentKafka")
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
			k.log.Warn("Не удалось прочитать сообщение", slog.String("err", err.Error()))
			continue
		}

		var fio controller.Fio
		err = json.Unmarshal(message.Value, &fio)
		if err != nil {
			k.log.Warn("Не удалось unmarshal данных", slog.String("err", err.Error()))
			continue
		}

		err = k.useCase.Enrichment(enrichment.Fio{
			Name:       fio.Name,
			Surname:    fio.Surname,
			Patronymic: fio.Patronymic,
		})
		if err != nil {
			k.log.Warn("Не удалось обогатить пользователя или добавить в базу")
			continue
		}

		k.log.Info("Успешно обогащён и добавлен в базу пользователь")

		reader.CommitMessages(context.Background(), message)
	}

}
