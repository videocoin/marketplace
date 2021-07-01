-- +goose Up
-- SQL in this section is executed when the migration is applied.
INSERT INTO tokens (symbol, address, image_url, name, decimals, eth_price, usd_price) VALUES ('DAI', '0x5ca6ba24df599d26b0c30b620233f0a11f7556fa', 'https://lh3.googleusercontent.com/9wRB1clkMysBIYUd3SYnDLCVeZ1biqgAXLHtIYeWv_oG5XPgTYayh47mZclHw81O6SfJhl3ZhICNApiKfGSeMIc', 'Dai Stablecoin', 18, 0.000445601026665, 1.000000000);
INSERT INTO tokens (symbol, address, image_url, name, decimals, eth_price, usd_price) VALUES ('ETH', '0x0000000000000000000000000000000000000000', 'https://lh3.googleusercontent.com/7hQyiGtBt8vmUTq4T0aIUhIhT00dPhnav87TuFQ5cLtjlk724JgXdjQjoH_CzYz-z37JpPuMFbRRQuyC7I9abyZRKA', 'Ether', 18, 1.0, 2232.23);
INSERT INTO tokens (symbol, address, image_url, name, decimals, eth_price, usd_price) VALUES ('MANA', '0x5c635b8e736d894c488f83e52739660728cc2ef1', 'https://lh3.googleusercontent.com/hsrIlncQIGqCWc1qQ7CUAIuVsNFwzuSmr_dCEpbUYnO_YO0VWoGZzCxVlDVSgcxATeuPcCQPdOh3t5PhXj_c8gAQ', 'Decentraland MANA', 18, 0.000516124439382, 1.160000000000000000);
INSERT INTO tokens (symbol, address, image_url, name, decimals, eth_price, usd_price) VALUES ('WETH', '0x9834f449937e29d5a26465d1a9cda2fe234570ec', 'https://lh3.googleusercontent.com/kPzD9AH9Xt4Vr7NXphLy2Yyf5ZWM0vfN_wMhJs0HWJpQjFZm0pcmZ880vcJQVLxPgdnOTEfUuYbkiaGxcTT_ZnCy', 'Wrapped Ether', 18, 1.000000000000000, 2533.679999999999836000);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
TRUNCATE TABLE tokens;
