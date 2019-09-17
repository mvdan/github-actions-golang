# GitHub Actions for Go

[GitHub Actions](https://github.com/features/actions) includes CI/CD for free
for Open Source repositories. This document contains information on making it
work well for [Go](https://github.com/features/actions). See them [in
action](https://github.com/mvdan/github-actions-golang/actions):

```
$ cat .github/workflows/test.yml
on: [push, pull_request]
name: Test
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.12.9, 1.13]
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

## Summary

Each workflow file has a number of jobs, which get run `on` specified events.

Each job runs on a configuration matrix. For example, we can test two major Go
versions on three operating systems.

Each job has a number of steps, such as installing Go, or checking out the
repository's code.

## FAQs

#### What about module support?

If your repository contains a `go.mod` file, Go 1.12 and later will already use
module mode by default. To turn it on explicitly, set `GO111MODULE=on`. To use
`GOPATH` mode, you'd need `$GOPATH/src/your/pkg/name` and `GO111MODULE=off`.

#### How do I set environent variables?

They can only be set for each step, as far as the documentation covers:

```
- name: Download Go dependencies with a custom proxy
  env:
    GOPROXY: "https://proxy.company.com"
  run: go mod download
```

On Go 1.13 and later, this can be simplified:

```
- name: Set Go env vars
  run: go env -w GOPROXY="https://proxy.company.com"
```

#### How do I set up caching between builds?

We haven't been able to find a simple way to accomplish this. It would be useful
to persist Go's module download and build caches.

#### How do I run a step conditionally?

You can use `if` conditionals, using their [custom expression
language](https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions):

```
- name: Run end-to-end tests on Linux
  if: github.event_name == 'push' && matrix.platform == 'ubuntu-latest'
  run: go run ./endtoend
```

#### How do I run multiline scripts?

```
- name: Series of commands
  run: |
    go test ./...
    go test -race ./...
```

#### Should I use two workflows, or two jobs on one workflow?

As far as we can tell, the only differences are in the UI and in how each
workflow can be triggered on a different set of events. Otherwise, there doesn't
seem to be a difference.

#### How do I set up a secret environment variable?

Follow [these steps](https://help.github.com/en/articles/virtual-environments-for-github-actions#creating-and-using-secrets-encrypted-variables)
to set up the secret in the repo's settings. After adding a secret like
`FOO_SECRET`, use it on a step as follows:

```
- name: Command that requires secret
  run: some-command
  env:
    FOO_SECRET: ${{ secrets.FOO_SECRET }}
```

#### How do I install private modules?

It's possible to install modules from private GitHub repositories without using
your own proxy. You'll need to add a
[personal access token](https://github.com/settings/tokens) as a secret
environment variable for this to work.

```yaml
- name: Configure git for private modules
  env:
    TOKEN: ${{ secrets.PERSONAL_ACCESS_TOKEN }}
  run: git config --global url."https://YOUR_GITHUB_USERNAME:${TOKEN}@github.com".insteadOf "https://github.com"
```

## Quick links

* Concepts, rate limits, joining the beta, etc: https://help.github.com/en/articles/about-github-actions

* Syntax and fields reference: https://help.github.com/en/articles/workflow-syntax-for-github-actions

* Environment reference: https://help.github.com/en/articles/virtual-environments-for-github-actions

* To report bugs: https://github.community/t5/GitHub-Actions/bd-p/actions

## Known bugs

* https://github.com/actions/setup-go/issues/14

The `setup-go` action doesn't set `PATH`, so currently it's not possible to `go
install` a program and run it directly. Until that's fixed, consider absolute
paths like `$(go env GOPATH)/bin/program`.

* https://github.community/t5/GitHub-Actions/git-config-core-autocrlf-should-default-to-false/m-p/30445

`git config core.autocrlf` defaults to true, so be careful about CRLF endings in
your plaintext `testdata` files on Windows. To work around this, set up the
following `.gitattributes`:

```
* -text
```

* https://github.community/t5/GitHub-Actions/LocalAppData-unset-on-Windows-when-using-Actions-for-CI/m-p/31349

`GOCACHE` won't be accessible on Windows by default, since `%LocalAppData%`
isn't defined. To temporarily work around the error below, set `GOCACHE`
manually:

> build cache is required, but could not be located: GOCACHE is not defined and
> %LocalAppData% is not defined

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
