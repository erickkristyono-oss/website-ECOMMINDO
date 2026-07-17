<?php

const JWT_TTL_SECONDS = 7 * 24 * 60 * 60; // 7 days

function base64url_encode(string $data): string
{
    return rtrim(strtr(base64_encode($data), '+/', '-_'), '=');
}

function base64url_decode(string $data): string
{
    $padded = str_pad($data, strlen($data) + (4 - strlen($data) % 4) % 4, '=');
    return base64_decode(strtr($padded, '-_', '+/'));
}

function jwt_generate(array $claims, string $secret): string
{
    $header = ['alg' => 'HS256', 'typ' => 'JWT'];
    $claims['iat'] = time();
    $claims['exp'] = time() + JWT_TTL_SECONDS;

    $segments = [
        base64url_encode(json_encode($header)),
        base64url_encode(json_encode($claims)),
    ];
    $signingInput = implode('.', $segments);
    $signature = hash_hmac('sha256', $signingInput, $secret, true);
    $segments[] = base64url_encode($signature);

    return implode('.', $segments);
}

function jwt_verify(string $token, string $secret): ?array
{
    $parts = explode('.', $token);
    if (count($parts) !== 3) {
        return null;
    }
    [$headerB64, $payloadB64, $sigB64] = $parts;

    $expectedSig = base64url_encode(hash_hmac('sha256', "{$headerB64}.{$payloadB64}", $secret, true));
    if (!hash_equals($expectedSig, $sigB64)) {
        return null;
    }

    $payload = json_decode(base64url_decode($payloadB64), true);
    if (!is_array($payload)) {
        return null;
    }
    if (isset($payload['exp']) && $payload['exp'] < time()) {
        return null;
    }

    return $payload;
}
