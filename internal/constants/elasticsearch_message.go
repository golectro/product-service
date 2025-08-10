package constants

import "golectro-product/internal/model"

var (
	FailedInsertProductToElasticsearch = model.Message{
		"en": "Failed to insert product into Elasticsearch",
		"id": "Gagal memasukkan produk ke Elasticsearch",
	}
	FailedDeleteProductFromElasticsearch = model.Message{
		"en": "Failed to delete product from Elasticsearch",
		"id": "Gagal menghapus produk dari Elasticsearch",
	}
	FailedUpdateProductInElasticsearch = model.Message{
		"en": "Failed to update product in Elasticsearch",
		"id": "Gagal memperbarui produk di Elasticsearch",
	}
	FailedDeleteImageFromElasticsearch = model.Message{
		"en": "Failed to delete image from Elasticsearch",
		"id": "Gagal menghapus gambar dari Elasticsearch",
	}
)
