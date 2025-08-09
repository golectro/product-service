package constants

import "golectro-product/internal/model"

var (
	SuccessGetProducts = model.Message{
		"en": "Successfully retrieved products",
		"id": "Berhasil mendapatkan produk",
	}
)

var (
	FailedGetProducts = model.Message{
		"en": "Failed to get products",
		"id": "Gagal mendapatkan produk",
	}
)
