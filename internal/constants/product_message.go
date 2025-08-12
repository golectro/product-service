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
	SuccessDeleteProductImage = model.Message{
		"en": "Successfully deleted product image",
		"id": "Berhasil menghapus gambar produk",
	}
	SuccessSearchProducts = model.Message{
		"en": "Successfully searched products",
		"id": "Berhasil mencari produk",
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
	FailedGetProductsByIDs = model.Message{
		"en": "Failed to get products by IDs",
		"id": "Gagal mendapatkan produk berdasarkan ID",
	}
	FailedDecreaseProductQuantity = model.Message{
		"en": "Failed to decrease product quantity",
		"id": "Gagal mengurangi jumlah produk",
	}
	InvalidSearchRequest = model.Message{
		"en": "Invalid search request",
		"id": "Permintaan pencarian tidak valid",
	}
	FailedSearchProducts = model.Message{
		"en": "Failed to search products",
		"id": "Gagal mencari produk",
	}
	ProductNotFound = model.Message{
		"en": "Product not found",
		"id": "Produk tidak ditemukan",
	}
	InsufficientProductQuantity = model.Message{
		"en": "Insufficient product quantity",
		"id": "Jumlah produk tidak mencukupi",
	}
	InvalidProductIDFormat = model.Message{
		"en": "Invalid product ID format",
		"id": "Format ID produk tidak valid",
	}
	FailedUploadProductImages = model.Message{
		"en": "Failed to upload product images",
		"id": "Gagal mengunggah gambar produk",
	}
	FailedGetPresignedURL = model.Message{
		"en": "Failed to get presigned URL",
		"id": "Gagal mendapatkan URL presigned",
	}
	FailedGetImageByID = model.Message{
		"en": "Failed to get image by ID",
		"id": "Gagal mendapatkan gambar berdasarkan ID",
	}
	FailedGetProductImages = model.Message{
		"en": "Failed to get product images",
		"id": "Gagal mendapatkan gambar produk",
	}
	ImageNotFound = model.Message{
		"en": "Image not found",
		"id": "Gambar tidak ditemukan",
	}
	FailedUpdateProduct = model.Message{
		"en": "Failed to update product",
		"id": "Gagal memperbarui produk",
	}
	FailedDeleteProduct = model.Message{
		"en": "Failed to delete product",
		"id": "Gagal menghapus produk",
	}
	FailedDeleteProductImage = model.Message{
		"en": "Failed to delete product image",
		"id": "Gagal menghapus gambar produk",
	}
	FailedDeleteImage = model.Message{
		"en": "Failed to delete image",
		"id": "Gagal menghapus gambar",
	}
	FailedCommitTransaction = model.Message{
		"en": "Failed to commit transaction",
		"id": "Gagal mengkomit transaksi",
	}
	FailedDeleteImageFromMinio = model.Message{
		"en": "Failed to delete image from Minio",
		"id": "Gagal menghapus gambar dari Minio",
	}
)
