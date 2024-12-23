// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// ToProto converts a hcldec.Spec to a protobuf-serialisable equivalent.
//
// This can then be used for gRPC communication over-the-wire for Packer plugins.
func ToProto(spec hcldec.Spec) (*Spec, error) {
	switch concreteSpec := spec.(type) {
	case *hcldec.AttrSpec:
		attr, err := attrSpecToProto(concreteSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to decode attribute spec: %s", err)
		}
		return &Spec{
			Block: &Spec_Attr{
				Attr: attr,
			},
		}, nil
	case *hcldec.BlockListSpec:
		blspec, err := blockListSpecToProto(concreteSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to decode block list spec: %s", err)
		}
		return &Spec{
			Block: &Spec_BlockList{
				BlockList: blspec,
			},
		}, nil
	case hcldec.ObjectSpec:
		objSpec, err := objectSpecToProto(&concreteSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to decode object spec: %s", err)
		}
		return &Spec{
			Block: &Spec_Object{
				Object: objSpec,
			},
		}, nil
	case *hcldec.BlockSpec:
		blockSpec, err := blockSpecToProto(concreteSpec)
		if err != nil {
			return nil, fmt.Errorf("failed to decode object spec: %s", err)
		}
		return &Spec{
			Block: &Spec_BlockValue{
				BlockValue: blockSpec,
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported hcldec.Spec type: %#v", reflect.TypeOf(spec).String())
}

func blockSpecToProto(bsp *hcldec.BlockSpec) (*Block, error) {
	nst, err := ToProto(bsp.Nested)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block nested spec: %s", err)
	}

	return &Block{
		Name:     bsp.TypeName,
		Required: bsp.Required,
		Nested:   nst,
	}, nil
}

func objectSpecToProto(hsp *hcldec.ObjectSpec) (*Object, error) {
	outSpec := &Object{
		Attributes: map[string]*Spec{},
	}

	for name, spec := range *hsp {
		pSpec, err := ToProto(spec)
		if err != nil {
			return nil, fmt.Errorf("failed to convert object attribute %q to hcldec.Spec: %s", name, err)
		}
		outSpec.Attributes[name] = pSpec
	}

	return outSpec, nil
}

func blockListSpecToProto(blspec *hcldec.BlockListSpec) (*BlockList, error) {
	nested, err := ToProto(blspec.Nested)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block list nested type: %s", err)
	}

	return &BlockList{
		Name:   blspec.TypeName,
		Nested: nested,
	}, nil
}

func attrSpecToProto(attrSpec *hcldec.AttrSpec) (*Attr, error) {
	convertedType, err := ctyTypeToProto(attrSpec.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ctyType for attribute %q: %s", attrSpec.Name, err)
	}

	return &Attr{
		Name:     attrSpec.Name,
		Required: attrSpec.Required,
		Type:     convertedType,
	}, nil
}

func ctyTypeToProto(cType cty.Type) (*CtyType, error) {
	if cType.IsPrimitiveType() {
		switch cType {
		case cty.Bool,
			cty.String,
			cty.Number:
			return &CtyType{
				TypeDef: &CtyType_Primitive{
					Primitive: &CtyPrimitive{
						TypeString: cType.GoString(),
					},
				},
			}, nil
		default:
			return nil, fmt.Errorf("Unknown primitive type: %s", cType.FriendlyName())
		}
	}

	if cType.IsListType() {
		el := cType.ElementType()
		elType, err := ctyTypeToProto(el)
		if err != nil {
			return nil, fmt.Errorf("failed to extract valid cty.Type from list element: %s", err)
		}
		return &CtyType{
			TypeDef: &CtyType_List{
				List: &CtyList{
					ElementType: elType,
				},
			},
		}, nil
	}

	// As per the specification, cty.Map are always a map from string to a cty type
	//
	// Therefore, we don't need to worry about other types than the element's
	if cType.IsMapType() {
		el := cType.MapElementType()
		elType, err := ctyTypeToProto(*el)
		if err != nil {
			return nil, fmt.Errorf("failed to extract valid cty.Type from map: %s", err)
		}
		return &CtyType{
			TypeDef: &CtyType_Map{
				Map: &CtyMap{
					ElementType: elType,
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported cty.Type conversion to protobuf-compatible structure: %+v", cType)
}

func (spec *Spec) FromProto() (hcldec.Spec, error) {
	switch realSpec := spec.Block.(type) {
	case *Spec_Attr:
		return protoArgToHCLDecSpec(realSpec.Attr)
	case *Spec_BlockList:
		return protoBlockListToHCLDecSpec(realSpec.BlockList)
	case *Spec_BlockValue:
		return protoBlockToHCLDecSpec(realSpec.BlockValue)
	case *Spec_Object:
		return protoObjectSpecToHCLDecSpec(realSpec.Object)
	}

	return nil, fmt.Errorf("unsupported spec type: %s", spec.String())
}

func protoObjectSpecToHCLDecSpec(protoSpec *Object) (*hcldec.ObjectSpec, error) {
	outSpec := hcldec.ObjectSpec{}

	for name, spec := range protoSpec.Attributes {
		attrSpec, err := spec.FromProto()
		if err != nil {
			return nil, fmt.Errorf("failed to decode object attribute %q: %s", name, err)
		}
		outSpec[name] = attrSpec
	}

	return &outSpec, nil
}

func protoBlockToHCLDecSpec(bl *Block) (*hcldec.BlockSpec, error) {
	nested, err := bl.Nested.FromProto()
	if err != nil {
		return nil, fmt.Errorf("failed to decode block nested type from proto: %s", err)
	}

	return &hcldec.BlockSpec{
		TypeName: bl.Name,
		Required: bl.Required,
		Nested:   nested,
	}, nil
}

func protoBlockListToHCLDecSpec(bll *BlockList) (*hcldec.BlockListSpec, error) {
	blSpec := bll.Nested
	nested, err := blSpec.FromProto()
	if err != nil {
		return nil, fmt.Errorf("failed to decode block list nested type from proto: %s", err)
	}

	return &hcldec.BlockListSpec{
		TypeName: bll.Name,
		Nested:   nested,
	}, nil
}

func protoArgToHCLDecSpec(attr *Attr) (*hcldec.AttrSpec, error) {
	relType, err := protoTypeToCtyType(attr.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to convert type of attribute %q: %s", attr.Name, err)
	}

	return &hcldec.AttrSpec{
		Name:     attr.Name,
		Required: attr.Required,
		Type:     relType,
	}, nil
}

func protoTypeToCtyType(protoType *CtyType) (cty.Type, error) {
	switch concrete := protoType.TypeDef.(type) {
	case *CtyType_Primitive:
		switch concrete.Primitive.TypeString {
		case "cty.String":
			return cty.String, nil
		case "cty.Bool":
			return cty.Bool, nil
		case "cty.Number":
			return cty.Number, nil
		}
	case *CtyType_List:
		elType, err := protoTypeToCtyType(concrete.List.ElementType)
		if err != nil {
			return cty.NilType, fmt.Errorf("failed to convert list element type: %s", err)
		}
		return cty.List(elType), nil
	case *CtyType_Map:
		elType, err := protoTypeToCtyType(concrete.Map.ElementType)
		if err != nil {
			return cty.NilType, fmt.Errorf("failed to convert map element type: %s", err)
		}
		return cty.Map(elType), nil
	}

	return cty.NilType, fmt.Errorf("unsupported cty.Type: %+v", protoType)
}
