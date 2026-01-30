-- Test User Seeder
-- Email: test@example.com
-- Password: password123
-- PIN: 123456
-- Initial Balance: 1,000,000,000 (1 billion in smallest unit)

INSERT INTO users (id, email, password_hash, pin_hash, created_at) 
VALUES (
    1,
    'test@example.com',
    '$2a$10$.hTZRP3v7nvbglq66sN5Z.YyA5rp9xxg.wMlNziIcBllRb6c6f5eW',
    '$2a$10$AncRTiktYRGln5GwXlzcW.e7Pyhy.yF4IJuJWhtGwoZzS2G0zD47C',
    NOW()
) ON DUPLICATE KEY UPDATE 
    email = VALUES(email),
    password_hash = VALUES(password_hash),
    pin_hash = VALUES(pin_hash);

INSERT INTO wallets (id, user_id, balance, created_at)
VALUES (
    1,
    1,
    1000000000,
    NOW()
) ON DUPLICATE KEY UPDATE
    balance = VALUES(balance);
