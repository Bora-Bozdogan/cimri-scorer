CREATE TABLE if not exists api_keys (
    key TEXT PRIMARY KEY,
    is_active BOOLEAN DEFAULT true,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO api_keys (key) VALUES ('test-key');

create table if not exists products (
    id serial primary key,
    name        TEXT NOT NULL,
    description TEXT DEFAULT 'No description has been provided',
    image_url   TEXT DEFAULT 'No image provided',
    popularity_score INT NOT NULL,
    urgency_score INT NOT NULL
);

create table if not exists merchants (
    id serial primary key,
    name text not null,
    key text not null,
    is_active BOOLEAN DEFAULT true
);

create table if not exists merchant_products (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    merchant_id int NOT NULL references merchants(id) on delete cascade,
    price INT NOT NULL,
    stock INT NOT NULL
);

with new_merchant as (
    insert into merchants
        (name, key, is_active)
        VALUES ('Amazon', 'amazon-key', true)
        returning id
),
new_product as (
insert into products
    (name, popularity_score, urgency_score)
    values ('IPhone 16', 5, 5)
    returning id
)
insert into merchant_products (product_id, merchant_id, price, stock)
       values (
        (select id from new_product),
        (select id from new_merchant),
               1500,
                50
              );

