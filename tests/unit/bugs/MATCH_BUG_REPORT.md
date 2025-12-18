# Match Expression Bug Report

## Summary
When a variable is assigned from a `match` expression, subsequent operations on that variable are parsed incorrectly, causing the infix operator and right operand (or function arguments) to be lost. The result is just the variable itself.

## Issue
- **Bug File**: `tests/unit/bugs/match_infix_bug_test.lang`
- **Status**: Open / Known Bug
- **Severity**: Critical (breaks code functionality)

## Reproduction

### Simple Case: String Concatenation
```lang
val = match "hello" {
    s: String -> s
    _ -> "Unknown"
}

result = val ++ " world"
// Expected: "hello world"
// Actual: "hello" (right operand and operator lost!)
```

### Function Call
```lang
fun echo(msg) { "Echo: " ++ msg }

msg = match "test" {
    s: String -> s
    _ -> "Unknown"
}

result = echo(msg)
// Expected: "Echo: test"
// Actual: "Echo: " (arguments lost!)
```

### Multiple Match Variables
```lang
x = match "x" { s: String -> s, _ -> "Unknown" }
y = match "y" { s: String -> s, _ -> "Unknown" }

result = x ++ y
// Expected: "xy"
// Actual: "x" (right operand lost!)
```

## Root Cause Analysis

The bug appears to be in the parser's handling of `match` expressions and the subsequent token stream. When a variable is assigned from a match expression:

1. The parser correctly parses the match expression
2. The variable is assigned the match result correctly
3. However, when that variable is used in subsequent operations, the parser seems to **lose the right operand or arguments**

This suggests the parser's state after match parsing leaves it in a position where:
- Infix operators (`++`) are not properly recognized, or
- Arguments in function calls are not properly parsed

### Related Code Locations
- `parser/expressions.go` - `parseMatchExpression()`
- `parser/patterns.go` - `parseMatchArm()`
- `parser/parser.go` - `disallowTrailingLambda` flag management
- Possibly: postfix/infix operator handling

## What's Lost

- **In string concatenation**: Right operand (e.g., `" world"`) is lost
- **In function calls**: Arguments after the first operand are lost
- **Pattern**: Any complex operation after `match` variable fails

## Workarounds

### ❌ Reassign the match result (DOESN'T WORK!)
```lang
temp = match "x" { s: String -> s, _ -> "Unknown" }
val = temp  // Reassign - still doesn't work!
result = val ++ " world"  // Still returns just "x"
```

### ✅ Only Working Workaround: Avoid match entirely
```lang
// This is the ONLY confirmed working approach - completely avoid match in operations
// Instead of:
// x = match ... { ... }
// result = x ++ "world"

// Use static values or values from non-match sources:
result = "hello" ++ " world"  // Works!
```

### ❌ Extra reassignments (DOESN'T WORK)
Even multiple reassignments don't help - the bug persists once a variable is tainted by match:
```lang
val = match "test" { s: String -> s, _ -> "Unknown" }
temp = val  // Extra reassignment - doesn't help
result = temp ++ "!"  // Still broken
```

## Test Status

### Expected Fail (Demonstrating the Bug)
- ⚠ `Bug: Match variable in string concatenation` - RIGHT OPERAND LOST
- ⚠ `Bug: Match variable in function call` - ARGUMENTS LOST
- ⚠ `Bug: Multiple operations on match variable` - RIGHT OPERAND LOST
- ⚠ `Bug: Match variable in sprintf` - ARGUMENTS LOST

### Should Pass (Workarounds)
- ✓ `Workaround: Avoid match in assignment chain` - Static values work correctly

## Next Steps

1. **Debug parser state**: Add debug output to trace `disallowTrailingLambda` and operator precedence after match
2. **Investigate infix parsing**: Check if `parseInfixExpression` properly handles variables from match
3. **Fix state management**: Ensure parser state is properly restored after match parsing
4. **Regression testing**: Once fixed, convert all `testExpectFail` to `testRun`
5. **Performance check**: Ensure fix doesn't break other parser functionality

## Additional Notes

- The bug is reproducible with ANY `match` expression followed by operations on the result
- The bug persists even with reassignment - once a variable is from match, it's "tainted"
- Direct operations (without going through match) work fine
- The bug affects BOTH infix operators AND function arguments
- This is a **CRITICAL parser bug** that breaks fundamental language functionality


