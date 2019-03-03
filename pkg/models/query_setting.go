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
		NfColId, NfColContentType, NfColTitle, NfColAddressInfo, NfColStatus,
	},
	TableTask: {
		TaskColNfId, TaskColTaskId, TaskColStatus, TaskColErrorCode,
	},
	TableAddress: {
		AddrColId, AddrColAddress, AddrColRemarks,
	},
	TableAddressList: {
		AddrColId, AddrLsName, AddrLsExtra,
	},
}

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	TableNotification: {
		NfColId, NfColStatus, NfColContentType,
	},
	TableTask: {
		TaskColTaskId, TaskColStatus, TaskColNfId,
	},
	TableAddress: {
		AddrColId, AddrColStatus, AddrColNotifyType,
	},
	TableAddressList: {
		AddrLsColId, AddrLsStatus,
	},
}
