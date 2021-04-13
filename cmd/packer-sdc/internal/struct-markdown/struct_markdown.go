package struct_markdown

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/fatih/structtag"
)

type Command struct {
}

func (cmd *Command) Help() string {
	return `
	Usage: stuct-markdown file
	`
}

func (cmd *Command) Run(args []string) int {
	if len(args) == 0 {
		// Default: process the file
		args = []string{os.Getenv("GOFILE")}
	}
	fname := args[0]

	absFilePath, err := filepath.Abs(fname)
	if err != nil {
		log.Printf(err.Error())
	}

	var projectRoot, docsFolder, filePath string

	for dir := filepath.Dir(absFilePath); len(dir) > 0 && projectRoot == ""; dir = filepath.Dir(dir) {
		base := filepath.Base(dir)
		log.Printf("base %s", base)
		if base == "packer" && filepath.Base(filepath.Dir(dir)) == "hashicorp" {
			projectRoot = dir
			filePath, _ = filepath.Rel(projectRoot, absFilePath)
			docsFolder = filepath.Join("website", "content", "partials")
		}
		if strings.HasPrefix(base, "packer-plugin-") {
			projectRoot = dir
			filePath, _ = filepath.Rel(projectRoot, absFilePath)
			docsFolder = filepath.Join("docs-partials")
		}
	}

	log.Printf("%s %v %s", projectRoot, filePath, docsFolder)

	if projectRoot == "" {
		log.Fatal("Failed to guess project ROOT. Is this a `packer-plugin-*` named project ?")
	}

	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("ReadFile: %+v", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, b, parser.ParseComments)
	if err != nil {
		log.Fatalf("ParseFile: %+v", err)
	}

	for _, decl := range f.Decls {
		typeDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		typeSpec, ok := typeDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		structDecl, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		fields := structDecl.Fields.List
		header := Struct{
			SourcePath: filePath,
			Name:       typeSpec.Name.Name,
			Filename:   typeSpec.Name.Name + ".mdx",
			Header:     strings.TrimSpace(typeDecl.Doc.Text()),
		}
		dataSourceOutput := Struct{
			SourcePath: filePath,
			Name:       typeSpec.Name.Name,
			Filename:   typeSpec.Name.Name + ".mdx",
		}
		required := Struct{
			SourcePath: filePath,
			Name:       typeSpec.Name.Name,
			Filename:   typeSpec.Name.Name + "-required.mdx",
		}
		notRequired := Struct{
			SourcePath: filePath,
			Name:       typeSpec.Name.Name,
			Filename:   typeSpec.Name.Name + "-not-required.mdx",
		}

		for _, field := range fields {
			if len(field.Names) == 0 || field.Tag == nil {
				continue
			}
			tag := field.Tag.Value[1:]
			tag = tag[:len(tag)-1]
			tags, err := structtag.Parse(tag)
			if err != nil {
				log.Fatalf("structtag.Parse(%s): err: %v", field.Tag.Value, err)
			}

			// Leave undocumented tags out of markdown. This is useful for
			// fields which exist for backwards compatability, or internal-use
			// only fields
			undocumented, _ := tags.Get("undocumented")
			if undocumented != nil {
				if undocumented.Name == "true" {
					continue
				}
			}
			mstr, err := tags.Get("mapstructure")
			if err != nil {
				continue
			}
			name := mstr.Name

			if name == "" {
				continue
			}

			var docs string
			if field.Doc != nil {
				docs = field.Doc.Text()
			} else {
				docs = strings.Join(camelcase.Split(field.Names[0].Name), " ")
			}

			if strings.Contains(docs, "TODO") {
				continue
			}
			fieldType := string(b[field.Type.Pos()-1 : field.Type.End()-1])
			fieldType = strings.ReplaceAll(fieldType, "*", `\*`)
			switch fieldType {
			case "time.Duration":
				fieldType = `duration string | ex: "1h5m2s"`
			case "config.Trilean":
				fieldType = `boolean`
			case "config.NameValues":
				fieldType = `[]{name string, value string}`
			case "config.KeyValues":
				fieldType = `[]{key string, value string}`
			}

			field := Field{
				Name: name,
				Type: fieldType,
				Docs: docs,
			}

			if typeSpec.Name.Name == "DatasourceOutput" {
				dataSourceOutput.Fields = append(dataSourceOutput.Fields, field)
				continue
			}

			if req, err := tags.Get("required"); err == nil && req.Value() == "true" {
				required.Fields = append(required.Fields, field)
			} else {
				notRequired.Fields = append(notRequired.Fields, field)
			}
		}

		dir := filepath.Join(projectRoot, docsFolder, filepath.Dir(filePath))
		log.Printf("dir: %q", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("mkdir of %q failed: %v", dir, err)
		}

		for _, str := range []Struct{header, dataSourceOutput, required, notRequired} {
			if len(str.Fields) == 0 && len(str.Header) == 0 {
				continue
			}
			outputPath := filepath.Join(dir, str.Filename)

			outputFile, err := os.Create(outputPath)
			if err != nil {
				log.Fatalf(err.Error())
				return 1
			}
			defer outputFile.Close()

			err = structDocsTemplate.Execute(outputFile, str)
			if err != nil {
				log.Printf(err.Error())
				return 1
			}
		}
	}

	return 0
}

func (cmd *Command) Synopsis() string {
	return "Generates a markdown file from the comments contained in a struct config."
}
