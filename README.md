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
