-- Сначала удаляем существующего админа если есть
DELETE FROM users WHERE email = 'admin@example.com';

-- Создаем нового админа с простым паролем: admin
INSERT INTO users (username, email, password_hash, is_admin, created_at, updated_at)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',  -- пароль: admin
    true,
    NOW(),
    NOW()
); 