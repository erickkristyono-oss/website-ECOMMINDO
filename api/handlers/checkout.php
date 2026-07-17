<?php

const BANK_ACCOUNTS = [
    ['bank' => 'Bank BCA', 'holder' => 'PT Ecommindo Jaya Persada', 'number' => '1234567890'],
    ['bank' => 'Bank Mandiri', 'holder' => 'PT Ecommindo Jaya Persada', 'number' => '9876543210'],
];

function generate_order_code(): string
{
    return 'EJP-' . (time() % 1000000) . str_pad((string) random_int(0, 999), 3, '0', STR_PAD_LEFT);
}

function handle_checkout(): void
{
    $userId = require_auth();
    $cart = cart_summary($userId);

    if (count($cart['items']) === 0) {
        json_error(400, 'Keranjang masih kosong.');
    }

    $pdo = db();
    $pdo->beginTransaction();
    try {
        $orderCode = generate_order_code();
        $status = 'Menunggu Konfirmasi Pembayaran';

        $stmt = $pdo->prepare('INSERT INTO orders (order_code, user_id, total, status) VALUES (?, ?, ?, ?)');
        $stmt->execute([$orderCode, $userId, $cart['total'], $status]);
        $orderId = (int) $pdo->lastInsertId();

        $itemStmt = $pdo->prepare('INSERT INTO order_items (order_id, service_id, service_name, price) VALUES (?, ?, ?, ?)');
        $items = [];
        foreach ($cart['items'] as $service) {
            $itemStmt->execute([$orderId, $service['id'], $service['name'], $service['price']]);
            $items[] = [
                'service_id' => $service['id'],
                'service_name' => $service['name'],
                'price' => $service['price'],
            ];
        }

        $clearStmt = $pdo->prepare('DELETE FROM cart_items WHERE user_id = ?');
        $clearStmt->execute([$userId]);

        $pdo->commit();
    } catch (Throwable $e) {
        $pdo->rollBack();
        json_error(500, 'Gagal membuat pesanan.');
    }

    json_response(201, [
        'order_code' => $orderCode,
        'total' => $cart['total'],
        'status' => $status,
        'items' => $items,
        'bank_accounts' => BANK_ACCOUNTS,
    ]);
}

function handle_list_orders(): void
{
    $userId = require_auth();

    $stmt = db()->prepare('SELECT id, order_code, total, status, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC');
    $stmt->execute([$userId]);
    $orders = $stmt->fetchAll();

    $itemStmt = db()->prepare('SELECT service_id, service_name, price FROM order_items WHERE order_id = ?');

    $result = [];
    foreach ($orders as $order) {
        $itemStmt->execute([$order['id']]);
        $items = array_map(function ($item) {
            return [
                'service_id' => $item['service_id'],
                'service_name' => $item['service_name'],
                'price' => (float) $item['price'],
            ];
        }, $itemStmt->fetchAll());

        $result[] = [
            'order_code' => $order['order_code'],
            'total' => (float) $order['total'],
            'status' => $order['status'],
            'items' => $items,
        ];
    }

    json_response(200, $result);
}
