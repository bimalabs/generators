package generators

import (
	"os"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
)

type (
	Generator interface {
		Generate(template Template, modulePath string, driver string)
	}

	Template struct {
		ApiPrefix             string
		PackageName           string
		Module                string
		ModuleLowercase       string
		ModulePlural          string
		ModulePluralLowercase string
		Columns               []FieldTemplate
	}

	ModuleJson struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	ModuleTemplate struct {
		Name   string
		Fields []FieldTemplate
	}

	FieldTemplate struct {
		Name           string
		NameUnderScore string
		ProtobufType   string
		GolangType     string
		Index          int
		IsRequired     bool
	}

	Factory struct {
		ApiPrefix  string
		Driver     string
		Pluralizer pluralize.Client
		Template   Template
		Generators []Generator
	}
)

func (f *Factory) Generate(module ModuleTemplate) {
	workDir, _ := os.Getwd()
	packageName := f.packageName(workDir)
	moduleName := strcase.ToCamel(module.Name)
	modulePlural := f.Pluralizer.Plural(module.Name)
	modulePluralLowercase := strcase.ToDelimited(modulePlural, '_')

	var modulePath strings.Builder

	modulePath.WriteString(workDir)
	modulePath.WriteString("/")
	modulePath.WriteString(modulePluralLowercase)

	f.Template.ApiPrefix = f.ApiPrefix
	f.Template.PackageName = packageName
	f.Template.Module = moduleName
	f.Template.ModuleLowercase = strcase.ToDelimited(module.Name, '_')
	f.Template.ModulePlural = modulePlural
	f.Template.ModulePluralLowercase = modulePluralLowercase
	f.Template.Columns = module.Fields

	os.MkdirAll(modulePath.String(), 0755)
	for _, generator := range f.Generators {
		generator.Generate(f.Template, modulePath.String(), f.Driver)
	}
}

func (f *Factory) packageName(workDir string) string {
	var path strings.Builder

	path.WriteString(workDir)
	path.WriteString("/go.mod")

	mod, err := os.ReadFile(path.String())
	if err != nil {
		panic(err)
	}

	return modfile.ModulePath(mod)
}
