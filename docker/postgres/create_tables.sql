CREATE TABLE FB_AD_INSIGHTS_PRODUCT (
    id SERIAL PRIMARY KEY,

    spend DOUBLE PRECISION DEFAULT 0.0,
    clicks BIGINT DEFAULT 0,
    impressions BIGINT DEFAULT 0,
    product_id BIGINT DEFAULT 0,
    date_start DATE DEFAULT NULL,

    CONSTRAINT metrics_product_date_unique
    UNIQUE (product_id, date_start)
);