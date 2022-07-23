package generators

import (
	"os"
	"strings"
	engine "text/template"

	"github.com/bimalabs/generators/templates"
)

type Proto struct {
}

func (g *Proto) Generate(template *Template, modulePath string, driver string) {
	var temp string
	if driver == "mongo" {
		temp = templates.MongoProto
	} else {
		temp = templates.GormProto
	}

	protoTemplate, err := engine.New("proto").Parse(temp)
	if err != nil {
		panic(err)
	}

	workDir, _ := os.Getwd()

	var path strings.Builder

	path.WriteString(workDir)
	path.WriteString("/protos/")
	path.WriteString(template.ModuleLowercase)
	path.WriteString(".proto")

	protoFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	protoTemplate.Execute(protoFile, template)
}
