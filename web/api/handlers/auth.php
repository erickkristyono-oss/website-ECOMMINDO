<?php

function user_to_array(array $user): array
{
    return [
        'id' => (int) $user['id'],
        'full_name' => $user['full_name'],
        'phone' => $user['phone'],
        'email' => $user['email'],
    ];
}

function find_user_by_email(string $email): ?array
{
    $stmt = db()->prepare('SELECT * FROM users WHERE email = ?');
    $stmt->execute([$email]);
    $user = $stmt->fetch();
    return $user ?: null;
}

function find_user_by_id(int $id): ?array
{
    $stmt = db()->prepare('SELECT * FROM users WHERE id = ?');
    $stmt->execute([$id]);
    $user = $stmt->fetch();
    return $user ?: null;
}

function handle_register(): void
{
    $body = read_json_body();
    $fullName = trim($body['full_name'] ?? '');
    $phone = trim($body['phone'] ?? '');
    $email = strtolower(trim($body['email'] ?? ''));
    $password = (string) ($body['password'] ?? '');

    if ($fullName === '' || $phone === '' || $email === '' || strlen($password) < 6) {
        json_error(400, 'Nama, nomor HP, email wajib diisi dan password minimal 6 karakter.');
    }

    if (find_user_by_email($email) !== null) {
        json_error(409, 'Email sudah terdaftar.');
    }

    $hash = password_hash($password, PASSWORD_BCRYPT);

    $stmt = db()->prepare('INSERT INTO users (full_name, phone, email, password_hash) VALUES (?, ?, ?, ?)');
    $stmt->execute([$fullName, $phone, $email, $hash]);
    $userId = (int) db()->lastInsertId();

    $token = jwt_generate(['user_id' => $userId, 'email' => $email], jwt_secret());

    json_response(201, [
        'token' => $token,
        'user' => ['id' => $userId, 'full_name' => $fullName, 'phone' => $phone, 'email' => $email],
    ]);
}

function handle_login(): void
{
    $body = read_json_body();
    $email = strtolower(trim($body['email'] ?? ''));
    $password = (string) ($body['password'] ?? '');

    $user = find_user_by_email($email);
    if ($user === null || !password_verify($password, $user['password_hash'])) {
        json_error(401, 'Email atau password salah.');
    }

    $token = jwt_generate(['user_id' => (int) $user['id'], 'email' => $user['email']], jwt_secret());

    json_response(200, [
        'token' => $token,
        'user' => user_to_array($user),
    ]);
}

function handle_me(): void
{
    $userId = require_auth();
    $user = find_user_by_id($userId);
    if ($user === null) {
        json_error(404, 'Pengguna tidak ditemukan.');
    }

    json_response(200, ['user' => user_to_array($user)]);
}
