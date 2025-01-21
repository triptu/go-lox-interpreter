# [Lox](https://craftinginterpreters.com/the-lox-language.html) Interpreter in Go

## Running the program

### Tokenize

Prints the tokens array for the source code.

```sh
./run.sh tokenize <filename>
```

### Parse

Parses the tokens array and prints the AST. This is not very clean, the AST is printed as a complicated looking string in the form `(<operator> <left> <right>)`. It's better to use the visualiser to see the AST.

```sh
./run.sh parse <filename>
```

### Visualise

Visualises the AST by creating a DOT file and then generating a PNG image with Graphviz.

```sh
./run.sh visualize <filename>
```

### Evaluate

Evaluates the AST and prints the result.

```sh
./run.sh evaluate <filename>
```


## Running tests

```sh
python test.py <suite_name> <filter_path>
```

example

```sh
> python test.py golox
> python test.py chap10_functions
```

`golox` runs all tests. Or else you can filter tests upto a specific chapter.

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
(* program is basically a list of statements *)
program        → declaration* EOF ;

(* declare variables, classes and functions *)
declaration    → classDecl
               | varDecl
               | funDecl
               | statement ;

classDecl      → "class" IDENTIFIER "{" function* "}" ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
funDecl        → "fun" function ;
function       → IDENTIFIER "(" parameters? ")" blockStmt ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;


statement      → exprStmt
               | printStmt
               | ifStmt
               | whileStmt
               | forStmt
               | returnStmt
               | blockStmt ;

exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;
ifStmt         → "if" "(" expression ")" statement
               ( "else" statement )? ;
whileStmt      → "while" "(" expression ")" statement ;
forStmt        → "for" "(" (varDecl | exprStmt | ";")
               expression? ";" 
               expression? ")" statement ;
blockStmt      → "{" declaration* "}" ;
returnStmt     → "return" expression? ";" ;

(* define expressions in order of precedence *)
expression     → assignment ;
(* a = 2 or breakfast.milk.sugar = 4 *)
assignment     → ( call "." )? IDENTIFIER "=" assignment
               | logic_or ;
logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | call ;
(* If there are no parentheses, this parses a bare primary expression. *)
(* Otherwise, there can be multiple layers of calls, like abc()() *)
(* and field access or both, like myClass.pqr().abc()() *)
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" 
               | IDENTIFIER ;

(* helper rules *)
arguments      → expression ( "," expression )* ;
```


