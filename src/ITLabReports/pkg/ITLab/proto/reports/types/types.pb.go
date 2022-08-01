// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.19.4
// source: rtuitlab/itlab/reports/types.proto

package types

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Assignees struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reporter    string `protobuf:"bytes,1,opt,name=reporter,proto3" json:"reporter,omitempty"`
	Implementer string `protobuf:"bytes,2,opt,name=implementer,proto3" json:"implementer,omitempty"`
}

func (x *Assignees) Reset() {
	*x = Assignees{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rtuitlab_itlab_reports_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Assignees) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Assignees) ProtoMessage() {}

func (x *Assignees) ProtoReflect() protoreflect.Message {
	mi := &file_rtuitlab_itlab_reports_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Assignees.ProtoReflect.Descriptor instead.
func (*Assignees) Descriptor() ([]byte, []int) {
	return file_rtuitlab_itlab_reports_types_proto_rawDescGZIP(), []int{0}
}

func (x *Assignees) GetReporter() string {
	if x != nil {
		return x.Reporter
	}
	return ""
}

func (x *Assignees) GetImplementer() string {
	if x != nil {
		return x.Implementer
	}
	return ""
}

type Report struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Text      string                 `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
	Assignees *Assignees             `protobuf:"bytes,4,opt,name=assignees,proto3" json:"assignees,omitempty"`
	Date      *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=date,proto3" json:"date,omitempty"`
}

func (x *Report) Reset() {
	*x = Report{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rtuitlab_itlab_reports_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Report) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Report) ProtoMessage() {}

func (x *Report) ProtoReflect() protoreflect.Message {
	mi := &file_rtuitlab_itlab_reports_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Report.ProtoReflect.Descriptor instead.
func (*Report) Descriptor() ([]byte, []int) {
	return file_rtuitlab_itlab_reports_types_proto_rawDescGZIP(), []int{1}
}

func (x *Report) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Report) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Report) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Report) GetAssignees() *Assignees {
	if x != nil {
		return x.Assignees
	}
	return nil
}

func (x *Report) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

var File_rtuitlab_itlab_reports_types_proto protoreflect.FileDescriptor

var file_rtuitlab_itlab_reports_types_proto_rawDesc = []byte{
	0x0a, 0x22, 0x72, 0x74, 0x75, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2f, 0x69, 0x74, 0x6c, 0x61, 0x62,
	0x2f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x72, 0x74, 0x75, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a,
	0x09, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x6d, 0x70,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x22, 0xb1, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x3f, 0x0a, 0x09, 0x61,
	0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x72, 0x74, 0x75, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x2e,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65,
	0x73, 0x52, 0x09, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x42, 0x5a, 0x5a, 0x33,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x52, 0x54, 0x55, 0x49, 0x54,
	0x4c, 0x61, 0x62, 0x2f, 0x49, 0x54, 0x4c, 0x61, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x3b, 0x74, 0x79,
	0x70, 0x65, 0x73, 0xaa, 0x02, 0x22, 0x52, 0x54, 0x55, 0x49, 0x54, 0x4c, 0x61, 0x62, 0x2e, 0x49,
	0x54, 0x4c, 0x61, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x73, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rtuitlab_itlab_reports_types_proto_rawDescOnce sync.Once
	file_rtuitlab_itlab_reports_types_proto_rawDescData = file_rtuitlab_itlab_reports_types_proto_rawDesc
)

func file_rtuitlab_itlab_reports_types_proto_rawDescGZIP() []byte {
	file_rtuitlab_itlab_reports_types_proto_rawDescOnce.Do(func() {
		file_rtuitlab_itlab_reports_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_rtuitlab_itlab_reports_types_proto_rawDescData)
	})
	return file_rtuitlab_itlab_reports_types_proto_rawDescData
}

var file_rtuitlab_itlab_reports_types_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rtuitlab_itlab_reports_types_proto_goTypes = []interface{}{
	(*Assignees)(nil),             // 0: rtuitlab.itlab.reports.Assignees
	(*Report)(nil),                // 1: rtuitlab.itlab.reports.Report
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_rtuitlab_itlab_reports_types_proto_depIdxs = []int32{
	0, // 0: rtuitlab.itlab.reports.Report.assignees:type_name -> rtuitlab.itlab.reports.Assignees
	2, // 1: rtuitlab.itlab.reports.Report.date:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_rtuitlab_itlab_reports_types_proto_init() }
func file_rtuitlab_itlab_reports_types_proto_init() {
	if File_rtuitlab_itlab_reports_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rtuitlab_itlab_reports_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Assignees); i {
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
		file_rtuitlab_itlab_reports_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Report); i {
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
			RawDescriptor: file_rtuitlab_itlab_reports_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rtuitlab_itlab_reports_types_proto_goTypes,
		DependencyIndexes: file_rtuitlab_itlab_reports_types_proto_depIdxs,
		MessageInfos:      file_rtuitlab_itlab_reports_types_proto_msgTypes,
	}.Build()
	File_rtuitlab_itlab_reports_types_proto = out.File
	file_rtuitlab_itlab_reports_types_proto_rawDesc = nil
	file_rtuitlab_itlab_reports_types_proto_goTypes = nil
	file_rtuitlab_itlab_reports_types_proto_depIdxs = nil
}