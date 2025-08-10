package constants

import "golectro-product/internal/model"

var (
	SuccessUpdateProduct = model.Message{
		"en": "Successfully updated product",
		"id": "Berhasil memperbarui produk",
	}
	SuccessDeleteProduct = model.Message{
		"en": "Successfully deleted product",
		"id": "Berhasil menghapus produk",
	}
)

var (
	WelcomeMessage = model.Message{
		"en": "Welcome to Golectro User API",
		"id": "Selamat datang di Golectro User API",
	}
	NotFound = model.Message{
		"en": "API not found",
		"id": "API tidak ditemukan",
	}
	FailedDataFromBody = model.Message{
		"en": "Failed to get data from body",
		"id": "Gagal mendapatkan data dari body",
	}
	FailedInputFormat = model.Message{
		"en": "Invalid input format",
		"id": "Format input tidak valid",
	}
	FailedValidationOccurred = model.Message{
		"en": "Validation error occurred",
		"id": "Terjadi kesalahan validasi",
	}
	InvalidToken = model.Message{
		"en": "Invalid token",
		"id": "Token tidak valid",
	}
	InvalidRequestData = model.Message{
		"en": "Invalid request data",
		"id": "Data permintaan tidak valid",
	}
	InternalServerError = model.Message{
		"en": "Internal server error",
		"id": "Kesalahan internal server",
	}
	TooManyRequests = model.Message{
		"en": "Too many requests, please try again later",
		"id": "Terlalu banyak permintaan, silakan coba lagi nanti",
	}
)
