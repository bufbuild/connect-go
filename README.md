reRPC
=====

reRPC is a small RPC framework built on [protocol buffers][protobuf] and
`net/http`. It generates code from API definition files so you can focus on
your application logic: no more artisanal query parameter parsing, no more
debates between PUT and PATCH, and no more hand-written clients.

reRPC servers and clients use [gRPC's][grpc] HTTP/2 protocol. reRPC servers
work seamlessly with clients generated by [any gRPC
implementation][grpc-implementations], command-line tools like [grpcurl][], and
proxies like [gRPC-Gateway][grpc-gateway] and [Envoy][envoy]. Similarly, reRPC
clients can call any gRPC server.

reRPC servers also support [Twirp's][twirp] HTTP/1.1 protocol. Of course,
clients generated by [any Twirp implementation][twirp-implementations] work
with reRPC servers. More importantly, Twirp's JSON variant is perfect for
debugging with cURL.

Sadly, nothing's free. To keep the implementation simple and expose the same
features over multiple protocols, reRPC only supports unary (request-response)
RPCs. There's an [open issue][streaming-issue] for discussion of streaming
support.

For more on reRPC, including a walkthrough and comparison to alternatives, see
the [docs][].

## A Small Example

Curious what all this looks like in practice? Here's a small h2c server:

```go
package main

import (
  "net/http"

  "golang.org/x/net/http2"
  "golang.org/x/net/http2/h2c"

  pingpb "github.com/akshayjshah/rerpc/internal/ping/v1test" // generated
)

type PingServer struct {
  pingpb.UnimplementedPingServiceReRPC // returns errors from all methods
}

func main() {
  ping := &PingServer{}
  mux := http.NewServeMux()
  mux.Handle(pingpb.NewPingHandlerReRPC(ping))
  handler := h2c.NewHandler(mux, &http2.Server{})
  http.ListenAndServe(":8081", handler)
}
```

With that server running, you can make requests with a gRPC client or with
cURL:

```bash
$ curl --request POST \
  --header "Content-Type: application/json" \
  http://localhost:8081/internal.ping.v1test.PingService/Ping

{"code":"unimplemented","msg":"internal.ping.v1test.PingService.Ping isn't implemented"}
```

You can find production-ready examples of [servers][prod-server] and
[clients][prod-client] in the API documentation.

## Status

This is the earliest of early alphas: APIs *will* break before the first stable
release.

## Support and Versioning

reRPC supports:

* The [two most recent major releases][go-support-policy] of Go.
* Version 3 of the protocol buffer language ([proto3][]).
* [APIv2][] of protocol buffers in Go (`google.golang.org/protobuf`).

Within those parameters, reRPC follows semantic versioning.

That said, please remember that reRPC is one person's labor of love.
(Well...love and frustration. Mostly love.) It'll probably take me a few days
to respond to issues and pull requests. If you're using reRPC in production,
I'd love your [help maintaining this project][maintainers-issue].

## Legal

Offered under the [MIT license][license]. This is a personal project developed
in my spare time - it's not endorsed by, supported by, or (as far as I know)
used by my current or former employers.

[APIv2]: https://blog.golang.org/protobuf-apiv2
[docs]: https://github.com/akshayjshah/rerpc/wiki
[envoy]: https://www.envoyproxy.io/
[godoc]: https://pkg.go.dev/github.com/akshayjshah/rerpc
[go-support-policy]: https://golang.org/doc/devel/release#policy
[grpc-gateway]: https://grpc-ecosystem.github.io/grpc-gateway/
[grpc]: https://grpc.io/
[grpc-implementations]: https://grpc.io/docs/languages/
[grpcurl]: https://github.com/fullstorydev/grpcurl
[license]: https://github.com/akshayjshah/rerpc/blob/main/LICENSE.txt
[maintainers-issue]: https://github.com/akshayjshah/rerpc/issues/2
[prod-client]: https://pkg.go.dev/github.com/akshayjshah/rerpc#example-Client
[prod-server]: https://pkg.go.dev/github.com/akshayjshah/rerpc#example-package
[proto3]: https://cloud.google.com/apis/design/proto3
[protobuf]: https://developers.google.com/protocol-buffers
[streaming-issue]: https://github.com/akshayjshah/rerpc/issues/1
[twirp]: https://twitchtv.github.io/twirp/
[twirp-implementations]: https://github.com/twitchtv/twirp#implementations-in-other-languages
