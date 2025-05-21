-- Добавляем поле slug в таблицу services
ALTER TABLE services ADD COLUMN IF NOT EXISTS slug VARCHAR(100) UNIQUE;

-- Обновляем существующие услуги, добавляя slug
UPDATE services SET slug = 'manicure_without' WHERE name LIKE '%маникюр без покрытия%';
UPDATE services SET slug = 'manicure_with' WHERE name LIKE '%маникюр с покрытием%';
UPDATE services SET slug = 'nail_extensions' WHERE name LIKE '%наращивание%';
UPDATE services SET slug = 'pedicure_without' WHERE name LIKE '%педикюр без покрытия%';
UPDATE services SET slug = 'pedicure_with' WHERE name LIKE '%педикюр с покрытием%';
UPDATE services SET slug = 'hand_care' WHERE name LIKE '%уход за кожей рук%';
UPDATE services SET slug = 'haircuts' WHERE name LIKE '%стриж%';
UPDATE services SET slug = 'styling' WHERE name LIKE '%уклад%';
UPDATE services SET slug = 'hair_care' WHERE name LIKE '%уход за волос%';
UPDATE services SET slug = 'curling' WHERE name LIKE '%биозавив%';
UPDATE services SET slug = 'eyelash_extensions' WHERE name LIKE '%наращивание ресниц%';
UPDATE services SET slug = 'lamination_lashes' WHERE name LIKE '%ламинирование ресниц%';
UPDATE services SET slug = 'lash_correction' WHERE name LIKE '%коррекция ресниц%';
UPDATE services SET slug = 'lash_coloring' WHERE name LIKE '%окрашивание ресниц%';
UPDATE services SET slug = 'brow_correction' WHERE name LIKE '%коррекция бров%';
UPDATE services SET slug = 'brow_tinting' WHERE name LIKE '%окрашивание бров%';
UPDATE services SET slug = 'brow_lamination' WHERE name LIKE '%ламинирование бров%';
UPDATE services SET slug = 'microblading' WHERE name LIKE '%микроблейдинг%';

-- Делаем поле slug обязательным
ALTER TABLE services ALTER COLUMN slug SET NOT NULL; 