drop table if exists product cascade;
CREATE TABLE public.product (
    product_id bigserial NOT NULL,
    sku text NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    category text NOT NULL,
    etalase text NOT NULL, 
    images jsonb NOT NULL,
    weight float4 NOT NULL,
    price integer NOT NULL,
    create_time timestamp NOT NULL DEFAULT now(),
    CONSTRAINT product_pkey PRIMARY KEY (product_id)
);

CREATE TABLE public.product_review (
    product_id bigint NOT NULL,
    review_id bigserial NOT NULL,
    rating smallint NOT NULL,
    review_comment text NOT NULL,
    CONSTRAINT product_review_pkey PRIMARY KEY (review_id),
    CONSTRAINT product_review_fkey FOREIGN KEY (product_id) REFERENCES public.product (product_id)
);