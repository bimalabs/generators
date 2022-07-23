package templates

import _ "embed"

var (
	// go:embed mongo/dic.tpl
	MongoDic string

	// go:embed mongo/model.tpl
	MongoModel string

	// go:embed mongo/module.tpl
	MongoModule string

	// go:embed mongo/proto.tpl
	MongoProto string

	// go:embed mongo/server.tpl
	MongoServer string
)

func Template() []string {
	return []string{MongoDic, MongoModel, MongoModule, MongoProto, MongoServer}
}
