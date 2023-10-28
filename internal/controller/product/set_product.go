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

func (p *ctrl) UpdateProduct(ctx context.Context, input UpdateProductInput) (err error) {
	if err = input.Data.ValidateComplete(); err != nil {
		return err
	}
	if err = prd.UpdateProduct(ctx, convertUpdateData(input)); err != nil {
		return err
	}
	return nil
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

func convertUpdateData(input UpdateProductInput) core.UpdateInput {
	return core.UpdateInput{
		ProductID: input.ProductID,
		Data:      convertData(input.Data),
	}
}
