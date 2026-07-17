<?php

function service_to_array(array $row): array
{
    return [
        'id' => $row['id'],
        'name' => $row['name'],
        'short_desc' => $row['short_desc'],
        'description' => $row['description'],
        'price' => (float) $row['price'],
        'icon' => $row['icon'],
    ];
}

function find_service_by_id(string $id): ?array
{
    $stmt = db()->prepare('SELECT * FROM services WHERE id = ?');
    $stmt->execute([$id]);
    $row = $stmt->fetch();
    return $row ?: null;
}

function handle_list_services(): void
{
    $rows = db()->query('SELECT * FROM services ORDER BY id')->fetchAll();
    json_response(200, array_map('service_to_array', $rows));
}
