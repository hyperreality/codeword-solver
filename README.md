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
> ./codewordsolver 112.2
llama
```

Or drop into a prompt:

```
> ./codewordsolver
Enter pattern: ..a122.1r
chauffeur

Enter pattern: 
```

Supply a custom dictionary file:

```
> ./codewordsolver --dict spanish.txt 112.2
```

