package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

type GoExample struct {
	Request    *plugin.CodeGeneratorRequest
	Response   *plugin.CodeGeneratorResponse
	Parameters map[string]string
}

type LocationMessage struct {
	Location        *descriptor.SourceCodeInfo_Location
	Message         *descriptor.DescriptorProto
	LeadingComments []string
}

func (runner *GoExample) PrintParameters(w io.Writer) {
	const padding = 3
	tw := tabwriter.NewWriter(w, 0, 0, padding, ' ', tabwriter.TabIndent)
	fmt.Fprintf(tw, "Parameters:\n")
	for k, v := range runner.Parameters {
		fmt.Fprintf(tw, "%s:\t%s\n", k, v)
	}
	fmt.Fprintln(tw, "")
	tw.Flush()
}

func (runner *GoExample) getLocationMessage() map[string][]*LocationMessage {

	ret := make(map[string][]*LocationMessage)
	for index, filename := range runner.Request.FileToGenerate {
		locationMessages := make([]*LocationMessage, 0)
		proto := runner.Request.ProtoFile[index]
		desc := proto.GetSourceCodeInfo()
		locations := desc.GetLocation()
		for _, location := range locations {
			// I would encourage developers to read the documentation about paths as I might have misunderstood this
			// I am trying to process message types which I understand to be `4` and only at the root level which I understand
			// to be path len == 2
			if len(location.GetPath()) > 2 {
				continue
			}

			leadingComments := strings.Split(location.GetLeadingComments(), "\n")
			if len(location.GetPath()) > 1 && location.GetPath()[0] == int32(4) {
				message := proto.GetMessageType()[location.GetPath()[1]]
				println(message.GetName())
				locationMessages = append(locationMessages, &LocationMessage{
					Message:  message,
					Location: location,
					// Because we are only parsing messages here at the root level we will not get field comments
					LeadingComments: leadingComments[:len(leadingComments)-1],
				})
			}
		}
		ret[filename] = locationMessages
	}
	return ret
}

func (runner *GoExample) CreateMarkdownFile(filename string, messages []*LocationMessage) error {
	// Create a file and append it to the output files

	var outfileName string
	var content string
	outfileName = strings.Replace(filename, ".proto", ".md", -1)
	var mdFile plugin.CodeGeneratorResponse_File
	mdFile.Name = &outfileName
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("# %s\n", outfileName))
	for _, locationMessage := range messages {
		buf.WriteString(fmt.Sprintf("\n## %s\n", locationMessage.Message.GetName()))
		buf.WriteString(fmt.Sprintf("### %s\n", "Leading Comments"))
		for _, comment := range locationMessage.LeadingComments {
			buf.WriteString(fmt.Sprintf("%s\n", comment))
		}
		if len(locationMessage.Message.NestedType) > 0 {
			buf.WriteString(fmt.Sprintf("### %s\n", "Nested Messages"))
			for _, nestedMessage := range locationMessage.Message.NestedType {
				buf.WriteString(fmt.Sprintf("#### %s\n", nestedMessage.GetName()))
				buf.WriteString(fmt.Sprintf("#### %s\n", "Fields"))
				for _, field := range nestedMessage.Field {
					buf.WriteString(fmt.Sprintf("%s - %s\n", field.GetName(), field.GetLabel()))
				}
			}
		}
		for _, field := range locationMessage.Message.Field {
			buf.WriteString(fmt.Sprintf("%s - %s\n", field.GetName(), field.GetLabel()))
		}
	}
	content = buf.String()
	mdFile.Content = &content
	runner.Response.File = append(runner.Response.File, &mdFile)
	os.Stderr.WriteString(fmt.Sprintf("Created File: %s", filename))
	return nil
}

func (runner *GoExample) generateMessageMarkdown() error {
	// This convenience method will return a structure of some types that I use
	fileLocationMessageMap := runner.getLocationMessage()
	for filename, locationMessages := range fileLocationMessageMap {
		runner.CreateMarkdownFile(filename, locationMessages)
	}
	return nil
}

func (runner *GoExample) generateCode() error {
	// Initialize the output file slice
	files := make([]*plugin.CodeGeneratorResponse_File, 0)
	runner.Response.File = files

	{
		err := runner.generateMessageMarkdown()
		if err != nil {
			return err
		}
	}
	return nil
}

type serviceDescriptor struct {
	Name string
}

type microGenerator struct {
	req    *plugin.CodeGeneratorRequest
	rsp    *plugin.CodeGeneratorResponse
	params map[string]string
}

func newMicroGenerator(
	req *plugin.CodeGeneratorRequest,
	rsp *plugin.CodeGeneratorResponse,
	params string,
) *microGenerator {
	m := map[string]string{}
	groupkv := strings.Split(params, ",")
	for _, element := range groupkv {
		kv := strings.Split(element, "=")
		if len(kv) > 1 {
			m[kv[0]] = kv[1]
		}
	}
	return &microGenerator{
		req:    req,
		rsp:    rsp,
		params: m,
	}
}

func fileWriteln(buf *bytes.Buffer, format string, args ...interface{}) {
	buf.WriteString(fmt.Sprintf(format, args...) + "\n")
}

func lcFirst(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

func (g *microGenerator) generateFile(filename string, proto *descriptorpb.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	file := &plugin.CodeGeneratorResponse_File{
		Name: &filename,
	}
	buf := new(bytes.Buffer)
	fileWriteln(buf, "// Code generated by protoc-gen-micro. DO NOT EDIT.")
	fileWriteln(buf, "")
	fileWriteln(buf, "package %s", *proto.Package)
	fileWriteln(buf, "")
	fileWriteln(buf, "import (")
	fileWriteln(buf, "  \"context\"")
	fileWriteln(buf, "  \"github.com/ofavor/micro-lite/server\"")
	fileWriteln(buf, "  \"github.com/ofavor/micro-lite/client\"")
	fileWriteln(buf, ")")
	fileWriteln(buf, "")
	services := proto.GetService()
	pkgName := *proto.Package
	for _, service := range services {
		svcName := *service.Name
		// service interface start
		fileWriteln(buf, "type %sService interface {", *service.Name)
		// service methods
		for _, method := range service.Method {
			it := strings.Replace(*method.InputType, "."+pkgName+".", "", 1)
			ot := strings.Replace(*method.OutputType, "."+pkgName+".", "", 1)
			fileWriteln(buf, "  %s(ctx context.Context, in *%s, opts ...client.CallOption) (*%s, error)", *method.Name, it, ot)
		}
		// service interface end
		fileWriteln(buf, "}")
		fileWriteln(buf, "")

		// service struct
		fileWriteln(buf, "type %sService struct {", lcFirst(svcName)) // should lower case
		fileWriteln(buf, "  serviceName string")
		fileWriteln(buf, "  c client.Client")
		fileWriteln(buf, "}")
		fileWriteln(buf, "")

		// service creator
		fileWriteln(buf, "func New%sService(name string, c client.Client) %sService {", svcName, svcName)
		fileWriteln(buf, "  return &%sService {", lcFirst(svcName)) // should lower case
		fileWriteln(buf, "    serviceName: name,")
		fileWriteln(buf, "    c: c,")
		fileWriteln(buf, "  }")
		fileWriteln(buf, "}")
		fileWriteln(buf, "")

		// service method impl
		for _, method := range service.Method {
			it := strings.Replace(*method.InputType, "."+pkgName+".", "", 1)
			ot := strings.Replace(*method.OutputType, "."+pkgName+".", "", 1)
			fileWriteln(buf, "func (s *%sService)%s(ctx context.Context, in *%s, opts ...client.CallOption) (*%s, error) {", lcFirst(svcName), *method.Name, it, ot) // should lower case
			fileWriteln(buf, "  req := client.NewRequest(s.serviceName, \"%s.%s\", in)", svcName, *method.Name)
			fileWriteln(buf, "  rsp := new(%s)", ot)
			fileWriteln(buf, "  err := s.c.Call(ctx, req, rsp, opts...)")
			fileWriteln(buf, "  if err != nil {")
			fileWriteln(buf, "    return nil, err")
			fileWriteln(buf, "  }")
			fileWriteln(buf, "  return rsp, nil")
			fileWriteln(buf, "}")
			fileWriteln(buf, "")
		}

		// handler interface start
		fileWriteln(buf, "type %sHandler interface {", *service.Name)
		// service methods
		for _, method := range service.Method {
			it := strings.Replace(*method.InputType, "."+pkgName+".", "", 1)
			ot := strings.Replace(*method.OutputType, "."+pkgName+".", "", 1)
			fileWriteln(buf, "  %s(ctx context.Context, in *%s, out *%s) error", *method.Name, it, ot)
		}
		// handler interface end
		fileWriteln(buf, "}")
		fileWriteln(buf, "")

		// handler register
		fileWriteln(buf, "func Register%sHandler(s server.Server, h %sHandler) {", svcName, svcName)
		fileWriteln(buf, "  hdr := server.NewHandler(\"%s\", h)", svcName)
		fileWriteln(buf, "  s.Handle(hdr)")
		fileWriteln(buf, "}")
		fileWriteln(buf, "")

	}
	content := buf.String()
	file.Content = &content
	return file, nil
}

func (g *microGenerator) generate() error {
	g.rsp.File = make([]*plugin.CodeGeneratorResponse_File, 0)
	for index, filename := range g.req.FileToGenerate {
		// fmt.Println(">>>>>>>", index, filename)
		proto := g.req.ProtoFile[index]
		folder := ""
		gpkg := *proto.Options.GoPackage
		if len(gpkg) > 0 {
			arr := strings.Split(gpkg, ";")
			if len(arr) > 0 {
				folder = arr[0]
			}
		}
		// TODO
		services := proto.GetService()
		if len(services) == 0 {
			continue
		}
		outfileName := path.Join(folder, strings.Replace(filename, ".proto", ".micro.pb.go", -1))
		file, err := g.generateFile(outfileName, proto)
		if err != nil {
			panic(err)
		}
		g.rsp.File = append(g.rsp.File, file)
		os.Stderr.WriteString(fmt.Sprintf("Created File: %s\n", outfileName))
	}
	return nil
}

func main() {
	// os.Stdin will contain data which will unmarshal into the following object:
	// https://godoc.org/github.com/golang/protobuf/protoc-gen-go/plugin#CodeGeneratorRequest
	req := &plugin.CodeGeneratorRequest{}
	rsp := &plugin.CodeGeneratorResponse{}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	// You must use the requests unmarshal method to handle this type
	if err := proto.Unmarshal(data, req); err != nil {
		panic(err)
	}

	generator := newMicroGenerator(req, rsp, req.GetParameter())
	if err := generator.generate(); err != nil {
		panic(err)
	}

	marshalled, err := proto.Marshal(rsp)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(marshalled)
}
