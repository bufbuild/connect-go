// protoc-gen-connect-go is a plugin for the protocol buffer compiler that
// generates Go code. To use it, build this program and make it available on
// your PATH as protoc-gen-connect-go.
//
// The 'connect-go' suffix becomes part of the arguments for the protocol buffer
// compiler, so you'll invoke it like this:
//	 protoc --connect-go_out=. path/to/file.proto
//
// This generates service definitions for the protocol buffer services defined
// by file.proto. As invoked above, the output will be written to:
//	 path/to/file_connect.pb.go
// If you'd prefer to write the output elsewhere, set '--connect-go_opt' as
// described in
// https://developers.google.com/protocol-buffers/docs/reference/go-generated.
package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/bufbuild/connect"
)

func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	version := flagSet.Bool("version", false, "Print the version and exit.")
	flagSet.Parse(os.Args[1:])
	if *version {
		fmt.Printf("%s %s\n", os.Args[0], connect.Version)
		return
	}
	protogen.Options{}.Run(
		func(plugin *protogen.Plugin) error {
			plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
			for _, f := range plugin.Files {
				if f.Generate {
					generate(plugin, f)
				}
			}
			return nil
		},
	)
}
