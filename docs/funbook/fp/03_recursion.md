# 03. Рекурсия и TCO

## Задача
Использовать рекурсию эффективно без переполнения стека.

## Простая рекурсия

```rust
fun factorial(n: Int) -> Int {
    if n <= 1 { 1 } else { n * factorial(n - 1) }
}

print(factorial(5))   // 120
print(factorial(10))  // 3628800
```

## Проблема: переполнение стека

```rust
// Эта версия НЕ оптимизирована - создаёт n стек-фреймов
fun badSum(n: Int) -> Int {
    if n == 0 { 0 } else { n + badSum(n - 1) }
}

// При больших n получим stack overflow!

```

## TCO (Tail Call Optimization)

Funxy оптимизирует хвостовые вызовы — когда рекурсивный вызов является последней операцией.

```rust
// Хвостовая рекурсия с аккумулятором
fun factorialTCO(n: Int, acc: Int) -> Int {
    if n <= 1 { acc } else { factorialTCO(n - 1, n * acc) }
}

fun factorial(n: Int) -> Int { factorialTCO(n, 1) }

print(factorial(20))  // 2432902008176640000 - работает!
```

## Объяснение TCO

```rust
// НЕ хвостовой вызов (после рекурсии есть операция *)
fun bad(n: Int) -> Int {
    if n <= 1 { 1 } else { n * bad(n - 1) }
}
//                         ^^^ умножение ПОСЛЕ рекурсии

// Хвостовой вызов (рекурсия - последняя операция)
fun good(n: Int, acc: Int) -> Int {
    if n <= 1 { acc } else { good(n - 1, n * acc) }
}
//                          ^^^ ничего после рекурсии

```

## Примеры с TCO

### Сумма списка

```rust
fun sumList(xs, acc) {
    match xs {
        [] -> acc
        [head, tail...] -> sumList(tail, acc + head)
    }
}

total = sumList([1, 2, 3, 4, 5], 0)
print(total)  // 15
```

### Длина списка

```rust
fun listLength(xs, acc: Int) -> Int {
    match xs {
        [] -> acc
        [_, tail...] -> listLength(tail, acc + 1)
    }
}

print(listLength([1, 2, 3, 4, 5], 0))  // 5
```

### Reverse списка

```rust
fun myReverse(xs, acc) {
    match xs {
        [] -> acc
        [head, tail...] -> myReverse(tail, [head] ++ acc)
    }
}

print(myReverse([1, 2, 3, 4, 5], []))  // [5, 4, 3, 2, 1]
```

### Фибоначчи

```rust
// TCO версия - линейная сложность
fun fibTCO(n: Int, a: Int, b: Int) -> Int {
    if n == 0 { a } else { fibTCO(n - 1, b, a + b) }
}

fun fib(n: Int) -> Int { fibTCO(n, 0, 1) }

print(fib(40))  // 102334155 - мгновенно!
```

## Pattern Matching + Рекурсия

```rust
type Tree = Leaf(Int) | Node((Tree, Tree))

fun treeSum(t: Tree) -> Int {
    match t {
        Leaf(n) -> n
        Node((left, right)) -> treeSum(left) + treeSum(right)
    }
}

tree = Node((
    Node((Leaf(1), Leaf(2))),
    Node((Leaf(3), Leaf(4)))
))
print(treeSum(tree))  // 10
```

## Взаимная рекурсия

```rust
fun isEven(n: Int) -> Bool {
    if n == 0 { true } else { isOdd(n - 1) }
}

fun isOdd(n: Int) -> Bool {
    if n == 0 { false } else { isEven(n - 1) }
}

print(isEven(100))  // true
print(isOdd(99))    // true
```

## Практический пример: обход дерева

```rust
import "lib/list" (map, foldl, flatten)

type FileTree = File((String, Int))
              | Dir((String, List<FileTree>))

fun totalSize(tree: FileTree) -> Int {
    match tree {
        File((name, size)) -> size
        Dir((name, children)) -> foldl(fun(acc, c) -> acc + totalSize(c), 0, children)
    }
}

fun findLargeFiles(tree: FileTree, threshold: Int) -> List<String> {
    match tree {
        File((name, size)) -> if size > threshold { [name] } else { [] }
        Dir((name, children)) -> flatten(map(fun(c) -> findLargeFiles(c, threshold), children))
    }
}

fs = Dir(("root", [
    File(("a.txt", 100)),
    Dir(("sub", [
        File(("b.txt", 500)),
        File(("c.txt", 50))
    ])),
    File(("d.txt", 200))
]))

print(totalSize(fs))              // 850
print(findLargeFiles(fs, 150))    // ["b.txt", "d.txt"]
```
