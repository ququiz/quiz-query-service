// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: quiz_grpc.proto

package pb

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

type GetQuestionReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuestionId string `protobuf:"bytes,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
	QuizId     string `protobuf:"bytes,2,opt,name=quiz_id,json=quizId,proto3" json:"quiz_id,omitempty"`
}

func (x *GetQuestionReq) Reset() {
	*x = GetQuestionReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quiz_grpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuestionReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuestionReq) ProtoMessage() {}

func (x *GetQuestionReq) ProtoReflect() protoreflect.Message {
	mi := &file_quiz_grpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuestionReq.ProtoReflect.Descriptor instead.
func (*GetQuestionReq) Descriptor() ([]byte, []int) {
	return file_quiz_grpc_proto_rawDescGZIP(), []int{0}
}

func (x *GetQuestionReq) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *GetQuestionReq) GetQuizId() string {
	if x != nil {
		return x.QuizId
	}
	return ""
}

type GetQuestionRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrectChoiceId      string `protobuf:"bytes,1,opt,name=correct_choice_id,json=correctChoiceId,proto3" json:"correct_choice_id,omitempty"`
	CorrectEssayAnswerId string `protobuf:"bytes,2,opt,name=correct_essay_answer_id,json=correctEssayAnswerId,proto3" json:"correct_essay_answer_id,omitempty"`
	Weight               uint64 `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *GetQuestionRes) Reset() {
	*x = GetQuestionRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quiz_grpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuestionRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuestionRes) ProtoMessage() {}

func (x *GetQuestionRes) ProtoReflect() protoreflect.Message {
	mi := &file_quiz_grpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuestionRes.ProtoReflect.Descriptor instead.
func (*GetQuestionRes) Descriptor() ([]byte, []int) {
	return file_quiz_grpc_proto_rawDescGZIP(), []int{1}
}

func (x *GetQuestionRes) GetCorrectChoiceId() string {
	if x != nil {
		return x.CorrectChoiceId
	}
	return ""
}

func (x *GetQuestionRes) GetCorrectEssayAnswerId() string {
	if x != nil {
		return x.CorrectEssayAnswerId
	}
	return ""
}

func (x *GetQuestionRes) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

type GetQuizParticipantsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuizId string `protobuf:"bytes,1,opt,name=quiz_id,json=quizId,proto3" json:"quiz_id,omitempty"`
}

func (x *GetQuizParticipantsReq) Reset() {
	*x = GetQuizParticipantsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quiz_grpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuizParticipantsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuizParticipantsReq) ProtoMessage() {}

func (x *GetQuizParticipantsReq) ProtoReflect() protoreflect.Message {
	mi := &file_quiz_grpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuizParticipantsReq.ProtoReflect.Descriptor instead.
func (*GetQuizParticipantsReq) Descriptor() ([]byte, []int) {
	return file_quiz_grpc_proto_rawDescGZIP(), []int{2}
}

func (x *GetQuizParticipantsReq) GetQuizId() string {
	if x != nil {
		return x.QuizId
	}
	return ""
}

type GetQuizParticipantRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserIds  []string `protobuf:"bytes,1,rep,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
	QuizName string   `protobuf:"bytes,2,opt,name=quiz_name,json=quizName,proto3" json:"quiz_name,omitempty"`
}

func (x *GetQuizParticipantRes) Reset() {
	*x = GetQuizParticipantRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_quiz_grpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuizParticipantRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuizParticipantRes) ProtoMessage() {}

func (x *GetQuizParticipantRes) ProtoReflect() protoreflect.Message {
	mi := &file_quiz_grpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuizParticipantRes.ProtoReflect.Descriptor instead.
func (*GetQuizParticipantRes) Descriptor() ([]byte, []int) {
	return file_quiz_grpc_proto_rawDescGZIP(), []int{3}
}

func (x *GetQuizParticipantRes) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

func (x *GetQuizParticipantRes) GetQuizName() string {
	if x != nil {
		return x.QuizName
	}
	return ""
}

var File_quiz_grpc_proto protoreflect.FileDescriptor

var file_quiz_grpc_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1c, 0x71, 0x75, 0x71, 0x75, 0x69, 0x7a, 0x2e, 0x6c, 0x69, 0x6e, 0x74, 0x61, 0x6e,
	0x67, 0x2e, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x63, 0x22,
	0x4a, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x71, 0x75, 0x69, 0x7a, 0x49, 0x64, 0x22, 0x8b, 0x01, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x2a,
	0x0a, 0x11, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x63, 0x68, 0x6f, 0x69, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x72, 0x72, 0x65,
	0x63, 0x74, 0x43, 0x68, 0x6f, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x17, 0x63, 0x6f,
	0x72, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x65, 0x73, 0x73, 0x61, 0x79, 0x5f, 0x61, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x63, 0x6f, 0x72,
	0x72, 0x65, 0x63, 0x74, 0x45, 0x73, 0x73, 0x61, 0x79, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0x31, 0x0a, 0x16, 0x47, 0x65, 0x74,
	0x51, 0x75, 0x69, 0x7a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x71, 0x75, 0x69, 0x7a, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x15,
	0x47, 0x65, 0x74, 0x51, 0x75, 0x69, 0x7a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x73,
	0x12, 0x1b, 0x0a, 0x09, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x71, 0x75, 0x69, 0x7a, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0x86, 0x02,
	0x0a, 0x10, 0x51, 0x75, 0x69, 0x7a, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x6f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x2c, 0x2e, 0x71, 0x75, 0x71, 0x75, 0x69, 0x7a,
	0x2e, 0x6c, 0x69, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x2e, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x2e, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x2c, 0x2e, 0x71, 0x75, 0x71, 0x75, 0x69, 0x7a, 0x2e, 0x6c,
	0x69, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x2e, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x2e, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x12, 0x80, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x51, 0x75, 0x69, 0x7a, 0x50,
	0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x34, 0x2e, 0x71, 0x75,
	0x71, 0x75, 0x69, 0x7a, 0x2e, 0x6c, 0x69, 0x6e, 0x74, 0x61, 0x6e, 0x67, 0x2e, 0x71, 0x75, 0x69,
	0x7a, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x51, 0x75,
	0x69, 0x7a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x71, 0x1a, 0x33, 0x2e, 0x71, 0x75, 0x71, 0x75, 0x69, 0x7a, 0x2e, 0x6c, 0x69, 0x6e, 0x74, 0x61,
	0x6e, 0x67, 0x2e, 0x71, 0x75, 0x69, 0x7a, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x63,
	0x2e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x69, 0x7a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x42, 0x17, 0x5a, 0x15, 0x71, 0x75, 0x69, 0x7a, 0x2d, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_quiz_grpc_proto_rawDescOnce sync.Once
	file_quiz_grpc_proto_rawDescData = file_quiz_grpc_proto_rawDesc
)

func file_quiz_grpc_proto_rawDescGZIP() []byte {
	file_quiz_grpc_proto_rawDescOnce.Do(func() {
		file_quiz_grpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_quiz_grpc_proto_rawDescData)
	})
	return file_quiz_grpc_proto_rawDescData
}

var file_quiz_grpc_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_quiz_grpc_proto_goTypes = []interface{}{
	(*GetQuestionReq)(nil),         // 0: ququiz.lintang.quiz_query.pc.GetQuestionReq
	(*GetQuestionRes)(nil),         // 1: ququiz.lintang.quiz_query.pc.GetQuestionRes
	(*GetQuizParticipantsReq)(nil), // 2: ququiz.lintang.quiz_query.pc.GetQuizParticipantsReq
	(*GetQuizParticipantRes)(nil),  // 3: ququiz.lintang.quiz_query.pc.GetQuizParticipantRes
}
var file_quiz_grpc_proto_depIdxs = []int32{
	0, // 0: ququiz.lintang.quiz_query.pc.QuizQueryService.GetQuestionDetail:input_type -> ququiz.lintang.quiz_query.pc.GetQuestionReq
	2, // 1: ququiz.lintang.quiz_query.pc.QuizQueryService.GetQuizParticipants:input_type -> ququiz.lintang.quiz_query.pc.GetQuizParticipantsReq
	1, // 2: ququiz.lintang.quiz_query.pc.QuizQueryService.GetQuestionDetail:output_type -> ququiz.lintang.quiz_query.pc.GetQuestionRes
	3, // 3: ququiz.lintang.quiz_query.pc.QuizQueryService.GetQuizParticipants:output_type -> ququiz.lintang.quiz_query.pc.GetQuizParticipantRes
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_quiz_grpc_proto_init() }
func file_quiz_grpc_proto_init() {
	if File_quiz_grpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_quiz_grpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuestionReq); i {
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
		file_quiz_grpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuestionRes); i {
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
		file_quiz_grpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuizParticipantsReq); i {
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
		file_quiz_grpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuizParticipantRes); i {
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
			RawDescriptor: file_quiz_grpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_quiz_grpc_proto_goTypes,
		DependencyIndexes: file_quiz_grpc_proto_depIdxs,
		MessageInfos:      file_quiz_grpc_proto_msgTypes,
	}.Build()
	File_quiz_grpc_proto = out.File
	file_quiz_grpc_proto_rawDesc = nil
	file_quiz_grpc_proto_goTypes = nil
	file_quiz_grpc_proto_depIdxs = nil
}
