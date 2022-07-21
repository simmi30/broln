package lnrpc

import "regexp"

var (
	// brolnClientStreamingURIs is a list of all broln RPCs that use a request-
	// streaming interface. Those request-streaming RPCs need to be handled
	// differently in the WebsocketProxy because of how the request body
	// parsing is implemented in the grpc-gateway library. Unfortunately
	// there is no straightforward way of obtaining this information on
	// runtime so we need to keep a hard coded list here.
	brolnClientStreamingURIs = []*regexp.Regexp{
		regexp.MustCompile("^/v1/channels/acceptor$"),
		regexp.MustCompile("^/v1/channels/transaction-stream$"),
		regexp.MustCompile("^/v2/router/htlcinterceptor$"),
	}
)
