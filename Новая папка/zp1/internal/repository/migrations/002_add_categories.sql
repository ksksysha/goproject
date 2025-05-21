-- Добавляем таблицу категорий
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE
);

-- Добавляем категории
INSERT INTO categories (name, slug) VALUES
    ('Ногтевой сервис', 'nails'),
    ('Волосы', 'hair'),
    ('Ресницы', 'lashes'),
    ('Брови', 'brows')
ON CONFLICT (slug) DO NOTHING;

-- Добавляем колонку category_id в таблицу services
ALTER TABLE services ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

-- Обновляем существующие услуги, связывая их с категориями
UPDATE services SET category_id = (SELECT id FROM categories WHERE slug = 'nails')
WHERE name LIKE '%маникюр%' OR name LIKE '%педикюр%' OR name LIKE '%ногт%';

UPDATE services SET category_id = (SELECT id FROM categories WHERE slug = 'hair')
WHERE name LIKE '%волос%' OR name LIKE '%стриж%' OR name LIKE '%уклад%';

UPDATE services SET category_id = (SELECT id FROM categories WHERE slug = 'lashes')
WHERE name LIKE '%ресниц%';

UPDATE services SET category_id = (SELECT id FROM categories WHERE slug = 'brows')
WHERE name LIKE '%бров%';

-- Делаем category_id обязательным
ALTER TABLE services ALTER COLUMN category_id SET NOT NULL; 