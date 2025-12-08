# Iteration 6: Custom Types

This iteration introduces powerful type system features: Type Aliases, Records, and Algebraic Data Types (ADTs).

## Type Aliases

Aliases create a new name for an existing type. They are interchangeable with the original type but help with readability.

```rust
type alias Money = Float

fun format(m: Money) -> String {
    "$" ++ show(m)
}

print(format(19.99))  // "$19.99"
```

## Records

Records are structural types with named fields.

```rust
// Type definition
type Point = { x: Int, y: Int }

// Literal
p = { x: 10, y: 20 }

// Access
print(p.x)  // 10
```

## Methods

You can define methods on custom types (Records or ADTs) using the `fun (receiver: Type) name()` syntax.

```rust
type Point = { x: Int, y: Int }

// Define a method 'dist_sq' on Point
fun (p: Point) dist_sq() -> Int {
    p.x * p.x + p.y * p.y
}

p: Point = { x: 3, y: 4 }
print(p.dist_sq()) // 25
```

## Algebraic Data Types (ADTs)

ADTs allow you to define types that can be one of several variants.

Parameters for generic ADTs are listed in angle brackets after the type name.

```rust
// Option is built-in, but custom ADTs look like:
type MyOption<T> = MySome(T) | MyNone

x = MySome(42)
y = MyNone
print(x)  // MySome(42)
```

Recursive ADTs are supported (e.g., for Lists or Trees).

```rust
// List is built-in, but custom recursive ADTs look like:
type MyList<T> = MyCons((T, MyList<T>)) | MyEmpty

list = MyCons((1, MyCons((2, MyEmpty))))
print(list)
```

## Проверка типов во время выполнения

### typeOf

Функция `typeOf(value, Type) -> Bool` проверяет тип значения:

```rust
x = 10
typeOf(x, Int)      // true
typeOf(x, String)   // false

name = "Alice"
typeOf(name, String)  // true
```

### Параметризованные типы

Для проверки параметризованных типов используйте **круглые скобки** (не угловые!):

```rust
type MyOption<T> = Yes T | NoVal

o = Yes(42)

// Проверка без параметра — любой MyOption
typeOf(o, MyOption)       // true

// Проверка с параметром — конкретный MyOption<Int>
typeOf(o, MyOption(Int))  // true
typeOf(o, MyOption(String))  // false
```

**Важно:** В выражениях `Type(Param)`, а не `Type<Param>`:
```rust
// Правильно:
typeOf(list, List)
typeOf(opt, Option(Int))

// Ошибка синтаксиса:
typeOf(list, List<Int>)  // угловые скобки не работают!
```

### getType

Функция `getType(value) -> Type` возвращает тип значения:

```rust
x = 42
t1 = getType(x)
print(t1)  // type(Int)

f = fun(a: Int) -> a * 2
t2 = getType(f)
print(t2)  // type((Int) -> Int)
```

## Runtime Representation

*   **Aliases**: Resolved to underlying type at analysis time.
*   **Records**: Represented as `RecordInstance` with named fields.
*   **ADTs**: Represented as `DataInstance` with constructor name and field values.
