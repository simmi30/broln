# How to write a Ruby gRPC client for the Lightning Network Daemon

This section enumerates what you need to do to write a client that communicates
with `broln` in Ruby.

### Introduction

`broln` uses the `gRPC` protocol for communication with clients like `brolncli`.

`gRPC` is based on protocol buffers and as such, you will need to compile
the `broln` proto file in Ruby before you can use it to communicate with `broln`.

### Setup

Install gRPC rubygems:

```shell
⛰  gem install grpc
⛰  gem install grpc-tools
```

Clone the Google APIs repository:

```shell
⛰  git clone https://github.com/googleapis/googleapis.git
```

Fetch the `lightning.proto` file (or copy it from your local source directory):

```shell
⛰  curl -o lightning.proto -s https://raw.githubusercontent.com/lightningnetwork/broln/master/lnrpc/lightning.proto
```

Compile the proto file:

```shell
⛰  grpc_tools_ruby_protoc --proto_path googleapis:. --ruby_out=. --grpc_out=. lightning.proto
```

Two files will be generated in the current directory: 

* `lightning_pb.rb`
* `lightning_services_pb.rb`

### Examples

#### Simple client to display wallet balance

Every time you use the Ruby gRPC you need to require the `lightning_services_pb` file.

We assume that `broln` runs on the default `localhost:10019`.

We further assume you run `broln` with `--no-macaroons`.

Note that when an IP address is used to connect to the node (e.g. 192.168.1.21 instead of localhost) you need to add `--tlsextraip=192.168.1.21` to your `broln` configuration and re-generate the certificate (delete tls.cert and tls.key and restart broln).

```ruby
#!/usr/bin/env ruby

$:.unshift(File.dirname(__FILE__))

require 'grpc'
require 'lightning_services_pb'

# Due to updated ECDSA generated tls.cert we need to let gprc know that
# we need to use that cipher suite otherwise there will be a handhsake
# error when we communicate with the broln rpc server.
ENV['GRPC_SSL_CIPHER_SUITES'] = "HIGH+ECDSA"

certificate = File.read(File.expand_path("~/.broln/tls.cert"))
credentials = GRPC::Core::ChannelCredentials.new(certificate)
stub = Lnrpc::Lightning::Stub.new('127.0.0.1:10019', credentials)

response = stub.wallet_balance(Lnrpc::WalletBalanceRequest.new())
puts "Total balance: #{response.total_balance}"
```

This will show the `total_balance` of the wallet.

#### Streaming client for invoice payment updates

```ruby
#!/usr/bin/env ruby

$:.unshift(File.dirname(__FILE__))

require 'grpc'
require 'lightning_services_pb'

ENV['GRPC_SSL_CIPHER_SUITES'] = "HIGH+ECDSA"

certificate = File.read(File.expand_path("~/.broln/tls.cert"))
credentials = GRPC::Core::ChannelCredentials.new(certificate)
stub = Lnrpc::Lightning::Stub.new('127.0.0.1:10019', credentials)

stub.subscribe_invoices(Lnrpc::InvoiceSubscription.new) do |invoice|
  puts invoice.inspect
end
```

Now, create an invoice on your node:

```shell
⛰  brolncli addinvoice --amt=590
{
	"r_hash": <R_HASH>,
	"pay_req": <PAY_REQ>
}
```

Next send a payment to it from another node:

```shell
⛰  brolncli sendpayment --pay_req=<PAY_REQ>
```

You should now see the details of the settled invoice appear.

#### Using Macaroons

To authenticate using macaroons you need to include the macaroon in the metadata of the request.

```ruby
# broln admin macaroon is at ~/.broln/data/chain/brocoin/simnet/admin.macaroon on Linux and
# ~/Library/Application Support/broln/data/chain/brocoin/simnet/admin.macaroon on Mac
macaroon_binary = File.read(File.expand_path("~/.broln/data/chain/brocoin/simnet/admin.macaroon"))
macaroon = macaroon_binary.each_byte.map { |b| b.to_s(16).rjust(2,'0') }.join
```

The simplest approach to use the macaroon is to include the metadata in each request as shown below.

```ruby
stub.get_info(Lnrpc::GetInfoRequest.new, metadata: {macaroon: macaroon})
```

However, this can get tiresome to do for each request. We can use gRPC interceptors to add this metadata to each request automatically. Our interceptor class would look like this.

```ruby
class MacaroonInterceptor < GRPC::ClientInterceptor
  attr_reader :macaroon

  def initialize(macaroon)
    @macaroon = macaroon
    super
  end

  def request_response(request:, call:, method:, metadata:)
    metadata['macaroon'] = macaroon
    yield
  end

  def server_streamer(request:, call:, method:, metadata:)
    metadata['macaroon'] = macaroon
    yield
  end
end
```

And then we would include it when we create our stub like so.

```ruby
certificate = File.read(File.expand_path("~/.broln/tls.cert"))
credentials = GRPC::Core::ChannelCredentials.new(certificate)
macaroon_binary = File.read(File.expand_path("~/.broln/data/chain/brocoin/simnet/admin.macaroon"))
macaroon = macaroon_binary.each_byte.map { |b| b.to_s(16).rjust(2,'0') }.join

stub = Lnrpc::Lightning::Stub.new(
	'localhost:10019',
	credentials,
	interceptors: [MacaroonInterceptor.new(macaroon)]
)

# Now we don't need to pass the metadata on a request level
p stub.get_info(Lnrpc::GetInfoRequest.new)
```

#### Receive Large Responses

A GRPC::ResourceExhausted exception is raised when a server response is too large. In particular, this will happen with mainnet DescribeGraph calls. The solution is to raise the default limits by including a channel_args hash when creating our stub.

```ruby
stub = Lnrpc::Lightning::Stub.new(
  'localhost:10019',
  credentials,
  channel_args: {"grpc.max_receive_message_length" => 1024 * 1024 * 50}
)
```
