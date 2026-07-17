<?php

function cart_summary(int $userId): array
{
    $stmt = db()->prepare('SELECT service_id FROM cart_items WHERE user_id = ? ORDER BY created_at');
    $stmt->execute([$userId]);
    $serviceIds = $stmt->fetchAll(PDO::FETCH_COLUMN);

    $items = [];
    $total = 0.0;
    foreach ($serviceIds as $serviceId) {
        $service = find_service_by_id($serviceId);
        if ($service === null) {
            continue;
        }
        $items[] = service_to_array($service);
        $total += (float) $service['price'];
    }

    return ['items' => $items, 'total' => $total];
}

function handle_get_cart(): void
{
    $userId = require_auth();
    json_response(200, cart_summary($userId));
}

function handle_add_cart(): void
{
    $userId = require_auth();
    $body = read_json_body();
    $serviceId = trim($body['service_id'] ?? '');

    if ($serviceId === '') {
        json_error(400, 'service_id wajib diisi.');
    }

    if (find_service_by_id($serviceId) === null) {
        json_error(404, 'Layanan tidak ditemukan.');
    }

    $stmt = db()->prepare('INSERT IGNORE INTO cart_items (user_id, service_id) VALUES (?, ?)');
    $stmt->execute([$userId, $serviceId]);

    json_response(200, cart_summary($userId));
}

function handle_remove_cart(string $serviceId): void
{
    $userId = require_auth();

    $stmt = db()->prepare('DELETE FROM cart_items WHERE user_id = ? AND service_id = ?');
    $stmt->execute([$userId, $serviceId]);

    json_response(200, cart_summary($userId));
}
