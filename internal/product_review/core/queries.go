package core

const (
	// qCreateProductReview is query to insert product review
	qCreateProductReview = `
		INSERT INTO public.product_review
		(product_id, rating, review_comment)
		VALUES($1, $2, $3)
		RETURNING review_id
	`
)
