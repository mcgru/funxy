# Криптография и кодирование (lib/crypto)

Модуль `lib/crypto` предоставляет функции для кодирования, хеширования и HMAC.

```rust
import "lib/crypto" (*)
```

## Функции кодирования

### Base64

```rust
base64Encode(s: String) -> String
base64Decode(s: String) -> String
```

Кодирование и декодирование в Base64.

```rust
import "lib/crypto" (base64Encode, base64Decode)

encoded = base64Encode("Hello, World!")
print(encoded)  // SGVsbG8sIFdvcmxkIQ==

decoded = base64Decode(encoded)
print(decoded)  // Hello, World!

// Кодирование бинарных данных для передачи
data = "some binary data"
safe = base64Encode(data)
// передать safe по сети
original = base64Decode(safe)
```

### Hex (шестнадцатеричное)

```rust
hexEncode(s: String) -> String
hexDecode(s: String) -> String
```

Кодирование в шестнадцатеричное представление и обратно.

```rust
import "lib/crypto" (hexEncode, hexDecode)

hex = hexEncode("ABC")
print(hex)  // 414243

original = hexDecode(hex)
print(original)  // ABC
```

## Хеш-функции

Все хеш-функции возвращают результат в виде шестнадцатеричной строки.

### MD5

```rust
md5(s: String) -> String
```

**Внимание:** MD5 считается криптографически небезопасным. Используйте только для контрольных сумм, не для безопасности.

```rust
import "lib/crypto" (md5)

hash = md5("hello")
print(hash)  // 5d41402abc4b2a76b9719d911017c592
print(len(hash))  // 32 (hex chars)
```

### SHA1

```rust
sha1(s: String) -> String
```

**Внимание:** SHA1 считается устаревшим для криптографических целей.

```rust
import "lib/crypto" (sha1)

hash = sha1("hello")
print(hash)  // aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
print(len(hash))  // 40 (hex chars)
```

### SHA256

```rust
sha256(s: String) -> String
```

Рекомендуемый алгоритм для большинства задач.

```rust
import "lib/crypto" (sha256)

hash = sha256("hello")
print(hash)  // 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
print(len(hash))  // 64 (hex chars)

// Хеширование паролей (простой пример)
passwordHash = sha256("mypassword" ++ "salt123")
```

### SHA512

```rust
sha512(s: String) -> String
```

Более длинный хеш для повышенной безопасности.

```rust
import "lib/crypto" (sha512)

hash = sha512("hello")
print(len(hash))  // 128 (hex chars)
```

## HMAC

HMAC (Hash-based Message Authentication Code) — код аутентификации сообщения с использованием хеш-функции и секретного ключа.

### hmacSha256

```rust
hmacSha256(key: String, message: String) -> String
```

```rust
import "lib/crypto" (hmacSha256)

signature = hmacSha256("secret-key", "message to sign")
print(signature)  // 64 hex chars

// Проверка подписи
expectedSig = hmacSha256("secret-key", "message to sign")
if signature == expectedSig {
    print("Valid signature")
}
```

### hmacSha512

```rust
hmacSha512(key: String, message: String) -> String
```

```rust
import "lib/crypto" (hmacSha512)

signature = hmacSha512("secret-key", "message")
print(len(signature))  // 128 hex chars
```

## Практические примеры

### Простая подпись API-запроса

```rust
import "lib/crypto" (hmacSha256, sha256)
import "lib/time" (timeNow)

fun signRequest(apiKey: String, secretKey: String, body: String) -> String {
    timestamp = show(timeNow())
    payload = timestamp ++ body
    hmacSha256(secretKey, payload)
}

// Использование
apiKey = "my-api-key"
secretKey = "my-secret"
body = "{\"action\": \"buy\"}"

signature = signRequest(apiKey, secretKey, body)
print("X-Signature: ${signature}")
```

### Контрольная сумма файла

```rust
import "lib/crypto" (sha256)
import "lib/io" (fileRead)

fun fileChecksum(path: String) -> String {
    match fileRead(path) {
        Ok(content) -> sha256(content)
        Fail(_) -> ""
    }
}

checksum = fileChecksum("myfile.txt")
print("SHA256: ${checksum}")
```

### Генерация токена

```rust
import "lib/crypto" (sha256, base64Encode)
import "lib/time" (clockNs)

fun generateToken(userId: String) -> String {
    // Простой токен на основе времени и userId
    data = userId ++ show(clockNs())
    hash = sha256(data)
    base64Encode(hash)
}

token = generateToken("user123")
print("Token: ${token}")
```

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `base64Encode` | `String -> String` | Кодирование в Base64 |
| `base64Decode` | `String -> String` | Декодирование из Base64 |
| `hexEncode` | `String -> String` | Кодирование в hex |
| `hexDecode` | `String -> String` | Декодирование из hex |
| `md5` | `String -> String` | MD5 хеш (32 hex) |
| `sha1` | `String -> String` | SHA1 хеш (40 hex) |
| `sha256` | `String -> String` | SHA256 хеш (64 hex) |
| `sha512` | `String -> String` | SHA512 хеш (128 hex) |
| `hmacSha256` | `(String, String) -> String` | HMAC-SHA256 |
| `hmacSha512` | `(String, String) -> String` | HMAC-SHA512 |

## Рекомендации по безопасности

1. **Не используйте MD5 или SHA1** для криптографических целей (пароли, подписи).
2. **Используйте SHA256** как стандартный выбор для хеширования.
3. **Используйте HMAC** для аутентификации сообщений, а не простой хеш.
4. **Храните секретные ключи** в переменных окружения, не в коде.
5. **Для паролей** используйте специализированные алгоритмы (bcrypt, argon2) — они не включены в этот модуль.

