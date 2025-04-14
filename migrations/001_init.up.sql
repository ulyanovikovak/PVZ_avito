CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Пользователи (схемка на будущее)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE,
    password TEXT,
    role TEXT NOT NULL CHECK (role IN ('employee', 'moderator'))
);

-- Пункты выдачи заказов
CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL,
    city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

-- Прием
CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY,
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    date_time TIMESTAMP NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

-- Товары
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE
);
