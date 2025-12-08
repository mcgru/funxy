# UUID (lib/uuid)

Модуль `lib/uuid` предоставляет функции для генерации и работы с UUID.

```rust
import "lib/uuid" (*)
```

## Что такое UUID?

UUID (Universally Unique Identifier) — 128-битный идентификатор, гарантирующий уникальность без централизованной координации.

Формат: `550e8400-e29b-41d4-a716-446655440000` (36 символов)

## Версии UUID

### v4 — Случайный (по умолчанию)

Генерируется из случайных чисел. Самый популярный вариант.

```rust
import "lib/uuid" (uuidNew, uuidToString)

myUuid = uuidNew()
print(uuidToString(myUuid))  // "550e8400-e29b-41d4-a716-446655440000"
```

### v7 — Упорядоченный по времени

Содержит метку времени. Идеален для первичных ключей в БД — сохраняет хронологический порядок.

```rust
import "lib/uuid" (uuidV7, uuidToString)

// Каждый следующий UUID "больше" предыдущего
id1 = uuidV7()
id2 = uuidV7()
// id2 > id1 по времени создания
```

### v5 — Детерминированный

Генерируется из namespace + name через SHA-1. Один и тот же input всегда даёт один и тот же UUID.

```rust
import "lib/uuid" (uuidV5, uuidNamespaceDNS, uuidToString)

ns = uuidNamespaceDNS()

// Одинаковые входные данные = одинаковый результат
id1 = uuidV5(ns, "example.com")
id2 = uuidV5(ns, "example.com")
print(uuidToString(id1) == uuidToString(id2))  // true
```

#### Стандартные namespace

```rust
uuidNamespaceDNS()   // для доменных имён
uuidNamespaceURL()   // для URL
uuidNamespaceOID()   // для OID
uuidNamespaceX500()  // для X.500 DN
```

## Специальные UUID

```rust
import "lib/uuid" (uuidNil, uuidMax, uuidToString)

// Nil UUID (все нули)
print(uuidToString(uuidNil()))  // "00000000-0000-0000-0000-000000000000"

// Max UUID (все единицы)
print(uuidToString(uuidMax()))  // "ffffffff-ffff-ffff-ffff-ffffffffffff"
```

## Парсинг

```rust
import "lib/uuid" (uuidParse, uuidToString, uuidVersion)

// Поддерживаются разные форматы
match uuidParse("550e8400-e29b-41d4-a716-446655440000") {
    Ok(u) -> {
        print("Версия: " ++ show(uuidVersion(u)))
        print("UUID: " ++ uuidToString(u))
    }
    Fail(err) -> print("Ошибка: " ++ err)
}

// Также работают:
// "550e8400e29b41d4a716446655440000"     (без дефисов)
// "{550e8400-e29b-41d4-a716-446655440000}" (с фигурными скобками)
// "urn:uuid:550e8400-e29b-41d4-a716-446655440000" (URN)
// "550E8400-E29B-41D4-A716-446655440000"  (uppercase)
```

## Форматы вывода

```rust
import "lib/uuid" (*)

u = uuidNew()

// Стандартный (8-4-4-4-12)
print(uuidToString(u))        // "550e8400-e29b-41d4-a716-446655440000"

// Компактный (без дефисов)
print(uuidToStringCompact(u)) // "550e8400e29b41d4a716446655440000"

// URN
print(uuidToStringUrn(u))     // "urn:uuid:550e8400-e29b-41d4-a716-446655440000"

// С фигурными скобками
print(uuidToStringBraces(u))  // "{550e8400-e29b-41d4-a716-446655440000}"

// Uppercase
print(uuidToStringUpper(u))   // "550E8400-E29B-41D4-A716-446655440000"
```

## Конвертация в байты

```rust
import "lib/uuid" (uuidNew, uuidToBytes, uuidFromBytes, uuidToString)

original = uuidNew()

// UUID -> Bytes (16 байт)
bytes = uuidToBytes(original)
print(len(bytes))  // 16

// Bytes -> UUID
match uuidFromBytes(bytes) {
    Ok(restored) -> {
        print(uuidToString(original) == uuidToString(restored))  // true
    }
    Fail(err) -> print("Ошибка: " ++ err)
}
```

## Информация о UUID

```rust
import "lib/uuid" (uuidNew, uuidV7, uuidNil, uuidVersion, uuidIsNil)

// Версия UUID
print(uuidVersion(uuidNew()))  // 4
print(uuidVersion(uuidV7()))   // 7

// Проверка на nil
print(uuidIsNil(uuidNil()))    // true
print(uuidIsNil(uuidNew()))    // false
```

## Сравнение

UUID поддерживают сравнение на равенство:

```rust
import "lib/uuid" (uuidNew, uuidNil)

// Равенство
a = uuidNil()
b = uuidNil()
print(a == b)  // true

// Неравенство (разные случайные UUID)
c = uuidNew()
d = uuidNew()
print(c != d)  // true
```

## Практические примеры

### Генерация ID для сущности

```rust
import "lib/uuid" (uuidNew, uuidToString)

type User = {
    id: String,
    name: String,
    email: String
}

fun createUser(name: String, email: String) -> User {
    {
        id: uuidToString(uuidNew()),
        name: name,
        email: email
    }
}
```

### ID для БД (v7 для сортировки)

```rust
import "lib/uuid" (uuidV7, uuidToString)
import "lib/sql" (*)

fun insertRecord(db, name: String) -> Result<String, Int> {
    recordId = uuidToString(uuidV7())
    sqlExec(db, "INSERT INTO records (id, name) VALUES ($1, $2)", [
        SqlString(recordId),
        SqlString(name)
    ])
}
```

### Детерминированный ID по email

```rust
import "lib/uuid" (uuidV5, uuidNamespaceDNS, uuidToString)

fun userIdFromEmail(email: String) -> String {
    ns = uuidNamespaceDNS()
    uuidToString(uuidV5(ns, email))
}

// Всегда один и тот же ID для одного email
id1 = userIdFromEmail("user@example.com")
id2 = userIdFromEmail("user@example.com")
print(id1 == id2)  // true
```

## Когда какую версию использовать?

| Версия | Когда использовать |
|--------|-------------------|
| v4 | Общего назначения, когда нужна просто уникальность |
| v7 | Primary keys в БД (сортируемые по времени) |
| v5 | Когда нужен детерминированный ID по известным данным |

## Сводка

| Функция | Тип | Описание |
|---------|-----|----------|
| `uuidNew` | `() -> Uuid` | Случайный UUID (v4) |
| `uuidV4` | `() -> Uuid` | Алиас для uuidNew |
| `uuidV5` | `(Uuid, String) -> Uuid` | Детерминированный UUID |
| `uuidV7` | `() -> Uuid` | Упорядоченный по времени |
| `uuidNil` | `() -> Uuid` | Nil UUID |
| `uuidMax` | `() -> Uuid` | Max UUID |
| `uuidNamespaceDNS` | `() -> Uuid` | DNS namespace |
| `uuidNamespaceURL` | `() -> Uuid` | URL namespace |
| `uuidNamespaceOID` | `() -> Uuid` | OID namespace |
| `uuidNamespaceX500` | `() -> Uuid` | X.500 namespace |
| `uuidParse` | `(String) -> Result<String, Uuid>` | Парсинг строки |
| `uuidFromBytes` | `(Bytes) -> Result<String, Uuid>` | Из 16 байт |
| `uuidToString` | `(Uuid) -> String` | Стандартный формат |
| `uuidToStringCompact` | `(Uuid) -> String` | Без дефисов |
| `uuidToStringUrn` | `(Uuid) -> String` | URN формат |
| `uuidToStringBraces` | `(Uuid) -> String` | С фигурными скобками |
| `uuidToStringUpper` | `(Uuid) -> String` | Uppercase |
| `uuidToBytes` | `(Uuid) -> Bytes` | В 16 байт |
| `uuidVersion` | `(Uuid) -> Int` | Версия UUID |
| `uuidIsNil` | `(Uuid) -> Bool` | Проверка на nil |

