// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.27.0
// source: gateway/api/auth/auth.proto

package auth

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

type SecretType int32

const (
	SecretType_Unknown  SecretType = 0
	SecretType_Password SecretType = 1
	SecretType_Github   SecretType = 2
)

// Enum value maps for SecretType.
var (
	SecretType_name = map[int32]string{
		0: "Unknown",
		1: "Password",
		2: "Github",
	}
	SecretType_value = map[string]int32{
		"Unknown":  0,
		"Password": 1,
		"Github":   2,
	}
)

func (x SecretType) Enum() *SecretType {
	p := new(SecretType)
	*p = x
	return p
}

func (x SecretType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SecretType) Descriptor() protoreflect.EnumDescriptor {
	return file_gateway_api_auth_auth_proto_enumTypes[0].Descriptor()
}

func (SecretType) Type() protoreflect.EnumType {
	return &file_gateway_api_auth_auth_proto_enumTypes[0]
}

func (x SecretType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SecretType.Descriptor instead.
func (SecretType) EnumDescriptor() ([]byte, []int) {
	return file_gateway_api_auth_auth_proto_rawDescGZIP(), []int{0}
}

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Secret     []byte     `protobuf:"bytes,1,opt,name=Secret,proto3" json:"Secret,omitempty"`
	SecretType SecretType `protobuf:"varint,2,opt,name=secretType,proto3,enum=auth.SecretType" json:"secretType,omitempty"`
	User       string     `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_api_auth_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_api_auth_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_gateway_api_auth_auth_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetSecret() []byte {
	if x != nil {
		return x.Secret
	}
	return nil
}

func (x *LoginRequest) GetSecretType() SecretType {
	if x != nil {
		return x.SecretType
	}
	return SecretType_Unknown
}

func (x *LoginRequest) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gateway_api_auth_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gateway_api_auth_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_gateway_api_auth_auth_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

var File_gateway_api_auth_auth_proto protoreflect.FileDescriptor

var file_gateway_api_auth_auth_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x61,
	0x75, 0x74, 0x68, 0x22, 0x6c, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x30, 0x0a, 0x0a, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x10, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x0a, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x22, 0x25, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2a, 0x33, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77,
	0x6e, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x10, 0x02, 0x32, 0x41, 0x0a,
	0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x05,
	0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a,
	0x6f, 0x73, 0x68, 0x6d, 0x65, 0x72, 0x61, 0x6e, 0x64, 0x61, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e,
	0x61, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gateway_api_auth_auth_proto_rawDescOnce sync.Once
	file_gateway_api_auth_auth_proto_rawDescData = file_gateway_api_auth_auth_proto_rawDesc
)

func file_gateway_api_auth_auth_proto_rawDescGZIP() []byte {
	file_gateway_api_auth_auth_proto_rawDescOnce.Do(func() {
		file_gateway_api_auth_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_gateway_api_auth_auth_proto_rawDescData)
	})
	return file_gateway_api_auth_auth_proto_rawDescData
}

var file_gateway_api_auth_auth_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_gateway_api_auth_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gateway_api_auth_auth_proto_goTypes = []interface{}{
	(SecretType)(0),       // 0: auth.SecretType
	(*LoginRequest)(nil),  // 1: auth.LoginRequest
	(*LoginResponse)(nil), // 2: auth.LoginResponse
}
var file_gateway_api_auth_auth_proto_depIdxs = []int32{
	0, // 0: auth.LoginRequest.secretType:type_name -> auth.SecretType
	1, // 1: auth.AuthService.Login:input_type -> auth.LoginRequest
	2, // 2: auth.AuthService.Login:output_type -> auth.LoginResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_gateway_api_auth_auth_proto_init() }
func file_gateway_api_auth_auth_proto_init() {
	if File_gateway_api_auth_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gateway_api_auth_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRequest); i {
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
		file_gateway_api_auth_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginResponse); i {
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
			RawDescriptor: file_gateway_api_auth_auth_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gateway_api_auth_auth_proto_goTypes,
		DependencyIndexes: file_gateway_api_auth_auth_proto_depIdxs,
		EnumInfos:         file_gateway_api_auth_auth_proto_enumTypes,
		MessageInfos:      file_gateway_api_auth_auth_proto_msgTypes,
	}.Build()
	File_gateway_api_auth_auth_proto = out.File
	file_gateway_api_auth_auth_proto_rawDesc = nil
	file_gateway_api_auth_auth_proto_goTypes = nil
	file_gateway_api_auth_auth_proto_depIdxs = nil
}
