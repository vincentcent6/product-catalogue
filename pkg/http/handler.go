package http

import productController "github.com/vincentcent6/product-catalogue/internal/controller/product"

var (
	prdCtrl productController.Controller
)

func Init() {
	if prdCtrl == nil {
		prdCtrl = productController.New()
	}
}
