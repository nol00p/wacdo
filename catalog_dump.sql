ALTER TABLE menu_products DISABLE TRIGGER ALL;
ALTER TABLE option_values DISABLE TRIGGER ALL;
ALTER TABLE product_options DISABLE TRIGGER ALL;
ALTER TABLE menus DISABLE TRIGGER ALL;
ALTER TABLE products DISABLE TRIGGER ALL;
ALTER TABLE categories DISABLE TRIGGER ALL;
TRUNCATE menu_products, option_values, product_options, menus, products, categories CASCADE;
--
-- PostgreSQL database dump
--

\restrict B3BKotaz9a1yKVyUT9yDg42OjDfAbMVR55d0Eeb2NmMq7OCEtnsOLUAIRD6IL3i

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.categories (id, name, description, display_order, image_url, created_at, updated_at) VALUES (2, 'Drinks', '', 0, '', '2026-02-23 18:04:29.631397+01', '2026-02-23 18:04:29.631397+01');
INSERT INTO public.categories (id, name, description, display_order, image_url, created_at, updated_at) VALUES (3, 'Pizza', 'Traditional Neapolitan wood-fired pizzas', 1, 'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=800', '2026-03-20 12:04:39.393261+01', '2026-03-20 12:04:39.393261+01');


--
-- Data for Name: menus; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (21, 2, 'Vino Rosso della Casa', 'House red wine, glass 15cl', 4.5, 200, true, 1, '2026-03-20 12:12:48.609421+01', '2026-03-20 12:12:48.609421+01', 'https://images.unsplash.com/photo-1510812431401-41d2bd2722f3?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (22, 2, 'Vino Bianco della Casa', 'House white wine, glass 15cl', 4.5, 200, true, 1, '2026-03-20 12:12:48.612968+01', '2026-03-20 12:12:48.612968+01', 'https://images.unsplash.com/photo-1566995541428-f2246c17cda1?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (3, 3, 'Margherita', 'San Marzano tomato sauce, fior di latte mozzarella, fresh basil, extra virgin olive oil', 15, 100, true, 8, '2026-03-20 12:05:08.966505+01', '2026-03-20 12:05:08.966505+01', 'https://images.unsplash.com/photo-1574071318508-1cdbab80d002?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (4, 3, 'Marinara', 'San Marzano tomato sauce, garlic, oregano, extra virgin olive oil', 14, 100, true, 7, '2026-03-20 12:05:08.995126+01', '2026-03-20 12:05:08.995126+01', 'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (5, 3, 'Margherita DOP', 'San Marzano DOP tomatoes, buffalo mozzarella DOP, fresh basil, extra virgin olive oil', 17, 100, true, 8, '2026-03-20 12:05:08.999968+01', '2026-03-20 12:05:08.999968+01', 'https://images.unsplash.com/photo-1595854341625-f33ee10dbf94?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (6, 3, 'Diavola', 'San Marzano tomato sauce, fior di latte mozzarella, spicy salame piccante, fresh basil', 17, 100, true, 9, '2026-03-20 12:05:09.004133+01', '2026-03-20 12:05:09.004133+01', 'https://images.unsplash.com/photo-1628840042765-356cda07504e?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (7, 3, 'Quattro Formaggi', 'Fior di latte mozzarella, gorgonzola, parmigiano reggiano, smoked provola', 18, 100, true, 10, '2026-03-20 12:05:09.008058+01', '2026-03-20 12:05:09.008058+01', 'https://images.unsplash.com/photo-1573821663912-569905455b1c?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (8, 3, 'Napoli', 'San Marzano tomato sauce, fior di latte mozzarella, anchovies, capers, oregano', 16.5, 100, true, 8, '2026-03-20 12:05:09.012217+01', '2026-03-20 12:05:09.012217+01', 'https://images.unsplash.com/photo-1544982503-9f984c14501a?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (9, 3, 'Capricciosa', 'San Marzano tomato sauce, fior di latte mozzarella, cooked ham, mushrooms, artichokes, olives', 18, 100, true, 10, '2026-03-20 12:05:09.016371+01', '2026-03-20 12:05:09.016371+01', 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (10, 3, 'Salsiccia e Friarielli', 'Fior di latte mozzarella, Neapolitan sausage, friarielli (broccoli rabe), smoked provola', 18.5, 100, true, 10, '2026-03-20 12:05:09.019973+01', '2026-03-20 12:05:09.019973+01', 'https://images.unsplash.com/photo-1571407970349-bc81e7e96d47?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (11, 3, 'Bufalina', 'San Marzano tomato sauce, buffalo mozzarella, cherry tomatoes, fresh basil', 18, 100, true, 9, '2026-03-20 12:05:09.023809+01', '2026-03-20 12:05:09.023809+01', 'https://images.unsplash.com/photo-1604382354936-07c5d9983bd3?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (12, 3, 'Quattro Stagioni', 'San Marzano tomato sauce, fior di latte mozzarella, ham, mushrooms, artichokes, olives in quartered sections', 18.5, 100, true, 10, '2026-03-20 12:05:09.028272+01', '2026-03-20 12:05:09.028272+01', 'https://images.unsplash.com/photo-1593560708920-61dd98c46a4e?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (13, 2, 'Coca-Cola', 'Classic Coca-Cola, served chilled', 3.5, 200, true, 1, '2026-03-20 12:12:48.557136+01', '2026-03-20 12:12:48.557136+01', 'https://images.unsplash.com/photo-1629203851122-3726ecdf080e?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (14, 2, 'Aranciata San Pellegrino', 'Sparkling Italian orange soda by San Pellegrino', 3.5, 200, true, 1, '2026-03-20 12:12:48.583874+01', '2026-03-20 12:12:48.583874+01', 'https://images.unsplash.com/photo-1625772299848-391b6a87d7b3?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (15, 2, 'Acqua Minerale', 'Still or sparkling natural mineral water, 50cl', 2, 200, true, 1, '2026-03-20 12:12:48.587681+01', '2026-03-20 12:12:48.587681+01', 'https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (16, 2, 'Birra Peroni', 'Italian lager beer, 33cl bottle', 5, 200, true, 1, '2026-03-20 12:12:48.591322+01', '2026-03-20 12:12:48.591322+01', 'https://images.unsplash.com/photo-1608270586620-248524c67de9?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (17, 2, 'Birra Moretti', 'Premium Italian lager, 33cl bottle', 5, 200, true, 1, '2026-03-20 12:12:48.594909+01', '2026-03-20 12:12:48.594909+01', 'https://images.unsplash.com/photo-1535958636474-b021ee887b13?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (18, 2, 'Limonata San Pellegrino', 'Sparkling Italian lemon soda by San Pellegrino', 3.5, 200, true, 1, '2026-03-20 12:12:48.598595+01', '2026-03-20 12:12:48.598595+01', 'https://images.unsplash.com/photo-1621263764928-df1444c5e859?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (19, 2, 'Espresso', 'Traditional Italian espresso, single shot', 2, 200, true, 2, '2026-03-20 12:12:48.602142+01', '2026-03-20 12:12:48.602142+01', 'https://images.unsplash.com/photo-1510707577719-ae7c14805e3a?w=800');
INSERT INTO public.products (id, category_id, name, description, price, stock_quantity, is_available, preparation_time, created_at, updated_at, image_url) VALUES (20, 2, 'Limoncello', 'Homemade lemon liqueur from the Amalfi coast, 4cl', 5.5, 200, true, 1, '2026-03-20 12:12:48.605828+01', '2026-03-20 12:12:48.605828+01', 'https://images.unsplash.com/photo-1560512823-829485b8bf24?w=800');


--
-- Data for Name: menu_products; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: product_options; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (3, 3, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (4, 3, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (5, 4, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (6, 4, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (7, 5, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (8, 5, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (9, 6, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (10, 6, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (11, 7, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (12, 7, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (13, 8, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (14, 8, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (15, 9, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (16, 9, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (17, 10, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (18, 10, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (19, 11, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (20, 11, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (21, 12, 'Taille', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (22, 12, 'Supplements', 'multiple', false);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (23, 15, 'Type', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (24, 21, 'Format', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (25, 22, 'Format', 'single', true);
INSERT INTO public.product_options (id, product_id, name, is_unique, is_required) VALUES (26, 19, 'Type', 'single', true);


--
-- Data for Name: option_values; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (7, 3, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (8, 3, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (9, 3, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (10, 4, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (11, 4, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (12, 4, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (13, 4, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (14, 4, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (15, 4, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (16, 4, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (17, 4, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (18, 4, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (19, 4, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (20, 4, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (21, 4, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (22, 4, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (23, 4, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (24, 4, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (25, 5, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (26, 5, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (27, 5, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (28, 6, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (29, 6, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (30, 6, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (31, 6, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (32, 6, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (33, 6, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (34, 6, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (35, 6, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (36, 6, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (37, 6, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (38, 6, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (39, 6, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (40, 6, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (41, 6, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (42, 6, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (43, 7, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (44, 7, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (45, 7, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (46, 8, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (47, 8, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (48, 8, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (49, 8, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (50, 8, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (51, 8, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (52, 8, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (53, 8, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (54, 8, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (55, 8, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (56, 8, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (57, 8, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (58, 8, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (59, 8, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (60, 8, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (61, 9, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (62, 9, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (63, 9, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (64, 10, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (65, 10, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (66, 10, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (67, 10, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (68, 10, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (69, 10, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (70, 10, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (71, 10, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (72, 10, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (73, 10, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (74, 10, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (75, 10, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (76, 10, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (77, 10, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (78, 10, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (79, 11, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (80, 11, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (81, 11, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (82, 12, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (83, 12, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (84, 12, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (85, 12, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (86, 12, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (87, 12, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (88, 12, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (89, 12, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (90, 12, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (91, 12, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (92, 12, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (93, 12, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (94, 12, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (95, 12, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (96, 12, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (97, 13, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (98, 13, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (99, 13, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (100, 14, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (101, 14, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (102, 14, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (103, 14, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (104, 14, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (105, 14, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (106, 14, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (107, 14, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (108, 14, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (109, 14, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (110, 14, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (111, 14, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (112, 14, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (113, 14, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (114, 14, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (115, 15, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (116, 15, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (117, 15, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (118, 16, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (119, 16, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (120, 16, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (121, 16, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (122, 16, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (123, 16, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (124, 16, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (125, 16, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (126, 16, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (127, 16, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (128, 16, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (129, 16, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (130, 16, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (131, 16, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (132, 16, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (133, 17, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (134, 17, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (135, 17, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (136, 18, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (137, 18, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (138, 18, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (139, 18, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (140, 18, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (141, 18, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (142, 18, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (143, 18, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (144, 18, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (145, 18, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (146, 18, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (147, 18, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (148, 18, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (149, 18, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (150, 18, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (151, 19, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (152, 19, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (153, 19, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (154, 20, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (155, 20, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (156, 20, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (157, 20, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (158, 20, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (159, 20, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (160, 20, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (161, 20, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (162, 20, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (163, 20, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (164, 20, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (165, 20, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (166, 20, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (167, 20, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (168, 20, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (169, 21, 'Piccola (25cm)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (170, 21, 'Media (30cm)', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (171, 21, 'Grande (36cm)', 4.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (172, 22, 'Mozzarella di bufala', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (173, 22, 'Prosciutto crudo', 2);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (174, 22, 'Salame piccante', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (175, 22, 'Funghi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (176, 22, 'Carciofi', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (177, 22, 'Olive nere', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (178, 22, 'Acciughe', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (179, 22, 'Capperi', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (180, 22, 'Peperoni', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (181, 22, 'Salsiccia', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (182, 22, 'Friarielli', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (183, 22, 'Pomodorini', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (184, 22, 'Provola affumicata', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (185, 22, 'Melanzane', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (186, 22, 'Parmigiano reggiano', 1.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (187, 23, 'Naturale (still)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (188, 23, 'Frizzante (sparkling)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (189, 24, 'Verre (15cl)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (190, 24, 'Demi (50cl)', 8);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (191, 24, 'Bouteille (75cl)', 14.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (192, 25, 'Verre (15cl)', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (193, 25, 'Demi (50cl)', 8);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (194, 25, 'Bouteille (75cl)', 14.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (195, 26, 'Espresso', 0);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (196, 26, 'Doppio', 1);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (197, 26, 'Macchiato', 0.5);
INSERT INTO public.option_values (id, option_id, value, option_price) VALUES (198, 26, 'Lungo', 0.5);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.categories_id_seq', 3, true);


--
-- Name: menu_products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.menu_products_id_seq', 1, false);


--
-- Name: menus_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.menus_id_seq', 1, false);


--
-- Name: option_values_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.option_values_id_seq', 198, true);


--
-- Name: product_options_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.product_options_id_seq', 26, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.products_id_seq', 22, true);


--
-- PostgreSQL database dump complete
--

\unrestrict B3BKotaz9a1yKVyUT9yDg42OjDfAbMVR55d0Eeb2NmMq7OCEtnsOLUAIRD6IL3i

ALTER TABLE menu_products ENABLE TRIGGER ALL;
ALTER TABLE option_values ENABLE TRIGGER ALL;
ALTER TABLE product_options ENABLE TRIGGER ALL;
ALTER TABLE menus ENABLE TRIGGER ALL;
ALTER TABLE products ENABLE TRIGGER ALL;
ALTER TABLE categories ENABLE TRIGGER ALL;
