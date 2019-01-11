# Nao
A command line spaced repetition flashcard program that uses the SM2 algorithm. It is supposed to be the successor of [nclt](https://github.com/gRastello/nclt).

# Install
To install the latest version of Nao use `go get`.
```bash
$ go get github.com/gRastello/nao
```

Install the manpage.
```bash
$ cd $GOPATH/src/github.com/gRastello/nao
$ sudo ./installman
```

# Changelog

## v1.1.0
- Improved `nao info` readability.
- It is now possible to set the deck directory to a directory whose path contains multiple consecutive whitespaces through the `naorc` file.
- Commands can now be shortened to their first letter (i.e. `nao a deck1` and `nao add deck1` are now the same).
- `Nao` is now licensed under MIT license.
- Minor bugfixes.
