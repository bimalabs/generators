package generators

import (
	"os"
	"strings"
	engine "text/template"

	"github.com/bimalabs/generators/templates"
)

type Server struct {
}

func (g *Server) Generate(template Template, modulePath string, driver string) {
	var temp string
	if driver == "mongo" {
		temp = templates.MongoServer
	} else {
		temp = templates.GormServer
	}

	serverTemplate, err := engine.New("server").Parse(temp)
	if err != nil {
		panic(err)
	}

	var path strings.Builder

	path.WriteString(modulePath)
	path.WriteString("/server.go")

	serverFile, err := os.Create(path.String())
	if err != nil {
		panic(err)
	}

	serverTemplate.Execute(serverFile, template)
}
