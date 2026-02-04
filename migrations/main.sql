-- ====================================================================
-- ТЕСТОВАЯ БАЗА ДАННЫХ: ПРОДАЖА ТОВАРОВ
-- ====================================================================

-- Удаляем таблицы если существуют
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS customers CASCADE;

-- ====================================================================
-- СОЗДАНИЕ ТАБЛИЦ
-- ====================================================================

-- Таблица категорий товаров
CREATE TABLE categories (
                            category_id SERIAL PRIMARY KEY,
                            category_name VARCHAR(100) NOT NULL,
                            description TEXT
);

-- Таблица товаров
CREATE TABLE products (
                          product_id SERIAL PRIMARY KEY,
                          product_name VARCHAR(200) NOT NULL,
                          category_id INTEGER REFERENCES categories(category_id),
                          price NUMERIC(10, 2) NOT NULL,
                          cost NUMERIC(10, 2) NOT NULL,
                          stock_quantity INTEGER DEFAULT 0
);

-- Таблица покупателей
CREATE TABLE customers (
                           customer_id SERIAL PRIMARY KEY,
                           first_name VARCHAR(100) NOT NULL,
                           last_name VARCHAR(100) NOT NULL,
                           email VARCHAR(150) UNIQUE,
                           phone VARCHAR(20),
                           city VARCHAR(100),
                           registration_date DATE DEFAULT CURRENT_DATE
);

-- Таблица заказов
CREATE TABLE orders (
                        order_id SERIAL PRIMARY KEY,
                        customer_id INTEGER REFERENCES customers(customer_id),
                        order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        status VARCHAR(50) DEFAULT 'completed',
                        total_amount NUMERIC(12, 2),
                        payment_method VARCHAR(50)
);

-- Таблица позиций в заказах
CREATE TABLE order_items (
                             order_item_id SERIAL PRIMARY KEY,
                             order_id INTEGER REFERENCES orders(order_id),
                             product_id INTEGER REFERENCES products(product_id),
                             quantity INTEGER NOT NULL,
                             price NUMERIC(10, 2) NOT NULL,
                             discount NUMERIC(5, 2) DEFAULT 0
);

-- ====================================================================
-- ЗАПОЛНЕНИЕ ТЕСТОВЫМИ ДАННЫМИ
-- ====================================================================

-- Категории
INSERT INTO categories (category_name, description) VALUES
                                                        ('Электроника', 'Электронные устройства и гаджеты'),
                                                        ('Одежда', 'Мужская и женская одежда'),
                                                        ('Продукты питания', 'Продукты и напитки'),
                                                        ('Книги', 'Художественная и техническая литература'),
                                                        ('Спорт и отдых', 'Спортивные товары и туристическое снаряжение'),
                                                        ('Дом и сад', 'Товары для дома и садоводства');

-- Товары
INSERT INTO products (product_name, category_id, price, cost, stock_quantity) VALUES
-- Электроника
('Смартфон Samsung Galaxy S23', 1, 65000.00, 50000.00, 45),
('Ноутбук Lenovo ThinkPad', 1, 85000.00, 65000.00, 20),
('Наушники Sony WH-1000XM5', 1, 28000.00, 20000.00, 60),
('Планшет iPad Air', 1, 55000.00, 42000.00, 30),
('Умные часы Apple Watch', 1, 35000.00, 27000.00, 40),

-- Одежда
('Джинсы Levis 501', 2, 6500.00, 3500.00, 100),
('Куртка зимняя North Face', 2, 15000.00, 9000.00, 45),
('Футболка Nike', 2, 2500.00, 1200.00, 150),
('Кроссовки Adidas Ultraboost', 2, 12000.00, 7000.00, 80),
('Платье вечернее', 2, 8500.00, 4500.00, 35),

-- Продукты питания
('Кофе Lavazza 1кг', 3, 1800.00, 1200.00, 200),
('Шоколад Lindt 100г', 3, 450.00, 250.00, 300),
('Оливковое масло 1л', 3, 850.00, 500.00, 120),
('Чай зеленый 100г', 3, 350.00, 200.00, 250),
('Мед натуральный 500г', 3, 650.00, 400.00, 100),

-- Книги
('Мастер и Маргарита', 4, 550.00, 300.00, 80),
('SQL. Сборник рецептов', 4, 2500.00, 1500.00, 40),
('Python для анализа данных', 4, 3200.00, 2000.00, 50),
('Атлас мира', 4, 1800.00, 1000.00, 30),
('1984 Джордж Оруэлл', 4, 480.00, 250.00, 90),

-- Спорт и отдых
('Велосипед горный', 5, 35000.00, 25000.00, 15),
('Палатка туристическая 4-местная', 5, 12000.00, 8000.00, 25),
('Коврик для йоги', 5, 1500.00, 800.00, 70),
('Гантели 10кг пара', 5, 3500.00, 2000.00, 40),
('Рюкзак туристический 60л', 5, 8500.00, 5500.00, 35),

-- Дом и сад
('Пылесос Dyson V15', 6, 45000.00, 35000.00, 20),
('Кофеварка Delonghi', 6, 18000.00, 12000.00, 30),
('Набор посуды 12 предметов', 6, 5500.00, 3500.00, 50),
('Лейка садовая 10л', 6, 650.00, 350.00, 80),
('Секатор профессиональный', 6, 1800.00, 1000.00, 60);

-- Покупатели
INSERT INTO customers (first_name, last_name, email, phone, city, registration_date) VALUES
                                                                                         ('Иван', 'Иванов', 'ivan.ivanov@mail.ru', '+79161234567', 'Москва', '2023-01-15'),
                                                                                         ('Мария', 'Петрова', 'maria.petrova@gmail.com', '+79162345678', 'Санкт-Петербург', '2023-02-20'),
                                                                                         ('Алексей', 'Сидоров', 'alex.sidorov@yandex.ru', '+79163456789', 'Москва', '2023-03-10'),
                                                                                         ('Елена', 'Смирнова', 'elena.smirnova@mail.ru', '+79164567890', 'Казань', '2023-04-05'),
                                                                                         ('Дмитрий', 'Кузнецов', 'dmitry.kuznetsov@gmail.com', '+79165678901', 'Новосибирск', '2023-05-12'),
                                                                                         ('Ольга', 'Попова', 'olga.popova@yandex.ru', '+79166789012', 'Екатеринбург', '2023-06-18'),
                                                                                         ('Сергей', 'Волков', 'sergey.volkov@mail.ru', '+79167890123', 'Москва', '2023-07-22'),
                                                                                         ('Анна', 'Соколова', 'anna.sokolova@gmail.com', '+79168901234', 'Краснодар', '2023-08-30'),
                                                                                         ('Павел', 'Лебедев', 'pavel.lebedev@yandex.ru', '+79169012345', 'Челябинск', '2023-09-14'),
                                                                                         ('Наталья', 'Козлова', 'natalia.kozlova@mail.ru', '+79160123456', 'Ростов-на-Дону', '2023-10-08'),
                                                                                         ('Андрей', 'Новиков', 'andrey.novikov@gmail.com', '+79161234560', 'Самара', '2023-11-25'),
                                                                                         ('Татьяна', 'Морозова', 'tatiana.morozova@yandex.ru', '+79162345671', 'Омск', '2023-12-03'),
                                                                                         ('Максим', 'Васильев', 'maxim.vasilev@mail.ru', '+79163456782', 'Воронеж', '2024-01-17'),
                                                                                         ('Юлия', 'Зайцева', 'julia.zaitseva@gmail.com', '+79164567893', 'Пермь', '2024-02-28'),
                                                                                         ('Владимир', 'Федоров', 'vladimir.fedorov@yandex.ru', '+79165678904', 'Волгоград', '2024-03-15'),
                                                                                         ('Екатерина', 'Михайлова', 'ekaterina.mikhailova@mail.ru', '+79166789015', 'Саратов', '2024-04-22'),
                                                                                         ('Роман', 'Александров', 'roman.alexandrov@gmail.com', '+79167890126', 'Тюмень', '2024-05-30'),
                                                                                         ('Светлана', 'Егорова', 'svetlana.egorova@yandex.ru', '+79168901237', 'Ижевск', '2024-06-12'),
                                                                                         ('Игорь', 'Семенов', 'igor.semenov@mail.ru', '+79169012348', 'Уфа', '2024-07-19'),
                                                                                         ('Виктория', 'Титова', 'victoria.titova@gmail.com', '+79160123459', 'Ярославль', '2024-08-25');

-- Заказы и позиции заказов (генерируем разнообразные данные за период с января 2024 по январь 2025)
-- Используем функцию для генерации случайных дат

DO $$
DECLARE
v_order_id INTEGER;
    v_customer_id INTEGER;
    v_order_date TIMESTAMP;
    v_product_id INTEGER;
    v_quantity INTEGER;
    v_price NUMERIC(10,2);
    v_discount NUMERIC(5,2);
    v_total NUMERIC(12,2);
    v_payment_methods TEXT[] := ARRAY['Наличные', 'Карта', 'Онлайн перевод', 'Электронный кошелек'];
    i INTEGER;
    j INTEGER;
    items_count INTEGER;
BEGIN
    -- Генерируем 500 заказов
FOR i IN 1..500 LOOP
        -- Случайный покупатель
        v_customer_id := (random() * 19 + 1)::INTEGER;

        -- Случайная дата в диапазоне 2024-01-01 до 2025-01-30
        v_order_date := '2024-01-01'::TIMESTAMP + (random() * 395)::INTEGER * INTERVAL '1 day'
                        + (random() * 23)::INTEGER * INTERVAL '1 hour'
                        + (random() * 59)::INTEGER * INTERVAL '1 minute';

        -- Создаем заказ
INSERT INTO orders (customer_id, order_date, status, payment_method)
VALUES (
           v_customer_id,
           v_order_date,
           CASE WHEN random() < 0.95 THEN 'completed' ELSE 'cancelled' END,
           v_payment_methods[(random() * 3 + 1)::INTEGER]
       )
    RETURNING order_id INTO v_order_id;

-- Добавляем от 1 до 5 товаров в заказ
items_count := (random() * 4 + 1)::INTEGER;
        v_total := 0;

FOR j IN 1..items_count LOOP
            -- Случайный товар
            v_product_id := (random() * 29 + 1)::INTEGER;

            -- Количество от 1 до 5
            v_quantity := (random() * 4 + 1)::INTEGER;

            -- Получаем цену товара
SELECT price INTO v_price FROM products WHERE product_id = v_product_id;

-- Случайная скидка (0%, 5%, 10%, 15%)
v_discount := CASE
                WHEN random() < 0.7 THEN 0
                WHEN random() < 0.85 THEN 5
                WHEN random() < 0.95 THEN 10
                ELSE 15
END;

            -- Добавляем позицию в заказ
INSERT INTO order_items (order_id, product_id, quantity, price, discount)
VALUES (v_order_id, v_product_id, v_quantity, v_price, v_discount);

-- Считаем сумму
v_total := v_total + (v_price * v_quantity * (100 - v_discount) / 100);
END LOOP;

        -- Обновляем общую сумму заказа
UPDATE orders SET total_amount = v_total WHERE order_id = v_order_id;
END LOOP;
END $$;

-- ====================================================================
-- СОЗДАНИЕ ИНДЕКСОВ для оптимизации запросов
-- ====================================================================

CREATE INDEX idx_orders_date ON orders(order_date);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_product ON order_items(product_id);
CREATE INDEX idx_products_category ON products(category_id);

-- ====================================================================
-- ПОЛЕЗНЫЕ ПРЕДСТАВЛЕНИЯ (VIEWS)
-- ====================================================================

-- Детальная информация о продажах
CREATE VIEW sales_detailed AS
SELECT
    o.order_id,
    o.order_date,
    o.status,
    o.payment_method,
    c.customer_id,
    c.first_name || ' ' || c.last_name AS customer_name,
    c.city,
    p.product_id,
    p.product_name,
    cat.category_name,
    oi.quantity,
    oi.price,
    oi.discount,
    (oi.price * oi.quantity * (100 - oi.discount) / 100) AS item_total,
    (p.price - p.cost) * oi.quantity AS profit
FROM orders o
         JOIN customers c ON o.customer_id = c.customer_id
         JOIN order_items oi ON o.order_id = oi.order_id
         JOIN products p ON oi.product_id = p.product_id
         JOIN categories cat ON p.category_id = cat.category_id
WHERE o.status = 'completed';

-- ====================================================================
-- ПРИМЕРЫ ЗАПРОСОВ ДЛЯ ПРАКТИКИ
-- ====================================================================

-- Проверка данных
SELECT 'Всего категорий: ' || COUNT(*) FROM categories
UNION ALL
SELECT 'Всего товаров: ' || COUNT(*) FROM products
UNION ALL
SELECT 'Всего покупателей: ' || COUNT(*) FROM customers
UNION ALL
SELECT 'Всего заказов: ' || COUNT(*) FROM orders
UNION ALL
SELECT 'Всего позиций: ' || COUNT(*) FROM order_items;

COMMENT ON TABLE categories IS 'Категории товаров';
COMMENT ON TABLE products IS 'Товары с ценами и остатками';
COMMENT ON TABLE customers IS 'Покупатели';
COMMENT ON TABLE orders IS 'Заказы покупателей';
COMMENT ON TABLE order_items IS 'Позиции в заказах';
COMMENT ON VIEW sales_detailed IS 'Детальная информация о продажах с расчетными полями';