# GitHub Actions for Go

[GitHub Actions](https://github.com/features/actions) includes CI/CD for free
for Open Source repositories. This document contains information on making it
work well for [Go](https://go.dev/). See them [in
action](https://github.com/mvdan/github-actions-golang/actions):

```yaml
$ cat .github/workflows/test.yml
on: [push, pull_request]
name: Test
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.23.x, 1.24.x]
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}
    - run: go test ./...
```

## Summary

Each workflow file has a number of jobs, which get run `on` specified events,
and run concurrently with each other. You can have workflow [status badges](https://docs.github.com/en/actions/monitoring-and-troubleshooting-workflows/monitoring-workflows/adding-a-workflow-status-badge).

Each `job` runs on a configuration `matrix`. For example, we can test two major
Go versions on three operating systems.

Each job has a number of `steps`, such as installing Go, or checking out the
repository's code.

Note that `name` fields are optional.

## FAQs

#### How do I set environment variables?

They can be set up via `env` for an [entire
workflow](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions#env),
a job, or for each step:

```yaml
env:
  GOPROXY: "https://proxy.company.com"
jobs:
  [...]
```

#### How do I set environment variables at run-time?

You can use [environment files](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/workflow-commands-for-github-actions#environment-files)
to set environment variables or add an element to `$PATH`. For example:

```yaml
steps:
- name: Set env vars
  run: |
      echo "CGO_ENABLED=0" >> $GITHUB_ENV
      echo "${HOME}/goroot/bin" >> $GITHUB_PATH
```

Note that these take effect for future steps in the job.

#### How do I set up caching between builds?

Since v4, [actions/setup-go](https://github.com/actions/setup-go) caches `GOCACHE`
and `GOMODCACHE` automatically, using `go.sum` as the cache key.
You can turn that off via `cache: false`, and then you may also use your own
custom caching, for example to only keep `GOMODCACHE`:

```yaml
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

See [this guide](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/caching-dependencies-to-speed-up-workflows)
for more details.

#### How do I run a step conditionally?

You can use `if` conditionals, using their [custom expression
language](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs):

```yaml
- if: github.event_name == 'push' && matrix.os == 'ubuntu-latest'
  run: go run ./endtoend
```

#### How do I set up a custom build matrix?

You can [include extra matrix
jobs](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions#example-including-new-combinations),
and you can [exclude specific matrix
jobs](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions#example-excluding-configurations-from-a-matrix).

#### How do I run multiline scripts?

```yaml
- name: Series of commands
  run: |
    go test ./...
    go test -race ./...
```

#### Should I use two workflows, or two jobs on one workflow?

The biggest difference is the UI; workflow results are shown separately.
Grouping jobs in workflows can also be useful if one wants to customize the
workflow triggers, or to set up dependencies via
[needs](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions#jobsjob_idneeds).

#### How do I set up a secret environment variable?

Follow [these steps](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions)
to set up the secret in the repo's settings. After adding a secret like
`FOO_SECRET`, use it on a step as follows:

```yaml
- run: some-command
  env:
    FOO_SECRET: ${{ secrets.FOO_SECRET }}
```

#### How do I install private modules?

It's possible to install modules from private GitHub repositories without using
your own proxy. You'll need to add a
[personal access token](https://github.com/settings/tokens) as a secret
environment variable, as well as configure
[GOPRIVATE](https://go.dev/ref/mod#private-modules).
You can also directly used the token
[provided by GitHub](https://docs.github.com/en/enterprise-cloud@latest/actions/security-for-github-actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow)
in the workflow.
You can define anything as username in the URL, it is not taken into account by GitHub.

```yaml
- name: Configure git for private modules
  env:
    TOKEN: ${{ secrets.PERSONAL_ACCESS_TOKEN }}
  run: git config --global url."https://user:${TOKEN}@github.com".insteadOf "https://github.com"
```

```yaml
env:
  GOPRIVATE: "*.company.com"
jobs:
  [...]
```

#### How do I install Linux packages?

Use `sudo apt`, making sure to only run the step on Linux:

```yaml
- if: matrix.os == 'ubuntu-latest'
  run: sudo apt update && sudo apt install -y --no-install-recommends mypackage
```

## Quick links

* Concepts, rate limits, etc: https://docs.github.com/en/actions/writing-workflows

* Syntax and fields reference: https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions

* GitHub-hosted runners: https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners

## Caveats

* https://github.com/actions/checkout/issues/135

`git config core.autocrlf` defaults to true, so be careful about CRLF endings in
your plaintext `testdata` files on Windows. To work around this, set up the
following `.gitattributes`:

```gitattributes
* -text
```
