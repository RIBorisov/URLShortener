// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: proto/user_urls.proto

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

type URL struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortUrl    string `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	OriginalUrl string `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *URL) Reset() {
	*x = URL{}
	mi := &file_proto_user_urls_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *URL) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*URL) ProtoMessage() {}

func (x *URL) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_urls_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use URL.ProtoReflect.Descriptor instead.
func (*URL) Descriptor() ([]byte, []int) {
	return file_proto_user_urls_proto_rawDescGZIP(), []int{0}
}

func (x *URL) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *URL) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type SavedByUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SavedByUserRequest) Reset() {
	*x = SavedByUserRequest{}
	mi := &file_proto_user_urls_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SavedByUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SavedByUserRequest) ProtoMessage() {}

func (x *SavedByUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_urls_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SavedByUserRequest.ProtoReflect.Descriptor instead.
func (*SavedByUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_user_urls_proto_rawDescGZIP(), []int{1}
}

type SavedByUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*URL `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *SavedByUserResponse) Reset() {
	*x = SavedByUserResponse{}
	mi := &file_proto_user_urls_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SavedByUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SavedByUserResponse) ProtoMessage() {}

func (x *SavedByUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_user_urls_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SavedByUserResponse.ProtoReflect.Descriptor instead.
func (*SavedByUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_user_urls_proto_rawDescGZIP(), []int{2}
}

func (x *SavedByUserResponse) GetUrls() []*URL {
	if x != nil {
		return x.Urls
	}
	return nil
}

var File_proto_user_urls_proto protoreflect.FileDescriptor

var file_proto_user_urls_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x75, 0x72, 0x6c,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x12, 0x1b,
	0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c, 0x22, 0x14,
	0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x64, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x2f, 0x0a, 0x13, 0x53, 0x61, 0x76, 0x65, 0x64, 0x42, 0x79, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x04, 0x75,
	0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x55, 0x52, 0x4c, 0x52,
	0x04, 0x75, 0x72, 0x6c, 0x73, 0x42, 0x1d, 0x5a, 0x1b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_user_urls_proto_rawDescOnce sync.Once
	file_proto_user_urls_proto_rawDescData = file_proto_user_urls_proto_rawDesc
)

func file_proto_user_urls_proto_rawDescGZIP() []byte {
	file_proto_user_urls_proto_rawDescOnce.Do(func() {
		file_proto_user_urls_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_user_urls_proto_rawDescData)
	})
	return file_proto_user_urls_proto_rawDescData
}

var file_proto_user_urls_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_user_urls_proto_goTypes = []any{
	(*URL)(nil),                 // 0: URL
	(*SavedByUserRequest)(nil),  // 1: SavedByUserRequest
	(*SavedByUserResponse)(nil), // 2: SavedByUserResponse
}
var file_proto_user_urls_proto_depIdxs = []int32{
	0, // 0: SavedByUserResponse.urls:type_name -> URL
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_user_urls_proto_init() }
func file_proto_user_urls_proto_init() {
	if File_proto_user_urls_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_user_urls_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_user_urls_proto_goTypes,
		DependencyIndexes: file_proto_user_urls_proto_depIdxs,
		MessageInfos:      file_proto_user_urls_proto_msgTypes,
	}.Build()
	File_proto_user_urls_proto = out.File
	file_proto_user_urls_proto_rawDesc = nil
	file_proto_user_urls_proto_goTypes = nil
	file_proto_user_urls_proto_depIdxs = nil
}