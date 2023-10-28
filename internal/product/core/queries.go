package core

const (
	// qCreateProduct is query to insert product
	qCreateProduct = `
		INSERT INTO public.product
		(sku, title, description, category, etalase, images, weight, price)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING product_id, create_time
	`

	// qUpdateProduct is query to update product
	qUpdateProduct = `
		UPDATE public.product
		SET %v
		WHERE product_id = $1
		returning product_id
	`

	// qGetProduct is query to get product
	qGetProduct = `
		SELECT product_id, sku, title, description, category, etalase, images, weight, price
		FROM public.product
		WHERE %v
	`
)
