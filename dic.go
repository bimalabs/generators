package generators

import (
	"os"
	"strings"
	engine "text/template"

	"github.com/bimalabs/generators/templates"
)

type Dic struct {
}

func (g *Dic) Generate(template Template, modulePath string, driver string) {
	var temp string
	if driver == "mongo" {
		temp = templates.MongoDic
	} else {
		temp = templates.GormDic
	}

	dicTemplate, err := engine.New("dic").Parse(temp)
	if err != nil {
		panic(err)
	}

	var path strings.Builder
	path.WriteString(modulePath)
	path.WriteString("/dic.go")

	dicFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	dicTemplate.Execute(dicFile, template)
}
