package generators

import (
	"os"
	"strings"
	engine "text/template"

	"github.com/bimalabs/generators/templates"
)

type Model struct {
}

func (g *Model) Generate(template Template, modulePath string, driver string) {
	var temp string
	if driver == "mongo" {
		temp = templates.MongoModel
	} else {
		temp = templates.GormModel
	}

	modelTemplate, err := engine.New("model").Parse(temp)
	if err != nil {
		panic(err)
	}

	var path strings.Builder
	path.WriteString(modulePath)
	path.WriteString("/model.go")

	modelFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	err = modelTemplate.Execute(modelFile, template)
	if err != nil {
		panic(err)
	}
}
