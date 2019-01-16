// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"testing"
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/idutil"
)

func TestNewService(t *testing.T) {
	nfservice := NewService()
	logger.Infof(nil, "nfservice=%+v", nfservice)
}

func TestCreateNfWaddrs(t *testing.T) {
	nfservice := NewService()
	testAddrsStr := "johuo@yunify.com;513590612@qq.com"
	nf := &models.Notification{
		NotificationId: idutil.GetUuid(constants.NfPostIDPrifix),
		ContentType:    "Email",
		AddrsStr:       testAddrsStr,
		Title:          "Title Test2",
		Content:        "Content2",
		ShortContent:   "ShortContent2",
		ExporedDays:    5,
		Owner:          "Huojiao",
		Status:         "New",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	q := globalcfg.GetInstance().GetEtcd().NewQueue(constants.EmailQueue)
	nfId, err := nfservice.CreateNfWithAddrs(nf, q)
	if err != nil {
		logger.Criticalf(nil, "failed to TestCreateNfWaddrs, error: [%+v]", err)
	}
	logger.Debugf(nil, "success to TestCreateNfWaddrs, nfId:%+s", nfId)

}
