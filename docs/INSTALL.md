# Table of Contents
* [Installation](#installation)
  * [Installing a binary release](#installing-a-binary-release)
  * [Building a tagged version with Docker](#building-a-tagged-version-with-docker)
  * [Building a development version from source](#building-a-development-version-from-source)
    * [Preliminaries](#preliminaries-for-installing-from-source)
    * [Installing broln](#installing-broln-from-source)
* [Available Backend Operating Modes](#available-backend-operating-modes)
  * [brond Options](#brond-options)
  * [Neutrino Options](#neutrino-options)
  * [Brocoind Options](#brocoind-options)
  * [Using brond](#using-brond)
    * [Installing brond](#installing-brond)
    * [Starting brond](#starting-brond)
    * [Running broln using the brond backend](#running-broln-using-the-brond-backend)
  * [Using Neutrino](#using-neutrino)
  * [Using brocoind or litecoind](#using-brocoind-or-litecoind)
* [Creating a Wallet](#creating-a-wallet)
* [Macaroons](#macaroons)
* [Network Reachability](#network-reachability)
* [Simnet vs. Testnet Development](#simnet-vs-testnet-development)
* [Creating an broln.conf (Optional)](#creating-an-brolnconf-optional)

# Installation

There are multiple ways to install `broln`. For most users the easiest way is to
[download and install an official release binary](#installing-a-binary-release).
Those release binaries are always built with production in mind and have all
RPC subservers enabled.

More advanced users that want to build `broln` from source also have multiple
options. To build a tagged version, there is a docker build helper script that
allows users to
[build `broln` from source without needing to install `golang`](#building-a-tagged-version-with-docker).
That is also the preferred way to build and verify the reproducible builds that
are released by the team. See
[release.md for more information about reproducible builds](release.md).

Finally, there is the option to build `broln` fully manually. This requires more
tooling to be set up first but allows to produce non-production (debug,
development) builds.

## Installing a binary release

Downloading and installing an official release binary is recommended for use on
mainnet.
[Visit the release page on GitHub](https://github.com/brsuite/broln/releases)
and select the latest version that does not have the "Pre-release" label set
(unless you explicitly want to help test a Release Candidate, RC).

Choose the package that best fits your operating system and system architecture.
It is recommended to choose 64bit versions over 32bit ones, if your operating
system supports both.

Extract the package and place the two binaries (`broln` and `brolncli` or `broln.exe`
and `brolncli.exe` on Windows) somewhere where the operating system can find them.

## Building a tagged version with Docker

To use the Docker build helper, you need to have the following software
installed and set up on your machine:
 - Docker
 - `make`
 - `bash`

To build a specific git tag of `broln`, simply run the following steps (assuming
`v0.x.y-beta` is the tagged version to build):

```shell
⛰  git clone https://github.com/brsuite/broln
⛰  cd broln
⛰  git checkout v0.x.y-beta
⛰  make docker-release tag=v0.x.y-beta
```

This will create a directory called `broln-v0.x.y-beta` that contains the release
binaries for all operating system and architecture pairs. A single pair can also
be selected by specifying the `sys=linux-amd64` flag for example. See
[release.md for more information on reproducible builds](release.md).

## Building a development version from source

Building and installing `broln` from source is only recommended for advanced users
and/or developers. Running the latest commit from the `master` branch is not
recommended for mainnet. The `master` branch can at times be unstable and
running your node off of it can prevent it to go back to a previous, stable
version if there are database migrations present.

### Preliminaries for installing from source
  In order to work with [`broln`](https://github.com/brsuite/broln), the
  following build dependencies are required:

  * **Go:** `broln` is written in Go. To install, run one of the following commands:


    **Note**: The minimum version of Go supported is Go 1.16. We recommend that
    users use the latest version of Go, which at the time of writing is
    [`1.17.1`](https://blog.golang.org/go1.17.1).


    On Linux:

    (x86-64)
    ```
    wget https://dl.google.com/go/go1.17.1.linux-amd64.tar.gz
    sha256sum go1.17.1.linux-amd64.tar.gz | awk -F " " '{ print $1 }'
    ```

    The final output of the command above should be
    `dab7d9c34361dc21ec237d584590d72500652e7c909bf082758fb63064fca0ef`. If it
    isn't, then the target REPO HAS BEEN MODIFIED, and you shouldn't install
    this version of Go. If it matches, then proceed to install Go:
    ```
    sudo tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    ```

    (ARMv6)
    ```
    wget https://dl.google.com/go/go1.17.1.linux-armv6l.tar.gz
    sha256sum go1.17.1.linux-armv6l.tar.gz | awk -F " " '{ print $1 }'
    ```

    The final output of the command above should be
    `ed3e4dbc9b80353f6482c441d65b51808290e94ff1d15d56da5f4a7be7353758`. If it
    isn't, then the target REPO HAS BEEN MODIFIED, and you shouldn't install
    this version of Go. If it matches, then proceed to install Go:
    ```
    tar -C /usr/local -xzf go1.17.1.linux-armv6l.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    ```

    On Mac OS X:
    ```
    brew install go@1.17.1
    ```

    On FreeBSD:
    ```
    pkg install go
    ```

    Alternatively, one can download the pre-compiled binaries hosted on the
    [Golang download page](https://golang.org/dl/). If one seeks to install
    from source, then more detailed installation instructions can be found
    [here](https://golang.org/doc/install).

    At this point, you should set your `$GOPATH` environment variable, which
    represents the path to your workspace. By default, `$GOPATH` is set to
    `~/go`. You will also need to add `$GOPATH/bin` to your `PATH`. This ensures
    that your shell will be able to detect the binaries you install.

    ```shell
    ⛰  export GOPATH=~/gocode
    ⛰  export PATH=$PATH:$GOPATH/bin
    ```

    We recommend placing the above in your .bashrc or in a setup script so that
    you can avoid typing this every time you open a new terminal window.

  * **Go modules:** This project uses [Go modules](https://github.com/golang/go/wiki/Modules) 
    to manage dependencies as well as to provide *reproducible builds*.

    Usage of Go modules (with Go 1.13) means that you no longer need to clone
    `broln` into your `$GOPATH` for development purposes. Instead, your `broln`
    repo can now live anywhere!

### Installing broln from source

With the preliminary steps completed, to install `broln`, `brolncli`, and all
related dependencies run the following commands:
```shell
⛰  git clone https://github.com/brsuite/broln
⛰  cd broln
⛰  make install
```

The command above will install the current _master_ branch of `broln`. If you
wish to install a tagged release of `broln` (as the master branch can at times be
unstable), then [visit then release page to locate the latest
release](https://github.com/brsuite/broln/releases). Assuming the name
of the release is `v0.x.x`, then you can compile this release from source with
a small modification to the above command: 
```shell
⛰  git clone https://github.com/brsuite/broln
⛰  cd broln
⛰  git checkout v0.x.x
⛰  make install
```


**NOTE**: Our instructions still use the `$GOPATH` directory from prior
versions of Go, but with Go 1.13, it's now possible for `broln` to live
_anywhere_ on your file system.

For Windows WSL users, make will need to be referenced directly via
/usr/bin/make/, or alternatively by wrapping quotation marks around make,
like so:

```shell
⛰  /usr/bin/make && /usr/bin/make install

⛰  "make" && "make" install
```

On FreeBSD, use gmake instead of make.

Alternatively, if one doesn't wish to use `make`, then the `go` commands can be
used directly:
```shell
⛰  go install -v ./...
```

**Updating**

To update your version of `broln` to the latest version run the following
commands:
```shell
⛰  cd $GOPATH/src/github.com/brsuite/broln
⛰  git pull
⛰  make clean && make && make install
```

On FreeBSD, use gmake instead of make.

Alternatively, if one doesn't wish to use `make`, then the `go` commands can be
used directly:
```shell
⛰  cd $GOPATH/src/github.com/brsuite/broln
⛰  git pull
⛰  go install -v ./...
```

**Tests**

To check that `broln` was installed properly run the following command:
```shell
⛰   make check
```

This command requires `brocoind` (almost any version should do) to be available
in the system's `$PATH` variable. Otherwise some of the tests will fail.

# Available Backend Operating Modes

In order to run, `broln` requires, that the user specify a chain backend. At the
time of writing of this document, there are three available chain backends:
`brond`, `neutrino`, `brocoind`. All including neutrino can run on mainnet with
an out of the box `broln` instance. We don't require `--txindex` when running
with `brocoind` or `brond` but activating the `txindex` will generally make
`broln` run faster. Note that since version 0.13 pruned nodes are supported
although they cause performance penalty and higher network usage.

The set of arguments for each of the backend modes is as follows:

## brond Options
```text
brond:
      --brond.dir=                                             The base directory that contains the node's data, logs, configuration file, etc. (default: /Users/roasbeef/Library/Application Support/Brond)
      --brond.rpchost=                                         The daemon's rpc listening address. If a port is omitted, then the default port for the selected chain parameters will be used. (default: localhost)
      --brond.rpcuser=                                         Username for RPC connections
      --brond.rpcpass=                                         Password for RPC connections
      --brond.rpccert=                                         File containing the daemon's certificate file (default: /Users/roasbeef/Library/Application Support/Brond/rpc.cert)
      --brond.rawrpccert=                                      The raw bytes of the daemon's PEM-encoded certificate chain which will be used to authenticate the RPC connection.
```

## Neutrino Options
```text
neutrino:
  -a, --neutrino.addpeer=                                     Add a peer to connect with at startup
      --neutrino.connect=                                     Connect only to the specified peers at startup
      --neutrino.maxpeers=                                    Max number of inbound and outbound peers
      --neutrino.banduration=                                 How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second
      --neutrino.banthreshold=                                Maximum allowed ban score before disconnecting and banning misbehaving peers.
      --neutrino.useragentname=                               Used to help identify ourselves to other brocoin peers.
      --neutrino.useragentversion=                            Used to help identify ourselves to other brocoin peers.
```

## Brocoind Options
```text
brocoind:
      --brocoind.dir=                                         The base directory that contains the node's data, logs, configuration file, etc. (default: /Users/roasbeef/Library/Application Support/Brocoin)
      --brocoind.rpchost=                                     The daemon's rpc listening address. If a port is omitted, then the default port for the selected chain parameters will be used. (default: localhost)
      --brocoind.rpcuser=                                     Username for RPC connections
      --brocoind.rpcpass=                                     Password for RPC connections
      --brocoind.zmqpubrawblock=                              The address listening for ZMQ connections to deliver raw block notifications
      --brocoind.zmqpubrawtx=                                 The address listening for ZMQ connections to deliver raw transaction notifications
      --brocoind.estimatemode=                                The fee estimate mode. Must be either "ECONOMICAL" or "CONSERVATIVE". (default: CONSERVATIVE)
```

## Using brond

### Installing brond

On FreeBSD, use gmake instead of make.

To install brond, run the following commands:

Install **brond**:
```shell
⛰   make brond
```

Alternatively, you can install [`brond` directly from its
repo](https://github.com/brsuite/brond).

### Starting brond

Running the following command will create `rpc.cert` and default `brond.conf`.

```shell
⛰   brond --testnet --rpcuser=REPLACEME --rpcpass=REPLACEME
```
If you want to use `broln` on testnet, `brond` needs to first fully sync the
testnet blockchain. Depending on your hardware, this may take up to a few
hours. Note that adding `--txindex` is optional, as it will take longer to sync
the node, but then `broln` will generally operate faster as it can hit the index
directly, rather than scanning blocks or BIP 158 filters for relevant items.

(NOTE: It may take several minutes to find segwit-enabled peers.)

While `brond` is syncing you can check on its progress using brond's `getinfo`
RPC command:
```shell
⛰   bronctl --testnet --rpcuser=REPLACEME --rpcpass=REPLACEME getinfo
{
  "version": 120000,
  "protocolversion": 70002,
  "blocks": 1114996,
  "timeoffset": 0,
  "connections": 7,
  "proxy": "",
  "difficulty": 422570.58270815,
  "testnet": true,
  "relayfee": 0.00001,
  "errors": ""
}
```

Additionally, you can monitor brond's logs to track its syncing progress in real
time.

You can test your `brond` node's connectivity using the `getpeerinfo` command:
```shell
⛰   bronctl --testnet --rpcuser=REPLACEME --rpcpass=REPLACEME getpeerinfo | more
```

### Running broln using the brond backend

If you are on testnet, run this command after `brond` has finished syncing.
Otherwise, replace `--brocoin.testnet` with `--brocoin.simnet`. If you are
installing `broln` in preparation for the
[tutorial](https://dev.lightning.community/tutorial), you may skip this step.
```shell
⛰   broln --brocoin.active --brocoin.testnet --debuglevel=debug \
       --brond.rpcuser=kek --brond.rpcpass=kek --externalip=X.X.X.X
```

## Using Neutrino

In order to run `broln` in its light client mode, you'll need to locate a
full-node which is capable of serving this new light client mode. `broln` uses
[BIP 157](https://github.com/brocoin/bips/blob/master/bip-0157.mediawiki) and [BIP
158](https://github.com/brocoin/bips/blob/master/bip-0158.mediawiki) for its light client
mode.  A public instance of such a node can be found at
`faucet.lightning.community`.

To run broln in neutrino mode, run `broln` with the following arguments, (swapping
in `--brocoin.simnet` if needed), and also your own `brond` node if available:
```shell
⛰   broln --brocoin.active --brocoin.testnet --debuglevel=debug \
       --brocoin.node=neutrino --neutrino.connect=faucet.lightning.community
```


## Using brocoind or litecoind

The configuration for brocoind and litecoind are nearly identical, the
following steps can be mirrored with loss of generality to enable a litecoind
backend.  Setup will be described in regards to `brocoind`, but note that `broln`
uses a distinct `litecoin.node=litecoind` argument and analogous
subconfigurations prefixed by `litecoind`. Note that adding `--txindex` is
optional, as it will take longer to sync the node, but then `broln` will
generally operate faster as it can hit the index directly, rather than scanning
blocks or BIP 158 filters for relevant items.

To configure your brocoind backend for use with broln, first complete and verify
the following:

- Since `broln` uses
  [ZeroMQ](https://github.com/brocoin/brocoin/blob/master/doc/zmq.md) to
  interface with `brocoind`, *your `brocoind` installation must be compiled with
  ZMQ*. Note that if you installed `brocoind` from source and ZMQ was not present, 
  then ZMQ support will be disabled, and `broln` will quit on a `connection refused` error. 
  If you installed `brocoind` via Homebrew in the past ZMQ may not be included 
  ([this has now been fixed](https://github.com/Homebrew/homebrew-core/pull/23088) 
  in the latest Homebrew recipe for brocoin)
- Configure the `brocoind` instance for ZMQ with `--zmqpubrawblock` and
  `--zmqpubrawtx`. These options must each use their own unique address in order
  to provide a reliable delivery of notifications (e.g.
  `--zmqpubrawblock=tcp://127.0.0.1:28332` and
  `--zmqpubrawtx=tcp://127.0.0.1:28333`).
- Start `brocoind` running against testnet, and let it complete a full sync with
  the testnet chain (alternatively, use `--brocoind.regtest` instead).

Here's a sample `brocoin.conf` for use with broln:
```text
testnet=1
server=1
daemon=1
zmqpubrawblock=tcp://127.0.0.1:28332
zmqpubrawtx=tcp://127.0.0.1:28333
```

Once all of the above is complete, and you've confirmed `brocoind` is fully
updated with the latest blocks on testnet, run the command below to launch
`broln` with `brocoind` as your backend (as with `brocoind`, you can create an
`broln.conf` to save these options, more info on that is described further
below):

```shell
⛰   broln --brocoin.active --brocoin.testnet --debuglevel=debug \
       --brocoin.node=brocoind --brocoind.rpcuser=REPLACEME \
       --brocoind.rpcpass=REPLACEME \
       --brocoind.zmqpubrawblock=tcp://127.0.0.1:28332 \
       --brocoind.zmqpubrawtx=tcp://127.0.0.1:28333 \
       --externalip=X.X.X.X
```

*NOTE:*
- The auth parameters `rpcuser` and `rpcpass` parameters can typically be
  determined by `broln` for a `brocoind` instance running under the same user,
  including when using cookie auth. In this case, you can exclude them from the
  `broln` options entirely.
- If you DO choose to explicitly pass the auth parameters in your `broln.conf` or
  command line options for `broln` (`brocoind.rpcuser` and `brocoind.rpcpass` as
  shown in example command above), you must also specify the
  `brocoind.zmqpubrawblock` and `brocoind.zmqpubrawtx` options. Otherwise, `broln`
  will attempt to get the configuration from your `brocoin.conf`.
- You must ensure the same addresses are used for the `brocoind.zmqpubrawblock`
  and `brocoind.zmqpubrawtx` options passed to `broln` as for the `zmqpubrawblock`
  and `zmqpubrawtx` passed in the `brocoind` options respectively.
- When running broln and brocoind on the same Windows machine, ensure you use
  127.0.0.1, not localhost, for all configuration options that require a TCP/IP
  host address.  If you use "localhost" as the host name, you may see extremely
  slow inter-process-communication between broln and the brocoind backend.  If broln
  is experiencing this issue, you'll see "Waiting for chain backend to finish
  sync, start_height=XXXXXX" as the last entry in the console or log output, and
  broln will appear to hang.  Normal broln output will quickly show multiple
  messages like this as broln consumes blocks from brocoind.
- Don't connect more than two or three instances of `broln` to `brocoind`. With
  the default `brocoind` settings, having more than one instance of `broln`, or
  `broln` plus any application that consumes the RPC could cause `broln` to miss
  crucial updates from the backend.
- The default fee estimate mode in `brocoind` is CONSERVATIVE. You can set
  `brocoind.estimatemode=ECONOMICAL` to change it into ECONOMICAL. Futhermore,
  if you start `brocoind` in `regtest`, this configuration won't take any effect.


# Creating a wallet
If `broln` is being run for the first time, create a new wallet with:
```shell
⛰   brolncli create
```
This will prompt for a wallet password, and optionally a cipher seed
passphrase.

`broln` will then print a 24 word cipher seed mnemonic, which can be used to
recover the wallet in case of data loss. The user should write this down and
keep in a safe place.

More [information about managing wallets can be found in the wallet management
document](wallet.md).

# Macaroons

`broln`'s authentication system is called **macaroons**, which are decentralized
bearer credentials allowing for delegation, attenuation, and other cool
features. You can learn more about them in Alex Akselrod's [writeup on
Github](https://github.com/brsuite/broln/issues/20).

Running `broln` for the first time will by default generate the `admin.macaroon`,
`read_only.macaroon`, and `macaroons.db` files that are used to authenticate
into `broln`. They will be stored in the network directory (default:
`brolndir/data/chain/brocoin/mainnet`) so that it's possible to use a distinct
password for mainnet, testnet, simnet, etc. Note that if you specified an
alternative data directory (via the `--datadir` argument), you will have to
additionally pass the updated location of the `admin.macaroon` file into `brolncli`
using the `--macaroonpath` argument.

To disable macaroons for testing, pass the `--no-macaroons` flag into *both*
`broln` and `brolncli`.

# Network Reachability

If you'd like to signal to other nodes on the network that you'll accept
incoming channels (as peers need to connect inbound to initiate a channel
funding workflow), then the `--externalip` flag should be set to your publicly
reachable IP address.

# Simnet vs. Testnet Development

If you are doing local development, such as for the tutorial, you'll want to
start both `brond` and `broln` in the `simnet` mode. Simnet is similar to regtest
in that you'll be able to instantly mine blocks as needed to test `broln`
locally. In order to start either daemon in the `simnet` mode use `simnet`
instead of `testnet`, adding the `--brocoin.simnet` flag instead of the
`--brocoin.testnet` flag.

Another relevant command line flag for local testing of new `broln` developments
is the `--debughtlc` flag. When starting `broln` with this flag, it'll be able to
automatically settle a special type of HTLC sent to it. This means that you
won't need to manually insert invoices in order to test payment connectivity.
To send this "special" HTLC type, include the `--debugsend` command at the end
of your `sendpayment` commands.


There are currently two primary ways to run `broln`: one requires a local `brond`
instance with the RPC service exposed, and the other uses a fully integrated
light client powered by [neutrino](https://github.com/brsuite/neutrino).

# Creating an broln.conf (Optional)

Optionally, if you'd like to have a persistent configuration between `broln`
launches, allowing you to simply type `broln --brocoin.testnet --brocoin.active`
at the command line, you can create an `broln.conf`.

**On MacOS, located at:**
`/Users/[username]/Library/Application Support/broln/broln.conf`

**On Linux, located at:**
`~/.broln/broln.conf`

Here's a sample `broln.conf` for `brond` to get you started:
```text
[Application Options]
debuglevel=trace
maxpendingchannels=10

[Brocoin]
brocoin.active=1
```

Notice the `[Brocoin]` section. This section houses the parameters for the
Brocoin chain. `broln` also supports Litecoin testnet4 (but not both BTC and LTC
at the same time), so when working with Litecoin be sure to set to parameters
for Litecoin accordingly. See a more detailed sample config file available
[here](https://github.com/brsuite/broln/blob/master/sample-broln.conf)
and explore the other sections for node configuration, including `[Brond]`,
`[Brocoind]`, `[Neutrino]`, `[Ltcd]`, and `[Litecoind]` depending on which
chain and node type you're using.
