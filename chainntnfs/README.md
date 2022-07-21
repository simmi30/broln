chainntnfs
==========

[![Build Status](http://img.shields.io/travis/lightningnetwork/broln.svg)](https://travis-ci.org/lightningnetwork/broln) 
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/brsuite/broln/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/brsuite/broln/chainntnfs)

The chainntnfs package implements a set of interfaces which allow callers to
receive notifications in response to specific on-chain events. The set of
notifications available include: 

  * Notifications for each new block connected to the current best chain.
  * Notifications once a `txid` has reached a specified number of
    confirmations.
  * Notifications once a target outpoint (`txid:index`) has been spent.

These notifications are used within `broln` in order to properly handle the
workflows for: channel funding, cooperative channel closures, forced channel
closures, channel contract breaches, sweeping time-locked outputs, and finally
pruning the channel graph. 

This package is intentionally general enough to be applicable outside the
specific use cases within `broln` outlined above. The current sole concrete
implementation of the `ChainNotifier` interface depends on `brond`.

## Installation and Updating

```shell
⛰  go get -u github.com/brsuite/broln/chainntnfs
```
