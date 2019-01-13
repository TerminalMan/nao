# Nao
A command line spaced repetition flashcard program that uses the SM2 algorithm. It is supposed to be the successor of [nclt](https://github.com/gRastello/nclt). (Slow) development is done on the `develop` branch.

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

If you want to install a particular version or, for some reason, you wish to use the one in development just checkout the wanted tag/branch and run `go install` and reinstall the manpage for that version with the above command.

# Changelog

## v1.2.0
- New configurable option `maxinterval` that set the maximum interval in days between repetitions of the same flashcard. 
- New configurable option `noprompt` that stops `Nao` from prompting for input while doing `nao add`. Useful for piping into `nao add`.

## v1.1.0
- Improved `nao info` readability.
- It is now possible to set the deck directory to a directory whose path contains multiple consecutive whitespaces through the `naorc` file.
- Commands can now be shortened to their first letter (i.e. `nao a deck1` and `nao add deck1` are now the same).
- `Nao` is now licensed under MIT license.
- Minor bugfixes.
