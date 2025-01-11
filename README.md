# Lox Interpreter in Go


## Grammar for the lox language

### How to read

- Each rule is in the form `<rule_name> -> symbols`.
- There are two types of symbols:
  - Terminal symbols are the characters that make up the language. They're tokens from the language grammar. They're either in double quotes or caps if they're referring to a literal(`STRING`, `NUMBER`, `IDENTIFIER`).
  - Non-terminal symbols recursively refer to other rules. It leads to composition of the rules to make the grammar.
- There are some postfix and binary operators based on regex to simplify writing the grammar -
  - `+` means one or more of the preceding symbol
  - `*` means zero or more of the preceding symbol
  - `?` means zero or one of the preceding symbol
  - `|` means one of the symbols on either side
  - `(` and `)` are used to group symbols
- The lines are in order of precedence. Each rule only matches expressions at
  its precedence level or higher.

```ebnf
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;
```
