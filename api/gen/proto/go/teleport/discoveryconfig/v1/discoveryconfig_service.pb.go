// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: teleport/discoveryconfig/v1/discoveryconfig_service.proto

package discoveryconfigv1

import (
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

// ListDiscoveryConfigsRequest is a request for a paginated list of DiscoveryConfigs.
type ListDiscoveryConfigsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// page_size is the size of the page to request.
	PageSize int32 `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// next_token is the page token.
	NextToken string `protobuf:"bytes,2,opt,name=next_token,json=nextToken,proto3" json:"next_token,omitempty"`
}

func (x *ListDiscoveryConfigsRequest) Reset() {
	*x = ListDiscoveryConfigsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDiscoveryConfigsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDiscoveryConfigsRequest) ProtoMessage() {}

func (x *ListDiscoveryConfigsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDiscoveryConfigsRequest.ProtoReflect.Descriptor instead.
func (*ListDiscoveryConfigsRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{0}
}

func (x *ListDiscoveryConfigsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListDiscoveryConfigsRequest) GetNextToken() string {
	if x != nil {
		return x.NextToken
	}
	return ""
}

// ListDiscoveryConfigsResponse is the response for ListDiscoveryConfigsRequest.
type ListDiscoveryConfigsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DiscoveryConfigs is a list of DiscoveryConfigs.
	DiscoveryConfigs []*DiscoveryConfig `protobuf:"bytes,1,rep,name=discovery_configs,json=discoveryConfigs,proto3" json:"discovery_configs,omitempty"`
	// NextKey is the key for the next page of DiscoveryConfigs.
	NextKey string `protobuf:"bytes,2,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
	// TotalCount is the total number of discovery_config in all pages.
	TotalCount int32 `protobuf:"varint,3,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
}

func (x *ListDiscoveryConfigsResponse) Reset() {
	*x = ListDiscoveryConfigsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListDiscoveryConfigsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDiscoveryConfigsResponse) ProtoMessage() {}

func (x *ListDiscoveryConfigsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDiscoveryConfigsResponse.ProtoReflect.Descriptor instead.
func (*ListDiscoveryConfigsResponse) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListDiscoveryConfigsResponse) GetDiscoveryConfigs() []*DiscoveryConfig {
	if x != nil {
		return x.DiscoveryConfigs
	}
	return nil
}

func (x *ListDiscoveryConfigsResponse) GetNextKey() string {
	if x != nil {
		return x.NextKey
	}
	return ""
}

func (x *ListDiscoveryConfigsResponse) GetTotalCount() int32 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

// GetDiscoveryConfigRequest is a request for a specific DiscoveryConfig resource.
type GetDiscoveryConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name is the name of the DiscoveryConfig to be requested.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetDiscoveryConfigRequest) Reset() {
	*x = GetDiscoveryConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDiscoveryConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDiscoveryConfigRequest) ProtoMessage() {}

func (x *GetDiscoveryConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDiscoveryConfigRequest.ProtoReflect.Descriptor instead.
func (*GetDiscoveryConfigRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetDiscoveryConfigRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// CreateDiscoveryConfigRequest is the request to create the provided DiscoveryConfig.
type CreateDiscoveryConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DiscoveryConfig is the DiscoveryConfig to be created.
	DiscoveryConfig *DiscoveryConfig `protobuf:"bytes,1,opt,name=discovery_config,json=discoveryConfig,proto3" json:"discovery_config,omitempty"`
}

func (x *CreateDiscoveryConfigRequest) Reset() {
	*x = CreateDiscoveryConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateDiscoveryConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDiscoveryConfigRequest) ProtoMessage() {}

func (x *CreateDiscoveryConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDiscoveryConfigRequest.ProtoReflect.Descriptor instead.
func (*CreateDiscoveryConfigRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{3}
}

func (x *CreateDiscoveryConfigRequest) GetDiscoveryConfig() *DiscoveryConfig {
	if x != nil {
		return x.DiscoveryConfig
	}
	return nil
}

// UpdateDiscoveryConfigRequest is the request to update the provided DiscoveryConfig.
type UpdateDiscoveryConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DiscoveryConfig is the DiscoveryConfig to be updated.
	DiscoveryConfig *DiscoveryConfig `protobuf:"bytes,1,opt,name=discovery_config,json=discoveryConfig,proto3" json:"discovery_config,omitempty"`
}

func (x *UpdateDiscoveryConfigRequest) Reset() {
	*x = UpdateDiscoveryConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateDiscoveryConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDiscoveryConfigRequest) ProtoMessage() {}

func (x *UpdateDiscoveryConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDiscoveryConfigRequest.ProtoReflect.Descriptor instead.
func (*UpdateDiscoveryConfigRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateDiscoveryConfigRequest) GetDiscoveryConfig() *DiscoveryConfig {
	if x != nil {
		return x.DiscoveryConfig
	}
	return nil
}

// UpsertDiscoveryConfigRequest is the request to upsert the provided DiscoveryConfig.
type UpsertDiscoveryConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// DiscoveryConfig is the DiscoveryConfig to be upserted.
	DiscoveryConfig *DiscoveryConfig `protobuf:"bytes,1,opt,name=discovery_config,json=discoveryConfig,proto3" json:"discovery_config,omitempty"`
}

func (x *UpsertDiscoveryConfigRequest) Reset() {
	*x = UpsertDiscoveryConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertDiscoveryConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertDiscoveryConfigRequest) ProtoMessage() {}

func (x *UpsertDiscoveryConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertDiscoveryConfigRequest.ProtoReflect.Descriptor instead.
func (*UpsertDiscoveryConfigRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{5}
}

func (x *UpsertDiscoveryConfigRequest) GetDiscoveryConfig() *DiscoveryConfig {
	if x != nil {
		return x.DiscoveryConfig
	}
	return nil
}

// DeleteDiscoveryConfigRequest is a request for deleting a specific DiscoveryConfig resource.
type DeleteDiscoveryConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name is the name of the DiscoveryConfig to be deleted.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeleteDiscoveryConfigRequest) Reset() {
	*x = DeleteDiscoveryConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteDiscoveryConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteDiscoveryConfigRequest) ProtoMessage() {}

func (x *DeleteDiscoveryConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteDiscoveryConfigRequest.ProtoReflect.Descriptor instead.
func (*DeleteDiscoveryConfigRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteDiscoveryConfigRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// DeleteAllDiscoveryConfigsRequest is the request for deleting all DiscoveryConfigs.
type DeleteAllDiscoveryConfigsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteAllDiscoveryConfigsRequest) Reset() {
	*x = DeleteAllDiscoveryConfigsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAllDiscoveryConfigsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAllDiscoveryConfigsRequest) ProtoMessage() {}

func (x *DeleteAllDiscoveryConfigsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAllDiscoveryConfigsRequest.ProtoReflect.Descriptor instead.
func (*DeleteAllDiscoveryConfigsRequest) Descriptor() ([]byte, []int) {
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP(), []int{7}
}

var File_teleport_discoveryconfig_v1_discoveryconfig_service_proto protoreflect.FileDescriptor

var file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDesc = []byte{
	0x0a, 0x39, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x74, 0x65, 0x6c,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x31, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f,
	0x76, 0x31, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x1b, 0x4c, 0x69, 0x73, 0x74,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65,
	0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0xb5, 0x01, 0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x11, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x10, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12,
	0x19, 0x0a, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x2f, 0x0a, 0x19, 0x47,
	0x65, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x77, 0x0a, 0x1c,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x57, 0x0a, 0x10,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x0f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x77, 0x0a, 0x1c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x57, 0x0a, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0f, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x77,
	0x0a, 0x1c, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x57,
	0x0a, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x32, 0x0a, 0x1c, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x22, 0x0a, 0x20, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x32,
	0x8b, 0x07, 0x0a, 0x16, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x8b, 0x01, 0x0a, 0x14, 0x4c,
	0x69, 0x73, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x73, 0x12, 0x38, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e,
	0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x7a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x36,
	0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x80, 0x01, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x39,
	0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x80, 0x01, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x39, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x74,
	0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x80, 0x01, 0x0a, 0x15, 0x55,
	0x70, 0x73, 0x65, 0x72, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x12, 0x39, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x6a, 0x0a,
	0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x39, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x72, 0x0a, 0x19, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x12, 0x3d, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x62, 0x5a,
	0x60, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x61, 0x76,
	0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x76, 0x31, 0x3b,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x76,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescOnce sync.Once
	file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescData = file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDesc
)

func file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescGZIP() []byte {
	file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescOnce.Do(func() {
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescData)
	})
	return file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDescData
}

var file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_goTypes = []interface{}{
	(*ListDiscoveryConfigsRequest)(nil),      // 0: teleport.discoveryconfig.v1.ListDiscoveryConfigsRequest
	(*ListDiscoveryConfigsResponse)(nil),     // 1: teleport.discoveryconfig.v1.ListDiscoveryConfigsResponse
	(*GetDiscoveryConfigRequest)(nil),        // 2: teleport.discoveryconfig.v1.GetDiscoveryConfigRequest
	(*CreateDiscoveryConfigRequest)(nil),     // 3: teleport.discoveryconfig.v1.CreateDiscoveryConfigRequest
	(*UpdateDiscoveryConfigRequest)(nil),     // 4: teleport.discoveryconfig.v1.UpdateDiscoveryConfigRequest
	(*UpsertDiscoveryConfigRequest)(nil),     // 5: teleport.discoveryconfig.v1.UpsertDiscoveryConfigRequest
	(*DeleteDiscoveryConfigRequest)(nil),     // 6: teleport.discoveryconfig.v1.DeleteDiscoveryConfigRequest
	(*DeleteAllDiscoveryConfigsRequest)(nil), // 7: teleport.discoveryconfig.v1.DeleteAllDiscoveryConfigsRequest
	(*DiscoveryConfig)(nil),                  // 8: teleport.discoveryconfig.v1.DiscoveryConfig
	(*emptypb.Empty)(nil),                    // 9: google.protobuf.Empty
}
var file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_depIdxs = []int32{
	8,  // 0: teleport.discoveryconfig.v1.ListDiscoveryConfigsResponse.discovery_configs:type_name -> teleport.discoveryconfig.v1.DiscoveryConfig
	8,  // 1: teleport.discoveryconfig.v1.CreateDiscoveryConfigRequest.discovery_config:type_name -> teleport.discoveryconfig.v1.DiscoveryConfig
	8,  // 2: teleport.discoveryconfig.v1.UpdateDiscoveryConfigRequest.discovery_config:type_name -> teleport.discoveryconfig.v1.DiscoveryConfig
	8,  // 3: teleport.discoveryconfig.v1.UpsertDiscoveryConfigRequest.discovery_config:type_name -> teleport.discoveryconfig.v1.DiscoveryConfig
	0,  // 4: teleport.discoveryconfig.v1.DiscoveryConfigService.ListDiscoveryConfigs:input_type -> teleport.discoveryconfig.v1.ListDiscoveryConfigsRequest
	2,  // 5: teleport.discoveryconfig.v1.DiscoveryConfigService.GetDiscoveryConfig:input_type -> teleport.discoveryconfig.v1.GetDiscoveryConfigRequest
	3,  // 6: teleport.discoveryconfig.v1.DiscoveryConfigService.CreateDiscoveryConfig:input_type -> teleport.discoveryconfig.v1.CreateDiscoveryConfigRequest
	4,  // 7: teleport.discoveryconfig.v1.DiscoveryConfigService.UpdateDiscoveryConfig:input_type -> teleport.discoveryconfig.v1.UpdateDiscoveryConfigRequest
	5,  // 8: teleport.discoveryconfig.v1.DiscoveryConfigService.UpsertDiscoveryConfig:input_type -> teleport.discoveryconfig.v1.UpsertDiscoveryConfigRequest
	6,  // 9: teleport.discoveryconfig.v1.DiscoveryConfigService.DeleteDiscoveryConfig:input_type -> teleport.discoveryconfig.v1.DeleteDiscoveryConfigRequest
	7,  // 10: teleport.discoveryconfig.v1.DiscoveryConfigService.DeleteAllDiscoveryConfigs:input_type -> teleport.discoveryconfig.v1.DeleteAllDiscoveryConfigsRequest
	1,  // 11: teleport.discoveryconfig.v1.DiscoveryConfigService.ListDiscoveryConfigs:output_type -> teleport.discoveryconfig.v1.ListDiscoveryConfigsResponse
	8,  // 12: teleport.discoveryconfig.v1.DiscoveryConfigService.GetDiscoveryConfig:output_type -> teleport.discoveryconfig.v1.DiscoveryConfig
	8,  // 13: teleport.discoveryconfig.v1.DiscoveryConfigService.CreateDiscoveryConfig:output_type -> teleport.discoveryconfig.v1.DiscoveryConfig
	8,  // 14: teleport.discoveryconfig.v1.DiscoveryConfigService.UpdateDiscoveryConfig:output_type -> teleport.discoveryconfig.v1.DiscoveryConfig
	8,  // 15: teleport.discoveryconfig.v1.DiscoveryConfigService.UpsertDiscoveryConfig:output_type -> teleport.discoveryconfig.v1.DiscoveryConfig
	9,  // 16: teleport.discoveryconfig.v1.DiscoveryConfigService.DeleteDiscoveryConfig:output_type -> google.protobuf.Empty
	9,  // 17: teleport.discoveryconfig.v1.DiscoveryConfigService.DeleteAllDiscoveryConfigs:output_type -> google.protobuf.Empty
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_init() }
func file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_init() {
	if File_teleport_discoveryconfig_v1_discoveryconfig_service_proto != nil {
		return
	}
	file_teleport_discoveryconfig_v1_discoveryconfig_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDiscoveryConfigsRequest); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListDiscoveryConfigsResponse); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDiscoveryConfigRequest); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateDiscoveryConfigRequest); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateDiscoveryConfigRequest); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertDiscoveryConfigRequest); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteDiscoveryConfigRequest); i {
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
		file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAllDiscoveryConfigsRequest); i {
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
			RawDescriptor: file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_goTypes,
		DependencyIndexes: file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_depIdxs,
		MessageInfos:      file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_msgTypes,
	}.Build()
	File_teleport_discoveryconfig_v1_discoveryconfig_service_proto = out.File
	file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_rawDesc = nil
	file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_goTypes = nil
	file_teleport_discoveryconfig_v1_discoveryconfig_service_proto_depIdxs = nil
}
