# Grammar definition


## Expressions

|Grammar|Definition|
|-------|----------|
|expression| literal &#124; unary &#124; binary &#124; grouping|
|literal| NUMBER &#124; STRING &#124; "true" &#124; "false" &#124; "nil"|
|grouping| "(" expression ")"|
|unary| ( "-" &#124; "!" ) expression|
|binary| expresion operator expression|
|operator| "==" &#124; "!=" &#124; "<" &#124; "<=" &#124; ">" &#124; ">=" &#124; "+" &#124; "-" &#124; "*" &#124; "/" |


## Expression parsing grammar

### Key
* Terminal      - Match and consume a token
* Nonterminal   - Call to rule's function
* &#124;        - `if` or `switch`
* \* or +       - `for` loop
* ?             - `if`

|Grammar|Definition|
|-------|----------|
|expression| -> equality|
|equality| -> comparison ( ( "!=" &#124; "==" ) comparison )*|
|comparison| -> term ( ( ">" &#124; ">=" &#124; "<" &#124; "<=") term )*|
|term| -> factor ( ( "-" &#124; "+" ) factor )*|
|factor| -> unary ( ( "/" &#124; "\*" ) unary )*|
|unary| -> ( "!" &#124; "-") unary &#124; primary|
|primary| -> NUMBER &#124; STRING &#124; "true" &#124; "false" &#124; "nil" &#124; "(" expression ")"|

