# Словари (Map)

## Представление

Map в Funxy — это иммутабельный ассоциативный массив (хэш-таблица). Внутренне реализован как HAMT (Hash Array Mapped Trie) для эффективных операций над иммутабельными данными.

```rust
import "lib/map" (*)

scores = %{ "Alice" => 100, "Bob" => 85, "Charlie" => 92 }
// Тип: Map<String, Int>
```

## Синтаксис

### Создание Map

```rust
// Map literal
m = %{ "key1" => 1, "key2" => 2 }

// Пустой map (требует аннотации типа)
empty: Map<String, Int> = %{}

// Multiline с trailing comma
config = %{
    "host" => "localhost",
    "port" => "8080",
    "debug" => "1",
}
```

### Ключи разных типов

```rust
// String keys
stringMap = %{ "a" => 1, "b" => 2 }

// Int keys  
intMap = %{ 1 => "one", 2 => "two", 3 => "three" }

// Любой тип с Eq
tupleMap = %{ (0, 0) => "origin", (1, 0) => "x-axis" }
```

## Модуль lib/map

```rust
import "lib/map" (*)
```

### Доступ к значениям

```rust
import "lib/map" (*)

scores = %{ "Alice" => 100, "Bob" => 85 }

// Index access — возвращает Option<V>
scores["Alice"]              // Some(100)
scores["Unknown"]            // Zero

// mapGet — то же самое
mapGet(scores, "Alice")      // Some(100)
mapGet(scores, "Unknown")    // Zero

// mapGetOr — с default значением
mapGetOr(scores, "Alice", 0)   // 100
mapGetOr(scores, "Unknown", 0) // 0

// mapContains — проверка наличия
mapContains(scores, "Alice")   // true
mapContains(scores, "Dave")    // false
```

### Размер

```rust
import "lib/map" (mapSize)

scores = %{ "Alice" => 100, "Bob" => 85 }
mapSize(scores)              // 2
len(scores)                  // 2 (встроенный len тоже работает)
```

### Модификация (иммутабельная)

Все операции модификации возвращают **новый** Map, оригинал не меняется:

```rust
import "lib/map" (*)

scores = %{ "Alice" => 100, "Bob" => 85 }

// Добавить или обновить
scores2 = mapPut(scores, "Charlie", 92)
mapSize(scores)              // 2 (оригинал не изменён)
mapSize(scores2)             // 3

// Обновить существующий
scores3 = mapPut(scores, "Alice", 110)
mapGet(scores, "Alice")      // Some(100)  — оригинал
mapGet(scores3, "Alice")     // Some(110)  — новый

// Удалить ключ
scores4 = mapRemove(scores, "Bob")
mapSize(scores4)             // 1
mapContains(scores4, "Bob")  // false
```

### Слияние

```rust
import "lib/map" (*)

m1 = %{ "a" => 1, "b" => 2 }
m2 = %{ "b" => 20, "c" => 3 }

// Merge — второй map "выигрывает" при конфликте
merged = mapMerge(m1, m2)
mapGet(merged, "a")          // Some(1)   — из m1
mapGet(merged, "b")          // Some(20)  — из m2 (перезаписал)
mapGet(merged, "c")          // Some(3)   — из m2
```

### Итерация

```rust
import "lib/map" (*)

scores = %{ "Alice" => 100, "Bob" => 85, "Charlie" => 92 }

// Получить все ключи
keys = mapKeys(scores)       // ["Alice", "Bob", "Charlie"]

// Получить все значения
vals = mapValues(scores)     // [100, 85, 92]

// Получить пары (ключ, значение)
items = mapItems(scores)     // [("Alice", 100), ("Bob", 85), ...]
```

## Pattern Matching с Option

Так как mapGet возвращает Option<V>, используйте pattern matching:

```rust
import "lib/map" (*)

fun getScore(scores: Map<String, Int>, name: String) -> String {
    match mapGet(scores, name) {
        Some(s) -> "Score: " ++ show(s)
        Zero -> "Not found"
    }
}

// Пример использования
scores = %{ "Alice" => 100, "Bob" => 85 }
print(getScore(scores, "Alice"))     // Score: 100
print(getScore(scores, "Unknown"))   // Not found

// Или с mapGetOr
score = mapGetOr(scores, "Alice", 0)
```

## Практические примеры

### Подсчёт частоты

```rust
import "lib/map" (*)
import "lib/list" (foldl)

fun countFreq(xs: List<Char>) -> Map<Char, Int> {
    foldl(fun(m, x) -> {
        count = mapGetOr(m, x, 0)
        mapPut(m, x, count + 1)
    }, %{}, xs)
}

freq = countFreq(['a', 'b', 'a', 'c', 'a', 'b'])
mapGet(freq, 'a')            // Some(3)
mapGet(freq, 'b')            // Some(2)
mapGet(freq, 'c')            // Some(1)
```

### Группировка по ключу

```rust
import "lib/map" (*)
import "lib/list" (foldl, length)

fun groupByLen(xs: List<String>) -> Map<Int, List<String>> {
    foldl(fun(m, x) -> {
        k = length(x)
        existing = mapGetOr(m, k, [])
        mapPut(m, k, existing ++ [x])
    }, %{}, xs)
}

words = ["hi", "hello", "hey", "world", "ok"]
byLen = groupByLen(words)
mapGet(byLen, 2)             // Some(["hi", "ok"])
mapGet(byLen, 5)             // Some(["hello", "world"])
```

### Конфигурация

```rust
import "lib/map" (*)

// Загрузка конфигурации
defaultConfig = %{
    "host" => "localhost",
    "port" => "8080",
    "timeout" => "30",
}

userConfig = %{
    "port" => "3000",
    "debug" => "true",
}

// Merge: user overrides default
config = mapMerge(defaultConfig, userConfig)
mapGet(config, "host")       // Some("localhost") — default
mapGet(config, "port")       // Some("3000")      — overridden
mapGet(config, "debug")      // Some("true")      — user only
```

### Инверсия Map

```rust
import "lib/map" (*)
import "lib/list" (foldl)
import "lib/tuple" (fst, snd)

// Поменять ключи и значения местами
fun invert(m: Map<String, Int>) -> Map<Int, String> {
    foldl(fun(acc, kv) -> {
        mapPut(acc, snd(kv), fst(kv))
    }, %{}, mapItems(m))
}

m = %{ "a" => 1, "b" => 2, "c" => 3 }
inv = invert(m)
mapGet(inv, 1)               // Some("a")
mapGet(inv, 2)               // Some("b")
```

## Когда использовать Map

**Используйте Map когда:**
- Нужен быстрый поиск по ключу
- Данные часто читаются, но редко модифицируются
- Нужна иммутабельность (безопасность в многопоточном коде)
- Ключи разнородны или динамичны

**Используйте Record когда:**
- Фиксированный набор полей известен заранее
- Нужна типизация каждого поля
- Структура определена на этапе компиляции

```rust
import "lib/map" (mapGet)

// Record — статическая структура
type User = { name: String, age: Int, email: String }
user: User = { name: "Alice", age: 30, email: "a@b.com" }
user.name                    // Типизированный доступ

// Map — динамическая структура
userData = %{ "name" => "Alice", "age" => "30" }
mapGet(userData, "name")     // Option<String>
```

## Сводка lib/map

| Функция | Тип | Описание |
|---------|-----|----------|
| mapGet | (Map K V, K) -> Option V | Получить значение |
| mapGetOr | (Map K V, K, V) -> V | Получить или default |
| mapContains | (Map K V, K) -> Bool | Проверить наличие |
| mapSize | (Map K V) -> Int | Количество записей |
| mapPut | (Map K V, K, V) -> Map K V | Добавить/обновить |
| mapRemove | (Map K V, K) -> Map K V | Удалить ключ |
| mapMerge | (Map K V, Map K V) -> Map K V | Объединить |
| mapKeys | (Map K V) -> List K | Все ключи |
| mapValues | (Map K V) -> List V | Все значения |
| mapItems | (Map K V) -> List (K, V) | Все пары |

Также работает встроенный len(m) для размера и m[key] для доступа.

