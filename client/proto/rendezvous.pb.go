// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.18.0
// source: proto/rendezvous.proto

package proto

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

type WorkspaceInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mnemonic string `protobuf:"bytes,1,opt,name=mnemonic,proto3" json:"mnemonic,omitempty"`
}

func (x *WorkspaceInfoRequest) Reset() {
	*x = WorkspaceInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rendezvous_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkspaceInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkspaceInfoRequest) ProtoMessage() {}

func (x *WorkspaceInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rendezvous_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkspaceInfoRequest.ProtoReflect.Descriptor instead.
func (*WorkspaceInfoRequest) Descriptor() ([]byte, []int) {
	return file_proto_rendezvous_proto_rawDescGZIP(), []int{0}
}

func (x *WorkspaceInfoRequest) GetMnemonic() string {
	if x != nil {
		return x.Mnemonic
	}
	return ""
}

// Exchanged between Rendezvous nodes,
// as well as between Client and Rendezvous nodes
type WorkspaceInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Mnemonic        string            `protobuf:"bytes,2,opt,name=mnemonic,proto3" json:"mnemonic,omitempty"`
	WorkspaceOwners []*WorkspaceOwner `protobuf:"bytes,3,rep,name=workspace_owners,json=workspaceOwners,proto3" json:"workspace_owners,omitempty"`
	SecurityType    string            `protobuf:"bytes,4,opt,name=security_type,json=securityType,proto3" json:"security_type,omitempty"`
	// Types that are assignable to SecuritySettings:
	//	*WorkspaceInfo_ContactsWrapper
	//	*WorkspaceInfo_PasswordHash
	SecuritySettings isWorkspaceInfo_SecuritySettings `protobuf_oneof:"security_settings"`
}

func (x *WorkspaceInfo) Reset() {
	*x = WorkspaceInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rendezvous_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkspaceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkspaceInfo) ProtoMessage() {}

func (x *WorkspaceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rendezvous_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkspaceInfo.ProtoReflect.Descriptor instead.
func (*WorkspaceInfo) Descriptor() ([]byte, []int) {
	return file_proto_rendezvous_proto_rawDescGZIP(), []int{1}
}

func (x *WorkspaceInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *WorkspaceInfo) GetMnemonic() string {
	if x != nil {
		return x.Mnemonic
	}
	return ""
}

func (x *WorkspaceInfo) GetWorkspaceOwners() []*WorkspaceOwner {
	if x != nil {
		return x.WorkspaceOwners
	}
	return nil
}

func (x *WorkspaceInfo) GetSecurityType() string {
	if x != nil {
		return x.SecurityType
	}
	return ""
}

func (m *WorkspaceInfo) GetSecuritySettings() isWorkspaceInfo_SecuritySettings {
	if m != nil {
		return m.SecuritySettings
	}
	return nil
}

func (x *WorkspaceInfo) GetContactsWrapper() *ContactsWrapper {
	if x, ok := x.GetSecuritySettings().(*WorkspaceInfo_ContactsWrapper); ok {
		return x.ContactsWrapper
	}
	return nil
}

func (x *WorkspaceInfo) GetPasswordHash() string {
	if x, ok := x.GetSecuritySettings().(*WorkspaceInfo_PasswordHash); ok {
		return x.PasswordHash
	}
	return ""
}

type isWorkspaceInfo_SecuritySettings interface {
	isWorkspaceInfo_SecuritySettings()
}

type WorkspaceInfo_ContactsWrapper struct {
	ContactsWrapper *ContactsWrapper `protobuf:"bytes,5,opt,name=contacts_wrapper,json=contactsWrapper,proto3,oneof"`
}

type WorkspaceInfo_PasswordHash struct {
	PasswordHash string `protobuf:"bytes,6,opt,name=password_hash,json=passwordHash,proto3,oneof"`
}

func (*WorkspaceInfo_ContactsWrapper) isWorkspaceInfo_SecuritySettings() {}

func (*WorkspaceInfo_PasswordHash) isWorkspaceInfo_SecuritySettings() {}

// Defines a workspace owner
type WorkspaceOwner struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PublicKey     string `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Libp2PAddress string `protobuf:"bytes,2,opt,name=libp2p_address,json=libp2pAddress,proto3" json:"libp2p_address,omitempty"`
}

func (x *WorkspaceOwner) Reset() {
	*x = WorkspaceOwner{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rendezvous_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkspaceOwner) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkspaceOwner) ProtoMessage() {}

func (x *WorkspaceOwner) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rendezvous_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkspaceOwner.ProtoReflect.Descriptor instead.
func (*WorkspaceOwner) Descriptor() ([]byte, []int) {
	return file_proto_rendezvous_proto_rawDescGZIP(), []int{2}
}

func (x *WorkspaceOwner) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *WorkspaceOwner) GetLibp2PAddress() string {
	if x != nil {
		return x.Libp2PAddress
	}
	return ""
}

// ContactsWrapper is a wrapper object for
// the contact public keys array
type ContactsWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContactPublicKeys []string `protobuf:"bytes,1,rep,name=contact_public_keys,json=contactPublicKeys,proto3" json:"contact_public_keys,omitempty"`
}

func (x *ContactsWrapper) Reset() {
	*x = ContactsWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_rendezvous_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContactsWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContactsWrapper) ProtoMessage() {}

func (x *ContactsWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_proto_rendezvous_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContactsWrapper.ProtoReflect.Descriptor instead.
func (*ContactsWrapper) Descriptor() ([]byte, []int) {
	return file_proto_rendezvous_proto_rawDescGZIP(), []int{3}
}

func (x *ContactsWrapper) GetContactPublicKeys() []string {
	if x != nil {
		return x.ContactPublicKeys
	}
	return nil
}

var File_proto_rendezvous_proto protoreflect.FileDescriptor

var file_proto_rendezvous_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x7a, 0x76, 0x6f,
	0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x14, 0x57, 0x6f, 0x72, 0x6b,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x22, 0x9b, 0x02, 0x0a,
	0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x12, 0x3a,
	0x0a, 0x10, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x52, 0x0f, 0x77, 0x6f, 0x72, 0x6b, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65,
	0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x3d, 0x0a, 0x10, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x5f, 0x77, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x73, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x48, 0x00, 0x52, 0x0f, 0x63,
	0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x25,
	0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x48, 0x61, 0x73, 0x68, 0x42, 0x13, 0x0a, 0x11, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74,
	0x79, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x56, 0x0a, 0x0e, 0x57, 0x6f,
	0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x6c,
	0x69, 0x62, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x69, 0x62, 0x70, 0x32, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x22, 0x41, 0x0a, 0x0f, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x57, 0x72,
	0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74,
	0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x11, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x4b, 0x65, 0x79, 0x73, 0x32, 0x51, 0x0a, 0x14, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x15, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x08, 0x5a, 0x06, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_rendezvous_proto_rawDescOnce sync.Once
	file_proto_rendezvous_proto_rawDescData = file_proto_rendezvous_proto_rawDesc
)

func file_proto_rendezvous_proto_rawDescGZIP() []byte {
	file_proto_rendezvous_proto_rawDescOnce.Do(func() {
		file_proto_rendezvous_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_rendezvous_proto_rawDescData)
	})
	return file_proto_rendezvous_proto_rawDescData
}

var file_proto_rendezvous_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_rendezvous_proto_goTypes = []interface{}{
	(*WorkspaceInfoRequest)(nil), // 0: WorkspaceInfoRequest
	(*WorkspaceInfo)(nil),        // 1: WorkspaceInfo
	(*WorkspaceOwner)(nil),       // 2: WorkspaceOwner
	(*ContactsWrapper)(nil),      // 3: ContactsWrapper
}
var file_proto_rendezvous_proto_depIdxs = []int32{
	2, // 0: WorkspaceInfo.workspace_owners:type_name -> WorkspaceOwner
	3, // 1: WorkspaceInfo.contacts_wrapper:type_name -> ContactsWrapper
	0, // 2: WorkspaceInfoService.GetWorkspaceInfo:input_type -> WorkspaceInfoRequest
	1, // 3: WorkspaceInfoService.GetWorkspaceInfo:output_type -> WorkspaceInfo
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_rendezvous_proto_init() }
func file_proto_rendezvous_proto_init() {
	if File_proto_rendezvous_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_rendezvous_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkspaceInfoRequest); i {
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
		file_proto_rendezvous_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkspaceInfo); i {
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
		file_proto_rendezvous_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkspaceOwner); i {
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
		file_proto_rendezvous_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContactsWrapper); i {
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
	file_proto_rendezvous_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*WorkspaceInfo_ContactsWrapper)(nil),
		(*WorkspaceInfo_PasswordHash)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_rendezvous_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_rendezvous_proto_goTypes,
		DependencyIndexes: file_proto_rendezvous_proto_depIdxs,
		MessageInfos:      file_proto_rendezvous_proto_msgTypes,
	}.Build()
	File_proto_rendezvous_proto = out.File
	file_proto_rendezvous_proto_rawDesc = nil
	file_proto_rendezvous_proto_goTypes = nil
	file_proto_rendezvous_proto_depIdxs = nil
}
