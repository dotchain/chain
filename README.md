# chain

Chain is a distributed functional reactive language built on top of DOT

## Priorities

1. Code is easy to read and modify
    * Immutable values: Easier to reason with and composes better
    * Expression syntax: Lack of statements make syntax much simpler
    * No explicit types: Easier to work with
2. Built for efficiency of coding
    * Not meant for very large code bases
    * Expect lots of small snippets
    * Embed code in markdown and other places
3. Runtime performance and type safety are "profiles"
    * Manual configuration (such as choosing data types) possible
    * Separate from main algorithm concerns
    * Profile-guided optimizations are the preferred route
4. Data provenance is easily available
    * Like call stacks, but for data
    * Both forward and backwards looking (dependents and dependencies)
5. Meta-programming
    * AST rewriting via interpreted code written in chain itself
    * Debugging hints and helpers for such rewrites
6. FRP/Incremental computation

## Language

1. Immutable values.  No facility to mutate, no variables (just names)
2. Composition using objects and lists:
    * `(x = 2, y = 3)` is an object with fields `x` and `y`
    * `z.x` access the `x` field of `z`
    * `(1, 2, 3)` is a list
    * Special syntax for empty list and object: `(,)` and `(=)`
    * Special syntax for single element list `(1,)`
    * List has `item` for element access `(1, 2, 3).item(2)`
2. Lexical scoping for the most part: `(x = y + 2, y = 5)`
    * fields in objects are available to all expressions within
    * nesting is possible: `(x = (z = y + 1).z, y = 5)`
    * recursion is not allowed except with functions
    * functions are like so: `(inc(x) = x + 1, z: inc(2))`
3. Dynamic scoping for lambdas: `list.filter(value < 42)`
    * lists and objects expose `filter` method which takes an
    expression
    * `value` represents the current value in the list/object
    * `key` represents the index (list) or key (object)
    * these values always come from the dynamic scope
4. Everything is an expression, no special statements
5. Edits are possible via session and branch:
    * builds upon Streams (pull-based FRP)
    * `session(variable, expression)`
    * expression can use the dynamic scoped name `value` to refer to
    the sesion variable
    * any mutation of this `value` is recorded with the session
    * branches are possible with `branch(variable, expression)`
       * edits within expression are saved locally until call to
       `push(variable, return_value)` or `pull(,)`
6. Objects are not dictionaries.  They are expected to mostly have
fixed schema with no support for fetching fields or modifying them
dynamically like dictionaries.  i.e. `x.(f())` is not expected to
work.  The language does support the syntax and some extensions may
actually support this behavior but this is not the intention.
7. No special syntax for commeents and annotations
    * A global function `note(comment, expr) = expr` can be used
    * example: `note("calculate distance", sqrt(x*x + y*y))`
    * Configuration file is expected to be "patches" on the code and
    can refer to these `note` calls and replace them as part of the
    build. The language will also support these comments showing up in
    the call stack
9. Unicode and rich-text support is built in:
    * text:rich(text:bold("hello"), " ", text:italic("world")) for
    instance. Also regular unicode formatting will be honored.
    * exact library for this is t.b.d.  Callstacks and other objects
    can show these as well.
10. unit testing is part of the langauge.  Code be annotated to
define what tests are to be fired to validate it.  This is likely
similar to `note(...)` in that it will define a name and some file
elsewehere will define the actual tests to match that name.

## Compiler

The base parser is in Golang and it is very likely this will get
ported to chain itself (as it is relatively small).

## DOT integration

The plan is to build an interpreter for chain that uses DOT values as
well as implement session semantics using DOT Streams.

## Other interpreters and compilers

No specific plans yet for this.