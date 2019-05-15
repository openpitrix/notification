// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/util/jsonutil"
)

type AddressInfo map[string][]string

type AddressListIds []string

func DecodeAddressInfo(data string) (*AddressInfo, error) {
	addressInfo := new(AddressInfo)
	err := jsonutil.Decode([]byte(data), addressInfo)
	if err != nil {
		logger.Errorf(nil, "Decode [%s] into address info failed: %+v", data, err)
	}
	return addressInfo, err
}

func DecodeAddressListIds(data string) (*AddressListIds, error) {
	addressListIds := new(AddressListIds)
	err := jsonutil.Decode([]byte(data), addressListIds)
	if err != nil {
		logger.Errorf(nil, "Decode [%s] into address list ids failed: %+v", data, err)
	}
	return addressListIds, err
}
