# Tuples

Tuples are fixed-size, heterogeneous collections. Unlike lists, tuples can contain elements of different types, and their size is known at compile time.

## Creating Tuples

```rust
pair = (1, "hello")        // (Int, String)
triple = (1, 2, 3)         // (Int, Int, Int)
mixed = (true, 42, "test") // (Bool, Int, String)
nested = ((1, 2), (3, 4))  // ((Int, Int), (Int, Int))
```

## Accessing Elements

### Index Syntax (Recommended)

Access elements using `[]` syntax, same as lists:

```rust
pair = (1, "hello")

print(pair[0])   // 1
print(pair[1])   // "hello"

// Negative indexing
print(pair[-1])  // "hello" (last element)
print(pair[-2])  // 1
```

### Pattern Matching

Destructure tuples in assignments:

```rust
pair = (1, "hello")
(a, b) = pair
print(a)  // 1
print(b)  // "hello"

// In function parameters
fun printPair(p: (Int, String)) -> Nil {
    (x, y) = p
    print(x)
    print(y)
}
```

## lib/tuple Module

For more tuple operations, import `lib/tuple`:

```rust
import "lib/tuple" (*)
```

### Basic Access Functions

```rust
import "lib/tuple" (fst, snd, tupleGet)

pair = (1, "hello")

print(fst(pair))  // 1 — first element of pair
print(snd(pair))  // "hello" — second element of pair

triple = (10, 20, 30)
print(tupleGet(triple, 0))  // 10 — element by index
print(tupleGet(triple, 2))  // 30
```

### Transformation

```rust
import "lib/tuple" (tupleSwap, tupleDup)

// tupleSwap: (A, B) -> (B, A)
print(tupleSwap((1, "hello")))  // ("hello", 1)

// tupleDup: A -> (A, A)
print(tupleDup(42))  // (42, 42)
```

### Mapping

Apply functions to tuple elements:

```rust
import "lib/tuple" (mapFst, mapSnd, mapPair)

pair = (5, "test")

// mapFst: ((A) -> C, (A, B)) -> (C, B)
print(mapFst(fun(x) -> x * 2, pair))  // (10, "test")

// mapSnd: ((B) -> C, (A, B)) -> (A, C)
print(mapSnd(fun(s) -> s ++ "!", pair))  // (5, "test!")

// mapPair: ((A) -> C, (B) -> D, (A, B)) -> (C, D)
print(mapPair(fun(x) -> x * 2, fun(s) -> s ++ "!", pair))  // (10, "test!")
```

### Currying

Convert between tuple-taking and curried functions:

```rust
import "lib/tuple" (curry, uncurry)

// Function that takes a tuple
fun addPair(p: (Int, Int)) -> Int { p[0] + p[1] }

// curry: ((A, B) -> C) -> (A) -> (B) -> C
curriedAdd = curry(addPair)
print(curriedAdd(3)(4))  // 7

// Useful for partial application
add3 = curriedAdd(3)
print(add3(7))  // 10

// Curried function
fun addCurried(a: Int) -> (Int) -> Int { fun(b) -> a + b }

// uncurry: ((A) -> (B) -> C) -> (A, B) -> C
uncurriedAdd = uncurry(addCurried)
print(uncurriedAdd((10, 20)))  // 30
```

### Predicates

```rust
import "lib/tuple" (tupleBoth, tupleEither)

isPositive = fun(x: Int) -> x > 0

// tupleBoth: ((A) -> Bool, (A, A)) -> Bool
// True if predicate holds for both elements
print(tupleBoth(isPositive, (5, 10)))   // true
print(tupleBoth(isPositive, (5, -1)))   // false

// tupleEither: ((A) -> Bool, (A, A)) -> Bool
// True if predicate holds for at least one element
print(tupleEither(isPositive, (5, -1)))   // true
print(tupleEither(isPositive, (-5, -1)))  // false
```

## Working with Tuples and Other Functions

### With zip

```rust
import "lib/list" (zip)
import "lib/tuple" (fst, snd)

names = ["Alice", "Bob"]
ages = [30, 25]

// zip creates list of tuples
people = zip(names, ages)  // [("Alice", 30), ("Bob", 25)]

// Access tuple elements
firstPerson = people[0]    // ("Alice", 30)
print(fst(firstPerson))    // "Alice"
print(firstPerson[1])      // 30
```

### With Pipe Operator

```rust
import "lib/tuple" (tupleSwap, mapFst)

result = (1, 2) |> tupleSwap |> mapFst(fun(x) -> x * 10)
print(result)  // (20, 1)
```

## Type Annotations

Tuple types are written with parentheses:

```rust
pair: (Int, String) = (1, "hello")
triple: (Int, Int, Int) = (1, 2, 3)

fun process(data: (String, Int)) -> String {
    (name, value) = data
    name ++ ": " ++ show(value)
}
```

## Summary

| Function | Type | Description |
|----------|------|-------------|
| `tuple[i]` | `Tuple -> T` | Index access (syntax) |
| `fst` | `(A, B) -> A` | First element |
| `snd` | `(A, B) -> B` | Second element |
| `get` | `(Tuple, Int) -> T` | Element by index |
| `swap` | `(A, B) -> (B, A)` | Swap elements |
| `dup` | `A -> (A, A)` | Duplicate into pair |
| `mapFst` | `((A) -> C, (A, B)) -> (C, B)` | Map first |
| `mapSnd` | `((B) -> C, (A, B)) -> (A, C)` | Map second |
| `mapPair` | `((A) -> C, (B) -> D, (A, B)) -> (C, D)` | Map both |
| `curry` | `((A, B) -> C) -> (A) -> (B) -> C` | Curry tuple function |
| `uncurry` | `((A) -> (B) -> C) -> (A, B) -> C` | Uncurry to tuple |
| `both` | `((A) -> Bool, (A, A)) -> Bool` | Both satisfy |
| `either` | `((A) -> Bool, (A, A)) -> Bool` | Either satisfies |

## See Also

- [Lists](14_builtins.md) - List operations
- [Pattern Matching](07_pattern_matching.md) - Destructuring tuples
- [Functions](04_functions.md) - Partial application

