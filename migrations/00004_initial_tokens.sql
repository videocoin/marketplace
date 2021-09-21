-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO tokens (symbol, address, image_url, name, decimals, eth_price, usd_price) VALUES ('WETH', '0xc778417e063141139fce010982780140aa0cd5ab', 'https://lh3.googleusercontent.com/kPzD9AH9Xt4Vr7NXphLy2Yyf5ZWM0vfN_wMhJs0HWJpQjFZm0pcmZ880vcJQVLxPgdnOTEfUuYbkiaGxcTT_ZnCy', 'Wrapped Ether', 18, 1.000000000000000, 2533.679999999999836000);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
TRUNCATE TABLE tokens;
