package controller

import (
	"context"

	"github.com/vincentcent6/product-catalogue/internal/core"
)

func (p *ctrl) CreateProduct(ctx context.Context, input ProductInput) (result ProductInput, err error) {
	if err = input.ValidateComplete(); err != nil {
		return ProductInput{}, err
	}
	if err = prd.CreateProduct(ctx, convertData(input)); err != nil {
		return ProductInput{}, err
	}
	return input, nil
}

func convertData(input ProductInput) core.Data {
	images := []core.ProductImage{}
	for _, img := range input.Images {
		images = append(images, core.ProductImage{
			URL:         img.URL,
			Description: img.Description,
		})
	}

	return core.Data{
		SKU:         input.SKU,
		Title:       input.Title,
		Description: input.Description,
		Category:    input.Category,
		Etalase:     input.Etalase,
		Images:      images,
		Weight:      input.Weight,
		Price:       input.Price,
	}
}
