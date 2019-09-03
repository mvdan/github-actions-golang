# GitHub Actions for Go

[GitHub Actions](https://github.com/features/actions) includes CI/CD for free
for Open Source repositories. This document contains information on making it
work well for [Go](https://github.com/features/actions).

### Quickstart

```
on: [push, pull_request]
name: Test
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.11.9, 1.12.9]
        platform: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.platform }}
    steps:
    - name: Install Go
      uses: actions/setup-go@v1
      with:
        go-version: ${{ matrix.go-version }}
    - name: Checkout code
      uses: actions/checkout@v1
    - name: Test
      run: go test ./...
```

### Summary

Each YAML file under `.github/workflows/` is a separate workflow to run on a
series of events specified by `on`. `name` is its human readable name.

Each workflow file has a number of jobs.

Each job runs on a configuration matrix. For example, we can test two major Go
versions on the three major operating systems.

Each job has a number of steps. Installing Go and checking out the repository's
code are the essential ones.

### FAQs

##### What about module support?

If your repository contains a `go.mod` file, Go 1.12 and later will already use
module mode by default. To turn it on explicitly, set `GO111MODULE=on`. To use
`GOPATH` mode, you'd need `GOPATH/src/your/pkg/name` and `GO111MODULE=off`.

##### How do I set environent variables?

They can only be set for each step, as far as the documentation covers:

```
- name: Download Go dependencies with the module proxy
  env:
    GOPROXY: "https://proxy.golang.org"
  run: go mod download
```

On Go 1.13 and later, this can be simplified:

```
- name: Set Go env vars
  run: go env -w GOPROXY="https://proxy.golang.org"
```

##### How do I set up caching between builds?

We haven't been able to find a simple way to accomplish this. It would be useful
to persist Go's module download and build caches.

##### How do I run a step conditionally?

You can use `if` conditionals, using their [custom expression
language](https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions):

```
- name: Run end-to-end tests
  if: github.event_name == 'push'
  run: go run ./endtoend
```

##### How do I run multiline scripts?

```
- name: Series of commands
  run: |
    go test ./...
    go test -race ./...
```

##### Should I use two workflows, or two jobs on one workflow?

As far as we can tell, the only differences are in the UI and in how each
workflow can be triggered on a different set of events. Otherwise, there doesn't
seem to be a difference.

##### Should I use two workflows, or two jobs on one workflow?

### Quick links

* Concepts, rate limits, joining the beta, etc: https://help.github.com/en/articles/about-github-actions

* Syntax and fields reference: https://help.github.com/en/articles/workflow-syntax-for-github-actions

* Environment reference: https://help.github.com/en/articles/virtual-environments-for-github-actions

* To report bugs: https://github.community/t5/GitHub-Actions/bd-p/actions

### Known bugs

* https://github.community/t5/GitHub-Actions/git-config-core-autocrlf-should-default-to-false/m-p/30445

`git config core.autocrlf` defaults to true, so be careful about CRLF endings in
your plaintext `testdata` files on Windows. To work around this, set up the
following `.gitattributes`:

```
* -text
```

* https://github.community/t5/GitHub-Actions/TEMP-is-broken-on-Windows/m-p/30432/thread-id/427

`os.TempDir` on Windows will contain a short name, since `%TEMP%` also contains
it. Note that case sensitivity doesn't matter, and that `os.Open` should still
work; but some programs not treaing short names might break.

```
$ echo %USERPROFILE%  # i.e. $HOME
C:\Users\runneradmin
$ echo %TEMP%         # a shortened version of $HOME/AppData/Local/Temp
C:\Users\RUNNER~1\AppData\Local\Temp
```
