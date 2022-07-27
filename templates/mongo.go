package templates

import "fmt"

var (
	MongoDic string = `package {{.ModulePluralLowercase}}

import (
    "github.com/bimalabs/framework/v4"
    "github.com/sarulabs/dingo/v4"
)

var Dic = []dingo.Def{
    {
        Name:  "module:{{.ModuleLowercase}}",
        Scope: bima.Application,
        Build: (*Module)(nil),
        Params: dingo.Params{
            "Module": dingo.Service("bima:module"),
        },
    },
    {
        Name:  "module:{{.ModuleLowercase}}:server",
        Scope: bima.Application,
        Build: (*Server)(nil),
        Params: dingo.Params{
            "Server": dingo.Service("bima:server"),
            "Module": dingo.Service("module:{{.ModuleLowercase}}"),
        },
    },
}
    `

	MongoInline   string = "`bson:\",inline\"`"
	MongoRequired string = "`bson:\"{{.NameUnderScore}}\" {{if .IsRequired}}validate:\"required\"{{end}}`"
	MongoModel    string = fmt.Sprintf(`package {{.ModulePluralLowercase}}

import (
    "context"
    "time"

    "github.com/bimalabs/framework/v4/configs"
)

type {{.Module}} struct {
    configs.MongoBase %s
{{range .Columns}}
    {{.Name}} {{.GolangType}} %s
{{end}}
}

func (m *{{.Module}}) CollectionName() string {
    return "{{.ModuleLowercase}}"
}

func (m *{{.Module}}) Creating(context.Context) error {
    m.CreatedAt = time.Now()

    return nil
}

func (m *{{.Module}}) Updating(context.Context) error {
    m.UpdatedAt = time.Now()

    return nil
}
`, MongoInline, MongoRequired)

	MongoModule string = `package {{.ModulePluralLowercase}}

import (
    "context"
    "fmt"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "github.com/bimalabs/framework/v4"
    "github.com/bimalabs/framework/v4/loggers"
    "github.com/bimalabs/framework/v4/paginations"
    "github.com/bimalabs/framework/v4/utils"
    "github.com/goccy/go-json"
    "github.com/jinzhu/copier"
    "{{.PackageName}}/protos/builds"
)

type Module struct {
    *bima.Module
    grpcs.Unimplemented{{.Module}}sServer
}

func (m *Module) GetPaginated(ctx context.Context, r *grpcs.PaginationRequest) (*grpcs.PaginationResponse, error) {
    model := {{.Module}}{}
    reqeust := paginations.Request{}

    m.Paginator.Model = &model
    m.Paginator.Table = model.CollectionName()

    copier.Copy(&reqeust, r)
    m.Paginator.Handle(reqeust)

    records := make([]*grpcs.{{.Module}}, 0, m.Paginator.Limit)
    metadata := m.Handler.Paginate(*m.Paginator, &records)

    return &grpcs.PaginationResponse{
        Code: http.StatusOK,
        Data: records,
        Meta: &grpcs.PaginationMetadata{
            Record:   int32(metadata.Record),
            Page:     int32(metadata.Page),
            Previous: int32(metadata.Previous),
            Next:     int32(metadata.Next),
            Limit:    int32(metadata.Limit),
            Total:    int32(metadata.Total),
        },
    }, nil
}

func (m *Module) Create(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
    v := {{.Module}}{}
    copier.Copy(&v, r)

    if message, err := utils.Validate(&v); err != nil {
        loggers.Logger.Error(ctx, string(message))

        return nil, status.Error(codes.InvalidArgument, string(message))
    }

    if err := m.Handler.Create(&v); err != nil {
        loggers.Logger.Error(ctx, err.Error())

        return nil, status.Error(codes.Internal, "Internal server error")
    }

    r.Id = v.ID.Hex()

    return r, nil
}

func (m *Module) Update(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
    v := {{.Module}}{}
    hold := v
    copier.Copy(&v, r)

    if message, err := utils.Validate(&v); err != nil {
        loggers.Logger.Error(ctx, string(message))

        return nil, status.Error(codes.InvalidArgument, string(message))
    }

    if err := m.Handler.Bind(&hold, r.Id); err != nil {
        loggers.Logger.Error(ctx, err.Error())

        return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
    }

    v.SetID(hold.GetID())
    v.CreatedAt = hold.CreatedAt
    if err := m.Handler.Update(&v, v.ID.Hex()); err != nil {
        loggers.Logger.Error(ctx, err.Error())

        return nil, status.Error(codes.Internal, "Internal server error")
    }

    m.Cache.Invalidate(r.Id)

    return r, nil
}

func (m *Module) Get(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
    var v {{.Module}}
    if data, found := m.Cache.Get(r.Id); found {
        err := json.Unmarshal(data, r)
        if err == nil {
            return r, nil
        }
    } else {
        if err := m.Handler.Bind(&v, r.Id); err != nil {
            loggers.Logger.Error(ctx, err.Error())

            return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
        }
    }

    copier.Copy(r, &v)

    data, err := json.Marshal(r)
    if err == nil {
        m.Cache.Set(r.Id, data)
    }

    return r, nil
}

func (m *Module) Delete(ctx context.Context, r *grpcs.{{.Module}}) (*grpcs.{{.Module}}, error) {
    ctx = context.WithValue(ctx, "scope", "{{.ModuleLowercase}}")
    v := {{.Module}}{}
    if err := m.Handler.Bind(&v, r.Id); err != nil {
        loggers.Logger.Error(ctx, err.Error())

        return nil, status.Error(codes.NotFound, fmt.Sprintf("Data with ID '%s' not found.", r.Id))
    }

    m.Handler.Delete(&v, r.Id)
    m.Cache.Invalidate(r.Id)

    return &grpcs.{{.Module}}{}, nil
}
`

	MongoProto string = `syntax = "proto3";

package grpcs;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "bima/pagination.proto";

option go_package = ".;grpcs";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
        title: "{{.Module}} Service";
        version: "{{.ApiVersion}}";
    };
    consumes: "application/json";
    produces: "application/json";
    security_definitions: {
        security: {
        key: "bearer";
        value: {
                type: TYPE_API_KEY;
                in: IN_HEADER;
                name: "Authorization";
                description: "Authentication token, prefixed by Bearer: Bearer (token)";
            }
        }
    };
    security: {
        security_requirement: {
            key: "bearer";
        }
    };
};

message {{.Module}} {
    string id = 1;
{{range .Columns}}
    {{.ProtobufType}} {{.NameUnderScore}} = {{.Index}};
{{end}}
}

service {{.Module}}s {
    rpc GetPaginated (PaginationRequest) returns (PaginationResponse) {
        option (google.api.http) = {
            get: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}"
        };
    }

    rpc Create ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            post: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}"
            body: "*"
        };
    }

    rpc Update ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            put: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
            body: "*"

            additional_bindings {
                patch: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
                body: "*"
            }
        };
    }

    rpc Get ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            get: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
        };
    }

    rpc Delete ({{.Module}}) returns ({{.Module}}) {
        option (google.api.http) = {
            delete: "/api/{{.ApiVersion}}/{{.ModulePluralLowercase}}/{id}"
        };
    }
}
`

	MongoServer string = `package {{.ModulePluralLowercase}}

import (
    "context"

    "github.com/bimalabs/framework/v4"
    "{{.PackageName}}/protos/builds"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
)

type Server struct {
    *bima.Server
    Module *Module
}

func (s *Server) Register(gs *grpc.Server) {
    grpcs.Register{{.Module}}sServer(gs, s.Module)
}

func (s *Server) Handle(context context.Context, server *runtime.ServeMux, client *grpc.ClientConn) error {
    return grpcs.Register{{.Module}}sHandler(context, server, client)
}
`
)
