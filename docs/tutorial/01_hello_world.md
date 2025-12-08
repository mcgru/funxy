# Hello World

Первая программа на Funxy.

## Вывод на экран

Функция `print` выводит значения в stdout:

```rust
print("Hello, World!")
```

## Строковые литералы

**Обычные строки** — в двойных кавычках:

```rust
message = "Hello, World!"
print(message)
```

**Многострочные строки** — в обратных кавычках:

```rust
text = `This is a
multi-line
string`

json = `{"name": "test", "value": 42}`
print(json)
```

Особенности raw-строк:
- Могут содержать переносы строк
- Не обрабатывают escape-последовательности
- Удобны для JSON, SQL, шаблонов

## Запуск

```bash
./funxy hello.lang
```

## Тесты

См. `tests/hello_world.lang`
