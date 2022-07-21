/*
Package lntest provides testing utilities for the broln repository.

This package contains infrastructure for integration tests that launch full broln
nodes in a controlled environment and interact with them via RPC. Using a
NetworkHarness, a test can launch multiple broln nodes, open channels between
them, create defined network topologies, and anything else that is possible with
RPC commands.
*/
package lntest
