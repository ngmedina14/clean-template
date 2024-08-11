-- +goose Up

INSERT INTO users (id, name, username, email, password, created_at, updated_at) VALUES
(1, 'User 1', 'User1', 'User1@gmail.com', 'abcdefgh123', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
(2, 'User 2', 'User2', 'User2@gmail.com', 'abcdefgh123', '2024-01-02 00:00:00', '2024-01-02 00:00:00'),
(3, 'User 3', 'User3', 'User3@gmail.com', 'abcdefgh123', '2024-01-03 00:00:00', '2024-01-03 00:00:00'),
(4, 'User 4', 'User4', 'User4@gmail.com', 'abcdefgh123', '2024-01-04 00:00:00', '2024-01-04 00:00:00'),
(5, 'User 5', 'User5', 'User5@gmail.com', 'abcdefgh123', '2024-01-05 00:00:00', '2024-01-05 00:00:00');

-- +goose Down
DELETE FROM users WHERE id IN ('1', '2', '3', '4', '5');