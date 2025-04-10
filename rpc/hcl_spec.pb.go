// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: hcl_spec.proto

package rpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CtyType is any of the cty types that can be used for a Attribute.
//
// Bodies aren't an issue since they're encompassing a bunch of different
// attributes, which end-up referencing a type from this structure.
type CtyType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to TypeDef:
	//
	//	*CtyType_Primitive
	//	*CtyType_List
	//	*CtyType_Map
	TypeDef isCtyType_TypeDef `protobuf_oneof:"typeDef"`
}

func (x *CtyType) Reset() {
	*x = CtyType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CtyType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CtyType) ProtoMessage() {}

func (x *CtyType) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CtyType.ProtoReflect.Descriptor instead.
func (*CtyType) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{0}
}

func (m *CtyType) GetTypeDef() isCtyType_TypeDef {
	if m != nil {
		return m.TypeDef
	}
	return nil
}

func (x *CtyType) GetPrimitive() *CtyPrimitive {
	if x, ok := x.GetTypeDef().(*CtyType_Primitive); ok {
		return x.Primitive
	}
	return nil
}

func (x *CtyType) GetList() *CtyList {
	if x, ok := x.GetTypeDef().(*CtyType_List); ok {
		return x.List
	}
	return nil
}

func (x *CtyType) GetMap() *CtyMap {
	if x, ok := x.GetTypeDef().(*CtyType_Map); ok {
		return x.Map
	}
	return nil
}

type isCtyType_TypeDef interface {
	isCtyType_TypeDef()
}

type CtyType_Primitive struct {
	Primitive *CtyPrimitive `protobuf:"bytes,1,opt,name=primitive,proto3,oneof"`
}

type CtyType_List struct {
	List *CtyList `protobuf:"bytes,2,opt,name=list,proto3,oneof"`
}

type CtyType_Map struct {
	Map *CtyMap `protobuf:"bytes,3,opt,name=map,proto3,oneof"`
}

func (*CtyType_Primitive) isCtyType_TypeDef() {}

func (*CtyType_List) isCtyType_TypeDef() {}

func (*CtyType_Map) isCtyType_TypeDef() {}

// CtyPrimitive is any of the cty.Type that match the `IsPrimitiveType` function
// i.e. either Number, Bool or String.
type CtyPrimitive struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TypeString string `protobuf:"bytes,1,opt,name=typeString,proto3" json:"typeString,omitempty"`
}

func (x *CtyPrimitive) Reset() {
	*x = CtyPrimitive{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CtyPrimitive) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CtyPrimitive) ProtoMessage() {}

func (x *CtyPrimitive) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CtyPrimitive.ProtoReflect.Descriptor instead.
func (*CtyPrimitive) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{1}
}

func (x *CtyPrimitive) GetTypeString() string {
	if x != nil {
		return x.TypeString
	}
	return ""
}

// CtyList is a list of a cty.Type
type CtyList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ElementType *CtyType `protobuf:"bytes,1,opt,name=elementType,proto3" json:"elementType,omitempty"`
}

func (x *CtyList) Reset() {
	*x = CtyList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CtyList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CtyList) ProtoMessage() {}

func (x *CtyList) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CtyList.ProtoReflect.Descriptor instead.
func (*CtyList) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{2}
}

func (x *CtyList) GetElementType() *CtyType {
	if x != nil {
		return x.ElementType
	}
	return nil
}

// CtyMap is a map from one type to another
type CtyMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ElementType *CtyType `protobuf:"bytes,1,opt,name=elementType,proto3" json:"elementType,omitempty"`
}

func (x *CtyMap) Reset() {
	*x = CtyMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CtyMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CtyMap) ProtoMessage() {}

func (x *CtyMap) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CtyMap.ProtoReflect.Descriptor instead.
func (*CtyMap) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{3}
}

func (x *CtyMap) GetElementType() *CtyType {
	if x != nil {
		return x.ElementType
	}
	return nil
}

// HCL2Spec matches what Packer already consumes from plugins in order to describe
// their contents' schema, and lets Packer decode the configuration provided by
// the user to cty values, and detect problems with the contents before executing them.
//
// These are sent over-the-wire over gRPC, much like the old system did using gob
// encoding and standard go RPC servers.
type HCL2Spec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TypeSpec map[string]*Spec `protobuf:"bytes,1,rep,name=TypeSpec,proto3" json:"TypeSpec,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *HCL2Spec) Reset() {
	*x = HCL2Spec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HCL2Spec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HCL2Spec) ProtoMessage() {}

func (x *HCL2Spec) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HCL2Spec.ProtoReflect.Descriptor instead.
func (*HCL2Spec) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{4}
}

func (x *HCL2Spec) GetTypeSpec() map[string]*Spec {
	if x != nil {
		return x.TypeSpec
	}
	return nil
}

// A Spec is any kind of object that can convert losslessly to any of the hcldec.Spec types.
type Spec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Block:
	//
	//	*Spec_Object
	//	*Spec_Attr
	//	*Spec_BlockValue
	//	*Spec_BlockList
	Block isSpec_Block `protobuf_oneof:"block"`
}

func (x *Spec) Reset() {
	*x = Spec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Spec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Spec) ProtoMessage() {}

func (x *Spec) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Spec.ProtoReflect.Descriptor instead.
func (*Spec) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{5}
}

func (m *Spec) GetBlock() isSpec_Block {
	if m != nil {
		return m.Block
	}
	return nil
}

func (x *Spec) GetObject() *Object {
	if x, ok := x.GetBlock().(*Spec_Object); ok {
		return x.Object
	}
	return nil
}

func (x *Spec) GetAttr() *Attr {
	if x, ok := x.GetBlock().(*Spec_Attr); ok {
		return x.Attr
	}
	return nil
}

func (x *Spec) GetBlockValue() *Block {
	if x, ok := x.GetBlock().(*Spec_BlockValue); ok {
		return x.BlockValue
	}
	return nil
}

func (x *Spec) GetBlockList() *BlockList {
	if x, ok := x.GetBlock().(*Spec_BlockList); ok {
		return x.BlockList
	}
	return nil
}

type isSpec_Block interface {
	isSpec_Block()
}

type Spec_Object struct {
	Object *Object `protobuf:"bytes,1,opt,name=object,proto3,oneof"`
}

type Spec_Attr struct {
	Attr *Attr `protobuf:"bytes,2,opt,name=attr,proto3,oneof"`
}

type Spec_BlockValue struct {
	BlockValue *Block `protobuf:"bytes,3,opt,name=block_value,json=blockValue,proto3,oneof"`
}

type Spec_BlockList struct {
	BlockList *BlockList `protobuf:"bytes,4,opt,name=block_list,json=blockList,proto3,oneof"`
}

func (*Spec_Object) isSpec_Block() {}

func (*Spec_Attr) isSpec_Block() {}

func (*Spec_BlockValue) isSpec_Block() {}

func (*Spec_BlockList) isSpec_Block() {}

// Attr spec type reads the value of an attribute in the current body
// and returns that value as its result. It also creates validation constraints
// for the given attribute name and its value.
type Attr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type     *CtyType `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Required bool     `protobuf:"varint,3,opt,name=required,proto3" json:"required,omitempty"`
}

func (x *Attr) Reset() {
	*x = Attr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attr) ProtoMessage() {}

func (x *Attr) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attr.ProtoReflect.Descriptor instead.
func (*Attr) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{6}
}

func (x *Attr) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Attr) GetType() *CtyType {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *Attr) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

// Block spec type applies one nested spec block to the contents of a
// block within the current body and returns the result of that spec. It also
// creates validation constraints for the given block type name.
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Required bool   `protobuf:"varint,2,opt,name=required,proto3" json:"required,omitempty"`
	Nested   *Spec  `protobuf:"bytes,3,opt,name=nested,proto3" json:"nested,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{7}
}

func (x *Block) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Block) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

func (x *Block) GetNested() *Spec {
	if x != nil {
		return x.Nested
	}
	return nil
}

// BlockList spec type is similar to `Block`, but it accepts zero or
// more blocks of a specified type rather than requiring zero or one. The
// result is a JSON array with one entry per block of the given type.
type BlockList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Nested *Spec  `protobuf:"bytes,2,opt,name=nested,proto3" json:"nested,omitempty"`
}

func (x *BlockList) Reset() {
	*x = BlockList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockList) ProtoMessage() {}

func (x *BlockList) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockList.ProtoReflect.Descriptor instead.
func (*BlockList) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{8}
}

func (x *BlockList) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *BlockList) GetNested() *Spec {
	if x != nil {
		return x.Nested
	}
	return nil
}

// Object spec type is the most commonly used at the root of a spec file.
// Its result is a JSON object whose properties are set based on any nested
// spec blocks:
type Object struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attributes map[string]*Spec `protobuf:"bytes,1,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Object) Reset() {
	*x = Object{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hcl_spec_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Object) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Object) ProtoMessage() {}

func (x *Object) ProtoReflect() protoreflect.Message {
	mi := &file_hcl_spec_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Object.ProtoReflect.Descriptor instead.
func (*Object) Descriptor() ([]byte, []int) {
	return file_hcl_spec_proto_rawDescGZIP(), []int{9}
}

func (x *Object) GetAttributes() map[string]*Spec {
	if x != nil {
		return x.Attributes
	}
	return nil
}

var File_hcl_spec_proto protoreflect.FileDescriptor

var file_hcl_spec_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x68, 0x63, 0x6c, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x80, 0x01, 0x0a, 0x07, 0x43, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2d, 0x0a, 0x09,
	0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x43, 0x74, 0x79, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x48, 0x00,
	0x52, 0x09, 0x70, 0x72, 0x69, 0x6d, 0x69, 0x74, 0x69, 0x76, 0x65, 0x12, 0x1e, 0x0a, 0x04, 0x6c,
	0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x43, 0x74, 0x79, 0x4c,
	0x69, 0x73, 0x74, 0x48, 0x00, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x03, 0x6d,
	0x61, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43, 0x74, 0x79, 0x4d, 0x61,
	0x70, 0x48, 0x00, 0x52, 0x03, 0x6d, 0x61, 0x70, 0x42, 0x09, 0x0a, 0x07, 0x74, 0x79, 0x70, 0x65,
	0x44, 0x65, 0x66, 0x22, 0x2e, 0x0a, 0x0c, 0x43, 0x74, 0x79, 0x50, 0x72, 0x69, 0x6d, 0x69, 0x74,
	0x69, 0x76, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x79, 0x70, 0x65, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x79, 0x70, 0x65, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x22, 0x35, 0x0a, 0x07, 0x43, 0x74, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2a,
	0x0a, 0x0b, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x43, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x65,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0x34, 0x0a, 0x06, 0x43, 0x74,
	0x79, 0x4d, 0x61, 0x70, 0x12, 0x2a, 0x0a, 0x0b, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x43, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0b, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x22, 0x83, 0x01, 0x0a, 0x08, 0x48, 0x43, 0x4c, 0x32, 0x53, 0x70, 0x65, 0x63, 0x12, 0x33, 0x0a,
	0x08, 0x54, 0x79, 0x70, 0x65, 0x53, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x48, 0x43, 0x4c, 0x32, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x53,
	0x70, 0x65, 0x63, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x54, 0x79, 0x70, 0x65, 0x53, 0x70,
	0x65, 0x63, 0x1a, 0x42, 0x0a, 0x0d, 0x54, 0x79, 0x70, 0x65, 0x53, 0x70, 0x65, 0x63, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xa7, 0x01, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x21, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x07, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x48, 0x00, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x61, 0x74, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x05, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x48, 0x00, 0x52, 0x04, 0x61, 0x74, 0x74, 0x72, 0x12,
	0x29, 0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x00, 0x52, 0x0a,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x00, 0x52, 0x09, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x07, 0x0a, 0x05, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x22, 0x54, 0x0a, 0x04, 0x41, 0x74, 0x74, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x43, 0x74, 0x79,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x56, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12,
	0x1d, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x05, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0x3e,
	0x0a, 0x09, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x05, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x22, 0x87,
	0x01, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x37, 0x0a, 0x0a, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x1a, 0x44, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x06, 0x5a, 0x04, 0x72, 0x70, 0x63, 0x2f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hcl_spec_proto_rawDescOnce sync.Once
	file_hcl_spec_proto_rawDescData = file_hcl_spec_proto_rawDesc
)

func file_hcl_spec_proto_rawDescGZIP() []byte {
	file_hcl_spec_proto_rawDescOnce.Do(func() {
		file_hcl_spec_proto_rawDescData = protoimpl.X.CompressGZIP(file_hcl_spec_proto_rawDescData)
	})
	return file_hcl_spec_proto_rawDescData
}

var file_hcl_spec_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_hcl_spec_proto_goTypes = []interface{}{
	(*CtyType)(nil),      // 0: CtyType
	(*CtyPrimitive)(nil), // 1: CtyPrimitive
	(*CtyList)(nil),      // 2: CtyList
	(*CtyMap)(nil),       // 3: CtyMap
	(*HCL2Spec)(nil),     // 4: HCL2Spec
	(*Spec)(nil),         // 5: Spec
	(*Attr)(nil),         // 6: Attr
	(*Block)(nil),        // 7: Block
	(*BlockList)(nil),    // 8: BlockList
	(*Object)(nil),       // 9: Object
	nil,                  // 10: HCL2Spec.TypeSpecEntry
	nil,                  // 11: Object.AttributesEntry
}
var file_hcl_spec_proto_depIdxs = []int32{
	1,  // 0: CtyType.primitive:type_name -> CtyPrimitive
	2,  // 1: CtyType.list:type_name -> CtyList
	3,  // 2: CtyType.map:type_name -> CtyMap
	0,  // 3: CtyList.elementType:type_name -> CtyType
	0,  // 4: CtyMap.elementType:type_name -> CtyType
	10, // 5: HCL2Spec.TypeSpec:type_name -> HCL2Spec.TypeSpecEntry
	9,  // 6: Spec.object:type_name -> Object
	6,  // 7: Spec.attr:type_name -> Attr
	7,  // 8: Spec.block_value:type_name -> Block
	8,  // 9: Spec.block_list:type_name -> BlockList
	0,  // 10: Attr.type:type_name -> CtyType
	5,  // 11: Block.nested:type_name -> Spec
	5,  // 12: BlockList.nested:type_name -> Spec
	11, // 13: Object.attributes:type_name -> Object.AttributesEntry
	5,  // 14: HCL2Spec.TypeSpecEntry.value:type_name -> Spec
	5,  // 15: Object.AttributesEntry.value:type_name -> Spec
	16, // [16:16] is the sub-list for method output_type
	16, // [16:16] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_hcl_spec_proto_init() }
func file_hcl_spec_proto_init() {
	if File_hcl_spec_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hcl_spec_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CtyType); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CtyPrimitive); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CtyList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CtyMap); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HCL2Spec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Spec); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_hcl_spec_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Object); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_hcl_spec_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*CtyType_Primitive)(nil),
		(*CtyType_List)(nil),
		(*CtyType_Map)(nil),
	}
	file_hcl_spec_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*Spec_Object)(nil),
		(*Spec_Attr)(nil),
		(*Spec_BlockValue)(nil),
		(*Spec_BlockList)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hcl_spec_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_hcl_spec_proto_goTypes,
		DependencyIndexes: file_hcl_spec_proto_depIdxs,
		MessageInfos:      file_hcl_spec_proto_msgTypes,
	}.Build()
	File_hcl_spec_proto = out.File
	file_hcl_spec_proto_rawDesc = nil
	file_hcl_spec_proto_goTypes = nil
	file_hcl_spec_proto_depIdxs = nil
}
