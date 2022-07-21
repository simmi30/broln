Makefile
========

To build, verify, and install `broln` from source, use the following
commands:
```shell
⛰  make
⛰  make check
⛰  make install
```

The command `make check` requires `brocoind` (almost any version should do) to
be available in the system's `$PATH` variable. Otherwise some of the tests will
fail.

Developers
==========

This document specifies all commands available from `broln`'s `Makefile`.
The commands included handle:
- Installation of all go-related dependencies.
- Compilation and installation of `broln` and `brolncli`.
- Compilation and installation of `brond` and `bronctl`.
- Running unit and integration suites.
- Testing, debugging, and flake hunting.
- Formatting and linting.

Commands
========

- [`all`](#scratch)
- [`brond`](#brond)
- [`build`](#build)
- [`check`](#check)
- [`clean`](#clean)
- [`default`](#default)
- [`dep`](#dep)
- [`flake-unit`](#flake-unit)
- [`flakehunter`](#flakehunter)
- [`fmt`](#fmt)
- [`install`](#install)
- [`itest`](#itest)
- [`lint`](#lint)
- [`list`](#list)
- [`rpc`](#rpc)
- [`scratch`](#scratch)
- [`travis`](#travis)
- [`unit`](#unit)
- [`unit-cover`](#unit-cover)
- [`unit-race`](#unit-race)

`all`
-----
Compiles, tests, and installs `broln` and `brolncli`. Equivalent to 
[`scratch`](#scratch) [`check`](#check) [`install`](#install).

`brond`
------
Ensures that the [`github.com/brsuite/brond`][brond] repository is checked out
locally. Lastly, installs the version of 
[`github.com/brsuite/brond`][brond] specified in `Gopkg.toml`

`build`
-------
Compiles the current source and vendor trees, creating `./broln` and
`./brolncli`.

`check`
-------
Installs the version of [`github.com/brsuite/brond`][brond] specified
in `Gopkg.toml`, then runs the unit tests followed by the integration
tests.

Related: [`unit`](#unit) [`itest`](#itest)

`clean`
-------
Removes compiled versions of both `./broln` and `./brolncli`, and removes the
`vendor` tree.

`default`
---------
Alias for [`scratch`](#scratch).

`flake-unit`
------------
Runs the unit test endlessly until a failure is detected.

Arguments:
- `pkg=<package>` 
- `case=<testcase>`
- `timeout=<timeout>`

Related: [`unit`](#unit)

`flakehunter`
-------------
Runs the itegration test suite endlessly until a failure is detected.

Arguments:
- `icase=<itestcase>`
- `timeout=<timeout>`

Related: [`itest`](#itest)

`fmt`
-----
Runs `go fmt` on the entire project. 

`install`
---------
Copies the compiled `broln` and `brolncli` binaries into `$GOPATH/bin`.

`itest`
-------
Installs the version of [`github.com/brsuite/brond`][brond] specified in
`Gopkg.toml`, builds the `./broln` and `./brolncli` binaries, then runs the
integration test suite.

Arguments:
- `icase=<itestcase>` (the snake_case version of the testcase name field in the testCases slice (i.e. sweep_coins), not the test func name)
- `timeout=<timeout>`

`itest-parallel`
------
Does the same as `itest` but splits the total set of tests into
`NUM_ITEST_TRANCHES` tranches (currently set to 6 by default, can be overwritten
by setting `tranches=Y`) and runs them in parallel.

Arguments:
- `icase=<itestcase>`: The snake_case version of the testcase name field in the
  testCases slice (i.e. `sweep_coins`, not the test func name) or any regular
  expression describing a set of tests.
- `timeout=<timeout>`
- `tranches=<number_of_tranches>`: The number of parts/tranches to split the
  total set of tests into.
- `parallel=<number_of_threads>`: The number of threads to run in parallel. Must
  be greater or equal to `tranches`, otherwise undefined behavior is expected.

`flakehunter-parallel`
------
Runs the test specified by `icase` simultaneously `parallel` (default=6) times
until an error occurs. Useful for hunting flakes.

Example:
```shell
⛰  make flakehunter-parallel icase='(data_loss_protection|channel_backup)' backend=neutrino
```

`lint`
------
Ensures that [`gopkg.in/alecthomas/gometalinter.v1`][gometalinter] is
installed, then lints the project.

`list`
------
Lists all known make targets.

`rpc`
-----
Compiles the `lnrpc` proto files.

`scratch`
---------
Compiles all dependencies and builds the `./broln` and `./brolncli` binaries.
Equivalent to [`lint`](#lint) [`brond`](#brond)
[`unit-race`](#unit-race).

`unit`
------
Runs the unit test suite. By default, this will run all known unit tests.

Arguments:
- `pkg=<package>` 
- `case=<testcase>`
- `timeout=<timeout>`
- `log="stdlog[ <log-level>]"` prints logs to stdout
  - `<log-level>` can be `info` (default), `debug`, `trace`, `warn`, `error`, `critical`, or `off`

`unit-cover`
------------
Runs the unit test suite with test coverage, compiling the statisitics in
`profile.cov`.

Arguments:
- `pkg=<package>` 
- `case=<testcase>`
- `timeout=<timeout>`
- `log="stdlog[ <log-level>]"` prints logs to stdout
  - `<log-level>` can be `info` (default), `debug`, `trace`, `warn`, `error`, `critical`, or `off`

Related: [`unit`](#unit)

`unit-race`
-----------
Runs the unit test suite with go's race detector.

Arguments:
- `pkg=<package>` 
- `case=<testcase>`
- `timeout=<timeout>`
- `log="stdlog[ <log-level>]"` prints logs to stdout
  - `<log-level>` can be `info` (default), `debug`, `trace`, `warn`, `error`, `critical`, or `off`

Related: [`unit`](#unit)

[brond]: https://github.com/brsuite/brond (github.com/brsuite/brond")
[gometalinter]: https://gopkg.in/alecthomas/gometalinter.v1 (gopkg.in/alecthomas/gometalinter.v1)
