# Iteration 7: Pattern Matching

Pattern matching is a powerful mechanism for checking a value against a pattern and, if it matches, destructuring it into constituent parts.

## Basic Matching

The `match` expression compares a value against a series of patterns.

```rust
x = 1
match x {
    1 -> print("One"),
    2 -> print("Two"),
    _ -> print("Other") // Wildcard pattern
}
```

## Literals

You can match against literals like integers, booleans, and strings.

```rust
match true {
    true -> print("True"),
    false -> print("False")
}

match "hello" {
    "hello" -> print("Greeting"),
    _ -> print("Unknown")
}
```

## Variables (Bindings)

Identifier patterns bind the matched value to a variable.

```rust
match 42 {
    val -> print(val) // val is 42
}
```

## Tuples

You can destructure tuples.

```rust
pair = (10, 20)
match pair {
    (0, 0) -> print("Origin"),
    (x, 0) -> print("X axis"),
    (0, y) -> print("Y axis"),
    (x, y) -> print("Point")
}
```

You can also use the spread operator `...` in tuple patterns to match the rest of the elements.

```rust
t = (1, 2, 3, 4)
match t {
    (head, tail...) -> {
        print(head) // 1
        print(tail) // (2, 3, 4)
    }
}
```

## Lists

Lists can be matched structurally.

```rust
l = [1, 2, 3]
match l {
    [] -> print("Empty"),
    [head, tail...] -> {
        print(head) // 1
        print(tail) // [2, 3]
    }
}
```

Fixed size matching:

```rust
match [1, 2] {
    [a, b] -> print(a + b),
    _ -> print("Not a pair")
}
```

## Records

You can match against record fields.

```rust
r = { x: 10, y: 20, z: 30 }
match r {
    { x: 0, y: 0 } -> print("Origin"), // Partial match on fields (Row Polymorphism)
    { x: val } -> print(val)
}
```

## Constructors (ADTs)

Pattern matching is the primary way to work with ADTs.

```rust
// Using built-in Option type
opt = Some(42)
match opt {
    Some(val) -> print(val)
    Zero -> print("Nothing")
}
```

## Nested Patterns

Patterns can be nested arbitrarily deep.

```rust
// Using built-in List type
list = [1, 2, 3]

match list {
    [1, x, rest...] -> print(x)  // Matches [1, 2, 3], x binds to 2
    _ -> print("No match")
}
```

Nested spread patterns:

```rust
match (1, [2, 3]) {
    (x, [y, z...]) -> {
        print(x) // 1
        print(y) // 2
        print(z) // [3]
    }
    _ -> print("No match")
}
```

## Guard Patterns

Guards add conditions to patterns using `if`. The arm is only executed if both the pattern matches AND the guard evaluates to `true`.

### Basic Guards

```rust
fun classify(n: Int) -> String {
    match n {
        x if x > 0 -> "positive"
        x if x < 0 -> "negative"
        _ -> "zero"
    }
}

print(classify(5))   // "positive"
print(classify(-3))  // "negative"
print(classify(0))   // "zero"
```

### FizzBuzz Example

Guards are perfect for FizzBuzz-style logic:

```rust
fun fizzbuzz(n: Int) -> String {
    match n {
        x if x % 15 == 0 -> "FizzBuzz"
        x if x % 3 == 0 -> "Fizz"
        x if x % 5 == 0 -> "Buzz"
        x -> show(x)
    }
}
```

### Guards with Destructuring

Guards can use variables bound by the pattern:

```rust
fun comparePair(pair: (Int, Int)) -> String {
    match pair {
        (a, b) if a == b -> "equal"
        (a, b) if a > b -> "first is greater"
        (a, b) -> "second is greater"
    }
}
```

### Guards with Lists

```rust
fun findFirstPositive(xs: List<Int>) -> Option<Int> {
    match xs {
        [] -> Zero
        [x, rest...] if x > 0 -> Some(x)
        [_, rest...] -> findFirstPositive(rest)
    }
}

print(findFirstPositive([-1, -2, 3, 4]))  // Some(3)
print(findFirstPositive([-1, -2]))        // Zero
```

### Guards with ADTs

```rust
type Shape = Circle Float | Rectangle Float Float

fun area(s: Shape) -> Float {
    match s {
        Circle(r) if r <= 0.0 -> 0.0
        Circle(r) -> 3.14159 * r * r
        Rectangle(w, h) if w <= 0.0 || h <= 0.0 -> 0.0
        Rectangle(w, h) -> w * h
    }
}
```

### Complex Guard Expressions

Guards can use any boolean expression:

```rust
fun inRange(n: Int, low: Int, high: Int) -> Bool {
    match n {
        x if x >= low && x <= high -> true
        _ -> false
    }
}
```

### Important Notes

1. **Guard must be Bool**: The guard expression must evaluate to a boolean.

```rust
match 5 {
    n if n + 1 -> "oops"  // ERROR: guard expression must be Bool, got Int
    _ -> "ok"
}
```

2. **Guards don't affect exhaustiveness**: The compiler cannot prove that guards cover all cases. Always include a catch-all `_` pattern when using guards.

3. **Order matters**: Arms are checked top-to-bottom. Put more specific guards before general ones.

## String Patterns with Captures

String patterns allow you to match strings with dynamic segments and capture their values. This is especially useful for URL routing and parsing.

### Basic Capture

Use `{name}` to capture a segment (up to the next `/` or end of string):

```rust
path = "/hello/world"
match path {
    "/hello/{name}" -> print("Hello " ++ name)  // Hello world
    _ -> print("Not found")
}
```

### Multiple Captures

Capture multiple segments in one pattern:

```rust
path = "/users/42/posts/123"
match path {
    "/users/{userId}/posts/{postId}" -> {
        print("User: " ++ userId)    // User: 42
        print("Post: " ++ postId)    // Post: 123
    }
    _ -> print("Not found")
}
```

### Greedy Capture

Use `{name...}` to capture the entire remaining path:

```rust
path = "/static/css/main/style.css"
match path {
    "/static/{file...}" -> print("Serving: " ++ file)  // Serving: css/main/style.css
    _ -> print("Not found")
}
```

### HTTP Routing Example

Combine with tuple matching for HTTP routing:

```rust
import "lib/http" (*)
import "lib/json" (jsonEncode)

fun handler(req) {
    match (req.method, req.path) {
        ("GET", "/") -> {
            status: 200,
            body: "Welcome!"
        }
        ("GET", "/users/{id}") -> {
            status: 200,
            body: jsonEncode({ userId: id })
        }
        ("GET", "/users/{id}/posts/{postId}") -> {
            status: 200,
            body: "User " ++ id ++ ", Post " ++ postId
        }
        ("GET", "/static/{file...}") -> {
            status: 200,
            body: "File: " ++ file
        }
        ("POST", "/users") -> {
            status: 201,
            body: "Created"
        }
        _ -> {
            status: 404,
            body: "Not found"
        }
    }
}

httpServe(8080, handler)
```

### Pattern Priority

Literal patterns take precedence over capture patterns:

```rust
match "/exact/match" {
    "/exact/match" -> "literal"      // This matches first
    "/exact/{name}" -> "captured"
    _ -> "other"
}
// Result: "literal"
```

### Empty Captures

Captures can be empty:

```rust
match "/prefix/" {
    "/prefix/{suffix}" -> "got: [" ++ suffix ++ "]"  // got: []
    _ -> "no match"
}
