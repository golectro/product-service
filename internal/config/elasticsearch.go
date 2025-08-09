package config

import (
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewElasticSearch(viper *viper.Viper, log *logrus.Logger) *elasticsearch.Client {
	indexesStr := viper.GetString("ELASTICSEARCH_INDEXES")
	if indexesStr == "" {
		log.Fatal("ELASTICSEARCH_INDEXES is required in environment variables")
	}

	indexes := strings.Split(indexesStr, ",")
	for i := range indexes {
		indexes[i] = strings.TrimSpace(indexes[i])
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			viper.GetString("ELASTICSEARCH_URL"),
		},
		Username: viper.GetString("ELASTICSEARCH_USERNAME"),
		Password: viper.GetString("ELASTICSEARCH_PASSWORD"),
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch client: %v", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Failed to ping Elasticsearch: %v", err)
	}
	defer res.Body.Close()
	log.Info("Elasticsearch connected")

	for _, index := range indexes {
		if index == "" {
			continue
		}

		existsRes, err := es.Indices.Exists([]string{index})
		if err != nil {
			log.Fatalf("Failed to check if index '%s' exists: %v", index, err)
		}
		defer existsRes.Body.Close()

		switch existsRes.StatusCode {
		case 404:
			createRes, err := es.Indices.Create(index)
			if err != nil {
				log.Fatalf("Failed to create index '%s': %v", index, err)
			}
			defer createRes.Body.Close()

			if createRes.IsError() {
				var e map[string]any
				_ = json.NewDecoder(createRes.Body).Decode(&e)
				log.Fatalf("Elasticsearch index creation error for '%s': %v", index, e)
			}

			log.Infof("Created new Elasticsearch index: %s", index)

		case 200:
			log.Infof("Elasticsearch index already exists: %s", index)

		default:
			log.Fatalf("Unexpected response checking index '%s': %s", index, existsRes.String())
		}
	}

	return es
}
