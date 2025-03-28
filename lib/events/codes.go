/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package events

import apievents "github.com/gravitational/teleport/api/types/events"

// Event describes an audit log event.
type Event struct {
	// Name is the event name.
	Name string
	// Code is the unique event code.
	Code string
}

// There is no strict algorithm for picking an event code, however existing
// event codes are currently loosely categorized as follows:
//
//   - Teleport event codes start with "T" and belong in this const block.
//
//   - Related events are grouped starting with the same number.
//     eg: All user related events are grouped under 1xxx.
//
//   - Suffix code with one of these letters: I (info), W (warn), E (error).
//
// After defining an event code, make sure to keep
// `web/packages/teleport/src/services/audit/types.ts` in sync.
const (
	// UserLocalLoginCode is the successful local user login event code.
	UserLocalLoginCode = "T1000I"
	// UserLocalLoginFailureCode is the unsuccessful local user login event code.
	UserLocalLoginFailureCode = "T1000W"
	// UserSSOLoginCode is the successful SSO user login event code.
	UserSSOLoginCode = "T1001I"
	// UserSSOLoginFailureCode is the unsuccessful SSO user login event code.
	UserSSOLoginFailureCode = "T1001W"
	// UserCreateCode is the user create event code.
	UserCreateCode = "T1002I"
	// UserUpdateCode is the user update event code.
	UserUpdateCode = "T1003I"
	// UserDeleteCode is the user delete event code.
	UserDeleteCode = "T1004I"
	// UserPasswordChangeCode is an event code for when user changes their own password.
	UserPasswordChangeCode = "T1005I"
	// MFADeviceAddEventCode is an event code for users adding MFA devices.
	MFADeviceAddEventCode = "T1006I"
	// MFADeviceDeleteEventCode is an event code for users deleting MFA devices.
	MFADeviceDeleteEventCode = "T1007I"
	// RecoveryCodesGenerateCode is an event code for generation of recovery codes.
	RecoveryCodesGenerateCode = "T1008I"
	// RecoveryCodeUseSuccessCode is an event code for when a
	// recovery code was used successfully.
	RecoveryCodeUseSuccessCode = "T1009I"
	// RecoveryCodeUseFailureCode is an event code for when a
	// recovery code was not used successfully.
	RecoveryCodeUseFailureCode = "T1009W"
	// UserSSOTestFlowLoginCode is the successful SSO test flow user login event code.
	UserSSOTestFlowLoginCode = "T1010I"
	// UserSSOTestFlowLoginFailureCode is the unsuccessful SSO test flow user login event code.
	UserSSOTestFlowLoginFailureCode = "T1011W"
	// UserHeadlessLoginRequestedCode is an event code for when headless login attempt was requested.
	UserHeadlessLoginRequestedCode = "T1012I"
	// UserHeadlessLoginApprovedCode is an event code for when headless login attempt was successfully approved.
	UserHeadlessLoginApprovedCode = "T1013I"
	// UserHeadlessLoginApprovedFailureCode is an event code for when headless login was approved with an error.
	UserHeadlessLoginApprovedFailureCode = "T1013W"
	// UserHeadlessLoginRejectedCode is an event code for when headless login attempt was rejected.
	UserHeadlessLoginRejectedCode = "T1014W"

	// BillingCardCreateCode is an event code for when a user creates a new credit card.
	BillingCardCreateCode = "TBL00I"
	// BillingCardDeleteCode is an event code for when a user deletes a credit card.
	BillingCardDeleteCode = "TBL01I"
	// BillingCardUpdateCode is an event code for when a user updates an existing credit card.
	BillingCardUpdateCode = "TBL02I"
	// BillingInformationUpdateCode is an event code for when a user updates their billing info.
	BillingInformationUpdateCode = "TBL03I"

	// SessionRejectedCode is an event code for when a user's attempt to create an
	// session/connection has been rejected.
	SessionRejectedCode = "T1006W"

	// SessionStartCode is the session start event code.
	SessionStartCode = "T2000I"
	// SessionJoinCode is the session join event code.
	SessionJoinCode = "T2001I"
	// TerminalResizeCode is the terminal resize event code.
	TerminalResizeCode = "T2002I"
	// SessionLeaveCode is the session leave event code.
	SessionLeaveCode = "T2003I"
	// SessionEndCode is the session end event code.
	SessionEndCode = "T2004I"
	// SessionUploadCode is the session upload event code.
	SessionUploadCode = "T2005I"
	// SessionDataCode is the session data event code.
	SessionDataCode = "T2006I"
	// AppSessionStartCode is the application session start code.
	AppSessionStartCode = "T2007I"
	// AppSessionChunkCode is the application session chunk create code.
	AppSessionChunkCode = "T2008I"
	// AppSessionRequestCode is the application request/response code.
	AppSessionRequestCode = "T2009I"
	// SessionConnectCode is the session connect event code.
	SessionConnectCode = "T2010I"
	// AppSessionEndCode is the application session end event code.
	AppSessionEndCode = "T2011I"
	// SessionRecordingAccessCode is the session recording view data event code.
	SessionRecordingAccessCode = "T2012I"
	// AppSessionDynamoDBRequestCode is the application request/response code.
	AppSessionDynamoDBRequestCode = "T2013I"

	// AppCreateCode is the app.create event code.
	AppCreateCode = "TAP03I"
	// AppUpdateCode is the app.update event code.
	AppUpdateCode = "TAP04I"
	// AppDeleteCode is the app.delete event code.
	AppDeleteCode = "TAP05I"

	// DatabaseSessionStartCode is the database session start event code.
	DatabaseSessionStartCode = "TDB00I"
	// DatabaseSessionStartFailureCode is the database session start failure event code.
	DatabaseSessionStartFailureCode = "TDB00W"
	// DatabaseSessionEndCode is the database session end event code.
	DatabaseSessionEndCode = "TDB01I"
	// DatabaseSessionQueryCode is the database query event code.
	DatabaseSessionQueryCode = "TDB02I"
	// DatabaseSessionQueryFailedCode is the database query failure event code.
	DatabaseSessionQueryFailedCode = "TDB02W"
	// DatabaseSessionMalformedPacketCode is the db.session.malformed_packet event code.
	DatabaseSessionMalformedPacketCode = "TDB06I"

	// PostgresParseCode is the db.session.postgres.statements.parse event code.
	PostgresParseCode = "TPG00I"
	// PostgresBindCode is the db.session.postgres.statements.bind event code.
	PostgresBindCode = "TPG01I"
	// PostgresExecuteCode is the db.session.postgres.statements.execute event code.
	PostgresExecuteCode = "TPG02I"
	// PostgresCloseCode is the db.session.postgres.statements.close event code.
	PostgresCloseCode = "TPG03I"
	// PostgresFunctionCallCode is the db.session.postgres.function event code.
	PostgresFunctionCallCode = "TPG04I"

	// MySQLStatementPrepareCode is the db.session.mysql.statements.prepare event code.
	MySQLStatementPrepareCode = "TMY00I"
	// MySQLStatementExecuteCode is the db.session.mysql.statements.execute event code.
	MySQLStatementExecuteCode = "TMY01I"
	// MySQLStatementSendLongDataCode is the db.session.mysql.statements.send_long_data event code.
	MySQLStatementSendLongDataCode = "TMY02I"
	// MySQLStatementCloseCode is the db.session.mysql.statements.close event code.
	MySQLStatementCloseCode = "TMY03I"
	// MySQLStatementResetCode is the db.session.mysql.statements.reset event code.
	MySQLStatementResetCode = "TMY04I"
	// MySQLStatementFetchCode is the db.session.mysql.statements.fetch event code.
	MySQLStatementFetchCode = "TMY05I"
	// MySQLStatementBulkExecuteCode is the db.session.mysql.statements.bulk_execute event code.
	MySQLStatementBulkExecuteCode = "TMY06I"
	// MySQLInitDBCode is the db.session.mysql.init_db event code.
	MySQLInitDBCode = "TMY07I"
	// MySQLCreateDBCode is the db.session.mysql.create_db event code.
	MySQLCreateDBCode = "TMY08I"
	// MySQLDropDBCode is the db.session.mysql.drop_db event code.
	MySQLDropDBCode = "TMY09I"
	// MySQLShutDownCode is the db.session.mysql.shut_down event code.
	MySQLShutDownCode = "TMY10I"
	// MySQLProcessKillCode is the db.session.mysql.process_kill event code.
	MySQLProcessKillCode = "TMY11I"
	// MySQLDebugCode is the db.session.mysql.debug event code.
	MySQLDebugCode = "TMY12I"
	// MySQLRefreshCode is the db.session.mysql.refresh event code.
	MySQLRefreshCode = "TMY13I"

	// SQLServerRPCRequestCode is the db.session.sqlserver.rpc_request event code.
	SQLServerRPCRequestCode = "TMS00I"

	// CassandraBatchEventCode is the db.session.cassandra.batch event code.
	CassandraBatchEventCode = "TCA01I"
	// CassandraPrepareEventCode is the db.session.cassandra.prepare event code.
	CassandraPrepareEventCode = "TCA02I"
	// CassandraExecuteEventCode is the db.session.cassandra.execute event code.
	CassandraExecuteEventCode = "TCA03I"
	// CassandraRegisterEventCode is the db.session.cassandra.register event code.
	CassandraRegisterEventCode = "TCA04I"

	// ElasticsearchRequestCode is the db.session.elasticsearch.request event code.
	ElasticsearchRequestCode = "TES00I"
	// ElasticsearchRequestFailureCode is the db.session.elasticsearch.request event failure code.
	ElasticsearchRequestFailureCode = "TES00E"

	// OpenSearchRequestCode is the db.session.opensearch.request event code.
	OpenSearchRequestCode = "TOS00I"
	// OpenSearchRequestFailureCode is the db.session.opensearch.request event failure code.
	OpenSearchRequestFailureCode = "TOS00E"

	// DynamoDBRequestCode is the db.session.dynamodb.request event code.
	DynamoDBRequestCode = "TDY01I"
	// DynamoDBRequestFailureCode is the db.session.dynamodb.request event failure code.
	// This is indicates that the database agent http transport failed to round trip the request.
	DynamoDBRequestFailureCode = "TDY01E"

	// DatabaseCreateCode is the db.create event code.
	DatabaseCreateCode = "TDB03I"
	// DatabaseUpdateCode is the db.update event code.
	DatabaseUpdateCode = "TDB04I"
	// DatabaseDeleteCode is the db.delete event code.
	DatabaseDeleteCode = "TDB05I"

	// DesktopSessionStartCode is the desktop session start event code.
	DesktopSessionStartCode = "TDP00I"
	// DesktopSessionStartFailureCode is event code for desktop sessions
	// that failed to start.
	DesktopSessionStartFailureCode = "TDP00W"
	// DesktopSessionEndCode is the desktop session end event code.
	DesktopSessionEndCode = "TDP01I"
	// DesktopClipboardSendCode is the desktop clipboard send code.
	DesktopClipboardSendCode = "TDP02I"
	// DesktopClipboardReceiveCode is the desktop clipboard receive code.
	DesktopClipboardReceiveCode = "TDP03I"
	// DesktopSharedDirectoryStartCode is the desktop directory start code.
	DesktopSharedDirectoryStartCode = "TDP04I"
	// DesktopSharedDirectoryStartFailureCode is the desktop directory start code
	// for when a start operation fails, or for when the internal cache state was corrupted
	// causing information loss, or for when the internal cache has exceeded its max size.
	DesktopSharedDirectoryStartFailureCode = "TDP04W"
	// DesktopSharedDirectoryReadCode is the desktop directory read code.
	DesktopSharedDirectoryReadCode = "TDP05I"
	// DesktopSharedDirectoryReadFailureCode is the desktop directory read code
	// for when a read operation fails, or for if the internal cache state was corrupted
	// causing information loss, or for when the internal cache has exceeded its max size.
	DesktopSharedDirectoryReadFailureCode = "TDP05W"
	// DesktopSharedDirectoryWriteCode is the desktop directory write code.
	DesktopSharedDirectoryWriteCode = "TDP06I"
	// DesktopSharedDirectoryWriteFailureCode is the desktop directory write code
	// for when a write operation fails, or for if the internal cache state was corrupted
	// causing information loss, or for when the internal cache has exceeded its max size.
	DesktopSharedDirectoryWriteFailureCode = "TDP06W"

	// SubsystemCode is the subsystem event code.
	SubsystemCode = "T3001I"
	// SubsystemFailureCode is the subsystem failure event code.
	SubsystemFailureCode = "T3001E"
	// ExecCode is the exec event code.
	ExecCode = "T3002I"
	// ExecFailureCode is the exec failure event code.
	ExecFailureCode = "T3002E"
	// PortForwardCode is the port forward event code.
	PortForwardCode = "T3003I"
	// PortForwardFailureCode is the port forward failure event code.
	PortForwardFailureCode = "T3003E"
	// SCPDownloadCode is the file download event code.
	SCPDownloadCode = "T3004I"
	// SCPDownloadFailureCode is the file download event failure code.
	SCPDownloadFailureCode = "T3004E"
	// SCPUploadCode is the file upload event code.
	SCPUploadCode = "T3005I"
	// SCPUploadFailureCode is the file upload failure event code.
	SCPUploadFailureCode = "T3005E"
	// ClientDisconnectCode is the client disconnect event code.
	ClientDisconnectCode = "T3006I"
	// AuthAttemptFailureCode is the auth attempt failure event code.
	AuthAttemptFailureCode = "T3007W"
	// X11ForwardCode is the x11 forward event code.
	X11ForwardCode = "T3008I"
	// X11ForwardFailureCode is the x11 forward failure event code.
	X11ForwardFailureCode = "T3008W"
	// KubeRequestCode is an event code for a generic kubernetes request.
	//
	// Note: some requests (like exec into a pod) use other codes (like
	// ExecCode).
	KubeRequestCode = "T3009I"

	// KubernetesClusterCreateCode is the kube.create event code.
	KubernetesClusterCreateCode = "T3010I"
	// KubernetesClusterUpdateCode is the kube.update event code.
	KubernetesClusterUpdateCode = "T3011I"
	// KubernetesClusterDeleteCode is the kube.delete event code.
	KubernetesClusterDeleteCode = "T3012I"

	// The following codes correspond to SFTP file operations.
	SFTPOpenCode            = "TS001I"
	SFTPOpenFailureCode     = "TS001E"
	SFTPCloseCode           = "TS002I"
	SFTPCloseFailureCode    = "TS002E"
	SFTPReadCode            = "TS003I"
	SFTPReadFailureCode     = "TS003E"
	SFTPWriteCode           = "TS004I"
	SFTPWriteFailureCode    = "TS004E"
	SFTPLstatCode           = "TS005I"
	SFTPLstatFailureCode    = "TS005E"
	SFTPFstatCode           = "TS006I"
	SFTPFstatFailureCode    = "TS006E"
	SFTPSetstatCode         = "TS007I"
	SFTPSetstatFailureCode  = "TS007E"
	SFTPFsetstatCode        = "TS008I"
	SFTPFsetstatFailureCode = "TS008E"
	SFTPOpendirCode         = "TS009I"
	SFTPOpendirFailureCode  = "TS009E"
	SFTPReaddirCode         = "TS010I"
	SFTPReaddirFailureCode  = "TS010E"
	SFTPRemoveCode          = "TS011I"
	SFTPRemoveFailureCode   = "TS011E"
	SFTPMkdirCode           = "TS012I"
	SFTPMkdirFailureCode    = "TS012E"
	SFTPRmdirCode           = "TS013I"
	SFTPRmdirFailureCode    = "TS013E"
	SFTPRealpathCode        = "TS014I"
	SFTPRealpathFailureCode = "TS014E"
	SFTPStatCode            = "TS015I"
	SFTPStatFailureCode     = "TS015E"
	SFTPRenameCode          = "TS016I"
	SFTPRenameFailureCode   = "TS016E"
	SFTPReadlinkCode        = "TS017I"
	SFTPReadlinkFailureCode = "TS017E"
	SFTPSymlinkCode         = "TS018I"
	SFTPSymlinkFailureCode  = "TS018E"
	SFTPLinkCode            = "TS019I"
	SFTPLinkFailureCode     = "TS019E"

	// SessionCommandCode is a session command code.
	SessionCommandCode = "T4000I"
	// SessionDiskCode is a session disk code.
	SessionDiskCode = "T4001I"
	// SessionNetworkCode is a session network code.
	SessionNetworkCode = "T4002I"

	// AccessRequestCreateCode is the the access request creation code.
	AccessRequestCreateCode = "T5000I"
	// AccessRequestUpdateCode is the access request state update code.
	AccessRequestUpdateCode = "T5001I"
	// AccessRequestReviewCode is the access review application code.
	AccessRequestReviewCode = "T5002I"
	// AccessRequestDeleteCode is the access request deleted code.
	AccessRequestDeleteCode = "T5003I"
	// AccessRequestResourceSearchCode is the access request resource search code.
	AccessRequestResourceSearchCode = "T5004I"

	// ResetPasswordTokenCreateCode is the token create event code.
	ResetPasswordTokenCreateCode = "T6000I"
	// RecoveryTokenCreateCode is the recovery token create event code.
	RecoveryTokenCreateCode = "T6001I"
	// PrivilegeTokenCreateCode is the privilege token create event code.
	PrivilegeTokenCreateCode = "T6002I"

	// TrustedClusterCreateCode is the event code for creating a trusted cluster.
	TrustedClusterCreateCode = "T7000I"
	// TrustedClusterDeleteCode is the event code for removing a trusted cluster.
	TrustedClusterDeleteCode = "T7001I"
	// TrustedClusterTokenCreateCode is the event code for creating new
	// provisioning token for a trusted cluster. Deprecated in favor of
	// [ProvisionTokenCreateEvent].
	TrustedClusterTokenCreateCode = "T7002I"

	// ProvisionTokenCreateCode is the event code for creating a provisioning
	// token, also known as Join Token. See
	// [github.com/gravitational/teleport/api/types.ProvisionToken].
	ProvisionTokenCreateCode = "TJT00I"

	// GithubConnectorCreatedCode is the Github connector created event code.
	GithubConnectorCreatedCode = "T8000I"
	// GithubConnectorDeletedCode is the Github connector deleted event code.
	GithubConnectorDeletedCode = "T8001I"

	// OIDCConnectorCreatedCode is the OIDC connector created event code.
	OIDCConnectorCreatedCode = "T8100I"
	// OIDCConnectorDeletedCode is the OIDC connector deleted event code.
	OIDCConnectorDeletedCode = "T8101I"

	// SAMLConnectorCreatedCode is the SAML connector created event code.
	SAMLConnectorCreatedCode = "T8200I"
	// SAMLConnectorDeletedCode is the SAML connector deleted event code.
	SAMLConnectorDeletedCode = "T8201I"

	// RoleCreatedCode is the role created event code.
	RoleCreatedCode = "T9000I"
	// RoleDeletedCode is the role deleted event code.
	RoleDeletedCode = "T9001I"

	// BotJoinCode is the 'bot.join' event code.
	BotJoinCode = "TJ001I"
	// InstanceJoinCode is the 'node.join' event code.
	InstanceJoinCode = "TJ002I"

	// LockCreatedCode is the lock created event code.
	LockCreatedCode = "TLK00I"
	// LockDeletedCode is the lock deleted event code.
	LockDeletedCode = "TLK01I"

	// CertificateCreateCode is the certificate issuance event code.
	CertificateCreateCode = "TC000I"

	// RenewableCertificateGenerationMismatchCode is the renewable cert
	// generation mismatch code.
	RenewableCertificateGenerationMismatchCode = "TCB00W"

	// UpgradeWindowStartUpdatedCode is the edit code of UpgradeWindowStartUpdateEvent.
	UpgradeWindowStartUpdatedCode = "TUW01I"

	// SSMRunSuccessCode is the discovery script success code.
	SSMRunSuccessCode = "TDS00I"
	// SSMRunFailCode is the discovery script success code.
	SSMRunFailCode = "TDS00W"

	// DeviceCreateCode is the device creation/registration code.
	DeviceCreateCode = "TV001I"
	// DeviceDeleteCode is the device deletion code.
	DeviceDeleteCode = "TV002I"
	// DeviceEnrollTokenCreateCode is the device enroll token creation code
	DeviceEnrollTokenCreateCode = "TV003I"
	// DeviceEnrollTokenSpentCode is the device enroll token spent code.
	DeviceEnrollTokenSpentCode = "TV004I"
	// DeviceEnrollCode is the device enrollment completion code.
	DeviceEnrollCode = "TV005I"
	// DeviceAuthenticateCode is the device authentication code.
	DeviceAuthenticateCode = "TV006I"
	// DeviceUpdateCode is the device update code.
	DeviceUpdateCode = "TV007I"

	// LoginRuleCreateCode is the login rule create code.
	LoginRuleCreateCode = "TLR00I"
	// LoginRuleDeleteCode is the login rule delete code.
	LoginRuleDeleteCode = "TLR01I"

	// SAMLIdPAuthAttemptCode is the SAML IdP auth attempt code.
	SAMLIdPAuthAttemptCode = "TSI000I"

	// SAMLIdPServiceProviderCreateCode is the SAML IdP service provider create code.
	SAMLIdPServiceProviderCreateCode = "TSI001I"

	// SAMLIdPServiceProviderCreateFailureCode is the SAML IdP service provider create failure code.
	SAMLIdPServiceProviderCreateFailureCode = "TSI001W"

	// SAMLIdPServiceProviderUpdateCode is the SAML IdP service provider update code.
	SAMLIdPServiceProviderUpdateCode = "TSI002I"

	// SAMLIdPServiceProviderUpdateFailureCode is the SAML IdP service provider update failure code.
	SAMLIdPServiceProviderUpdateFailureCode = "TSI002W"

	// SAMLIdPServiceProviderDeleteCode is the SAML IdP service provider delete code.
	SAMLIdPServiceProviderDeleteCode = "TSI003I"

	// SAMLIdPServiceProviderDeleteFailureCode is the SAML IdP service provider delete failure code.
	SAMLIdPServiceProviderDeleteFailureCode = "TSI003W"

	// SAMLIdPServiceProviderDeleteAllCode is the SAML IdP service provider delete all code.
	SAMLIdPServiceProviderDeleteAllCode = "TSI004I"

	// SAMLIdPServiceProviderDeleteAllFailureCode is the SAML IdP service provider delete all failure code.
	SAMLIdPServiceProviderDeleteAllFailureCode = "TSI004W"

	// OktaGroupsUpdateCode is the Okta groups updated code.
	OktaGroupsUpdateCode = "TOK001I"

	// OktaApplicationsUpdateCode is the Okta applications updated code.
	OktaApplicationsUpdateCode = "TOK002I"

	// OktaSyncFailureCode is the Okta synchronization failure code.
	OktaSyncFailureCode = "TOK003E"

	// OktaAssignmentProcessSuccessCode is the Okta assignment process success code.
	OktaAssignmentProcessSuccessCode = "TOK004I"

	// OktaAssignmentProcessFailureCode is the Okta assignment process failure code.
	OktaAssignmentProcessFailureCode = "TOK004E"

	// OktaAssignmentCleanupSuccessCode is the Okta assignment cleanup success code.
	OktaAssignmentCleanupSuccessCode = "TOK005I"

	// OktaAssignmentCleanupFailureCode is the Okta assignment cleanup failure code.
	OktaAssignmentCleanupFailureCode = "TOK005E"

	// AccessListCreateSuccessCode is the access list create success code.
	AccessListCreateSuccessCode = "TAL001I"

	// AccessListCreateFailureCode is the access list create failure code.
	AccessListCreateFailureCode = "TAL001E"

	// AccessListUpdateSuccessCode is the access list update success code.
	AccessListUpdateSuccessCode = "TAL002I"

	// AccessListUpdateFailureCode is the access list update failure code.
	AccessListUpdateFailureCode = "TAL002E"

	// AccessListDeleteSuccessCode is the access list delete success code.
	AccessListDeleteSuccessCode = "TAL003I"

	// AccessListDeleteFailureCode is the access list delete failure code.
	AccessListDeleteFailureCode = "TAL003E"

	// AccessListReviewSuccessCode is the access list review success code.
	AccessListReviewSuccessCode = "TAL004I"

	// AccessListReviewFailureCode is the access list review failure code.
	AccessListReviewFailureCode = "TAL004E"

	// AccessListMemberCreateSuccessCode is the access list member create success code.
	AccessListMemberCreateSuccessCode = "TAL005I"

	// AccessListMemberCreateFailureCode is the access list member create failure code.
	AccessListMemberCreateFailureCode = "TAL005E"

	// AccessListMemberUpdateSuccessCode is the access list member update success code.
	AccessListMemberUpdateSuccessCode = "TAL006I"

	// AccessListMemberUpdateFailureCode is the access list member update failure code.
	AccessListMemberUpdateFailureCode = "TAL006E"

	// AccessListMemberDeleteSuccessCode is the access list member delete success code.
	AccessListMemberDeleteSuccessCode = "TAL007I"

	// AccessListMemberDeleteFailureCode is the access list member delete failure code.
	AccessListMemberDeleteFailureCode = "TAL007E"

	// AccessListMemberDeleteAllForAccessListSuccessCode is the access list all member delete success code.
	AccessListMemberDeleteAllForAccessListSuccessCode = "TAL008I"

	// AccessListMemberDeleteAllForAccessListFailureCode is the access list member delete failure code.
	AccessListMemberDeleteAllForAccessListFailureCode = "TAL008E"

	// SecReportsAuditQueryRunCode is used when a custom Security Reports Query is run.
	SecReportsAuditQueryRunCode = "SRE001I"

	// SecReportsReportRunCode is used when a report in run.
	SecReportsReportRunCode = "SRE002I"

	// UnknownCode is used when an event of unknown type is encountered.
	UnknownCode = apievents.UnknownCode
)
