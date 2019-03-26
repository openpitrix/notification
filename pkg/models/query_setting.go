// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

var SearchWordColumnTable = []string{
	TableNotification,
	TableTask,
	TableAddress,
	TableAddressList,
}

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	TableNotification: {
		NfColId, NfColContentType, NfColOwner, NfColStatus,
	},
	TableTask: {
		TaskColTaskId, TaskColNfId, TaskColErrorCode, TaskColStatus,
	},
	TableAddress: {
		AddrColId, AddrColAddress, AddrColNotifyType, AddrColStatus,
	},
	TableAddressList: {
		AddrLsColId, AddrLsColName, AddrLsColExtra, AddrLsColStatus,
	},
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	TableNotification: {
		NfColId, NfColContentType, NfColTitle, NfColShortContent, NfColAddressInfo, NfColStatus, NfColOwner,
	},
	TableTask: {
		TaskColTaskId, TaskColNfId, TaskColStatus, TaskColErrorCode,
	},
	TableAddress: {
		AddrColId, AddrColAddress, AddrColNotifyType, AddrColStatus, AddrColRemarks,
	},
	TableAddressList: {
		AddrLsColId, AddrLsColName, AddrLsColName, AddrLsColExtra,
	},
}
