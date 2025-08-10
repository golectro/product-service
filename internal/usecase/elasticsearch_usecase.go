package usecase

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ElasticsearchUseCase struct {
	Elasticsearch *elasticsearch.Client
	Validate      *validator.Validate
	Log           *logrus.Logger
	Viper         *viper.Viper
}

func NewElasticsearchUsecase(elasticsearch *elasticsearch.Client, log *logrus.Logger, validate *validator.Validate, viper *viper.Viper) *ElasticsearchUseCase {
	return &ElasticsearchUseCase{
		Elasticsearch: elasticsearch,
		Validate:      validate,
		Log:           log,
		Viper:         viper,
	}
}

func (e *ElasticsearchUseCase) InsertDocument(id uuid.UUID, document any) error {
	if err := e.Validate.Struct(document); err != nil {
		e.Log.WithError(err).Error("Invalid document structure")
		return err
	}

	jsonBody, err := json.Marshal(document)
	if err != nil {
		e.Log.WithError(err).Error("Failed to marshal document to JSON")
		return err
	}

	res, err := e.Elasticsearch.Index(e.Viper.GetString("ELASTICSEARCH_INDEX"), bytes.NewReader(jsonBody), e.Elasticsearch.Index.WithDocumentID(id.String()))
	if err != nil {
		e.Log.WithError(err).Error("Failed to index document in Elasticsearch")
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		e.Log.Errorf("Error response from Elasticsearch: %s", res.String())
		return err
	}

	return nil
}

func (e *ElasticsearchUseCase) DeleteDocumentByID(id string) error {
	res, err := e.Elasticsearch.Delete(e.Viper.GetString("ELASTICSEARCH_INDEX"), id)
	if err != nil {
		e.Log.WithError(err).Error("Failed to delete document from Elasticsearch")
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		e.Log.Errorf("Error response from Elasticsearch: %s", res.String())
		return err
	}

	return nil
}
