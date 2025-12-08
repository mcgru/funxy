# 01. Hello World и вывод

## Задача
Вывести текст на экран.

## Решение

```rust
// Простой вывод
print("Hello, Funxy!")

// Вывод с переменными
name = "World"
print("Hello, " ++ name ++ "!")

// Интерполяция строк
print("Hello, ${name}!")

// Многострочный вывод
print("Line 1\nLine 2\nLine 3")
```

## Объяснение

- `print()` — выводит с переносом строки
- `write()` — выводит без переноса строки
- `++` — конкатенация строк
- `${...}` — интерполяция внутри строк
- `\n` — перенос строки

## Вариации

```rust
// Вывод без переноса строки
write("Loading")
write(".")
write(".")
write(".\n")

// Вывод типа значения
x = 42
print(getType(x))  // Int

// Вывод списка
print([1, 2, 3])  // [1, 2, 3]
```

## Форматирование чисел

```rust
import "lib/string" (stringPadLeft)

price = 42
print("Price: $" ++ stringPadLeft(show(price), 5, '0'))
// Price: $00042
```

