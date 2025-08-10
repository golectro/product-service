package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"

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

func (e *ElasticsearchUseCase) SearchProducts(query map[string]any) ([]map[string]any, int64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		e.Log.WithError(err).Error("Failed to encode search query")
		return nil, 0, err
	}

	res, err := e.Elasticsearch.Search(
		e.Elasticsearch.Search.WithIndex(e.Viper.GetString("ELASTICSEARCH_INDEX")),
		e.Elasticsearch.Search.WithBody(&buf),
		e.Elasticsearch.Search.WithTrackTotalHits(true),
		e.Elasticsearch.Search.WithPretty(),
	)
	if err != nil {
		e.Log.WithError(err).Error("Failed to execute search query")
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		e.Log.Errorf("Elasticsearch search error: %s", res.String())
		return nil, 0, fmt.Errorf("search request failed: %s", res.String())
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		e.Log.WithError(err).Error("Failed to decode search response")
		return nil, 0, err
	}

	hitsData, ok := result["hits"].(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("unexpected hits format")
	}

	totalObj, ok := hitsData["total"].(map[string]any)
	if !ok {
		return nil, 0, fmt.Errorf("unexpected total format")
	}
	totalValue := int64(totalObj["value"].(float64))

	var products []map[string]any
	if hitsArray, ok := hitsData["hits"].([]any); ok {
		for _, h := range hitsArray {
			if hit, ok := h.(map[string]any); ok {
				if src, ok := hit["_source"].(map[string]any); ok {
					products = append(products, src)
				}
			}
		}
	}

	return products, totalValue, nil
}
