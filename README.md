# Codeword Solver

CLI version of [codewordsolver.com](https://codewordsolver.com)

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

The solver also supports searching for intersecting patterns using the following syntax:

```
> ./codeword-solver .osmos m....o..sm 4 10
cosmos:metabolism
```

We are finding words which intersect on the 4th character of .osmos (m) and the 10th character of m....o..sm

If no positions are specified, it defaults to using the first letter of each word.
