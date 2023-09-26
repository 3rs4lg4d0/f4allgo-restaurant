CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO restaurant (id, name, city, state, street, zip) VALUES (1000, 'restaurant1', 'city1', 'state1', 'street1', 'zip1');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (1000, 1, 'item1.1', '13.14');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (1000, 2, 'item1.2', '14.15');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (1000, 3, 'item1.3', '15.16');

INSERT INTO restaurant (id, name, city, state, street, zip) VALUES (2000, 'restaurant2', 'city2', 'state2', 'street2', 'zip2');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (2000, 1, 'item2.1', '13.14');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (2000, 2, 'item2.2', '14.15');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (2000, 3, 'item2.3', '15.16');

INSERT INTO restaurant (id, name, city, state, street, zip) VALUES (3000, 'restaurant3', 'city3', 'state3', 'street3', 'zip3');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (3000, 1, 'item3.1', '13.14');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (3000, 2, 'item3.2', '14.15');
INSERT INTO menu_item (restaurant_id, id, name, price) VALUES (3000, 3, 'item3.3', '15.16');

INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload) VALUES (uuid_generate_v4(), 'restaurant', '1', 'RestaurantCreated', E'\\xDEADBEEF');
INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload) VALUES (uuid_generate_v4(), 'restaurant', '2', 'RestaurantUpdated', E'\\xDEADBEEF');
INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload) VALUES (uuid_generate_v4(), 'restaurant', '3', 'RestaurantDeleted', E'\\xDEADBEEF');
INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload) VALUES (uuid_generate_v4(), 'restaurant', '4', 'RestaurantCreated', E'\\xDEADBEEF');
INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload) VALUES (uuid_generate_v4(), 'restaurant', '5', 'RestaurantCreated', E'\\xDEADBEEF');
