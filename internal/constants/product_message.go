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
	SuccessGetProductByID = model.Message{
		"en": "Successfully retrieved product by ID",
		"id": "Berhasil mendapatkan produk berdasarkan ID",
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
	FailedToCreateProduct = model.Message{
		"en": "Failed to create product",
		"id": "Gagal membuat produk",
	}
	NoFilesUploaded = model.Message{
		"en": "No files uploaded",
		"id": "Tidak ada file yang diunggah",
	}
	InvalidProductID = model.Message{
		"en": "Invalid product ID",
		"id": "ID produk tidak valid",
	}
	FailedGetProductByID = model.Message{
		"en": "Failed to get product by ID",
		"id": "Gagal mendapatkan produk berdasarkan ID",
	}
	ProductNotFound = model.Message{
		"en": "Product not found",
		"id": "Produk tidak ditemukan",
	}
	InvalidProductIDFormat = model.Message{
		"en": "Invalid product ID format",
		"id": "Format ID produk tidak valid",
	}
	FailedUploadProductImages = model.Message{
		"en": "Failed to upload product images",
		"id": "Gagal mengunggah gambar produk",
	}
)
