package http

import (
	productController "github.com/vincentcent6/product-catalogue/internal/product/controller"
	productReviewController "github.com/vincentcent6/product-catalogue/internal/product_review/controller"
)

var (
	prdCtrl       productController.Controller
	prdReviewCtrl productReviewController.Controller
)

func Init() {
	if prdCtrl == nil {
		prdCtrl = productController.New()
	}
	if prdReviewCtrl == nil {
		prdReviewCtrl = productReviewController.New()
	}
}
