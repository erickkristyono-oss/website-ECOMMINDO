<?php
// Local-dev-only router for `php -S`, which does not read .htaccess.
// Mirrors the rewrite rule in web/api/.htaccess so the API behaves the
// same locally as it will on real Apache/LiteSpeed hosting.

$uri = urldecode(parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH));

if ($uri === '/') {
    $uri = '/index.html';
}

if (file_exists(__DIR__ . $uri) && !is_dir(__DIR__ . $uri)) {
    return false;
}

if (str_starts_with($uri, '/api/')) {
    $_GET['_route'] = substr($uri, strlen('/api/'));
    require __DIR__ . '/api/index.php';
    return true;
}

http_response_code(404);
echo 'Not Found';
