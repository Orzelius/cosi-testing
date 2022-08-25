// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.18.1
// source: v1alpha1/resource.proto

// Resource package defines protobuf serialization of COSI resources.

package v1alpha1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LabelTerm_Operation int32

const (
	// Label exists.
	LabelTerm_EXISTS LabelTerm_Operation = 0
	// Label value is equal.
	LabelTerm_EQUAL LabelTerm_Operation = 1
	// Label doesn't exist.
	LabelTerm_NOT_EXISTS LabelTerm_Operation = 2
)

// Enum value maps for LabelTerm_Operation.
var (
	LabelTerm_Operation_name = map[int32]string{
		0: "EXISTS",
		1: "EQUAL",
		2: "NOT_EXISTS",
	}
	LabelTerm_Operation_value = map[string]int32{
		"EXISTS":     0,
		"EQUAL":      1,
		"NOT_EXISTS": 2,
	}
)

func (x LabelTerm_Operation) Enum() *LabelTerm_Operation {
	p := new(LabelTerm_Operation)
	*p = x
	return p
}

func (x LabelTerm_Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LabelTerm_Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_v1alpha1_resource_proto_enumTypes[0].Descriptor()
}

func (LabelTerm_Operation) Type() protoreflect.EnumType {
	return &file_v1alpha1_resource_proto_enumTypes[0]
}

func (x LabelTerm_Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LabelTerm_Operation.Descriptor instead.
func (LabelTerm_Operation) EnumDescriptor() ([]byte, []int) {
	return file_v1alpha1_resource_proto_rawDescGZIP(), []int{3, 0}
}

// Metadata represents resource metadata.
//
// (namespace, type, id) is a recource pointer.
// (version) is a current resource version.
// (owner) is filled in for controller-managed resources with controller name.
// (phase) indicates whether resource is going through tear down phase.
// (finalizers) are attached controllers blocking teardown of the resource.
type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace  string                 `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Type       string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Id         string                 `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Version    string                 `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	Owner      string                 `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	Phase      string                 `protobuf:"bytes,6,opt,name=phase,proto3" json:"phase,omitempty"`
	Created    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created,proto3" json:"created,omitempty"`
	Updated    *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated,proto3" json:"updated,omitempty"`
	Finalizers []string               `protobuf:"bytes,9,rep,name=finalizers,proto3" json:"finalizers,omitempty"`
	Labels     map[string]string      `protobuf:"bytes,10,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1alpha1_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_v1alpha1_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_v1alpha1_resource_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Metadata) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Metadata) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metadata) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Metadata) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Metadata) GetPhase() string {
	if x != nil {
		return x.Phase
	}
	return ""
}

func (x *Metadata) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

func (x *Metadata) GetUpdated() *timestamppb.Timestamp {
	if x != nil {
		return x.Updated
	}
	return nil
}

func (x *Metadata) GetFinalizers() []string {
	if x != nil {
		return x.Finalizers
	}
	return nil
}

func (x *Metadata) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

// Spec defines content of the resource.
type Spec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Protobuf-serialized representation of the resource.
	ProtoSpec []byte `protobuf:"bytes,1,opt,name=proto_spec,json=protoSpec,proto3" json:"proto_spec,omitempty"`
	// YAML representation of the spec (optional).
	YamlSpec string `protobuf:"bytes,2,opt,name=yaml_spec,json=yamlSpec,proto3" json:"yaml_spec,omitempty"`
}

func (x *Spec) Reset() {
	*x = Spec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1alpha1_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Spec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Spec) ProtoMessage() {}

func (x *Spec) ProtoReflect() protoreflect.Message {
	mi := &file_v1alpha1_resource_proto_msgTypes[1]
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
	return file_v1alpha1_resource_proto_rawDescGZIP(), []int{1}
}

func (x *Spec) GetProtoSpec() []byte {
	if x != nil {
		return x.ProtoSpec
	}
	return nil
}

func (x *Spec) GetYamlSpec() string {
	if x != nil {
		return x.YamlSpec
	}
	return ""
}

// Resource is a combination of metadata and spec.
type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec     *Spec     `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1alpha1_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_v1alpha1_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_v1alpha1_resource_proto_rawDescGZIP(), []int{2}
}

func (x *Resource) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Resource) GetSpec() *Spec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// LabelTerm is an expression on a label.
type LabelTerm struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string              `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Op    LabelTerm_Operation `protobuf:"varint,2,opt,name=op,proto3,enum=cosi.resource.LabelTerm_Operation" json:"op,omitempty"`
	Value string              `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LabelTerm) Reset() {
	*x = LabelTerm{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1alpha1_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelTerm) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelTerm) ProtoMessage() {}

func (x *LabelTerm) ProtoReflect() protoreflect.Message {
	mi := &file_v1alpha1_resource_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelTerm.ProtoReflect.Descriptor instead.
func (*LabelTerm) Descriptor() ([]byte, []int) {
	return file_v1alpha1_resource_proto_rawDescGZIP(), []int{3}
}

func (x *LabelTerm) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *LabelTerm) GetOp() LabelTerm_Operation {
	if x != nil {
		return x.Op
	}
	return LabelTerm_EXISTS
}

func (x *LabelTerm) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// LabelQuery is a query on resource metadata labels.
//
// Terms are combined with AND.
type LabelQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Terms []*LabelTerm `protobuf:"bytes,1,rep,name=terms,proto3" json:"terms,omitempty"`
}

func (x *LabelQuery) Reset() {
	*x = LabelQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_v1alpha1_resource_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabelQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabelQuery) ProtoMessage() {}

func (x *LabelQuery) ProtoReflect() protoreflect.Message {
	mi := &file_v1alpha1_resource_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabelQuery.ProtoReflect.Descriptor instead.
func (*LabelQuery) Descriptor() ([]byte, []int) {
	return file_v1alpha1_resource_proto_rawDescGZIP(), []int{4}
}

func (x *LabelQuery) GetTerms() []*LabelTerm {
	if x != nil {
		return x.Terms
	}
	return nil
}

var File_v1alpha1_resource_proto protoreflect.FileDescriptor

var file_v1alpha1_resource_proto_rawDesc = []byte{
	0x0a, 0x17, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x63, 0x6f, 0x73, 0x69, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x96, 0x03, 0x0a, 0x08, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x68, 0x61, 0x73,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x68, 0x61, 0x73, 0x65, 0x12, 0x34,
	0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x07, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x69,
	0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x66, 0x69, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x73, 0x12, 0x3b, 0x0a, 0x06, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x6f, 0x73,
	0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x42, 0x0a, 0x04, 0x53, 0x70, 0x65, 0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x53, 0x70, 0x65, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x79, 0x61, 0x6d,
	0x6c, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x79, 0x61,
	0x6d, 0x6c, 0x53, 0x70, 0x65, 0x63, 0x22, 0x68, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x73, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x73, 0x69, 0x2e, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63,
	0x22, 0x9b, 0x01, 0x0a, 0x09, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x32, 0x0a, 0x02, 0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x63,
	0x6f, 0x73, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x61, 0x62,
	0x65, 0x6c, 0x54, 0x65, 0x72, 0x6d, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x02, 0x6f, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x32, 0x0a, 0x09, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x58, 0x49, 0x53, 0x54,
	0x53, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x10, 0x02, 0x22, 0x3c,
	0x0a, 0x0a, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x2e, 0x0a, 0x05,
	0x74, 0x65, 0x72, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f,
	0x73, 0x69, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x54, 0x65, 0x72, 0x6d, 0x52, 0x05, 0x74, 0x65, 0x72, 0x6d, 0x73, 0x42, 0x2e, 0x5a, 0x2c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x69, 0x2d,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_v1alpha1_resource_proto_rawDescOnce sync.Once
	file_v1alpha1_resource_proto_rawDescData = file_v1alpha1_resource_proto_rawDesc
)

func file_v1alpha1_resource_proto_rawDescGZIP() []byte {
	file_v1alpha1_resource_proto_rawDescOnce.Do(func() {
		file_v1alpha1_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_v1alpha1_resource_proto_rawDescData)
	})
	return file_v1alpha1_resource_proto_rawDescData
}

var file_v1alpha1_resource_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_v1alpha1_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_v1alpha1_resource_proto_goTypes = []interface{}{
	(LabelTerm_Operation)(0),      // 0: cosi.resource.LabelTerm.Operation
	(*Metadata)(nil),              // 1: cosi.resource.Metadata
	(*Spec)(nil),                  // 2: cosi.resource.Spec
	(*Resource)(nil),              // 3: cosi.resource.Resource
	(*LabelTerm)(nil),             // 4: cosi.resource.LabelTerm
	(*LabelQuery)(nil),            // 5: cosi.resource.LabelQuery
	nil,                           // 6: cosi.resource.Metadata.LabelsEntry
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_v1alpha1_resource_proto_depIdxs = []int32{
	7, // 0: cosi.resource.Metadata.created:type_name -> google.protobuf.Timestamp
	7, // 1: cosi.resource.Metadata.updated:type_name -> google.protobuf.Timestamp
	6, // 2: cosi.resource.Metadata.labels:type_name -> cosi.resource.Metadata.LabelsEntry
	1, // 3: cosi.resource.Resource.metadata:type_name -> cosi.resource.Metadata
	2, // 4: cosi.resource.Resource.spec:type_name -> cosi.resource.Spec
	0, // 5: cosi.resource.LabelTerm.op:type_name -> cosi.resource.LabelTerm.Operation
	4, // 6: cosi.resource.LabelQuery.terms:type_name -> cosi.resource.LabelTerm
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_v1alpha1_resource_proto_init() }
func file_v1alpha1_resource_proto_init() {
	if File_v1alpha1_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_v1alpha1_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
		file_v1alpha1_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_v1alpha1_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
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
		file_v1alpha1_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelTerm); i {
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
		file_v1alpha1_resource_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabelQuery); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1alpha1_resource_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1alpha1_resource_proto_goTypes,
		DependencyIndexes: file_v1alpha1_resource_proto_depIdxs,
		EnumInfos:         file_v1alpha1_resource_proto_enumTypes,
		MessageInfos:      file_v1alpha1_resource_proto_msgTypes,
	}.Build()
	File_v1alpha1_resource_proto = out.File
	file_v1alpha1_resource_proto_rawDesc = nil
	file_v1alpha1_resource_proto_goTypes = nil
	file_v1alpha1_resource_proto_depIdxs = nil
}
