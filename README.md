# verifyiface

`verifyiface` is a static code analyzer which checks that interface implementation is verified

## Install

You can get `verifyiface` by `go get` command.

```bash
$ go get -u github.com/dcu/verifyiface
```

## QuickStart

```bash
$ verifyiface package/...
```

## Analyzer

`verifyiface` checks that an implementation of an interface is verified as explained here: https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
