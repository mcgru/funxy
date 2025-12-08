# 02. Pipes и композиция

## Задача
Создавать читаемые цепочки трансформаций данных.

## Pipe оператор `|>`

```rust
import "lib/math" (absInt)

// Без pipe (вложенные вызовы, читать изнутри наружу)
result = show(absInt(-42))

// С pipe (линейно, читать слева направо)
result = -42 |> absInt |> show
print(result)  // "42"
```

## Объяснение

`x |> f` эквивалентно `f(x)`

Pipe передаёт левый операнд как первый аргумент правой функции.

## Трансформация списков

```rust
import "lib/list" (filter, map)

numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

// Императивно (многословно)
result1 = []
for n in numbers {
    if n % 2 == 0 {
        result1 = result1 ++ [n * n]
    }
}

// Функционально с pipe (элегантно!)
result2 = numbers
    |> filter(fun(n) -> n % 2 == 0)  // [2, 4, 6, 8, 10]
    |> map(fun(n) -> n * n)           // [4, 16, 36, 64, 100]

print(result2)  // [4, 16, 36, 64, 100]
```

## Агрегация

```rust
import "lib/list" (foldl)

sales = [100, 250, 180, 320, 90]

total = foldl(fun(acc, x) -> acc + x, 0, sales)
print(total)  // 940

average = total / len(sales)
print(average)  // 188
```

## Реальный пример: обработка данных

```rust
import "lib/string" (stringSplit)
import "lib/list" (map, filter)

// Простой парсер (без обработки ошибок для наглядности)
fun parseAge(s: String) -> Int {
    match read(s, Int) {
        Some(n) -> n
        Zero -> 0
    }
}

csv = "Alice,30,Engineer\nBob,25,Designer\nCarol,35,Manager"

employees = stringSplit(csv, "\n")
    |> map(fun(line) -> {
        parts = stringSplit(line, ",")
        {
            name: parts[0],
            age: parseAge(parts[1]),
            role: parts[2]
        }
    })
    |> filter(fun(e) -> e.age >= 30)
    |> map(fun(e) -> e.name)

print(employees)  // ["Alice", "Carol"]
```

## Композиция функций

```rust
// Создаём новые функции из существующих
fun double(x: Int) -> Int { x * 2 }
fun increment(x: Int) -> Int { x + 1 }
fun square(x: Int) -> Int { x * x }

// Применяем цепочку
result = 3 |> double |> increment |> square
print(result)  // ((3 * 2) + 1)² = 49
```

## Частичное применение с лямбдами

```rust
import "lib/list" (map, filter, all)

numbers = [1, 2, 3, 4, 5]

// Умножить все на 10
tens = numbers |> map(fun(x) -> x * 10)
print(tens)  // [10, 20, 30, 40, 50]

// Отфильтровать больше 3
big = numbers |> filter(fun(x) -> x > 3)
print(big)  // [4, 5]

// Проверить все ли положительные
allPositive = all(fun(x) -> x > 0, numbers)
print(allPositive)  // true
```

## Получение полей из списка объектов

```rust
import "lib/list" (map)

users = [
    { name: "Alice", age: 30 },
    { name: "Bob", age: 25 }
]

// Извлечение полей
names = users |> map(fun(u) -> u.name)
print(names)  // ["Alice", "Bob"]

ages = users |> map(fun(u) -> u.age)
print(ages)  // [30, 25]
```

## Комбинирование операций

```rust
import "lib/string" (stringToUpper)
import "lib/list" (filter, map, foldl)

words = ["hello", "world", "funxy", "rocks"]

result = words
    |> filter(fun(w) -> len(w) > 4)
    |> map(fun(w) -> stringToUpper(w))
    |> foldl(fun(acc, w) -> if acc == "" { w } else { acc ++ ", " ++ w }, "")

print(result)  // "HELLO, WORLD, FUNXY, ROCKS"
```

## Работа с Option

```rust
// Option встроен: Some T | Zero

fun mapOption(opt, f) {
    match opt {
        Some(x) -> Some(f(x))
        Zero -> Zero
    }
}

fun flatMapOption(opt, f) {
    match opt {
        Some(x) -> f(x)
        Zero -> Zero
    }
}

// Цепочка безопасных операций
result = Some(10)
    |> fun(o) -> mapOption(o, fun(x) -> x * 2)
    |> fun(o) -> flatMapOption(o, fun(x) -> if x > 15 { Some(x) } else { Zero })
    |> fun(o) -> mapOption(o, fun(x) -> "Result: " ++ show(x))

print(result)  // Some("Result: 20")
```

## Сравнение стилей

```rust
import "lib/list" (range, filter, map, foldl)

// Задача: найти сумму квадратов чётных чисел от 1 до 10

// Императивно
sum1 = 0
for i in range(1, 11) {
    if i % 2 == 0 {
        sum1 += i * i
    }
}
print(sum1)  // 220

// Функционально
sum2 = range(1, 11)
    |> filter(fun(x) -> x % 2 == 0)
    |> map(fun(x) -> x * x)
    |> foldl(fun(a, b) -> a + b, 0)
print(sum2)  // 220
```
