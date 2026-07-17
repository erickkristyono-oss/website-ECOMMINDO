<?php

// Returns the authenticated user's ID, or sends a 401 JSON error and exits.
function require_auth(): int
{
    $header = $_SERVER['HTTP_AUTHORIZATION'] ?? '';
    if ($header === '' && function_exists('apache_request_headers')) {
        $headers = apache_request_headers();
        $header = $headers['Authorization'] ?? '';
    }

    if (!str_starts_with($header, 'Bearer ')) {
        json_error(401, 'Silakan login terlebih dahulu.');
    }

    $token = substr($header, 7);
    $claims = jwt_verify($token, jwt_secret());
    if ($claims === null || !isset($claims['user_id'])) {
        json_error(401, 'Sesi tidak valid, silakan login kembali.');
    }

    return (int) $claims['user_id'];
}
