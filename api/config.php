<?php

// Reads an environment variable if set (works on hosts that support it),
// otherwise falls back to the second argument.
//
// On Hostinger Cloud Hosting, the simplest path is usually to just edit the
// fallback values below directly with the credentials shown in
// hPanel -> Databases, rather than relying on env vars.
function env_value(string $key, string $fallback): string
{
    $value = getenv($key);
    return ($value === false || $value === '') ? $fallback : $value;
}

function db(): PDO
{
    static $pdo = null;
    if ($pdo === null) {
        $host = env_value('DB_HOST', '127.0.0.1');
        $port = env_value('DB_PORT', '3306');
        $name = env_value('DB_NAME', 'ecommindo');
        $user = env_value('DB_USER', 'root');
        $pass = env_value('DB_PASSWORD', '');

        $dsn = "mysql:host={$host};port={$port};dbname={$name};charset=utf8mb4";
        $pdo = new PDO($dsn, $user, $pass, [
            PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,
            PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
        ]);
    }
    return $pdo;
}

function jwt_secret(): string
{
    return env_value('JWT_SECRET', 'ecommindo-dev-secret-change-me');
}
