# Записи (Records)

Записи (Records) — это структурные типы с именованными полями.

## Определение типа записи

```rust
// Именованный тип записи
type Point = { x: Int, y: Int }

// Тип с несколькими полями разных типов
type Person = { name: String, age: Int, active: Bool }

// Вложенные записи
type Rectangle = { topLeft: Point, bottomRight: Point }
```

## Создание записей

### Литералы записей

```rust
type Point = { x: Int, y: Int }
type Person = { name: String, age: Int, active: Bool }

// Анонимная запись (без объявления типа)
point = { x: 10, y: 20 }

// Запись с аннотацией типа
p: Point = { x: 10, y: 20 }

// Порядок полей не важен
person: Person = { age: 30, name: "Alice", active: true }
```

### Пустая запись

```rust
type Empty = {}

e: Empty = {}
```

## Доступ к полям

Используйте точечную нотацию:

```rust
p = { x: 10, y: 20 }

print(p.x)    // 10
print(p.y)    // 20

// Вложенный доступ
rect = { topLeft: { x: 0, y: 0 }, bottomRight: { x: 100, y: 50 } }
print(rect.topLeft.x)      // 0
print(rect.bottomRight.y)  // 50
```

## Изменение полей

Funxy позволяет изменять поля записей:

```rust
p = { x: 10, y: 20 }
p.x = 100
print(p.x)  // 100

// Вложенное изменение
rect = { tl: { x: 0, y: 0 }, br: { x: 10, y: 10 } }
rect.tl.y = 50
print(rect.tl.y)  // 50
```

## Обновление записи (spread)

Создание новой записи на основе существующей с изменением некоторых полей:

```rust
type Point = { x: Int, y: Int }

base: Point = { x: 1, y: 2 }

// Создать новую запись, изменив только x
updated = { ...base, x: 10 }
print(updated.x)  // 10
print(updated.y)  // 2 (из base)

// Оригинал не изменён
print(base.x)     // 1
```

## Деструктуризация записей

Записи поддерживают деструктуризацию для извлечения полей в переменные:

```rust
// Базовая деструктуризация
p = { x: 3, y: 4 }
{ x: a, y: b } = p
print(a)  // 3
print(b)  // 4
```

### С именованными типами

```rust
type Point = { x: Int, y: Int }

p: Point = { x: 10, y: 20 }
{ x: x, y: y } = p
print(x)  // 10
print(y)  // 20
```

### Частичная деструктуризация

Можно извлечь только нужные поля:

```rust
person = { name: "Alice", age: 30, city: "Moscow" }

// Только name
{ name: n } = person
print(n)  // Alice
```

### Вложенная деструктуризация

```rust
data = { 
    user: { name: "Bob", role: "admin" },
    count: 42
}

{ user: { name: userName, role: r }, count: c } = data
print(userName)  // Bob
print(r)         // admin
print(c)         // 42
```

### В функциях

```rust
import "lib/math" (sqrt)

type Point = { x: Int, y: Int }

fun magnitude(p: Point) -> Float {
    { x: x, y: y } = p
    sqrt(intToFloat(x * x + y * y))
}

p: Point = { x: 3, y: 4 }
print(magnitude(p))  // 5
```

## Номинальные vs Структурные типы

### Структурные типы (анонимные записи)

Без аннотации типа запись имеет структурный тип:

```rust
point1 = { x: 10, y: 20 }
print(getType(point1))  // { x: Int, y: Int }
```

### Номинальные типы (определение типа записи)

Когда вы определяете тип записи и используете аннотацию, запись имеет номинальный тип:

```rust
// Определение типа записи
type Point = { x: Int, y: Int }

p: Point = { x: 10, y: 20 }
print(getType(p))  // Point
```

**Важно:** `type Point = { ... }` — это **определение типа записи**, а не type alias. Type alias определяется с ключевым словом `alias`:

```rust
type Point = { x: Int, y: Int }

// Type alias — просто другое имя для существующего типа
type alias Coordinate = Point
```

### Почему это важно?

Номинальные типы нужны для:

1. **Extension methods** — методы привязаны к имени типа:

```rust
type Point = { x: Int, y: Int }

fun (p: Point) length() -> Int {
    p.x + p.y  // упрощённо
}

// Работает только с Point
p: Point = { x: 3, y: 4 }
print(p.length())  // 7

// НЕ работает с анонимной записью
anon = { x: 3, y: 4 }
// anon.length()  // Ошибка! У { x: Int, y: Int } нет метода length
```

2. **Trait instances** — реализации трейтов привязаны к имени типа:

```rust
type Point = { x: Int, y: Int }

instance Default Point {}

p: Point = default(Point)  // { x: 0, y: 0 }
```

3. **Ясность кода** — имена типов документируют намерение:

```rust
type Point = { x: Int, y: Int }

// Неясно что это
fun processAnon(data: { x: Int, y: Int }) -> Int { data.x + data.y }

// Ясно что это координата
fun processNamed(data: Point) -> Int { data.x + data.y }
```

## Extension Methods

Определение методов на записях:

```rust
type Point = { x: Int, y: Int }

// Метод без параметров
fun (p: Point) distanceFromOrigin() -> Int {
    p.x * p.x + p.y * p.y
}

// Метод с параметрами
fun (p: Point) add(other: Point) -> Point {
    { x: p.x + other.x, y: p.y + other.y }
}

// Использование
p1: Point = { x: 3, y: 4 }
p2: Point = { x: 1, y: 1 }

print(p1.distanceFromOrigin())  // 25
print(p1.add(p2).x)             // 4
```

## Pattern Matching

Анонимные записи поддерживают паттерн-матчинг:

```rust
// Pattern matching работает с анонимными записями
fun describe(p: { x: Int, y: Int }) -> String {
    match p {
        { x: 0, y: 0 } -> "origin"
        { x: 0, y: _ } -> "on Y axis"
        { x: _, y: 0 } -> "on X axis"
        { x: x, y: y } -> "point at (" ++ show(x) ++ ", " ++ show(y) ++ ")"
    }
}

p = { x: 0, y: 5 }
print(describe(p))  // "on Y axis"
```

Для именованных типов записей используйте доступ к полям:

```rust
type Point = { x: Int, y: Int }

fun describePoint(p: Point) -> String {
    if p.x == 0 && p.y == 0 { "origin" }
    else if p.x == 0 { "on Y axis" }
    else if p.y == 0 { "on X axis" }
    else { "point at (" ++ show(p.x) ++ ", " ++ show(p.y) ++ ")" }
}

p: Point = { x: 0, y: 5 }
print(describePoint(p))  // "on Y axis"
```

### Частичный матчинг (Row Polymorphism)

Можно матчить только часть полей:

```rust
r = { x: 10, y: 20, z: 30 }

match r {
    { x: val } -> print(val)  // 10 — остальные поля игнорируются
}
```

## Row Polymorphism

Функция может принимать записи с "дополнительными" полями. В отличие от некоторых языков, синтаксис `...` **не требуется** — row polymorphism работает автоматически:

```rust
// Функция ожидает запись с полем x (без "...")
fun getX(r: { x: Int }) -> Int {
    r.x
}

// Можно передать запись с любым количеством дополнительных полей
point = { x: 10, y: 20 }
print(getX(point))  // 10

extended = { x: 5, y: 6, z: 7, name: "test" }
print(getX(extended))  // 5

// Даже совершенно разные записи, если есть нужное поле
config = { x: 100, host: "localhost", port: 8080 }
print(getX(config))  // 100
```

> **Примечание**: Синтаксис `{ x: Int, ... }` для row polymorphism **не нужен** и не поддерживается. Язык автоматически допускает любые записи, содержащие как минимум указанные поля.

## Записи в функциях

### Как параметры

```rust
type Config = { host: String, port: Int, debug: Bool }

fun connect(cfg: Config) -> String {
    cfg.host ++ ":" ++ show(cfg.port)
}

config: Config = { host: "localhost", port: 8080, debug: true }
print(connect(config))  // "localhost:8080"
```

### Как возвращаемое значение

```rust
type Point = { x: Int, y: Int }

fun makePoint(x: Int, y: Int) -> Point {
    { x: x, y: y }
}

p = makePoint(10, 20)
print(p.x)  // 10
```

## Generics с записями

```rust
type Box<T> = { value: T }

fun makeBox<T>(v: T) -> Box<T> {
    { value: v }
}

intBox = makeBox(42)
print(intBox.value)  // 42

strBox = makeBox("hello")
print(strBox.value)  // "hello"
```

## Практические примеры

### Конфигурация

```rust
type ServerConfig = {
    host: String,
    port: Int,
    timeout: Int,
    debug: Bool
}

defaultConfig: ServerConfig = {
    host: "localhost",
    port: 8080,
    timeout: 30,
    debug: false
}

// Создать конфигурацию для production
prodConfig = { ...defaultConfig, host: "api.example.com", debug: false }
```

### Состояние в игре

```rust
type Player = { name: String, health: Int, x: Int, y: Int }

fun damage(p: Player, amount: Int) -> Player {
    { ...p, health: p.health - amount }
}

fun moveTo(p: Player, newX: Int, newY: Int) -> Player {
    { ...p, x: newX, y: newY }
}

player: Player = { name: "Hero", health: 100, x: 0, y: 0 }
player = damage(player, 20)
player = moveTo(player, 10, 5)
print(player.health)  // 80
print(player.x)       // 10
```

### DTO для API

```rust
type User = { id: Int, name: String, email: String }
type CreateUserRequest = { name: String, email: String }

fun createUser(req: CreateUserRequest, id: Int) -> User {
    { id: id, name: req.name, email: req.email }
}
```

## Когда использовать Records

**Используйте Records когда:**
- Структура данных фиксирована и известна заранее
- Нужна типизация каждого поля
- Хотите использовать extension methods
- Хотите реализовать трейты для типа

**Используйте Map когда:**
- Ключи динамические или неизвестны заранее
- Все значения одного типа
- Нужен быстрый поиск по произвольному ключу

**Используйте Tuple когда:**
- Нужна фиксированная коллекция без имён полей
- Данные временные (например, возврат нескольких значений)

## Сводка

| Операция | Синтаксис | Пример |
|----------|-----------|--------|
| Определение типа | `type Name = { field: Type }` | `type Point = { x: Int, y: Int }` |
| Создание | `{ field: value }` | `{ x: 10, y: 20 }` |
| С аннотацией | `name: Type = { ... }` | `p: Point = { x: 10, y: 20 }` |
| Доступ к полю | `record.field` | `p.x` |
| Изменение поля | `record.field = value` | `p.x = 100` |
| Spread/Update | `{ ...base, field: value }` | `{ ...p, x: 0 }` |
| Extension method | `fun (r: Type) name() { }` | `fun (p: Point) len() { }` |
| Pattern match | `match r { { f: v } -> ... }` | `match p { { x: 0 } -> "zero" }` |

## См. также

- [Custom Types](06_custom_types.md) — ADT и type aliases
- [Pattern Matching](07_pattern_matching.md) — деструктуризация записей
- [Traits](08_traits.md) — реализация трейтов для записей
- [Maps](24_maps.md) — динамические ассоциативные массивы

