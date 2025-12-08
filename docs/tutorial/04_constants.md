# Constants

Constants are named values that cannot be changed after they are defined. They are declared using the `:-` operator.

## Syntax

```rust
name :- value
name : Type :- value
```

## Examples

```rust
// Inferred type (Float)
pi :- 3.14159

// Explicit type
max_retries: Int :- 5

// String constant
app_name :- "My Application"

fun area(r: Float) -> Float {
    pi * r * r
}

print(area(2.0))
print(app_name)
```

## Semantics

*   **Immutability**: You cannot assign a new value to a constant.
*   **Single Definition**: A constant can only be defined once. You cannot redefine an existing variable as constant.
*   **Scope**: Constants obey standard scoping rules. Top-level constants are visible throughout the package.
*   **Evaluation**: Constants are evaluated once when the module is loaded.

## Difference from Variables

```rust
// Variable (mutable) — can be reassigned
x = 10
x = 20       // OK
print(x)     // 20

// Constant (immutable) — cannot be reassigned
y :- 10
print(y)     // 10
// y = 20   // Would cause: Error: cannot reassign constant 'y'

// Cannot redefine variable as constant
z = 10
// z :- 20  // Would cause: Error: redefinition of symbol 'z'
```

## Tuple Unpacking

Constants support tuple unpacking:

```rust
pair :- (1, "hello")
(a, b) :- pair      // a = 1, b = "hello"

// Nested unpacking
nested :- ((1, 2), 3)
((x, y), z) :- nested

// Wildcard for unused parts
(first, _) :- pair  // only binds first
```

## Note on Naming
Constants typically use lowerCamelCase or snake_case, just like variables. Capitalized identifiers are reserved for Types and Constructors.

