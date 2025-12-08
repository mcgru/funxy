# Iteration 5: Generics

Generics allow you to write code that works with different types, increasing flexibility and type safety.

## Naming Convention

**Important**: Type parameters must start with an **uppercase letter**.

```rust
// ✓ T, U are uppercase - correct
fun swap<T, U>(pair: (T, U)) -> (U, T) {
    (a, b) = pair
    (b, a)
}

// ✗ Lowercase type parameters would be an error:
// fun swap<t, u>(pair: (t, u)) -> (u, t) { ... }
print(swap((1, "hello")))  // ("hello", 1)
```

This follows the language-wide convention:
- **Uppercase**: types, constructors, traits, type parameters (`Int`, `Some`, `Order`, `T`)
- **Lowercase**: values, functions, variables (`myVar`, `calculate`, `x`)

## Generic Functions

You can define functions that accept type parameters using angle brackets `<T>`.

```rust
// Identity function working on any type T
fun id<T>(x: T) -> T {
    x
}

n = id(42)       // T is Int
s = id("hello")  // T is String
print(n)         // 42
print(s)         // "hello"
```

## Generic Types

Type declarations can also be generic. Type parameters are listed in angle brackets after the type name.

```rust
// A simple wrapper type
type Box<T> = { value: T }

b = { value: 10 }  // Box<Int>
print(b.value)     // 10
```

## Type Inference

The type system infers concrete types at call sites:

```rust
fun id<T>(x: T) -> T { x }

id(42)      // T inferred as Int
id("hello") // T inferred as String
```

