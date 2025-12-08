# 04. Pattern Matching

## Задача
Элегантно обрабатывать разные случаи без цепочек if-else.

## Базовое решение

```rust
fun describe(x: Int) -> String {
    match x {
        0 -> "zero"
        1 -> "one"
        n if n < 0 -> "negative"
        n if n > 100 -> "big"
        _ -> "some number"
    }
}

print(describe(0))    // zero
print(describe(-5))   // negative
print(describe(150))  // big
print(describe(42))   // some number
```

## Объяснение

- `match` проверяет значение по порядку
- `_` — wildcard, ловит всё остальное
- `n if condition` — guard pattern с условием
- Возвращает значение (это выражение!)

## Деструктуризация кортежей

```rust
fun process(pair: (Int, Int)) -> String {
    match pair {
        (0, 0) -> "origin"
        (0, y) -> "on Y axis at " ++ show(y)
        (x, 0) -> "on X axis at " ++ show(x)
        (x, y) -> "point (" ++ show(x) ++ ", " ++ show(y) ++ ")"
    }
}

print(process((0, 0)))   // origin
print(process((0, 5)))   // on Y axis at 5
print(process((3, 4)))   // point (3, 4)
```

## Деструктуризация списков

```rust
import "lib/list" (length)

fun listInfo(xs) {
    match xs {
        [] -> "empty"
        [x] -> "single element"
        [x, y] -> "two elements"
        [head, tail...] -> "starts with something, " ++ show(length(tail)) ++ " more"
    }
}

print(listInfo([]))           // empty
print(listInfo([1]))          // single element
print(listInfo([1, 2]))       // two elements
print(listInfo([1, 2, 3, 4])) // starts with something, 3 more
```

Синтаксис spread: `tail...` (после переменной).

## ADT (Algebraic Data Types)

```rust
// Конструктор с одним аргументом: Circle(Float)
// Конструктор с несколькими: Rectangle((Float, Float)) — кортеж
type Shape = Circle(Float) | Rectangle((Float, Float))

fun area(s: Shape) -> Float {
    match s {
        Circle(r) -> 3.14159 * r * r
        Rectangle((w, h)) -> w * h
    }
}

print(area(Circle(5.0)))             // 78.53975
print(area(Rectangle((4.0, 3.0))))   // 12
```

## Option (встроенный)

`Option<T>` — встроенный тип: `Some T | Zero`.

```rust
fun safeDivide(a: Int, b: Int) -> Option<Int> {
    if b == 0 { Zero } else { Some(a / b) }
}

fun showResult(opt: Option<Int>) -> String {
    match opt {
        Some(x) -> "Result: " ++ show(x)
        Zero -> "Cannot divide by zero"
    }
}

print(showResult(safeDivide(10, 2)))  // Result: 5
print(showResult(safeDivide(10, 0)))  // Cannot divide by zero
```

## FizzBuzz

```rust
import "lib/list" (range)

fun fizzbuzz(n: Int) -> String {
    match (n % 3, n % 5) {
        (0, 0) -> "FizzBuzz"
        (0, _) -> "Fizz"
        (_, 0) -> "Buzz"
        _ -> show(n)
    }
}

for i in range(1, 16) {
    print(fizzbuzz(i))
}
```

## Вложенный matching

`Result<E, A>` — встроенный тип: `Ok(A) | Fail(E)`.

```rust
// Result<E, A> = Ok(A) | Fail(E) — встроенный

fun processResult(r: Result<String, Option<Int>>) -> String {
    match r {
        Ok(Some(x)) -> "Got value: " ++ show(x)
        Ok(Zero) -> "Ok but empty"
        Fail(e) -> "Error: " ++ e
    }
}

print(processResult(Ok(Some(42))))   // Got value: 42
print(processResult(Ok(Zero)))        // Ok but empty
print(processResult(Fail("oops")))    // Error: oops
```

## Guard Patterns (условия)

```rust
fun grade(score: Int) -> String {
    match score {
        s if s >= 90 -> "A"
        s if s >= 80 -> "B"
        s if s >= 70 -> "C"
        s if s >= 60 -> "D"
        _ -> "F"
    }
}

print(grade(95))  // A
print(grade(72))  // C
print(grade(55))  // F
```

## Комбинация условий

```rust
fun classify(x: Int) -> String {
    match x {
        n if n > 0 && n % 2 == 0 -> "positive even"
        n if n > 0 -> "positive odd"
        n if n < 0 -> "negative"
        _ -> "zero"
    }
}

print(classify(4))   // positive even
print(classify(3))   // positive odd
print(classify(-5))  // negative
print(classify(0))   // zero
```

## Сводка синтаксиса

| Паттерн | Пример | Описание |
|---------|--------|----------|
| Литерал | `0 ->` | Точное совпадение |
| Переменная | `x ->` | Привязка к переменной |
| Wildcard | `_ ->` | Любое значение |
| Кортеж | `(x, y) ->` | Деструктуризация |
| Список | `[x, y] ->` | Фиксированная длина |
| Список spread | `[h, t...] ->` | Голова и хвост |
| ADT | `Circle(r) ->` | Деструктуризация ADT |
| Guard | `x if x > 0 ->` | С условием |
