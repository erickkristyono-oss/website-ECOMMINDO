<?php

require __DIR__ . '/config.php';
require __DIR__ . '/lib/jwt.php';
require __DIR__ . '/lib/response.php';
require __DIR__ . '/lib/auth.php';
require __DIR__ . '/handlers/auth.php';
require __DIR__ . '/handlers/services.php';
require __DIR__ . '/handlers/cart.php';
require __DIR__ . '/handlers/checkout.php';

header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

$method = $_SERVER['REQUEST_METHOD'];
$route = trim($_GET['_route'] ?? '', '/');
$segments = $route === '' ? [] : explode('/', $route);

try {
    // POST /api/auth/register
    if ($method === 'POST' && $segments === ['auth', 'register']) {
        handle_register();
    }
    // POST /api/auth/login
    elseif ($method === 'POST' && $segments === ['auth', 'login']) {
        handle_login();
    }
    // GET /api/auth/me
    elseif ($method === 'GET' && $segments === ['auth', 'me']) {
        handle_me();
    }
    // GET /api/services
    elseif ($method === 'GET' && $segments === ['services']) {
        handle_list_services();
    }
    // GET /api/cart
    elseif ($method === 'GET' && $segments === ['cart']) {
        handle_get_cart();
    }
    // POST /api/cart
    elseif ($method === 'POST' && $segments === ['cart']) {
        handle_add_cart();
    }
    // DELETE /api/cart/{serviceId}
    elseif ($method === 'DELETE' && count($segments) === 2 && $segments[0] === 'cart') {
        handle_remove_cart($segments[1]);
    }
    // POST /api/checkout
    elseif ($method === 'POST' && $segments === ['checkout']) {
        handle_checkout();
    }
    // GET /api/orders
    elseif ($method === 'GET' && $segments === ['orders']) {
        handle_list_orders();
    } else {
        json_error(404, 'Endpoint tidak ditemukan.');
    }
} catch (Throwable $e) {
    json_error(500, 'Terjadi kesalahan pada server.');
}
