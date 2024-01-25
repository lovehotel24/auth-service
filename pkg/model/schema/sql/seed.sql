INSERT INTO users (user_id, name, phone, role, password_hash, date_created, date_updated)
VALUES ('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', 'ADMIN',
        '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2024-01-24 00:00:00',
        '2024-01-25 00:00:00'),
       ('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Guest Gopher', 'guest@example.com', 'USER',
        '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2024-01-24 00:00:00',
        '2024-01-25 00:00:00') ON CONFLICT DO NOTHING;