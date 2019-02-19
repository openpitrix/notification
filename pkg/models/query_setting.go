// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const (
	DefaultSelectLimit = 200
)

var SearchWordColumnTable = []string{
	TableNotification,
	TableTask,
}

// columns that can be search through sql 'like' operator
var SearchColumns = map[string][]string{
	TableNotification: {
		NfColId, NfColContentType, NfColTitle, NfColAddressInfo, NfColStatus,
	},
	TableTask: {
		TaskColNfId, TaskColTaskId, TaskColStatus, TaskColErrorCode,
	},
}

// columns that can be search through sql '=' operator
var IndexedColumns = map[string][]string{
	TableNotification: {
		NfColId, NfColStatus,
	},
	TableTask: {
		TaskColNfId, TaskColStatus,
	},
}
