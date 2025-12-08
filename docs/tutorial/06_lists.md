# Списки

## Представление списков

Списки в Funxy — это иммутабельные последовательности элементов одного типа. Внутренне реализованы как persistent vector (HAMT-based) для эффективных операций.

```rust
xs = [1, 2, 3, 4, 5]    // Тип: List<Int>
names = ["Alice", "Bob"] // Тип: List<String>
empty: List<Int> = []    // Пустой список с типом
```

## Базовые операции

### Доступ по индексу

```rust
xs = [10, 20, 30, 40]
xs[0]      // 10
xs[2]      // 30
xs[-1]     // 40 (с конца)
```

### Конкатенация и cons

```rust
// Конкатенация списков
[1, 2] ++ [3, 4]        // [1, 2, 3, 4]

// Cons — добавление в начало (::)
0 :: [1, 2, 3]          // [0, 1, 2, 3]
```

### Длина

```rust
len([1, 2, 3])          // 3
length([1, 2, 3])       // 3 (из lib/list)
```

## Multiline списки

Поддерживается многострочный синтаксис с trailing comma:

```rust
numbers = [
    1,
    2,
    3,
    4,
    5,
]
```

## Модуль lib/list

```rust
import "lib/list" (*)
```

### Доступ к элементам

```rust
import "lib/list" (*)

xs = [10, 20, 30, 40, 50]

// Первый/последний
head(xs)                // Some(10)
headOr(xs, 0)           // 10 (или default)
last(xs)                // Some(50)
lastOr(xs, 0)           // 50

// По индексу
nth(xs, 2)              // Some(30)
nthOr(xs, 10, -1)       // -1 (default)
```

### Срезы

```rust
import "lib/list" (*)

xs = [1, 2, 3, 4, 5]

tail(xs)                // Some([2, 3, 4, 5]) — всё кроме первого
init(xs)                // Some([1, 2, 3, 4]) — всё кроме последнего

take(xs, 3)             // [1, 2, 3]
drop(xs, 2)             // [3, 4, 5]
slice(xs, 1, 4)         // [2, 3, 4]

takeWhile(fun(x) -> x < 4, xs)  // [1, 2, 3]
dropWhile(fun(x) -> x < 3, xs)  // [3, 4, 5]
```

### Проверки

```rust
import "lib/list" (*)

xs = [1, 2, 3, 4, 5]

isEmpty([])             // true
isEmpty(xs)             // false
length(xs)              // 5
contains(xs, 3)         // true
any(fun(x) -> x > 4, xs)   // true
all(fun(x) -> x > 0, xs)   // true
```

### Поиск

```rust
import "lib/list" (*)

xs = [10, 20, 30, 40, 50]

indexOf(xs, 30)                   // Some(2)
indexOf(xs, 99)                   // Zero

find(fun(x) -> x > 25, xs)        // Some(30)
findIndex(fun(x) -> x > 25, xs)   // Some(2)
```

### Трансформации

```rust
import "lib/list" (*)

xs = [1, 2, 3, 4, 5]

// Map
map(fun(x) -> x * 2, xs)          // [2, 4, 6, 8, 10]

// Filter
filter(fun(x) -> x % 2 == 0, xs)  // [2, 4]

// Reverse
reverse(xs)                       // [5, 4, 3, 2, 1]

// Unique (remove duplicates)
unique([1, 2, 2, 3, 1])           // [1, 2, 3]

// Sort
sort([3, 1, 4, 1, 5])             // [1, 1, 3, 4, 5]
words = ["hello", "a", "world"]
sortBy(words, fun(a, b) -> length(a) - length(b))
```

### Свёртки (Fold)

```rust
import "lib/list" (*)

xs = [1, 2, 3, 4, 5]

// Left fold: ((((0 + 1) + 2) + 3) + 4) + 5 = 15
foldl((+), 0, xs)                 // 15

// Right fold: 1 + (2 + (3 + (4 + (5 + 0)))) = 15
foldr((+), 0, xs)                 // 15

// Пример: конкатенация строк
foldl((++), "", ["a", "b", "c"])  // "abc"
```

### Комбинирование

```rust
import "lib/list" (*)

// Zip — объединить два списка в пары
zip([1, 2, 3], ["a", "b", "c"])   // [(1, "a"), (2, "b"), (3, "c")]

// Unzip — разделить пары
pairs: List<(Int, String)> = [(1, "a"), (2, "b")]
unzip(pairs)                      // ([1, 2], ["a", "b"])

// Concat — объединить списки
concat([[1, 2], [3], [4, 5]])     // [1, 2, 3, 4, 5]

// Flatten — то же что concat
flatten([[1, 2], [3, 4]])         // [1, 2, 3, 4]

// Partition — разделить по условию
partition(fun(x) -> x % 2 == 0, [1, 2, 3, 4, 5])
// ([2, 4], [1, 3, 5])

// forEach — выполнить функцию для каждого элемента (побочные эффекты)
[1, 2, 3] |> forEach(fun(x) -> print(x))
// prints: 1, 2, 3
// returns: Nil
```

### Генерация

```rust
// Range — диапазон [start, end)
range(1, 5)                       // [1, 2, 3, 4]
range(0, 10)                      // [0, 1, 2, ..., 9]
```

## Pipe оператор

Списковые функции отлично работают с pipe:

```rust
import "lib/list" (*)

result = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    |> filter(fun(x) -> x % 2 == 0)   // [2, 4, 6, 8, 10]
    |> map(fun(x) -> x * x)           // [4, 16, 36, 64, 100]
    |> fun(xs) -> take(xs, 3)         // [4, 16, 36]
    |> foldl((+), 0)                   // 56
```

## Pattern Matching

### Деструктуризация

```rust
fun sum(xs: List<Int>) -> Int {
    match xs {
        [] -> 0
        [x, rest...] -> x + sum(rest)
    }
}

sum([1, 2, 3, 4, 5])  // 15
```

### Spread patterns

```rust
fun first3(xs: List<Int>) -> List<Int> {
    match xs {
        [a, b, c, rest...] -> [a, b, c]
        _ -> xs
    }
}

first3([1, 2, 3, 4, 5])  // [1, 2, 3]
first3([1, 2])           // [1, 2]
```

### Конкретные паттерны

```rust
import "lib/list" (length)

fun describe(xs: List<Int>) -> String {
    match xs {
        [] -> "empty"
        [x] -> "single: ${x}"
        [x, y] -> "pair: ${x}, ${y}"
        [x, rest...] -> "starts with ${x}, has ${length(rest)} more"
    }
}
```

## Spread в литералах

```rust
xs = [1, 2, 3]
ys = [4, 5]

// Spread в начале
[0, ...xs]              // [0, 1, 2, 3]

// Spread в конце
[...xs, 4]              // [1, 2, 3, 4]

// Несколько spread
[...xs, ...ys]          // [1, 2, 3, 4, 5]
```

## Практические примеры

### Сумма и среднее

```rust
import "lib/list" (length, foldl)

sum = fun(xs: List<Int>) -> Int {
    foldl((+), 0, xs)
}

average = fun(xs: List<Float>) -> Float {
    foldl((+), 0.0, xs) / intToFloat(length(xs))
}

sum([1, 2, 3, 4, 5])      // 15
average([1.0, 2.0, 3.0])  // 2.0
```

### Максимум и минимум

```rust
import "lib/list" (foldl, head, tail)

maximum = fun(xs: List<Int>) -> Int {
    foldl(fun(a, b) -> if a > b { a } else { b }, head(xs), tail(xs))
}

maximum([3, 1, 4, 1, 5, 9])  // 9
```

### Группировка

```rust
import "lib/list" (filter, unique, map)

fun groupBy<K, V>(keyFn: (V) -> K, xs: List<V>) -> List<(K, List<V>)> {
    keys = unique(map(keyFn, xs))
    map(fun(k) -> (k, filter(fun(x) -> keyFn(x) == k, xs)), keys)
}

// Группировка по чётности
nums = [1, 2, 3, 4, 5, 6]
groupBy(fun(x) -> x % 2, nums)
// [(1, [1, 3, 5]), (0, [2, 4, 6])]
```

### Обработка данных

```rust
import "lib/list" (*)
import "lib/string" (*)

// Подсчёт слов в тексте
fun countWords(text: String) -> Int {
    length(stringWords(text))
}

// Найти самое длинное слово
fun longestWord(text: String) -> Option<String> {
    stringWords(text)
        |> fun(ws) -> sortBy(ws, fun(a, b) -> length(b) - length(a))
        |> head
}

// Частота символов
fun charFreq(s: String) -> List<(Char, Int)> {
    chars = unique(s)
    map(fun(c) -> (c, length(filter(fun(x) -> x == c, s))), chars)
}
```

## Сводка lib/list

| Функция | Описание |
|---------|----------|
| `head`, `headOr` | Первый элемент |
| `last`, `lastOr` | Последний элемент |
| `nth`, `nthOr` | Элемент по индексу |
| `tail`, `init` | Без первого/последнего |
| `take`, `drop` | Взять/пропустить N элементов |
| `slice` | Срез по индексам |
| `takeWhile`, `dropWhile` | По условию |
| `isEmpty`, `length` | Проверка и длина |
| `contains`, `indexOf` | Поиск элемента |
| `find`, `findIndex` | Поиск по условию |
| `any`, `all` | Проверки по условию |
| `map`, `filter` | Трансформации |
| `foldl`, `foldr` | Свёртки |
| `reverse`, `sort`, `sortBy` | Упорядочивание |
| `unique` | Удаление дубликатов |
| `concat`, `flatten` | Объединение |
| `zip`, `unzip` | Комбинирование |
| `partition` | Разделение |
| `forEach` | Побочные эффекты |
| `range` | Генерация |

