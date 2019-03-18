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

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	TableNotification: {
		NfColContentType, NfColTitle, NfColAddressInfo, NfColStatus, NfColOwner,
	},
	TableTask: {
		TaskColTaskId, TaskColStatus, TaskColErrorCode,
	},
	TableAddress: {
		AddrColAddress, AddrColRemarks,
	},
	TableAddressList: {
		AddrLsColName, AddrLsColExtra,
	},
}

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	TableNotification: {
		NfColId, NfColStatus, NfColContentType, NfColOwner,
	},
	TableTask: {
		TaskColTaskId, TaskColStatus, TaskColNfId,
	},
	TableAddress: {
		AddrColId, AddrColStatus, AddrColNotifyType,
	},
	TableAddressList: {
		AddrLsColId, AddrLsColStatus,
	},
}
