// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// mapstructure-to-hcl2 fills the gaps between hcl2 and mapstructure for Packer
//
// By generating a struct that the HCL2 ecosystem understands making use of
// mapstructure tags.
//
// Packer heavily uses the mapstructure decoding library to load/parse user
// config files. Packer now needs to move to HCL2.
//
// Here are a few differences/gaps betweens hcl2 and mapstructure:
//
//   - in HCL2 all basic struct fields (string/int/struct) that are not pointers
//     are required ( must be set ). In mapstructure everything is optional.
//
//   - mapstructure allows to 'squash' fields
//     (ex: Field CommonStructType `mapstructure:",squash"`) this allows to
//     decorate structs and reuse configuration code. HCL2 parsing libs don't have
//     anything similar.
//
// mapstructure-to-hcl2 will parse Packer's config files and generate the HCL2
// compliant code that will allow to not change any of the current builders in
// order to softly move to HCL2.
package mapstructure_to_hcl2

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/types"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var (
	cmdPrefix = "mapstructure-to-hcl2"

	//go:embed README.md
	readme string
)

const HCL_LABEL_INDEX_KEY = "hcllabelindex"

type Command struct {
	typeNames string
	output    string
}

// Usage is a replacement usage function for the flags package.
func (cmd *Command) Help() string {
	cmd.Flags().Usage()
	return "\n" + readme
}

func (cmd *Command) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet(cmdPrefix, flag.ExitOnError)
	fs.StringVar(&cmd.typeNames, "type", "", "comma-separated list of type names; must be set")
	fs.StringVar(&cmd.output, "output", "", "output file name; default srcdir/<type>_hcl2.go")
	return fs
}

func (cmd *Command) Synopsis() string {
	return "[//go:generate command] Generates the code necessary for a packer-plugin to 'speak' HCL2."
}

func (cmd *Command) Run(args []string) int {
	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		log.Fatalf("unable to parse flags: %s", err)
	}
	args = fs.Args()

	log.SetFlags(0)
	log.SetPrefix(cmdPrefix + ": ")
	if len(cmd.typeNames) == 0 {
		cmd.Help()
		return 2
	}
	typeNames := strings.Split(cmd.typeNames, ",")

	// We accept either one directory or a list of files. Which do we have?
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}
	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		goFile = args[0]
	}
	outputPath := goFile[:len(goFile)-2] + "hcl2spec.go"

	log.SetPrefix(fmt.Sprintf(cmdPrefix+": %s.%v: ", os.Getenv("GOPACKAGE"), typeNames))

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, args...)
	if err != nil {
		log.Fatalf("package.Load: %v", err)
	}
	if len(pkgs) != 1 {
		log.Printf("error: %d packages found", len(pkgs))
		return 1
	}
	topPkg := pkgs[0]
	sort.Strings(typeNames)

	var structs []StructDef
	usedImports := map[NamePath]*types.Package{}

	for id, obj := range topPkg.TypesInfo.Defs {
		if obj == nil {
			continue
		}
		t := obj.Type()
		nt, isANamedType := t.(*types.Named)
		if !isANamedType {
			continue
		}
		if nt.Obj().Pkg() != topPkg.Types {
			// Sometimes a struct embeds another struct named the same. ex:
			// builder/osc/bsuvolume.BlockDevice. This makes sure the type is
			// defined in topPkg.
			continue
		}
		ut := nt.Underlying()
		utStruct, utOk := ut.(*types.Struct)
		if !utOk {
			continue
		}

		pos := sort.SearchStrings(typeNames, id.Name)
		if pos >= len(typeNames) || typeNames[pos] != id.Name {
			continue // not a struct we care about
		}
		// Sometimes we see the underlying struct for a similar named type, which results
		// in an incorrect FlatMap. If the type names are not exactly the same skip.
		if nt.Obj().Name() != id.Name {
			continue // not the struct we are looking for
		}
		// make sure each type is found once where somehow sometimes they can be found twice
		typeNames = append(typeNames[:pos], typeNames[pos+1:]...)
		flatenedStruct, err := getMapstructureSquashedStruct(obj.Pkg(), utStruct)
		if err != nil {
			log.Printf("%s.%s: %s", obj.Pkg().Name(), obj.Id(), err)
			return 1
		}

		flatenedStruct, err = addTagsToStruct(flatenedStruct)
		if err != nil {
			log.Printf("%s.%s: %s", obj.Pkg().Name(), obj.Id(), err)
			return 1
		}

		newStructName := "Flat" + id.Name
		structs = append(structs, StructDef{
			OriginalStructName: id.Name,
			FlatStructName:     newStructName,
			Struct:             flatenedStruct,
		})

		for k, v := range getUsedImports(flatenedStruct) {
			if _, found := usedImports[k]; !found {
				usedImports[k] = v
			}
		}
	}

	out := bytes.NewBuffer(nil)

	fmt.Fprintf(out, `// Code generated by "packer-sdc %s"; DO NOT EDIT.`, cmdPrefix)
	fmt.Fprintf(out, "\n\npackage %s\n", topPkg.Name)

	delete(usedImports, NamePath{topPkg.Name, topPkg.PkgPath})
	usedImports[NamePath{"hcldec", "github.com/hashicorp/hcl/v2/hcldec"}] = types.NewPackage("hcldec", "github.com/hashicorp/hcl/v2/hcldec")
	usedImports[NamePath{"cty", "github.com/zclconf/go-cty/cty"}] = types.NewPackage("cty", "github.com/zclconf/go-cty/cty")
	outputImports(out, usedImports)

	sort.Slice(structs, func(i int, j int) bool {
		return structs[i].OriginalStructName < structs[j].OriginalStructName
	})
	for _, flatenedStruct := range structs {
		fmt.Fprintf(out, "\n// %s is an auto-generated flat version of %s.", flatenedStruct.FlatStructName, flatenedStruct.OriginalStructName)
		fmt.Fprintf(out, "\n// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.")
		fmt.Fprintf(out, "\ntype %s struct {\n", flatenedStruct.FlatStructName)
		outputStructFields(out, flatenedStruct.Struct)
		fmt.Fprint(out, "}\n")

		fmt.Fprintf(out, "\n// FlatMapstructure returns a new %s.", flatenedStruct.FlatStructName)
		fmt.Fprintf(out, "\n// %s is an auto-generated flat version of %s.", flatenedStruct.FlatStructName, flatenedStruct.OriginalStructName)
		fmt.Fprintf(out, "\n// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.")
		fmt.Fprintf(out, "\nfunc (*%s) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {", flatenedStruct.OriginalStructName)
		fmt.Fprintf(out, "\nreturn new(%s)", flatenedStruct.FlatStructName)
		fmt.Fprint(out, "\n}\n")

		fmt.Fprintf(out, "\n// HCL2Spec returns the hcl spec of a %s.", flatenedStruct.OriginalStructName)
		fmt.Fprintf(out, "\n// This spec is used by HCL to read the fields of %s.", flatenedStruct.OriginalStructName)
		fmt.Fprintf(out, "\n// The decoded values from this spec will then be applied to a %s.", flatenedStruct.FlatStructName)
		fmt.Fprintf(out, "\nfunc (*%s) HCL2Spec() map[string]hcldec.Spec {\n", flatenedStruct.FlatStructName)
		outputStructHCL2SpecBody(out, flatenedStruct.Struct)
		fmt.Fprint(out, "}\n")
	}

	for impt := range usedImports {
		if strings.ContainsAny(impt.Path, "/") {
			out = bytes.NewBuffer(bytes.ReplaceAll(out.Bytes(),
				[]byte(impt.Path+"."),
				[]byte(impt.Name+".")))
		}
	}

	// avoid needing to import current pkg; there's probably a better way.
	out = bytes.NewBuffer(bytes.ReplaceAll(out.Bytes(),
		[]byte(topPkg.PkgPath+"."),
		nil))

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("os.Create: %v", err)
	}

	_, err = outputFile.Write(goFmt(outputFile.Name(), out.Bytes()))
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	return 0
}

type StructDef struct {
	OriginalStructName string
	FlatStructName     string
	Struct             *types.Struct
}

// outputStructHCL2SpecBody writes the map[string]hcldec.Spec that defines the HCL spec of a
// struct. Based on the layout of said struct.
// If a field of s is a struct then the HCL2Spec() function of that struct will be called, otherwise a
// cty.Type is outputed.
func outputStructHCL2SpecBody(w io.Writer, s *types.Struct) {
	fmt.Fprintf(w, "s := map[string]hcldec.Spec{\n")

	for i := 0; i < s.NumFields(); i++ {
		field, tag := s.Field(i), s.Tag(i)
		st, _ := structtag.Parse(tag)
		ctyTag, _ := st.Get("cty")
		fmt.Fprintf(w, "	\"%s\": ", ctyTag.Name)
		outputHCL2SpecField(w, ctyTag.Name, field.Type(), st)
		fmt.Fprintln(w, `,`)
	}

	fmt.Fprintln(w, `}`)
	fmt.Fprintln(w, `return s`)
}

// outputHCL2SpecField is called on each field of a struct.
// outputHCL2SpecField writes the values of the `map[string]hcldec.Spec` map
// supposed to define the HCL spec of a struct.
func outputHCL2SpecField(w io.Writer, accessor string, fieldType types.Type, tags *structtag.Tags) {
	if m2h, err := tags.Get(cmdPrefix); err == nil && m2h.HasOption("self-defined") {
		fmt.Fprintf(w, `(&%s{}).HCL2Spec()`, fieldType.String())
		return
	}
	spec, _ := goFieldToCtyType(accessor, fieldType, tags)
	switch spec := spec.(type) {
	case string:
		fmt.Fprint(w, spec)
	default:
		fmt.Fprintf(w, `%#v`, spec)
	}

}

// goFieldToCtyType is a recursive method that returns a cty.Type (or a string) based on the fieldType.
// goFieldToCtyType returns the values of the `map[string]hcldec.Spec` map
// supposed to define the HCL spec of a struct.
// To allow it to be recursive, the method returns two values: an interface that can either be
// a cty.Type or a string. The second argument is used for recursion and is the
// type that will be used by the parent. For example when fieldType is a []string; a
// recursive goFieldToCtyType call will return a cty.String.
func goFieldToCtyType(accessor string, fieldType types.Type, tags *structtag.Tags) (interface{}, cty.Type) {
	switch f := fieldType.(type) {
	case *types.Pointer:
		return goFieldToCtyType(accessor, f.Elem(), tags)
	case *types.Basic:
		if f.Kind() == types.String {
			hcl, err1 := tags.Get("hcl")
			hcllabelindex, err2 := tags.Get(HCL_LABEL_INDEX_KEY)
			if err1 == nil && err2 == nil && hcl.HasOption("label") {
				index, err := strconv.Atoi(hcllabelindex.Name)
				if err != nil {
					panic(err)
				}
				return &hcldec.BlockLabelSpec{
					Index: index,
					Name:  accessor,
				}, cty.NilType
			}
		}

		ctyType := basicKindToCtyType(f.Kind())

		return &hcldec.AttrSpec{
			Name:     accessor,
			Type:     ctyType,
			Required: hasTrueRequiredStructTag(tags),
		}, ctyType
	case *types.Map:
		return &hcldec.AttrSpec{
			Name:     accessor,
			Type:     cty.Map(cty.String), // for now everything can be simplified to a map[string]string
			Required: hasTrueRequiredStructTag(tags),
		}, cty.Map(cty.String)
	case *types.Named:
		// Named is the relative type when of a field with a struct.
		// E.g. SourceAmiFilter    *common.FlatAmiFilterOptions
		// SourceAmiFilter will become a block with nested elements from the struct itself.
		underlyingType := f.Underlying()
		switch underlyingType.(type) {
		case *types.Struct:
			// A struct returns NilType because its HCL2Spec is written in the related file
			// and we don't need to write it again.
			reqString := `false`
			if hasTrueRequiredStructTag(tags) {
				reqString = `true`
			}
			return fmt.Sprintf(
				`&hcldec.BlockSpec{TypeName: "%s", `+
					`Nested: hcldec.ObjectSpec((*%s)(nil).HCL2Spec()), `+
					`Required: %s}`,
				accessor, f.String(), reqString,
			), cty.NilType
		default:
			return goFieldToCtyType(accessor, underlyingType, tags)
		}
	case *types.Slice:
		elem := f.Elem()
		if ptr, isPtr := elem.(*types.Pointer); isPtr {
			elem = ptr.Elem()
		}
		switch elem := elem.(type) {
		case *types.Named:
			// A Slice of Named is the relative type of a filed with a slice of structs.
			// E.g. LaunchMappings []common.FlatBlockDevice
			// LaunchMappings will validate more than one block with nested elements.
			b := bytes.NewBuffer(nil)
			underlyingType := elem.Underlying()
			switch underlyingType.(type) {
			case *types.Struct:
				fmt.Fprintf(b, `hcldec.ObjectSpec((*%s)(nil).HCL2Spec())`, elem.String())
			}
			minCount := 0
			if hasTrueRequiredStructTag(tags) {
				minCount = 1
			}
			return fmt.Sprintf(
				`&hcldec.BlockListSpec{TypeName: "%s", Nested: %s, MinItems: %d}`,
				accessor, b.String(), minCount,
			), cty.NilType
		default:
			_, specType := goFieldToCtyType(accessor, elem, tags)
			if specType == cty.NilType {
				return goFieldToCtyType(accessor, elem.Underlying(), tags)
			}
			return &hcldec.AttrSpec{
				Name:     accessor,
				Type:     cty.List(specType),
				Required: hasTrueRequiredStructTag(tags),
			}, cty.List(specType)
		}
	}
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, `%#v`, &hcldec.AttrSpec{
		Name:     accessor,
		Type:     basicKindToCtyType(types.Bool),
		Required: hasTrueRequiredStructTag(tags),
	})
	fmt.Fprintf(b, `/* TODO(azr): could not find type */`)
	return b.String(), cty.NilType
}

func basicKindToCtyType(kind types.BasicKind) cty.Type {
	switch kind {
	case types.Bool:
		return cty.Bool
	case types.String:
		return cty.String
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.Float32, types.Float64,
		types.Complex64, types.Complex128:
		return cty.Number
	case types.Invalid:
		return cty.String // TODO(azr): fix that beforehand ?
	default:
		log.Printf("Un handled basic kind: %d", kind)
		return cty.String
	}
}

func outputStructFields(w io.Writer, s *types.Struct) {
	for i := 0; i < s.NumFields(); i++ {
		field, tag := s.Field(i), s.Tag(i)
		st, err := structtag.Parse(tag)
		if err == nil {
			// Remove hcllabelindex from the printout because it is not needed
			// in the generated struct.
			st.Delete(HCL_LABEL_INDEX_KEY)
		}
		fieldNameStr := field.String()
		fieldNameStr = strings.Replace(fieldNameStr, "field ", "", 1)
		fmt.Fprintf(w, "	%s `%s`\n", fieldNameStr, st.String())
	}
}

type NamePath struct {
	Name, Path string
}

func outputImports(w io.Writer, imports map[NamePath]*types.Package) {
	if len(imports) == 0 {
		return
	}
	// naive implementation
	pkgs := []NamePath{}
	for k := range imports {
		pkgs = append(pkgs, k)
	}
	sort.Slice(pkgs, func(i int, j int) bool {
		return pkgs[i].Path < pkgs[j].Path
	})

	fmt.Fprint(w, "import (\n")
	for _, pkg := range pkgs {
		if pkg.Name == pkg.Path || strings.HasSuffix(pkg.Path, "/"+pkg.Name) {
			fmt.Fprintf(w, "	\"%s\"\n", pkg.Path)
		} else {
			fmt.Fprintf(w, "	%s \"%s\"\n", pkg.Name, pkg.Path)
		}
	}
	fmt.Fprint(w, ")\n")
}

func getUsedImports(s *types.Struct) map[NamePath]*types.Package {
	res := map[NamePath]*types.Package{}
	for i := 0; i < s.NumFields(); i++ {
		fieldType := s.Field(i).Type()
		if p, ok := fieldType.(*types.Pointer); ok {
			fieldType = p.Elem()
		}
		if p, ok := fieldType.(*types.Slice); ok {
			fieldType = p.Elem()
		}
		namedType, ok := fieldType.(*types.Named)
		if !ok {
			continue
		}
		pkg := namedType.Obj().Pkg()
		if pkg == nil {
			continue
		}
		res[NamePath{pkg.Name(), pkg.Path()}] = pkg
	}
	return res
}

func isCtyStringOrStringPointer(field *types.Var) bool {
	switch f := field.Type().(type) {
	case *types.Basic:
		if f.Kind() == types.String {
			return true
		}
	case *types.Pointer:
		switch fp := f.Elem().(type) {
		case *types.Basic:
			if fp.Kind() == types.String {
				return true
			}
		}
	}
	return false
}

func addTagsToStruct(s *types.Struct) (*types.Struct, error) {
	var hclLabelIndex = 0

	vars, tags := structFields(s)
	for i := range tags {
		field, tag := vars[i], tags[i]
		ctyAccessor := ToSnakeCase(field.Name())
		var hclOptions []string
		st, err := structtag.Parse(tag)
		if err == nil {
			if ms, err := st.Get("mapstructure"); err == nil && ms.Name != "" {
				ctyAccessor = ms.Name
			}
			if hcl, err := st.Get("hcl"); err == nil && hcl.HasOption("label") {
				if !isCtyStringOrStringPointer(field) {
					return nil, fmt.Errorf("field %q has an `hcl:\",label\"` struct tag but is not a string or string pointer", ctyAccessor)
				}
				hclOptions = append(hclOptions, "label")
				st.Set(&structtag.Tag{Key: HCL_LABEL_INDEX_KEY, Name: fmt.Sprintf("%d", hclLabelIndex)})
				hclLabelIndex++
				if required, err := st.Get("required"); err != nil || required.Name != "true" {
					return nil, fmt.Errorf("field %q has an `hcl:\",label\"` struct tag, but has a malformed or missing `required:\"true\"` struct tag", ctyAccessor)
				}
			}
		}
		_ = st.Set(&structtag.Tag{Key: "cty", Name: ctyAccessor})
		_ = st.Set(&structtag.Tag{Key: "hcl", Name: ctyAccessor, Options: hclOptions})
		tags[i] = st.String()
	}

	vars, tags, err := uniqueTags("cty", vars, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to add tag to struct: %s", err)
	}

	return types.NewStruct(vars, tags), nil
}

func uniqueTags(tagName string, fields []*types.Var, tags []string) ([]*types.Var, []string, error) {
	outVars := []*types.Var{}
	outTags := []string{}
	uniqueTags := map[string]bool{}
	for i := range fields {
		field, tag := fields[i], tags[i]
		structtag, _ := structtag.Parse(tag)
		h, err := structtag.Get(tagName)
		if err == nil {
			if uniqueTags[h.Name] {
				return nil, nil, fmt.Errorf("field %q: duplicate tag %q", field.Name(), tagName)
			}
			uniqueTags[h.Name] = true
		}
		outVars = append(outVars, field)
		outTags = append(outTags, tag)
	}
	return outVars, outTags, nil
}

// getMapstructureSquashedStruct will return the same struct but embedded
// fields with a `mapstructure:",squash"` tag will be un-nested.
func getMapstructureSquashedStruct(topPkg *types.Package, utStruct *types.Struct) (*types.Struct, error) {
	res := &types.Struct{}
	for i := 0; i < utStruct.NumFields(); i++ {
		field, tag := utStruct.Field(i), utStruct.Tag(i)
		if !field.Exported() {
			continue
		}
		if _, ok := field.Type().(*types.Signature); ok {
			continue // ignore funcs
		}
		structtag, err := structtag.Parse(tag)
		if err != nil {
			log.Printf("could not parse field tag %s of : %v", tag, err)
			continue
		}

		// Contains mapstructure-to-hcl2 tag
		if ms, err := structtag.Get("mapstructure-to-hcl2"); err == nil {
			// Stop if is telling to skip it
			if ms.HasOption("skip") {
				continue
			}
		}

		// Contains mapstructure tag
		if ms, err := structtag.Get("mapstructure"); err == nil {
			// Squash structs
			if ms.HasOption("squash") {
				ot := field.Type()
				uot := ot.Underlying()
				utStruct, utOk := uot.(*types.Struct)
				if !utOk {
					continue
				}

				sqStr, err := getMapstructureSquashedStruct(topPkg, utStruct)
				if err != nil {
					return nil, err
				}

				res, err = squashStructs(res, sqStr)
				if err != nil {
					return nil, err
				}

				continue
			}
		}

		if field.Pkg() != topPkg {
			field = types.NewField(field.Pos(), topPkg, field.Name(), field.Type(), field.Embedded())
		}
		if p, isPointer := field.Type().(*types.Pointer); isPointer {
			// in order to make the following switch simpler we 'unwrap' this
			// pointer all structs are going to be made pointers anyways.
			field = types.NewField(field.Pos(), field.Pkg(), field.Name(), p.Elem(), field.Embedded())
		}
		switch f := field.Type().(type) {
		case *types.Named:
			switch f.String() {
			case "time.Duration":
				field = types.NewField(field.Pos(), field.Pkg(), field.Name(), types.NewPointer(types.Typ[types.String]), field.Embedded())
			case "github.com/hashicorp/packer-plugin-sdk/template/config.Trilean": // TODO(azr): unhack this situation
				field = types.NewField(field.Pos(), field.Pkg(), field.Name(), types.NewPointer(types.Typ[types.Bool]), field.Embedded())
			case "github.com/hashicorp/packer/provisioner/powershell.ExecutionPolicy": // TODO(azr): unhack this situation
				field = types.NewField(field.Pos(), field.Pkg(), field.Name(), types.NewPointer(types.Typ[types.String]), field.Embedded())
			default:
				if str, isStruct := f.Underlying().(*types.Struct); isStruct {
					obj := flattenNamed(f, str)
					field = types.NewField(field.Pos(), field.Pkg(), field.Name(), obj, field.Embedded())
					field = makePointer(field)
				}
				if slice, isSlice := f.Underlying().(*types.Slice); isSlice {
					if f, fNamed := slice.Elem().(*types.Named); fNamed {
						if str, isStruct := f.Underlying().(*types.Struct); isStruct {
							// this is a slice of named structs; we want to change
							// the struct ref to a 'FlatStruct'.
							obj := flattenNamed(f, str)
							slice := types.NewSlice(obj)
							field = types.NewField(field.Pos(), field.Pkg(), field.Name(), slice, field.Embedded())
						}
					}
				}
				if _, isBasic := f.Underlying().(*types.Basic); isBasic {
					field = makePointer(field)
				}
			}
		case *types.Slice:
			if f, fNamed := f.Elem().(*types.Named); fNamed {
				if str, isStruct := f.Underlying().(*types.Struct); isStruct {
					obj := flattenNamed(f, str)
					field = types.NewField(field.Pos(), field.Pkg(), field.Name(), types.NewSlice(obj), field.Embedded())
				}
			}
		case *types.Basic:
			// since everything is optional, everything must be a pointer
			// non optional fields should be non pointers.
			field = makePointer(field)
		}
		res, err = addFieldToStruct(res, field, tag)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func flattenNamed(f *types.Named, underlying types.Type) *types.Named {
	obj := f.Obj()
	obj = types.NewTypeName(obj.Pos(), obj.Pkg(), "Flat"+obj.Name(), obj.Type())
	return types.NewNamed(obj, underlying, nil)
}

func makePointer(field *types.Var) *types.Var {
	return types.NewField(field.Pos(), field.Pkg(), field.Name(), types.NewPointer(field.Type()), field.Embedded())
}

func addFieldToStruct(s *types.Struct, field *types.Var, tag string) (*types.Struct, error) {
	sf, st := structFields(s)

	vars, tags, err := uniqueFields(append(sf, field), append(st, tag))
	if err != nil {
		return nil, err
	}

	str := types.NewStruct(vars, tags)
	return str, nil
}

func squashStructs(a, b *types.Struct) (*types.Struct, error) {
	va, ta := structFields(a)
	vb, tb := structFields(b)

	vars, tags, err := uniqueFields(append(va, vb...), append(ta, tb...))
	if err != nil {
		return nil, fmt.Errorf("failed to squash struct: %s", err)
	}

	str := types.NewStruct(vars, tags)
	return str, nil
}

func uniqueFields(fields []*types.Var, tags []string) ([]*types.Var, []string, error) {
	outVars := []*types.Var{}
	outTags := []string{}
	fieldNames := map[string]bool{}
	for i := range fields {
		field, tag := fields[i], tags[i]
		if fieldNames[field.Name()] {
			return nil, nil, fmt.Errorf("duplicate field %q", field.Name())
		}
		fieldNames[field.Name()] = true
		outVars = append(outVars, field)
		outTags = append(outTags, tag)
	}
	return outVars, outTags, nil
}

func structFields(s *types.Struct) (vars []*types.Var, tags []string) {
	for i := 0; i < s.NumFields(); i++ {
		field, tag := s.Field(i), s.Tag(i)
		vars = append(vars, field)
		tags = append(tags, tag)
	}
	return vars, tags
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func goFmt(filename string, b []byte) []byte {
	fb, err := imports.Process(filename, b, nil)
	if err != nil {
		log.Printf("formatting err: %v", err)
		return b
	}
	return fb
}

func hasTrueRequiredStructTag(st *structtag.Tags) bool {
	if st == nil {
		return false
	}

	requiredTag, err := st.Get("required")
	return err == nil && requiredTag.Name == "true"
}
