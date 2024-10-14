# Codeword Solver

Solves codeword puzzles that are found in newspapers, e.g. https://puzzles.independent.co.uk/games/independent-codeword

### Build

```
go build
```

## Usage

Enter the letters that you already have, with . or ? or * to represent unique letters, and numbers for _repeated_ unknown letters.

Either provide a pattern as a command line argument:

```
> ./codeword-solver 112.2
llama
```

Or drop into a prompt:

```
> ./codeword-solver
Enter pattern: ..a122.1r
chauffeur

Enter pattern: 
```

Supply a custom dictionary file:

```
> ./codeword-solver --dict spanish.txt 112.2
```

### Search for intersecting patterns

The solver also supports searching for multiple patterns together using the following semantics:

```
> ./codeword-solver .os1os 1....o..s1
cosmos:metabolism
```

