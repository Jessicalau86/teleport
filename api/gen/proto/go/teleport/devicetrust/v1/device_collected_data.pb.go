// Copyright 2022 Gravitational, Inc
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
// source: teleport/devicetrust/v1/device_collected_data.proto

package devicetrustv1

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

// DeviceCollectedData contains information gathered from the device during
// various ceremonies.
// Gathered information must match, within reason, the original registration
// data and previous instances of collected data.
type DeviceCollectedData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Time of data collection, set by the client.
	// Required.
	CollectTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=collect_time,json=collectTime,proto3" json:"collect_time,omitempty"`
	// Time of data collection, as received by the server.
	// System managed.
	RecordTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=record_time,json=recordTime,proto3" json:"record_time,omitempty"`
	// Device operating system.
	// Required.
	OsType OSType `protobuf:"varint,3,opt,name=os_type,json=osType,proto3,enum=teleport.devicetrust.v1.OSType" json:"os_type,omitempty"`
	// Device serial number used to match the device with the inventory.
	// This field is one of the three following
	// values in this precedence:
	// - reported_asset_tag
	// - system_serial_number
	// - base_board_serial_number
	// Required.
	SerialNumber string `protobuf:"bytes,4,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	// Non-descriptive model identifier.
	// Example: "MacBookPro9,2".
	ModelIdentifier string `protobuf:"bytes,5,opt,name=model_identifier,json=modelIdentifier,proto3" json:"model_identifier,omitempty"`
	// OS version number, without the leading 'v'.
	// Example: "13.2.1".
	OsVersion string `protobuf:"bytes,6,opt,name=os_version,json=osVersion,proto3" json:"os_version,omitempty"`
	// OS build identifier. Augments the os_version.
	// May match either the DeviceProfile os_build or os_build_supplemental.
	// Example: "22D68" or "22F770820d".
	OsBuild string `protobuf:"bytes,7,opt,name=os_build,json=osBuild,proto3" json:"os_build,omitempty"`
	// OS username (distinct from the Teleport user).
	OsUsername string `protobuf:"bytes,8,opt,name=os_username,json=osUsername,proto3" json:"os_username,omitempty"`
	// Jamf binary version, without the leading 'v'.
	// Example: "9.27" or "10.44.1-t1677509507".
	JamfBinaryVersion string `protobuf:"bytes,9,opt,name=jamf_binary_version,json=jamfBinaryVersion,proto3" json:"jamf_binary_version,omitempty"`
	// Unmodified output of `/usr/bin/profiles status -type enrollment`.
	// Used to verify the presence of an enrollment profile.
	MacosEnrollmentProfiles string `protobuf:"bytes,10,opt,name=macos_enrollment_profiles,json=macosEnrollmentProfiles,proto3" json:"macos_enrollment_profiles,omitempty"`
	// The asset tag of the device as reported by the BIOS DMI Type 3. Tools
	// used by customers to manage their fleet may set this value.
	ReportedAssetTag string `protobuf:"bytes,11,opt,name=reported_asset_tag,json=reportedAssetTag,proto3" json:"reported_asset_tag,omitempty"`
	// The serial number of the "system" as reported by the BIOS DMI Type 1.
	// This field can be empty if no value has been configured.
	SystemSerialNumber string `protobuf:"bytes,12,opt,name=system_serial_number,json=systemSerialNumber,proto3" json:"system_serial_number,omitempty"`
	// The serial number of the "base board" as reported by BIOS DMI Type 2.
	// This field can be empty if no value has been configured.
	BaseBoardSerialNumber string `protobuf:"bytes,13,opt,name=base_board_serial_number,json=baseBoardSerialNumber,proto3" json:"base_board_serial_number,omitempty"`
	// If during the collection of this device data, the device performed a TPM
	// platform attestation (e.g during enrollment or authentication), then this
	// field holds the record of this attestation. This allows the state of the
	// device to be compared to historical state, and allows for the platform
	// attestations to be revalidated at a later date.
	//
	// This field is not explicitly sent up by the client, and any DCD sent by a
	// client including this field should be rejected. The server should inject
	// this field once verifying that the submitted platform attestation during
	// the enrollment or authentication.
	//
	// System managed.
	TpmPlatformAttestation *TPMPlatformAttestation `protobuf:"bytes,14,opt,name=tpm_platform_attestation,json=tpmPlatformAttestation,proto3" json:"tpm_platform_attestation,omitempty"`
}

func (x *DeviceCollectedData) Reset() {
	*x = DeviceCollectedData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teleport_devicetrust_v1_device_collected_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceCollectedData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceCollectedData) ProtoMessage() {}

func (x *DeviceCollectedData) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_devicetrust_v1_device_collected_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceCollectedData.ProtoReflect.Descriptor instead.
func (*DeviceCollectedData) Descriptor() ([]byte, []int) {
	return file_teleport_devicetrust_v1_device_collected_data_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceCollectedData) GetCollectTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CollectTime
	}
	return nil
}

func (x *DeviceCollectedData) GetRecordTime() *timestamppb.Timestamp {
	if x != nil {
		return x.RecordTime
	}
	return nil
}

func (x *DeviceCollectedData) GetOsType() OSType {
	if x != nil {
		return x.OsType
	}
	return OSType_OS_TYPE_UNSPECIFIED
}

func (x *DeviceCollectedData) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *DeviceCollectedData) GetModelIdentifier() string {
	if x != nil {
		return x.ModelIdentifier
	}
	return ""
}

func (x *DeviceCollectedData) GetOsVersion() string {
	if x != nil {
		return x.OsVersion
	}
	return ""
}

func (x *DeviceCollectedData) GetOsBuild() string {
	if x != nil {
		return x.OsBuild
	}
	return ""
}

func (x *DeviceCollectedData) GetOsUsername() string {
	if x != nil {
		return x.OsUsername
	}
	return ""
}

func (x *DeviceCollectedData) GetJamfBinaryVersion() string {
	if x != nil {
		return x.JamfBinaryVersion
	}
	return ""
}

func (x *DeviceCollectedData) GetMacosEnrollmentProfiles() string {
	if x != nil {
		return x.MacosEnrollmentProfiles
	}
	return ""
}

func (x *DeviceCollectedData) GetReportedAssetTag() string {
	if x != nil {
		return x.ReportedAssetTag
	}
	return ""
}

func (x *DeviceCollectedData) GetSystemSerialNumber() string {
	if x != nil {
		return x.SystemSerialNumber
	}
	return ""
}

func (x *DeviceCollectedData) GetBaseBoardSerialNumber() string {
	if x != nil {
		return x.BaseBoardSerialNumber
	}
	return ""
}

func (x *DeviceCollectedData) GetTpmPlatformAttestation() *TPMPlatformAttestation {
	if x != nil {
		return x.TpmPlatformAttestation
	}
	return nil
}

var File_teleport_devicetrust_v1_device_collected_data_proto protoreflect.FileDescriptor

var file_teleport_devicetrust_v1_device_collected_data_proto_rawDesc = []byte{
	0x0a, 0x33, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x74, 0x72, 0x75, 0x73, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x74, 0x72, 0x75, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x25, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x74, 0x72, 0x75, 0x73, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x74, 0x72, 0x75, 0x73, 0x74, 0x2f, 0x76, 0x31, 0x2f,
	0x74, 0x70, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe6, 0x05, 0x0a, 0x13, 0x44, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x3d, 0x0a, 0x0c, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x0b, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x3b, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x38, 0x0a,
	0x07, 0x6f, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f,
	0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x74, 0x72, 0x75, 0x73, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x53, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x06, 0x6f, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x29, 0x0a, 0x10,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x73, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x73, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x73, 0x5f, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x73, 0x42, 0x75, 0x69, 0x6c,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6f, 0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6f, 0x73, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x6a, 0x61, 0x6d, 0x66, 0x5f, 0x62, 0x69, 0x6e, 0x61, 0x72,
	0x79, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x11, 0x6a, 0x61, 0x6d, 0x66, 0x42, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x19, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x5f, 0x65, 0x6e, 0x72, 0x6f,
	0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x45, 0x6e, 0x72, 0x6f,
	0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x2c,
	0x0a, 0x12, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x5f, 0x74, 0x61, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x64, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x61, 0x67, 0x12, 0x30, 0x0a, 0x14,
	0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x37,
	0x0a, 0x18, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x5f, 0x73, 0x65, 0x72,
	0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x15, 0x62, 0x61, 0x73, 0x65, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x53, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x69, 0x0a, 0x18, 0x74, 0x70, 0x6d, 0x5f, 0x70,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x74, 0x72, 0x75, 0x73, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x54, 0x50, 0x4d, 0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x41,
	0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x16, 0x74, 0x70, 0x6d, 0x50,
	0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x41, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x5a, 0x5a, 0x58, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x67, 0x72, 0x61, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f, 0x74,
	0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x74, 0x72, 0x75, 0x73, 0x74, 0x2f, 0x76, 0x31,
	0x3b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x74, 0x72, 0x75, 0x73, 0x74, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_teleport_devicetrust_v1_device_collected_data_proto_rawDescOnce sync.Once
	file_teleport_devicetrust_v1_device_collected_data_proto_rawDescData = file_teleport_devicetrust_v1_device_collected_data_proto_rawDesc
)

func file_teleport_devicetrust_v1_device_collected_data_proto_rawDescGZIP() []byte {
	file_teleport_devicetrust_v1_device_collected_data_proto_rawDescOnce.Do(func() {
		file_teleport_devicetrust_v1_device_collected_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_teleport_devicetrust_v1_device_collected_data_proto_rawDescData)
	})
	return file_teleport_devicetrust_v1_device_collected_data_proto_rawDescData
}

var file_teleport_devicetrust_v1_device_collected_data_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_teleport_devicetrust_v1_device_collected_data_proto_goTypes = []interface{}{
	(*DeviceCollectedData)(nil),    // 0: teleport.devicetrust.v1.DeviceCollectedData
	(*timestamppb.Timestamp)(nil),  // 1: google.protobuf.Timestamp
	(OSType)(0),                    // 2: teleport.devicetrust.v1.OSType
	(*TPMPlatformAttestation)(nil), // 3: teleport.devicetrust.v1.TPMPlatformAttestation
}
var file_teleport_devicetrust_v1_device_collected_data_proto_depIdxs = []int32{
	1, // 0: teleport.devicetrust.v1.DeviceCollectedData.collect_time:type_name -> google.protobuf.Timestamp
	1, // 1: teleport.devicetrust.v1.DeviceCollectedData.record_time:type_name -> google.protobuf.Timestamp
	2, // 2: teleport.devicetrust.v1.DeviceCollectedData.os_type:type_name -> teleport.devicetrust.v1.OSType
	3, // 3: teleport.devicetrust.v1.DeviceCollectedData.tpm_platform_attestation:type_name -> teleport.devicetrust.v1.TPMPlatformAttestation
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_teleport_devicetrust_v1_device_collected_data_proto_init() }
func file_teleport_devicetrust_v1_device_collected_data_proto_init() {
	if File_teleport_devicetrust_v1_device_collected_data_proto != nil {
		return
	}
	file_teleport_devicetrust_v1_os_type_proto_init()
	file_teleport_devicetrust_v1_tpm_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_teleport_devicetrust_v1_device_collected_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceCollectedData); i {
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
			RawDescriptor: file_teleport_devicetrust_v1_device_collected_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_devicetrust_v1_device_collected_data_proto_goTypes,
		DependencyIndexes: file_teleport_devicetrust_v1_device_collected_data_proto_depIdxs,
		MessageInfos:      file_teleport_devicetrust_v1_device_collected_data_proto_msgTypes,
	}.Build()
	File_teleport_devicetrust_v1_device_collected_data_proto = out.File
	file_teleport_devicetrust_v1_device_collected_data_proto_rawDesc = nil
	file_teleport_devicetrust_v1_device_collected_data_proto_goTypes = nil
	file_teleport_devicetrust_v1_device_collected_data_proto_depIdxs = nil
}
