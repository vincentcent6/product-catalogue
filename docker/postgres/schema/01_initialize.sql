drop table if exists product cascade;
CREATE TABLE public.product (
    product_id bigserial NOT NULL,
    sku integer NOT NULL,
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