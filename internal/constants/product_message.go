package constants

import "golectro-product/internal/model"

var (
	SuccessGetProducts = model.Message{
		"en": "Successfully retrieved products",
		"id": "Berhasil mendapatkan produk",
	}
	SuccessCreateProduct = model.Message{
		"en": "Successfully created product",
		"id": "Berhasil membuat produk",
	}
)

var (
	FailedGetProducts = model.Message{
		"en": "Failed to get products",
		"id": "Gagal mendapatkan produk",
	}
	FailedCreateProduct = model.Message{
		"en": "Failed to create product",
		"id": "Gagal membuat produk",
	}
	AccessDenied = model.Message{
		"en": "Access denied only for admin",
		"id": "Akses ditolak hanya untuk admin",
	}
)
