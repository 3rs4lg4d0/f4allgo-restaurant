CREATE TABLE restaurant (
    id     BIGSERIAL    PRIMARY KEY,
    name   VARCHAR(255) NOT NULL,
    city   VARCHAR(255) NOT NULL,
    state  VARCHAR(255) NOT NULL,
    street VARCHAR(255) NOT NULL,
    zip    VARCHAR(255) NOT NULL
);

CREATE TABLE menu_item (
	restaurant_id BIGINT         NOT NULL,
	id            INTEGER        NOT NULL,
	name          VARCHAR(255)   NOT NULL,
	price         VARCHAR(10) NOT NULL,
    PRIMARY KEY (restaurant_id, id)
);

ALTER TABLE menu_item ADD CONSTRAINT fk_restaurant_id FOREIGN KEY (restaurant_id) REFERENCES restaurant(id);
