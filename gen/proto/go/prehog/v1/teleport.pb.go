// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
// source: prehog/v1/teleport.proto

package prehogv1

import (
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

type UserActivityReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// randomly generated UUID for this specific report, 16 bytes (in string order)
	ReportUuid []byte `protobuf:"bytes,1,opt,name=report_uuid,json=reportUuid,proto3" json:"report_uuid,omitempty"`
	// anonymized, 32 bytes (HMAC-SHA-256)
	ClusterName []byte `protobuf:"bytes,2,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	// anonymized, 32 bytes (HMAC-SHA-256)
	ReporterHostid []byte `protobuf:"bytes,3,opt,name=reporter_hostid,json=reporterHostid,proto3" json:"reporter_hostid,omitempty"`
	// beginning of the time window for this data; ending is not specified but is
	// intended to be at most 15 minutes
	StartTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Records   []*UserActivityRecord  `protobuf:"bytes,5,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *UserActivityReport) Reset() {
	*x = UserActivityReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prehog_v1_teleport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserActivityReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserActivityReport) ProtoMessage() {}

func (x *UserActivityReport) ProtoReflect() protoreflect.Message {
	mi := &file_prehog_v1_teleport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserActivityReport.ProtoReflect.Descriptor instead.
func (*UserActivityReport) Descriptor() ([]byte, []int) {
	return file_prehog_v1_teleport_proto_rawDescGZIP(), []int{0}
}

func (x *UserActivityReport) GetReportUuid() []byte {
	if x != nil {
		return x.ReportUuid
	}
	return nil
}

func (x *UserActivityReport) GetClusterName() []byte {
	if x != nil {
		return x.ClusterName
	}
	return nil
}

func (x *UserActivityReport) GetReporterHostid() []byte {
	if x != nil {
		return x.ReporterHostid
	}
	return nil
}

func (x *UserActivityReport) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *UserActivityReport) GetRecords() []*UserActivityRecord {
	if x != nil {
		return x.Records
	}
	return nil
}

type UserActivityRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// anonymized, 32 bytes (HMAC-SHA-256)
	UserName []byte `protobuf:"bytes,1,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	// counter of user.login events
	Logins uint64 `protobuf:"varint,2,opt,name=logins,proto3" json:"logins,omitempty"`
	// counter of session.start events (non-Kube)
	SshSessions uint64 `protobuf:"varint,3,opt,name=ssh_sessions,json=sshSessions,proto3" json:"ssh_sessions,omitempty"`
	// counter of app.session.start events (non-TCP)
	AppSessions uint64 `protobuf:"varint,4,opt,name=app_sessions,json=appSessions,proto3" json:"app_sessions,omitempty"`
	// counter of session.start events (only Kube)
	KubeSessions uint64 `protobuf:"varint,5,opt,name=kube_sessions,json=kubeSessions,proto3" json:"kube_sessions,omitempty"`
	// counter of db.session.start events
	DbSessions uint64 `protobuf:"varint,6,opt,name=db_sessions,json=dbSessions,proto3" json:"db_sessions,omitempty"`
	// counter of windows.desktop.session.start events
	DesktopSessions uint64 `protobuf:"varint,7,opt,name=desktop_sessions,json=desktopSessions,proto3" json:"desktop_sessions,omitempty"`
	// counter of app.session.start events (only TCP)
	AppTcpSessions uint64 `protobuf:"varint,8,opt,name=app_tcp_sessions,json=appTcpSessions,proto3" json:"app_tcp_sessions,omitempty"`
	// counter of port events (both SSH and Kube)
	//
	// Deprecated: Marked as deprecated in prehog/v1/teleport.proto.
	SshPortSessions uint64 `protobuf:"varint,9,opt,name=ssh_port_sessions,json=sshPortSessions,proto3" json:"ssh_port_sessions,omitempty"`
	// counter of kube.request events
	KubeRequests uint64 `protobuf:"varint,10,opt,name=kube_requests,json=kubeRequests,proto3" json:"kube_requests,omitempty"`
	// counter of sftp events
	SftpEvents uint64 `protobuf:"varint,11,opt,name=sftp_events,json=sftpEvents,proto3" json:"sftp_events,omitempty"`
	// counter of port events (only SSH)
	SshPortV2Sessions uint64 `protobuf:"varint,12,opt,name=ssh_port_v2_sessions,json=sshPortV2Sessions,proto3" json:"ssh_port_v2_sessions,omitempty"`
	// counter of port events (only Kube)
	KubePortSessions uint64 `protobuf:"varint,13,opt,name=kube_port_sessions,json=kubePortSessions,proto3" json:"kube_port_sessions,omitempty"`
}

func (x *UserActivityRecord) Reset() {
	*x = UserActivityRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prehog_v1_teleport_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserActivityRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserActivityRecord) ProtoMessage() {}

func (x *UserActivityRecord) ProtoReflect() protoreflect.Message {
	mi := &file_prehog_v1_teleport_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserActivityRecord.ProtoReflect.Descriptor instead.
func (*UserActivityRecord) Descriptor() ([]byte, []int) {
	return file_prehog_v1_teleport_proto_rawDescGZIP(), []int{1}
}

func (x *UserActivityRecord) GetUserName() []byte {
	if x != nil {
		return x.UserName
	}
	return nil
}

func (x *UserActivityRecord) GetLogins() uint64 {
	if x != nil {
		return x.Logins
	}
	return 0
}

func (x *UserActivityRecord) GetSshSessions() uint64 {
	if x != nil {
		return x.SshSessions
	}
	return 0
}

func (x *UserActivityRecord) GetAppSessions() uint64 {
	if x != nil {
		return x.AppSessions
	}
	return 0
}

func (x *UserActivityRecord) GetKubeSessions() uint64 {
	if x != nil {
		return x.KubeSessions
	}
	return 0
}

func (x *UserActivityRecord) GetDbSessions() uint64 {
	if x != nil {
		return x.DbSessions
	}
	return 0
}

func (x *UserActivityRecord) GetDesktopSessions() uint64 {
	if x != nil {
		return x.DesktopSessions
	}
	return 0
}

func (x *UserActivityRecord) GetAppTcpSessions() uint64 {
	if x != nil {
		return x.AppTcpSessions
	}
	return 0
}

// Deprecated: Marked as deprecated in prehog/v1/teleport.proto.
func (x *UserActivityRecord) GetSshPortSessions() uint64 {
	if x != nil {
		return x.SshPortSessions
	}
	return 0
}

func (x *UserActivityRecord) GetKubeRequests() uint64 {
	if x != nil {
		return x.KubeRequests
	}
	return 0
}

func (x *UserActivityRecord) GetSftpEvents() uint64 {
	if x != nil {
		return x.SftpEvents
	}
	return 0
}

func (x *UserActivityRecord) GetSshPortV2Sessions() uint64 {
	if x != nil {
		return x.SshPortV2Sessions
	}
	return 0
}

func (x *UserActivityRecord) GetKubePortSessions() uint64 {
	if x != nil {
		return x.KubePortSessions
	}
	return 0
}

type SubmitUsageReportsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// at most 10 in a single RPC, each shouldn't exceed 128KiB or so
	UserActivity []*UserActivityReport `protobuf:"bytes,1,rep,name=user_activity,json=userActivity,proto3" json:"user_activity,omitempty"`
}

func (x *SubmitUsageReportsRequest) Reset() {
	*x = SubmitUsageReportsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prehog_v1_teleport_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitUsageReportsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitUsageReportsRequest) ProtoMessage() {}

func (x *SubmitUsageReportsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_prehog_v1_teleport_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitUsageReportsRequest.ProtoReflect.Descriptor instead.
func (*SubmitUsageReportsRequest) Descriptor() ([]byte, []int) {
	return file_prehog_v1_teleport_proto_rawDescGZIP(), []int{2}
}

func (x *SubmitUsageReportsRequest) GetUserActivity() []*UserActivityReport {
	if x != nil {
		return x.UserActivity
	}
	return nil
}

type SubmitUsageReportsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// randomly generated UUID for this specific batch, 16 bytes (in string order)
	BatchUuid []byte `protobuf:"bytes,1,opt,name=batch_uuid,json=batchUuid,proto3" json:"batch_uuid,omitempty"`
}

func (x *SubmitUsageReportsResponse) Reset() {
	*x = SubmitUsageReportsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prehog_v1_teleport_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitUsageReportsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitUsageReportsResponse) ProtoMessage() {}

func (x *SubmitUsageReportsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_prehog_v1_teleport_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitUsageReportsResponse.ProtoReflect.Descriptor instead.
func (*SubmitUsageReportsResponse) Descriptor() ([]byte, []int) {
	return file_prehog_v1_teleport_proto_rawDescGZIP(), []int{3}
}

func (x *SubmitUsageReportsResponse) GetBatchUuid() []byte {
	if x != nil {
		return x.BatchUuid
	}
	return nil
}

var File_prehog_v1_teleport_proto protoreflect.FileDescriptor

var file_prehog_v1_teleport_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x72, 0x65, 0x68,
	0x6f, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf5, 0x01, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x55, 0x75, 0x69, 0x64, 0x12, 0x21,
	0x0a, 0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x5f, 0x68, 0x6f,
	0x73, 0x74, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x48, 0x6f, 0x73, 0x74, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x22, 0xff,
	0x03, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x73,
	0x68, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0b, 0x73, 0x73, 0x68, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x21, 0x0a,
	0x0c, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x61, 0x70, 0x70, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x6b, 0x75, 0x62, 0x65, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x62, 0x5f, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x64, 0x62, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x64, 0x65, 0x73, 0x6b, 0x74, 0x6f,
	0x70, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0f, 0x64, 0x65, 0x73, 0x6b, 0x74, 0x6f, 0x70, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x5f, 0x74, 0x63, 0x70, 0x5f, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x61, 0x70, 0x70,
	0x54, 0x63, 0x70, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2e, 0x0a, 0x11, 0x73,
	0x73, 0x68, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0f, 0x73, 0x73, 0x68, 0x50,
	0x6f, 0x72, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x6b,
	0x75, 0x62, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0c, 0x6b, 0x75, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x66, 0x74, 0x70, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x66, 0x74, 0x70, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x2f, 0x0a, 0x14, 0x73, 0x73, 0x68, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x76, 0x32,
	0x5f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x11, 0x73, 0x73, 0x68, 0x50, 0x6f, 0x72, 0x74, 0x56, 0x32, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10,
	0x6b, 0x75, 0x62, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x22, 0x5f, 0x0a, 0x19, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x42, 0x0a,
	0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x22, 0x3b, 0x0a, 0x1a, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x62, 0x61, 0x74, 0x63, 0x68, 0x55, 0x75, 0x69, 0x64, 0x32, 0x7f,
	0x0a, 0x18, 0x54, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x63, 0x0a, 0x12, 0x53, 0x75,
	0x62, 0x6d, 0x69, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73,
	0x12, 0x24, 0x2e, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62,
	0x6d, 0x69, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0xa6, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2e, 0x76,
	0x31, 0x42, 0x0d, 0x54, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x41, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67,
	0x72, 0x61, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x6c,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x72, 0x65,
	0x68, 0x6f, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x50, 0x72,
	0x65, 0x68, 0x6f, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x50, 0x72, 0x65, 0x68, 0x6f, 0x67,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x50, 0x72, 0x65, 0x68, 0x6f, 0x67, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x50, 0x72,
	0x65, 0x68, 0x6f, 0x67, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_prehog_v1_teleport_proto_rawDescOnce sync.Once
	file_prehog_v1_teleport_proto_rawDescData = file_prehog_v1_teleport_proto_rawDesc
)

func file_prehog_v1_teleport_proto_rawDescGZIP() []byte {
	file_prehog_v1_teleport_proto_rawDescOnce.Do(func() {
		file_prehog_v1_teleport_proto_rawDescData = protoimpl.X.CompressGZIP(file_prehog_v1_teleport_proto_rawDescData)
	})
	return file_prehog_v1_teleport_proto_rawDescData
}

var file_prehog_v1_teleport_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_prehog_v1_teleport_proto_goTypes = []interface{}{
	(*UserActivityReport)(nil),         // 0: prehog.v1.UserActivityReport
	(*UserActivityRecord)(nil),         // 1: prehog.v1.UserActivityRecord
	(*SubmitUsageReportsRequest)(nil),  // 2: prehog.v1.SubmitUsageReportsRequest
	(*SubmitUsageReportsResponse)(nil), // 3: prehog.v1.SubmitUsageReportsResponse
	(*timestamppb.Timestamp)(nil),      // 4: google.protobuf.Timestamp
}
var file_prehog_v1_teleport_proto_depIdxs = []int32{
	4, // 0: prehog.v1.UserActivityReport.start_time:type_name -> google.protobuf.Timestamp
	1, // 1: prehog.v1.UserActivityReport.records:type_name -> prehog.v1.UserActivityRecord
	0, // 2: prehog.v1.SubmitUsageReportsRequest.user_activity:type_name -> prehog.v1.UserActivityReport
	2, // 3: prehog.v1.TeleportReportingService.SubmitUsageReports:input_type -> prehog.v1.SubmitUsageReportsRequest
	3, // 4: prehog.v1.TeleportReportingService.SubmitUsageReports:output_type -> prehog.v1.SubmitUsageReportsResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_prehog_v1_teleport_proto_init() }
func file_prehog_v1_teleport_proto_init() {
	if File_prehog_v1_teleport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_prehog_v1_teleport_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserActivityReport); i {
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
		file_prehog_v1_teleport_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserActivityRecord); i {
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
		file_prehog_v1_teleport_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitUsageReportsRequest); i {
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
		file_prehog_v1_teleport_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitUsageReportsResponse); i {
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
			RawDescriptor: file_prehog_v1_teleport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_prehog_v1_teleport_proto_goTypes,
		DependencyIndexes: file_prehog_v1_teleport_proto_depIdxs,
		MessageInfos:      file_prehog_v1_teleport_proto_msgTypes,
	}.Build()
	File_prehog_v1_teleport_proto = out.File
	file_prehog_v1_teleport_proto_rawDesc = nil
	file_prehog_v1_teleport_proto_goTypes = nil
	file_prehog_v1_teleport_proto_depIdxs = nil
}
