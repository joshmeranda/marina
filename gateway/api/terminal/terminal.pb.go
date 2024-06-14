// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.27.0
// source: gateway/api/terminal/terminal.proto

package terminal

import (
	core "github.com/joshmeranda/marina/gateway/api/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TerminalSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image string `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
}

func (x *TerminalSpec) Reset() {
	*x = TerminalSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_api_terminal_terminal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TerminalSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TerminalSpec) ProtoMessage() {}

func (x *TerminalSpec) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_api_terminal_terminal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TerminalSpec.ProtoReflect.Descriptor instead.
func (*TerminalSpec) Descriptor() ([]byte, []int) {
	return file_gateway_api_terminal_terminal_proto_rawDescGZIP(), []int{0}
}

func (x *TerminalSpec) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

type TerminalCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name *core.NamespacedName `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Spec *TerminalSpec        `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *TerminalCreateRequest) Reset() {
	*x = TerminalCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_api_terminal_terminal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TerminalCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TerminalCreateRequest) ProtoMessage() {}

func (x *TerminalCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_api_terminal_terminal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TerminalCreateRequest.ProtoReflect.Descriptor instead.
func (*TerminalCreateRequest) Descriptor() ([]byte, []int) {
	return file_gateway_api_terminal_terminal_proto_rawDescGZIP(), []int{1}
}

func (x *TerminalCreateRequest) GetName() *core.NamespacedName {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *TerminalCreateRequest) GetSpec() *TerminalSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

type TerminalCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod   *core.NamespacedName `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	Token []byte               `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Host  string               `protobuf:"bytes,3,opt,name=host,proto3" json:"host,omitempty"`
}

func (x *TerminalCreateResponse) Reset() {
	*x = TerminalCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_api_terminal_terminal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TerminalCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TerminalCreateResponse) ProtoMessage() {}

func (x *TerminalCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_api_terminal_terminal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TerminalCreateResponse.ProtoReflect.Descriptor instead.
func (*TerminalCreateResponse) Descriptor() ([]byte, []int) {
	return file_gateway_api_terminal_terminal_proto_rawDescGZIP(), []int{2}
}

func (x *TerminalCreateResponse) GetPod() *core.NamespacedName {
	if x != nil {
		return x.Pod
	}
	return nil
}

func (x *TerminalCreateResponse) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *TerminalCreateResponse) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

type TerminalDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name *core.NamespacedName `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *TerminalDeleteRequest) Reset() {
	*x = TerminalDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_api_terminal_terminal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TerminalDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TerminalDeleteRequest) ProtoMessage() {}

func (x *TerminalDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_api_terminal_terminal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TerminalDeleteRequest.ProtoReflect.Descriptor instead.
func (*TerminalDeleteRequest) Descriptor() ([]byte, []int) {
	return file_gateway_api_terminal_terminal_proto_rawDescGZIP(), []int{3}
}

func (x *TerminalDeleteRequest) GetName() *core.NamespacedName {
	if x != nil {
		return x.Name
	}
	return nil
}

var File_gateway_api_terminal_terminal_proto protoreflect.FileDescriptor

var file_gateway_api_terminal_terminal_proto_rawDesc = []byte{
	0x0a, 0x23, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x65,
	0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x24, 0x0a, 0x0c, 0x54, 0x65, 0x72,
	0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x53, 0x70, 0x65, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x22,
	0x6d, 0x0a, 0x15, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2e, 0x54, 0x65, 0x72, 0x6d,
	0x69, 0x6e, 0x61, 0x6c, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x22, 0x6a,
	0x0a, 0x16, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x03, 0x70, 0x6f, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x22, 0x41, 0x0a, 0x15, 0x54, 0x65,
	0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0xb5, 0x01,
	0x0a, 0x0f, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x55, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x65, 0x72, 0x6d, 0x69,
	0x6e, 0x61, 0x6c, 0x12, 0x1f, 0x2e, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2e, 0x54,
	0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2e,
	0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x12, 0x1f, 0x2e, 0x74, 0x65, 0x72,
	0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2e, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x73, 0x68, 0x6d, 0x65, 0x72, 0x61, 0x6e, 0x64, 0x61, 0x2f,
	0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x74, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_gateway_api_terminal_terminal_proto_rawDescOnce sync.Once
	file_gateway_api_terminal_terminal_proto_rawDescData = file_gateway_api_terminal_terminal_proto_rawDesc
)

func file_gateway_api_terminal_terminal_proto_rawDescGZIP() []byte {
	file_gateway_api_terminal_terminal_proto_rawDescOnce.Do(func() {
		file_gateway_api_terminal_terminal_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_api_terminal_terminal_proto_rawDescData)
	})
	return file_gateway_api_terminal_terminal_proto_rawDescData
}

var file_gateway_api_terminal_terminal_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gateway_api_terminal_terminal_proto_goTypes = []interface{}{
	(*TerminalSpec)(nil),           // 0: terminal.TerminalSpec
	(*TerminalCreateRequest)(nil),  // 1: terminal.TerminalCreateRequest
	(*TerminalCreateResponse)(nil), // 2: terminal.TerminalCreateResponse
	(*TerminalDeleteRequest)(nil),  // 3: terminal.TerminalDeleteRequest
	(*core.NamespacedName)(nil),    // 4: core.NamespacedName
	(*emptypb.Empty)(nil),          // 5: google.protobuf.Empty
}
var file_gateway_api_terminal_terminal_proto_depIdxs = []int32{
	4, // 0: terminal.TerminalCreateRequest.name:type_name -> core.NamespacedName
	0, // 1: terminal.TerminalCreateRequest.spec:type_name -> terminal.TerminalSpec
	4, // 2: terminal.TerminalCreateResponse.pod:type_name -> core.NamespacedName
	4, // 3: terminal.TerminalDeleteRequest.name:type_name -> core.NamespacedName
	1, // 4: terminal.TerminalService.CreateTerminal:input_type -> terminal.TerminalCreateRequest
	3, // 5: terminal.TerminalService.DeleteTerminal:input_type -> terminal.TerminalDeleteRequest
	2, // 6: terminal.TerminalService.CreateTerminal:output_type -> terminal.TerminalCreateResponse
	5, // 7: terminal.TerminalService.DeleteTerminal:output_type -> google.protobuf.Empty
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_gateway_api_terminal_terminal_proto_init() }
func file_gateway_api_terminal_terminal_proto_init() {
	if File_gateway_api_terminal_terminal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gateway_api_terminal_terminal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TerminalSpec); i {
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
		file_gateway_api_terminal_terminal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TerminalCreateRequest); i {
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
		file_gateway_api_terminal_terminal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TerminalCreateResponse); i {
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
		file_gateway_api_terminal_terminal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TerminalDeleteRequest); i {
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
			RawDescriptor: file_gateway_api_terminal_terminal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gateway_api_terminal_terminal_proto_goTypes,
		DependencyIndexes: file_gateway_api_terminal_terminal_proto_depIdxs,
		MessageInfos:      file_gateway_api_terminal_terminal_proto_msgTypes,
	}.Build()
	File_gateway_api_terminal_terminal_proto = out.File
	file_gateway_api_terminal_terminal_proto_rawDesc = nil
	file_gateway_api_terminal_terminal_proto_goTypes = nil
	file_gateway_api_terminal_terminal_proto_depIdxs = nil
}
