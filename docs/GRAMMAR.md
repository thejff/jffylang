# Grammar definition


## Expressions

|Grammar    |Definition                                                                                                     |
|-------    |----------                                                                                                     |
|expression | literal &#124; unary &#124; binary &#124; grouping                                                            |
|literal    | NUMBER &#124; STRING &#124; "true" &#124; "false" &#124; "nil"                                                |
|grouping   | "(" expression ")"                                                                                            |
|unary      | ( "-" &#124; "!" ) expression                                                                                 |
|binary     | expresion operator expression                                                                                 |
|operator   | "==" &#124; "!=" &#124; "<" &#124; "<=" &#124; ">" &#124; ">=" &#124; "+" &#124; "-" &#124; "*" &#124; "/"    |


## Expression parsing grammar

### Key
* Terminal      - Match and consume a token
* Nonterminal   - Call to rule's function
* &#124;        - `if` or `switch`
* \* or +       - `for` loop
* ?             - `if`

|Grammar    |Definition                                                                                     |
|-------    |----------                                                                                     |
|expression | -> assignment                                                                                 |
|assignment | -> IDENTIFIER "=" assignment  &#124; logic_or ;                                                |
|logic_or   | -> logic_and ( "or" logic_and )* ; |
|logic_and  | -> equality ( "and" equality )* ; |
|equality   | -> comparison ( ( "!=" &#124; "==" ) comparison )*                                            |
|comparison | -> term ( ( ">" &#124; ">=" &#124; "<" &#124; "<=") term )*                                   |
|term       | -> factor ( ( "-" &#124; "+" ) factor )*                                                      |
|factor     | -> unary ( ( "/" &#124; "\*" ) unary )*                                                       |
|unary      | -> ( "!" &#124; "-") unary &#124; primary                                                     |
|primary    | -> NUMBER &#124; STRING &#124; "true" &#124; "false" &#124; "nil" &#124; "(" expression ")"   |


## Definition grammar

|Grammar    |Definition                                                                                                         |
|-------    |----------                                                                                                         |
|program    | -> declaration* EOF                                                                                               |
|declaration| -> varDecl &#124; statement                                                                                       |
|statement  | -> exprStmt &#124; forStmt &#124; ifStmt &#124; printStmt &#124; whileStmt &#124; "break;" &#124; "continue;" &#124; block ;|
|whileStmt  | -> "while" expression block ; |
|forStmt    | -> "for" ( varDecl &#124; exprStmt &#124; ";" ) expression? ";" expression? block ; |
|ifStmt     | -> "if" expression block "else" block ? ;
|printStmt  | -> "print" ( STRING &#124; NUMBER ) ";"
|block      | -> "{" declaration* "}" ;  
|varDecl    | -> "var" IDENTIFIER ( "=" expression )? ";" ;                                                                     |
|primary    | -> "true" &#124; "false" &#124; "nil" &#124; NUMBER &#124; STRING &#124; "(" expression ")" &#124; IDENTIFIER ;   |

