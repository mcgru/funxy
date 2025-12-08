# 03. Коллекции: списки и Map

## Задача
Работать с коллекциями данных: создавать, трансформировать, искать, агрегировать.

---

## Списки (List)

### Создание

```rust
import "lib/list" (range)

// Литерал
nums = [1, 2, 3, 4, 5]
empty = []

// Через range (start включён, end не включён)
r = range(1, 10)   // [1, 2, 3, ..., 9]
r2 = range(1, 11)  // [1, 2, 3, ..., 10]
print(r)
print(r2)
```

### Строки — это списки символов

```rust
import "lib/list" (head, tail, take, drop, reverse, map, filter, length)
import "lib/char" (charToUpper)

s = "hello"

// Индексирование
print(s[0])        // 'h'
print(s[4])        // 'o'

// Функции списков работают со строками
print(head(s))     // 'h'
print(tail(s))     // "ello"
print(take(s, 3))  // "hel"
print(drop(s, 2))  // "llo"
print(reverse(s))  // "olleh"
print(length(s))   // 5

// map — преобразовать каждый символ
upper = map(fun(c) -> charToUpper(c), s)
print(upper)       // "HELLO"

// filter — оставить символы по условию
vowels = filter(fun(c) -> c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u', s)
print(vowels)      // "eo"
```

### Доступ к элементам

```rust
import "lib/list" (head, tail, last, nth, headOr, nthOr)

nums = [1, 2, 3, 4, 5]

// По индексу
print(nums[0])   // 1
print(nums[2])   // 3

// Безопасный доступ
print(head(nums))        // 1 (первый элемент)
print(last(nums))        // 5 (последний)
print(nth(nums, 2))      // 3 (по индексу)

// С дефолтным значением
print(headOr([], 0))      // 0
print(nthOr(nums, 10, 0)) // 0
```

### Подсписки

```rust
import "lib/list" (tail, init, take, drop, slice)

nums = [1, 2, 3, 4, 5]

print(tail(nums))        // [2, 3, 4, 5] (без первого)
print(init(nums))        // [1, 2, 3, 4] (без последнего)
print(take(nums, 3))     // [1, 2, 3]
print(drop(nums, 2))     // [3, 4, 5]
print(slice(nums, 1, 4)) // [2, 3, 4]
```

### Трансформации

```rust
import "lib/list" (map, filter, reverse, unique, flatten)

nums = [1, 2, 3, 4, 5]

// map — применить функцию ко всем
doubled = map(fun(x) -> x * 2, nums)
print(doubled)  // [2, 4, 6, 8, 10]

// filter — оставить по условию
evens = filter(fun(x) -> x % 2 == 0, nums)
print(evens)  // [2, 4]

// reverse
print(reverse(nums))  // [5, 4, 3, 2, 1]

// unique — убрать дубликаты
print(unique([1, 2, 2, 3, 3, 3]))  // [1, 2, 3]

// flatten — развернуть вложенные списки
print(flatten([[1, 2], [3, 4], [5]]))  // [1, 2, 3, 4, 5]
```

### Pipes — элегантные цепочки

```rust
import "lib/list" (map, filter, take)

result = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
    |> filter(fun(x) -> x % 2 == 0)    // [2, 4, 6, 8, 10]
    |> map(fun(x) -> x * x)             // [4, 16, 36, 64, 100]
    |> fun(xs) -> take(xs, 3)           // [4, 16, 36]

print(result)  // [4, 16, 36]
```

Для функций с несколькими аргументами в pipes используйте лямбду-обёртку.

### Свёртка (fold)

```rust
import "lib/list" (foldl, foldr)

nums = [1, 2, 3, 4, 5]

// Сумма
sum = foldl(fun(acc, x) -> acc + x, 0, nums)
print(sum)  // 15

// Произведение
product = foldl(fun(acc, x) -> acc * x, 1, nums)
print(product)  // 120

// Максимум
max = foldl(fun(acc, x) -> if x > acc { x } else { acc }, nums[0], nums)
print(max)  // 5

// Конкатенация строк
words = ["Hello", " ", "World"]
sentence = foldl(fun(acc, w) -> acc ++ w, "", words)
print(sentence)  // "Hello World"
```

### Поиск

```rust
import "lib/list" (find, findIndex, indexOf, contains, any, all)

nums = [1, 2, 3, 4, 5]

// contains — есть ли элемент
print(contains(nums, 3))  // true

// find — найти по условию
print(find(fun(x) -> x > 3, nums))       // Some(4)
print(findIndex(fun(x) -> x > 3, nums))  // Some(3)
print(indexOf(nums, 3))                  // Some(2)

// any/all — предикаты
print(any(fun(x) -> x > 4, nums))   // true (есть хоть один > 4)
print(all(fun(x) -> x > 0, nums))   // true (все > 0)
```

### Условные операции

```rust
import "lib/list" (takeWhile, dropWhile, partition)

nums = [1, 2, 3, 4, 5, 4, 3, 2, 1]

// takeWhile — брать пока условие true
print(takeWhile(fun(x) -> x < 4, nums))  // [1, 2, 3]

// dropWhile — пропускать пока условие true
print(dropWhile(fun(x) -> x < 4, nums))  // [4, 5, 4, 3, 2, 1]

// partition — разделить на два списка
(small, big) = partition(fun(x) -> x < 3, nums)
print(small)  // [1, 2, 2, 1]
print(big)    // [3, 4, 5, 4, 3]
```

### Комбинирование

```rust
import "lib/list" (zip, unzip, concat)

// concat — объединить два списка
print(concat([1, 2], [3, 4]))  // [1, 2, 3, 4]

// или через ++
print([1, 2] ++ [3, 4])  // [1, 2, 3, 4]

// zip — объединить попарно
names = ["Alice", "Bob"]
ages = [30, 25]
pairs = zip(names, ages)
print(pairs)  // [("Alice", 30), ("Bob", 25)]

// unzip — обратно
(ns, as) = unzip(pairs)
print(ns)  // ["Alice", "Bob"]
print(as)  // [30, 25]
```

### Сортировка

```rust
import "lib/list" (sort, sortBy)

nums = [3, 1, 4, 1, 5, 9, 2, 6]

// sort (требует Order trait)
print(sort(nums))  // [1, 1, 2, 3, 4, 5, 6, 9]

// sortBy — с кастомным компаратором
desc = sortBy(nums, fun(a, b) -> b - a)
print(desc)  // [9, 6, 5, 4, 3, 2, 1, 1]

// Сортировка записей
users = [
    { name: "Alice", age: 30 },
    { name: "Bob", age: 25 },
    { name: "Carol", age: 35 }
]

byAge = sortBy(users, fun(a, b) -> a.age - b.age)
for u in byAge {
    print(u.name)  // Bob, Alice, Carol
}
```

### Итерация

```rust
import "lib/list" (forEach, zip, range)

// for loop
for x in [1, 2, 3] {
    print(x)
}

// for с индексом (через zip + range)
for pair in zip(range(0, 100), ["a", "b", "c"]) {
    (i, x) = pair
    print("Index " ++ show(i) ++ ": " ++ x)
}

// forEach — для side effects, возвращает Nil
forEach(fun(x) -> print(x), [10, 20, 30])
```

---

## Map (ассоциативный массив)

Map — это иммутабельный ассоциативный массив с ключами и значениями одного типа.

### Создание

```rust
import "lib/map" (mapSize)

// Литерал Map: %{ key => value }
scores = %{ "Alice" => 100, "Bob" => 85, "Carol" => 92 }
print(mapSize(scores))  // 3

// Пустой map (требует аннотации типа)
empty: Map<String, Int> = %{}

// Ключи разных типов
intKeys = %{ 1 => "one", 2 => "two", 3 => "three" }
```

Все значения в Map должны быть одного типа.

### Доступ

```rust
import "lib/map" (mapGet, mapGetOr, mapContains)

scores = %{ "Alice" => 100, "Bob" => 85 }

// mapGet — возвращает Option
print(mapGet(scores, "Alice"))    // Some(100)
print(mapGet(scores, "Unknown"))  // Zero

// mapGetOr — с дефолтным значением
print(mapGetOr(scores, "Alice", 0))    // 100
print(mapGetOr(scores, "Unknown", 0))  // 0

// mapContains — проверка наличия
print(mapContains(scores, "Alice"))  // true
print(mapContains(scores, "Dave"))   // false
```

### Модификация (иммутабельная)

Все операции возвращают новый Map, оригинал не меняется:

```rust
import "lib/map" (mapPut, mapRemove, mapMerge, mapGet, mapSize, mapContains)

scores = %{ "Alice" => 100, "Bob" => 85 }

// Добавить/обновить (возвращает новую мапу)
updated = mapPut(scores, "Charlie", 92)
print(mapSize(scores))   // 2 (оригинал не изменён)
print(mapSize(updated))  // 3

// Удалить ключ
smaller = mapRemove(scores, "Bob")
print(mapContains(smaller, "Bob"))  // false

// Объединить (второй выигрывает при конфликте)
m1 = %{ "a" => 1, "b" => 2 }
m2 = %{ "b" => 20, "c" => 3 }
merged = mapMerge(m1, m2)
print(mapGet(merged, "b"))  // Some(20) — из m2
```

### Итерация

```rust
import "lib/map" (mapKeys, mapValues, mapItems)

scores = %{ "Alice" => 100, "Bob" => 85, "Carol" => 92 }

// Ключи
print(mapKeys(scores))    // ["Bob", "Carol", "Alice"]

// Значения
print(mapValues(scores))  // [85, 92, 100]

// Пары (ключ, значение)
for pair in mapItems(scores) {
    (k, v) = pair
    print(k ++ ": " ++ show(v))
}
```

### Размер

```rust
import "lib/map" (mapSize)

scores = %{ "Alice" => 100, "Bob" => 85 }
print(mapSize(scores))  // 2
```

---

## Практические примеры

### Частотный словарь

```rust
import "lib/list" (foldl)
import "lib/map" (mapGetOr, mapPut)

fun frequency(items) {
    foldl(fun(acc, item) -> {
        count = mapGetOr(acc, item, 0)
        mapPut(acc, item, count + 1)
    }, %{}, items)
}

words = ["apple", "banana", "apple", "cherry", "banana", "apple"]
freq = frequency(words)
print(freq)
// %{"apple" => 3, "banana" => 2, "cherry" => 1}
```

### Группировка

```rust
import "lib/list" (foldl, length)
import "lib/map" (mapGetOr, mapPut)

fun groupByLength(words) {
    foldl(fun(acc, word) -> {
        key = length(word)
        existing = mapGetOr(acc, key, [])
        mapPut(acc, key, existing ++ [word])
    }, %{}, words)
}

words = ["hi", "hello", "hey", "world", "ok"]
byLen = groupByLength(words)
print(byLen)
// %{2 => ["hi", "ok"], 5 => ["hello", "world"], 3 => ["hey"]}
```

### Top-N элементов

```rust
import "lib/list" (sortBy, take)

fun topN(items, n, scoreFn) {
    sorted = sortBy(items, fun(a, b) -> scoreFn(b) - scoreFn(a))
    take(sorted, n)
}

products = [
    { name: "A", sales: 100 },
    { name: "B", sales: 500 },
    { name: "C", sales: 250 }
]

top2 = topN(products, 2, fun(p) -> p.sales)
for p in top2 {
    print(p.name ++ ": " ++ show(p.sales))
}
// B: 500
// C: 250
```

---

## Сводка функций lib/list

| Категория | Функции |
|-----------|---------|
| Доступ | head, headOr, last, lastOr, nth, nthOr |
| Подсписки | tail, init, take, drop, slice, takeWhile, dropWhile |
| Трансформации | map, filter, reverse, unique, flatten, sort, sortBy |
| Свёртка | foldl, foldr |
| Поиск | find, findIndex, indexOf, contains, any, all |
| Комбинирование | concat, zip, unzip, partition |
| Генерация | range |
| Итерация | forEach |
| Размер | length |

## Сводка функций lib/map

| Функция | Описание |
|---------|----------|
| mapGet | Получить значение (Option) |
| mapGetOr | Получить или дефолт |
| mapContains | Проверить наличие ключа |
| mapSize | Количество записей |
| mapPut | Добавить/обновить (новая Map) |
| mapRemove | Удалить ключ (новая Map) |
| mapMerge | Объединить две Map |
| mapKeys | Список ключей |
| mapValues | Список значений |
| mapItems | Список пар (ключ, значение) |
