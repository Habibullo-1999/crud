CREATE TABLE customers 
(
  id        BIGSERIAL PRIMARY Key,
  name      TEXT      NOT NULL,
  phone     TEXT      NOT NULL UNIQUE,
  password  TEXT      NOT NULL ,
  active    BOOLEAN   NOT NULL DEFAULT TRUE,
  created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE managers 
(
  id        BIGSERIAL PRIMARY KEY,
  name      TEXT      NOT NULL,
  salary    INTEGER   NOT NULL CHECK ( salary > 0 ),
  plan      INTEGER   NOT NULL CHECK ( salary > 0 ),
  boss_id   BIGINT    REFERENCES managers,
  departament TEXT,
  login     TEXT      NOT NULL UNIQUE,
  password  TEXT      NOT NULL, 
  active    BOOLEAN   NOT NULL DEFAULT TRUE,
  created   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers_tokens 
(
  token TEXT NOT NULL UNIQUE,
  customer_id BIGINT NOT NULL REFERENCES customers,
  expire TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
  created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)
ALTER TABLE customers
ADD active BOOLEAN NOT NULL DEFAULT TRUE


INSERT INTO customers_tokens(token,customer_id) VALUES('111',1)

-- DROP TABLE customers_tokens;
-- DROP TABLE managers;
-- DROP TABLE customers;

CREATE TABLE products 
(
    id      BIGSERIAL   PRIMARY KEY,
    name    TEXT        NOT NULL,
    price   INTEGER     NOT NULL CHECK (price > 0),
    qty     INTEGER     NOT NULL DEFAULT 0 CHECK (qty >=0),
    active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
)

CREATE TABLE users 
( 
   id       BIGSERIAL     PRIMARY KEY,
   name     TEXT          NOT NULL,
   phone    TEXT          NOT NULL UNIQUE,
   password TEXT          NOT NULL,
   roles    TEXT[]        NOT NULL DEFAULT '{}',
   active  BOOLEAN     NOT NULL DEFAULT TRUE,
   created TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
)

	SELECT id,name, price,qty FROM products WHERE active 


CREATE TABLE sales 
( 
    id           BIGSERIAL  PRIMARY KEY,
    manager_id   BIGINT     NOT NULL REFERENCES managers,
    customer_id  BIGINT     REFERENCES customers,
    created      TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sale_positions 
( 
    id           BIGSERIAL  PRIMARY KEY,
    sale_id      BIGINT     NOT NULL REFERENCES sales,
    product_id   BIGINT     NOT NULL REFERENCES products,
    name         TEXT       NOT NULL,
    price        INTEGER    NOT NULL CHECK ( price >= 0 ),
    qty          INTEGER    NOT NULL DEFAULT 1 CHECK ( qty > 0 ),
    created      TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
)

SELECT sp.id, sp.name, sp.price,sp.qty,sp.created 
FROM sale_positions sp 
JOIN sales s on s.id = sp.sale_id
WHERE s.customer_id = 1;

SELECT * FROM products


INSERT INTO sale_positions (
    sale_id,
    product_id,
    name,
    price,
    qty
  )
VALUES (
    '2',
    '1',
    'cofee',
    100,
    10
  );
  INsert into sale_positions(sale_id, product_id, name, price, qty) VALUES (1,1,'cofee',100,10) 


INsert into sales(manager_id, customer_id) VALUES (1,1)

INSERT INTO managers (
    name,
    salary,
    plan,
    boss_id,
    departament,
    login,
    password
  )
VALUES (
    'habib',
    10000,
    10000,
    1,
    1200,
    'habib',
    'h@b1b'
  );