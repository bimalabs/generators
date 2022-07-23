package templates

import _ "embed"

var (
	// go:embed gorm/dic.tpl
	GormDic string

	// go:embed gorm/model.tpl
	GormModel string

	// go:embed gorm/module.tpl
	GormModule string

	// go:embed gorm/proto.tpl
	GormProto string

	// go:embed gorm/server.tpl
	GormServer string
)
