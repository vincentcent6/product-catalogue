package core

const (
	// qCreateProduct is query to insert product
	qCreateProduct = `
		INSERT INTO public.product
		(sku, title, description, category, etalase, images, weight, price)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// qUpdateProduct is query to update product
	qUpdateProduct = `
		UPDATE public.product
		SET %v
		WHERE product_id = $1
		returning product_id
	`
)
