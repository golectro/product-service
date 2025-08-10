package constants

import "golectro-product/internal/model"

var (
	FailedInsertProductToElasticsearch = model.Message{
		"en": "Failed to insert product into Elasticsearch",
		"id": "Gagal memasukkan produk ke Elasticsearch",
	}
)
