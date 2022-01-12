// protoc-gen-go-rerpc is a plugin for the protocol buffer compiler that
// generates Go code. To use it, build this program and make it available on
// your PATH as protoc-gen-go-rerpc.
//
// The 'go-rerpc' suffix becomes part of the arguments for the protocol buffer
// compiler, so you'll invoke it like this:
//	 protoc --go-rerpc_out=. path/to/file.proto
//
// This generates service definitions for the protocol buffer services defined
// by file.proto. As invoked above, the output will be written to:
//	 path/to/file_rerpc.pb.go
// If you'd prefer to write the output elsewhere, set '--go-rerpc_opt' as
// described in
// https://developers.google.com/protocol-buffers/docs/reference/go-generated.
package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/rerpc/rerpc"
)

func main() {
	version := flag.Bool("version", false, "Print the version and exit.")
	flag.Parse()
	if *version {
		fmt.Printf("protoc-gen-go-rerpc %s\n", rerpc.Version)
		return
	}

	var flags flag.FlagSet
	// By default, behave like the gRPC and Twirp plugins and generate code in
	// the same package as protoc-gen-go's output (typically determined by the
	// "go_package" file option).
	//
	// Setting this flag generates our code into a separate package, importing
	// message types from the base package. This escape hatch lets us generate
	// code with shorter names, since we don't need to worry about collisions
	// with gRPC, Twirp, or whatever other plugins the user might run.
	externalGoTypes := flags.Bool(
		"external_go_types",
		false,
		"Generate RPC code in a separate package from the basic Go types.",
	)
	protogen.Options{ParamFunc: flags.Set}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if f.Generate {
				generate(gen, f, *externalGoTypes)
			}
		}
		return nil
	})
}