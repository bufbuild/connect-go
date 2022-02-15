package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/bufconnect/connect"
)

const (
	contextPackage = protogen.GoImportPath("context")
	httpPackage    = protogen.GoImportPath("net/http")
	stringsPackage = protogen.GoImportPath("strings")
	errorsPackage  = protogen.GoImportPath("errors")

	protoPackage = protogen.GoImportPath("google.golang.org/protobuf/proto")

	connectPackage      = protogen.GoImportPath("github.com/bufconnect/connect")
	connectProtoPackage = protogen.GoImportPath("github.com/bufconnect/connect/codec/protobuf")
	cstreamPackage      = protogen.GoImportPath("github.com/bufconnect/connect/clientstream")
	hstreamPackage      = protogen.GoImportPath("github.com/bufconnect/connect/handlerstream")
)

var (
	contextContext          = contextPackage.Ident("Context")
	contextCanceled         = contextPackage.Ident("Canceled")
	contextDeadlineExceeded = contextPackage.Ident("DeadlineExceeded")
	errorsIs                = errorsPackage.Ident("Is")
)

func generate(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_connect.pb.go"
	var path protogen.GoImportPath
	g := gen.NewGeneratedFile(filename, path)
	preamble(gen, file, g)
	content(file, g)
	return g
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	out := fmt.Sprintf("v%d.%d.%d", v.GetMajor(), v.GetMinor(), v.GetPatch())
	if s := v.GetSuffix(); s != "" {
		out += "-" + s
	}
	return out
}

func preamble(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Code generated by protoc-gen-go-connect. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-go-connect v", connect.Version)
	g.P("// - protoc              ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		wrap(g, file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
}

func content(file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}
	handshake(g)
	for _, svc := range file.Services {
		service(file, g, svc)
	}
}

func handshake(g *protogen.GeneratedFile) {
	wrap(g, "This is a compile-time assertion to ensure that this generated file ",
		"and the connect package are compatible. If you get a compiler error that this constant ",
		"isn't defined, this code was generated with a version of connect newer than the one ",
		"compiled into your binary. You can fix the problem by either regenerating this code ",
		"with an older version of connect or updating the connect version compiled into your binary.")
	g.P("const _ = ", connectPackage.Ident("IsAtLeastVersion0_0_1"))
	g.P()
}

type names struct {
	Base string

	Client             string
	ClientConstructor  string
	ClientImpl         string
	ClientExposeMethod string

	Server              string
	ServerConstructor   string
	UnimplementedServer string
}

func newNames(service *protogen.Service) names {
	base := service.GoName
	return names{
		Base: base,

		Client:            fmt.Sprintf("%sClient", base),
		ClientConstructor: fmt.Sprintf("New%sClient", base),
		ClientImpl:        fmt.Sprintf("%sClient", unexport(base)),

		Server:              fmt.Sprintf("%sHandler", base),
		ServerConstructor:   fmt.Sprintf("With%sHandler", base),
		UnimplementedServer: fmt.Sprintf("Unimplemented%sHandler", base),
	}
}

func service(file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	names := newNames(service)

	clientInterface(g, service, names)
	clientImplementation(g, service, names)

	serverInterface(g, service, names)
	serverConstructor(g, service, names)
	unimplementedServerImplementation(g, service, names)
}

func clientInterface(g *protogen.GeneratedFile, service *protogen.Service, names names) {
	wrap(g, names.Client, " is a client for the ", service.Desc.FullName(), " service.")
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		deprecated(g)
	}
	g.Annotate(names.Client, service.Location)
	g.P("type ", names.Client, " interface {")
	for _, method := range service.Methods {
		g.Annotate(names.Client+"."+method.GoName, method.Location)
		leadingComments(
			g,
			method.Comments.Leading,
			method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated(),
		)
		g.P(clientSignature(g, method, false /* named */))
	}
	g.P("}")
	g.P()
}

func clientSignature(g *protogen.GeneratedFile, method *protogen.Method, named bool) string {
	reqName := "req"
	ctxName := "ctx"
	if !named {
		reqName, ctxName = "", ""
	}
	if method.Desc.IsStreamingClient() && method.Desc.IsStreamingServer() {
		// bidi streaming
		return method.GoName + "(" + ctxName + " " + g.QualifiedGoIdent(contextContext) + ") " +
			"*" + g.QualifiedGoIdent(cstreamPackage.Ident("Bidirectional")) +
			"[" + g.QualifiedGoIdent(method.Input.GoIdent) + ", " + g.QualifiedGoIdent(method.Output.GoIdent) + "]"
	}
	if method.Desc.IsStreamingClient() {
		// client streaming
		return method.GoName + "(" + ctxName + " " + g.QualifiedGoIdent(contextContext) + ") " +
			"*" + g.QualifiedGoIdent(cstreamPackage.Ident("Client")) +
			"[" + g.QualifiedGoIdent(method.Input.GoIdent) + ", " + g.QualifiedGoIdent(method.Output.GoIdent) + "]"
	}
	if method.Desc.IsStreamingServer() {
		return method.GoName + "(" + ctxName + " " + g.QualifiedGoIdent(contextContext) +
			", " + reqName + " *" + g.QualifiedGoIdent(connectPackage.Ident("Request")) + "[" +
			g.QualifiedGoIdent(method.Input.GoIdent) + "]) " +
			"(*" + g.QualifiedGoIdent(cstreamPackage.Ident("Server")) +
			"[" + g.QualifiedGoIdent(method.Output.GoIdent) + "]" +
			", error)"
	}
	// unary; symmetric so we can re-use server templating
	return method.GoName + serverSignatureParams(g, method, named)
}

func procedureName(method *protogen.Method) string {
	return fmt.Sprintf(
		"%s.%s/%s",
		method.Parent.Desc.ParentFile().Package(),
		method.Parent.Desc.Name(),
		method.Desc.Name(),
	)
}

func reflectionName(service *protogen.Service) string {
	return fmt.Sprintf("%s.%s", service.Desc.ParentFile().Package(), service.Desc.Name())
}

func clientImplementation(g *protogen.GeneratedFile, service *protogen.Service, names names) {
	clientOption := connectPackage.Ident("ClientOption")

	// Client constructor.
	wrap(g, names.ClientConstructor, " constructs a client for the ", service.Desc.FullName(),
		" service. By default, it uses the binary protobuf codec.")
	g.P("//")
	wrap(g, "The URL supplied here should be the base URL for the gRPC server ",
		"(e.g., https://api.acme.com or https://acme.com/grpc).")
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		deprecated(g)
	}
	g.P("func ", names.ClientConstructor, " (baseURL string, doer ", connectPackage.Ident("Doer"),
		", opts ...", clientOption, ") (", names.Client, ", error) {")
	g.P("baseURL = ", stringsPackage.Ident("TrimRight"), `(baseURL, "/")`)
	g.P("opts = append([]", clientOption, "{")
	g.P(connectPackage.Ident("Codec"), "(", connectProtoPackage.Ident("NameBinary"), ", ",
		connectProtoPackage.Ident("NewBinary"), "()),")
	g.P("}, opts...)")
	g.P("var (")
	g.P("client ", names.ClientImpl)
	g.P("err error")
	g.P(")")
	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			g.P(
				"client.",
				unexport(method.GoName),
				", err = ",
				connectPackage.Ident("NewClientStream"),
				"(",
			)
			g.P("doer,")
			if method.Desc.IsStreamingClient() && method.Desc.IsStreamingServer() {
				g.P(connectPackage.Ident("StreamTypeBidirectional"), ",")
			} else if method.Desc.IsStreamingClient() {
				g.P(connectPackage.Ident("StreamTypeClient"), ",")
			} else {
				g.P(connectPackage.Ident("StreamTypeServer"), ",")
			}
			g.P("baseURL,")
			g.P(`"`, procedureName(method), `",`)
			g.P("opts...,")
			g.P(")")
		} else {
			g.P("client.", unexport(method.GoName), ", err = ", connectPackage.Ident("NewClientFunc"), "[", method.Input.GoIdent, ", ", method.Output.GoIdent, "](")
			g.P("doer,")
			g.P("baseURL,")
			g.P(`"`, procedureName(method), `",`)
			g.P("opts...,")
			g.P(")")
		}
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
	}
	g.P("return &client, nil")
	g.P("}")
	g.P()

	// Client struct.
	wrap(g, names.ClientImpl, " implements ", names.Client, ".")
	g.P("type ", names.ClientImpl, " struct {")
	typeSender := connectPackage.Ident("Sender")
	typeReceiver := connectPackage.Ident("Receiver")
	for _, method := range service.Methods {
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			g.P(unexport(method.GoName), " func(", contextContext, ") (", typeSender, ", ", typeReceiver, ")")
		} else {
			g.P(unexport(method.GoName), " func", serverSignatureParams(g, method, false /* named */))
		}
	}
	g.P("}")
	g.P()
	g.P("var _ ", names.Client, " = (*", names.ClientImpl, ")(nil) // verify interface implementation")
	g.P()
	for _, method := range service.Methods {
		clientMethod(g, service, method, names)
	}
}

func clientMethod(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method, names names) {
	receiver := names.ClientImpl
	isStreamingClient := method.Desc.IsStreamingClient()
	isStreamingServer := method.Desc.IsStreamingServer()
	wrap(g, method.GoName, " calls ", method.Desc.FullName(), ".")
	if method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
		g.P("//")
		deprecated(g)
	}
	g.P("func (c *", receiver, ") ", clientSignature(g, method, true /* named */), " {")

	if isStreamingClient || isStreamingServer {
		g.P("sender, receiver := c.", unexport(method.GoName), "(ctx)")
		if !isStreamingClient && isStreamingServer {
			// server streaming, we need to send the request.
			g.P("if err := sender.Send(req.Msg); err != nil {")
			g.P("_ = sender.Close(err)")
			g.P("_ = receiver.Close()")
			g.P("return nil, err")
			g.P("}")
			g.P("if err := sender.Close(nil); err != nil {")
			g.P("_ = receiver.Close()")
			g.P("return nil, err")
			g.P("}")
			g.P("return ", cstreamPackage.Ident("NewServer"),
				"[", method.Output.GoIdent, "]", "(receiver), nil")
		} else if isStreamingClient && !isStreamingServer {
			// client streaming
			g.P("return ", cstreamPackage.Ident("NewClient"),
				"[", method.Input.GoIdent, ", ", method.Output.GoIdent, "]", "(sender, receiver)")
		} else {
			// bidi streaming
			g.P("return ", cstreamPackage.Ident("NewBidirectional"),
				"[", method.Input.GoIdent, ", ", method.Output.GoIdent, "]", "(sender, receiver)")
		}
	} else {
		g.P("return c.", unexport(method.GoName), "(ctx, req)")
	}
	g.P("}")
	g.P()
}

func serverInterface(g *protogen.GeneratedFile, service *protogen.Service, names names) {
	wrap(g, names.Server, " is an implementation of the ", service.Desc.FullName(), " service.")
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		deprecated(g)
	}
	g.Annotate(names.Server, service.Location)
	g.P("type ", names.Server, " interface {")
	for _, method := range service.Methods {
		leadingComments(
			g,
			method.Comments.Leading,
			method.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated(),
		)
		g.Annotate(names.Server+"."+method.GoName, method.Location)
		g.P(serverSignature(g, method))
	}
	g.P("}")
	g.P()
}

func serverSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	return method.GoName + serverSignatureParams(g, method, false /* named */)
}

func serverSignatureParams(g *protogen.GeneratedFile, method *protogen.Method, named bool) string {
	ctxName := "ctx "
	reqName := "req "
	streamName := "stream "
	if !named {
		ctxName, reqName, streamName = "", "", ""
	}
	if method.Desc.IsStreamingClient() && method.Desc.IsStreamingServer() {
		// bidi streaming
		return "(" + ctxName + g.QualifiedGoIdent(contextContext) + ", " +
			streamName + "*" + g.QualifiedGoIdent(hstreamPackage.Ident("Bidirectional")) +
			"[" + g.QualifiedGoIdent(method.Input.GoIdent) + ", " + g.QualifiedGoIdent(method.Output.GoIdent) + "]" +
			") error"
	}
	if method.Desc.IsStreamingClient() {
		// client streaming
		return "(" + ctxName + g.QualifiedGoIdent(contextContext) + ", " +
			streamName + "*" + g.QualifiedGoIdent(hstreamPackage.Ident("Client")) +
			"[" + g.QualifiedGoIdent(method.Input.GoIdent) + ", " + g.QualifiedGoIdent(method.Output.GoIdent) + "]" +
			") error"
	}
	if method.Desc.IsStreamingServer() {
		// server streaming
		return "(" + ctxName + g.QualifiedGoIdent(contextContext) +
			", " + reqName + "*" + g.QualifiedGoIdent(connectPackage.Ident("Request")) + "[" +
			g.QualifiedGoIdent(method.Input.GoIdent) + "], " +
			streamName + "*" + g.QualifiedGoIdent(hstreamPackage.Ident("Server")) +
			"[" + g.QualifiedGoIdent(method.Output.GoIdent) + "]" +
			") error"
	}
	// unary
	return "(" + ctxName + g.QualifiedGoIdent(contextContext) +
		", " + reqName + "*" + g.QualifiedGoIdent(connectPackage.Ident("Request")) + "[" +
		g.QualifiedGoIdent(method.Input.GoIdent) + "]) " +
		"(*" + g.QualifiedGoIdent(connectPackage.Ident("Response")) + "[" +
		g.QualifiedGoIdent(method.Output.GoIdent) + "], error)"
}

func serverConstructor(g *protogen.GeneratedFile, service *protogen.Service, names names) {
	wrap(g, names.ServerConstructor, " wraps the service implementation in a connect.MuxOption,",
		" which can then be passed to connect.NewServeMux.")
	g.P("//")
	wrap(g, "By default, services support the gRPC and gRPC-Web protocols with ",
		"the binary protobuf and JSON codecs.")
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		deprecated(g)
	}
	handlerOption := connectPackage.Ident("HandlerOption")
	g.P("func ", names.ServerConstructor, "(svc ", names.Server, ", opts ...", handlerOption,
		") ", connectPackage.Ident("MuxOption"), " {")
	g.P("handlers := make([]", connectPackage.Ident("Handler"), ", 0, ", len(service.Methods), ")")
	g.P("opts = append([]", handlerOption, "{")
	g.P(connectPackage.Ident("Codec"), "(", connectProtoPackage.Ident("NameBinary"), ", ", connectProtoPackage.Ident("NewBinary"), "()", "),")
	g.P(connectPackage.Ident("Codec"), "(", connectProtoPackage.Ident("NameJSON"), ", ", connectProtoPackage.Ident("NewJSON"), "()", "),")
	g.P("}, opts...)")
	g.P()
	for _, method := range service.Methods {
		hname := unexport(string(method.Desc.Name()))

		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			g.P(hname, ", err := ", connectPackage.Ident("NewStreamingHandler"), "(")
			if method.Desc.IsStreamingServer() && method.Desc.IsStreamingClient() {
				g.P(connectPackage.Ident("StreamTypeBidirectional"), ",")
			} else if method.Desc.IsStreamingServer() {
				g.P(connectPackage.Ident("StreamTypeServer"), ",")
			} else {
				g.P(connectPackage.Ident("StreamTypeClient"), ",")
			}
			g.P(`"`, procedureName(method), `", // procedure name`)
			g.P(`"`, reflectionName(service), `", // reflection name`)
			g.P("func(ctx ", contextContext, ", sender ", connectPackage.Ident("Sender"),
				", receiver ", connectPackage.Ident("Receiver"), ") {")
			if method.Desc.IsStreamingServer() && method.Desc.IsStreamingClient() {
				// bidi streaming
				g.P("typed := ", hstreamPackage.Ident("NewBidirectional"),
					"[", method.Input.GoIdent, ", ", method.Output.GoIdent, "]", "(sender, receiver)")
			} else if method.Desc.IsStreamingClient() {
				// client streaming
				g.P("typed := ", hstreamPackage.Ident("NewClient"),
					"[", method.Input.GoIdent, ", ", method.Output.GoIdent, "]", "(sender, receiver)")
			} else {
				// server streaming
				g.P("typed := ", hstreamPackage.Ident("NewServer"),
					"[", method.Output.GoIdent, "]", "(sender)")
			}
			if method.Desc.IsStreamingServer() && !method.Desc.IsStreamingClient() {
				g.P("req, err := ", connectPackage.Ident("ReceiveRequest"), "[", method.Input.GoIdent, "]",
					"(receiver)")
				g.P("if err != nil {")
				g.P("_ = receiver.Close()")
				g.P("_ = sender.Close(err)")
				g.P("return")
				g.P("}")
				g.P("if err = receiver.Close(); err != nil {")
				g.P("_ = sender.Close(err)")
				g.P("return")
				g.P("}")
				g.P("err = svc.", method.GoName, "(ctx, req, typed)")
			} else {
				g.P("err := svc.", method.GoName, "(ctx, typed)")
				g.P("_ = receiver.Close()")
			}
			g.P("_ = sender.Close(err)")
			g.P("},")
			g.P("opts...,")
			g.P(")")
		} else {
			g.P(hname, ", err := ", connectPackage.Ident("NewUnaryHandler"), "(")
			g.P(`"`, procedureName(method), `", // procedure name`)
			g.P(`"`, reflectionName(service), `", // reflection name`)
			g.P("svc.", method.GoName, ",")
			g.P("opts...,")
			g.P(")")
		}
		g.P("if err != nil {")
		g.P("return ", connectPackage.Ident("WithHandlers"), "(nil, err)")
		g.P("}")
		g.P("handlers = append(handlers, *", hname, ")")
		g.P()
	}
	g.P("return ", connectPackage.Ident("WithHandlers"), "(handlers, nil)")
	g.P("}")
	g.P()
}

func unimplementedServerImplementation(g *protogen.GeneratedFile, service *protogen.Service, names names) {
	wrap(g, names.UnimplementedServer, " returns CodeUnimplemented from all methods.")
	g.P("type ", names.UnimplementedServer, " struct {}")
	g.P()
	g.P("var _ ", names.Server, " = (*", names.UnimplementedServer, ")(nil) // verify interface implementation")
	g.P()
	for _, method := range service.Methods {
		g.P("func (", names.UnimplementedServer, ") ", serverSignature(g, method), "{")
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			g.P("return ", connectPackage.Ident("Errorf"), "(", connectPackage.Ident("CodeUnimplemented"), `, "`, method.Desc.FullName(), ` isn't implemented")`)
		} else {
			g.P("return nil, ", connectPackage.Ident("Errorf"), "(", connectPackage.Ident("CodeUnimplemented"), `, "`, method.Desc.FullName(), ` isn't implemented")`)
		}
		g.P("}")
		g.P()
	}
	g.P()
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func deprecated(g *protogen.GeneratedFile) {
	g.P("// Deprecated: do not use.")
}

func leadingComments(g *protogen.GeneratedFile, comments protogen.Comments, isDeprecated bool) {
	if comments.String() != "" {
		g.P(strings.TrimSpace(comments.String()))
	}
	if isDeprecated {
		if comments.String() != "" {
			g.P("//")
		}
		deprecated(g)
	}
}
