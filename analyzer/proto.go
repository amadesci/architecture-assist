package analyzer

import (
	"os"
	"path/filepath"
)

func detectProtoContracts(projectRoot string) string {
	if _, err := os.Stat(filepath.Join(projectRoot, "*.proto")); err == nil {
		return DockerComposeEnv
	}
}

func parseProtoFile(path string) (GrpcService, error) {
	parser := proto.NewParser(os.Open(path))
	definition, _ := parser.Parse()
	visitor := &GrpcVisitor{}
	proto.Walk(definition, visitor)
	return visitor.Service, nil
}

type GrpcVisitor struct {
	proto.NoopVisitor
	Service GrpcService
}

func (v *GrpcVisitor) VisitService(s *proto.Service) {
	v.Service.ServiceName = s.Name
	for _, e := range s.Elements {
		if rpc, ok := e.(*proto.RPC); ok {
			v.Service.Methods = append(v.Service.Methods, GrpcMethod{
				Name:       rpc.Name,
				InputType:  rpc.RequestType,
				OutputType: rpc.ReturnsType,
			})
		}
	}
}
