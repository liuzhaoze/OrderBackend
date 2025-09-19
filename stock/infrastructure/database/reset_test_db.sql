USE stock_test;

DROP TABLE IF EXISTS items;

CREATE TABLE items (
    item_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    price_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO items (item_id, name, quantity, price_id) VALUES
    ('prod_SNjRcpjxpiazxk', 'Pencil', 100, 'price_1RSxyuPSQHt2xYB8XhSJRSVX'),
    ('prod_SNjQpQjNC8QuaD', 'Book', 200, 'price_1RSxy2PSQHt2xYB8uZzd0XQx');
