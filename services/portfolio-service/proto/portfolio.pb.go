// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.1
// source: portfolio.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetPortfolioRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPortfolioRequest) Reset() {
	*x = GetPortfolioRequest{}
	mi := &file_portfolio_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPortfolioRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPortfolioRequest) ProtoMessage() {}

func (x *GetPortfolioRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPortfolioRequest.ProtoReflect.Descriptor instead.
func (*GetPortfolioRequest) Descriptor() ([]byte, []int) {
	return file_portfolio_proto_rawDescGZIP(), []int{0}
}

func (x *GetPortfolioRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetPortfolioResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Holdings      []*Holding             `protobuf:"bytes,1,rep,name=holdings,proto3" json:"holdings,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPortfolioResponse) Reset() {
	*x = GetPortfolioResponse{}
	mi := &file_portfolio_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPortfolioResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPortfolioResponse) ProtoMessage() {}

func (x *GetPortfolioResponse) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPortfolioResponse.ProtoReflect.Descriptor instead.
func (*GetPortfolioResponse) Descriptor() ([]byte, []int) {
	return file_portfolio_proto_rawDescGZIP(), []int{1}
}

func (x *GetPortfolioResponse) GetHoldings() []*Holding {
	if x != nil {
		return x.Holdings
	}
	return nil
}

type Holding struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Symbol        string                 `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Quantity      float64                `protobuf:"fixed64,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	AveragePrice  float64                `protobuf:"fixed64,3,opt,name=average_price,json=averagePrice,proto3" json:"average_price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Holding) Reset() {
	*x = Holding{}
	mi := &file_portfolio_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Holding) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Holding) ProtoMessage() {}

func (x *Holding) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Holding.ProtoReflect.Descriptor instead.
func (*Holding) Descriptor() ([]byte, []int) {
	return file_portfolio_proto_rawDescGZIP(), []int{2}
}

func (x *Holding) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Holding) GetQuantity() float64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *Holding) GetAveragePrice() float64 {
	if x != nil {
		return x.AveragePrice
	}
	return 0
}

type UpdateHoldingsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Symbol        string                 `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Quantity      float64                `protobuf:"fixed64,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateHoldingsRequest) Reset() {
	*x = UpdateHoldingsRequest{}
	mi := &file_portfolio_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateHoldingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateHoldingsRequest) ProtoMessage() {}

func (x *UpdateHoldingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateHoldingsRequest.ProtoReflect.Descriptor instead.
func (*UpdateHoldingsRequest) Descriptor() ([]byte, []int) {
	return file_portfolio_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateHoldingsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdateHoldingsRequest) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *UpdateHoldingsRequest) GetQuantity() float64 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *UpdateHoldingsRequest) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type UpdateHoldingsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateHoldingsResponse) Reset() {
	*x = UpdateHoldingsResponse{}
	mi := &file_portfolio_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateHoldingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateHoldingsResponse) ProtoMessage() {}

func (x *UpdateHoldingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_portfolio_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateHoldingsResponse.ProtoReflect.Descriptor instead.
func (*UpdateHoldingsResponse) Descriptor() ([]byte, []int) {
	return file_portfolio_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateHoldingsResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_portfolio_proto protoreflect.FileDescriptor

var file_portfolio_proto_rawDesc = string([]byte{
	0x0a, 0x0f, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x22, 0x2e, 0x0a, 0x13,
	0x47, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x46, 0x0a, 0x14,
	0x47, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x08, 0x68, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c,
	0x69, 0x6f, 0x2e, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x08, 0x68, 0x6f, 0x6c, 0x64,
	0x69, 0x6e, 0x67, 0x73, 0x22, 0x62, 0x0a, 0x07, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x61, 0x76, 0x65, 0x72,
	0x61, 0x67, 0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x22, 0x7a, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x22, 0x32, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x6f,
	0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0xba, 0x01, 0x0a, 0x10, 0x50, 0x6f, 0x72,
	0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4f, 0x0a,
	0x0c, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x12, 0x1e, 0x2e,
	0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72,
	0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x72,
	0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55,
	0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x20, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x45, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x6b, 0x61, 0x6e, 0x38, 0x2f, 0x73, 0x77, 0x61, 0x70, 0x73,
	0x79, 0x6e, 0x63, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x6f, 0x72, 0x74, 0x66, 0x6f, 0x6c, 0x69, 0x6f, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_portfolio_proto_rawDescOnce sync.Once
	file_portfolio_proto_rawDescData []byte
)

func file_portfolio_proto_rawDescGZIP() []byte {
	file_portfolio_proto_rawDescOnce.Do(func() {
		file_portfolio_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_portfolio_proto_rawDesc), len(file_portfolio_proto_rawDesc)))
	})
	return file_portfolio_proto_rawDescData
}

var file_portfolio_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_portfolio_proto_goTypes = []any{
	(*GetPortfolioRequest)(nil),    // 0: portfolio.GetPortfolioRequest
	(*GetPortfolioResponse)(nil),   // 1: portfolio.GetPortfolioResponse
	(*Holding)(nil),                // 2: portfolio.Holding
	(*UpdateHoldingsRequest)(nil),  // 3: portfolio.UpdateHoldingsRequest
	(*UpdateHoldingsResponse)(nil), // 4: portfolio.UpdateHoldingsResponse
}
var file_portfolio_proto_depIdxs = []int32{
	2, // 0: portfolio.GetPortfolioResponse.holdings:type_name -> portfolio.Holding
	0, // 1: portfolio.PortfolioService.GetPortfolio:input_type -> portfolio.GetPortfolioRequest
	3, // 2: portfolio.PortfolioService.UpdateHoldings:input_type -> portfolio.UpdateHoldingsRequest
	1, // 3: portfolio.PortfolioService.GetPortfolio:output_type -> portfolio.GetPortfolioResponse
	4, // 4: portfolio.PortfolioService.UpdateHoldings:output_type -> portfolio.UpdateHoldingsResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_portfolio_proto_init() }
func file_portfolio_proto_init() {
	if File_portfolio_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_portfolio_proto_rawDesc), len(file_portfolio_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_portfolio_proto_goTypes,
		DependencyIndexes: file_portfolio_proto_depIdxs,
		MessageInfos:      file_portfolio_proto_msgTypes,
	}.Build()
	File_portfolio_proto = out.File
	file_portfolio_proto_goTypes = nil
	file_portfolio_proto_depIdxs = nil
}
