package models

import (
	"testing"

	"openpitrix.io/logger"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
)

func TestDecodeNotificationExtra(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()

	testExtra := "{\"ws_service\": \"op\",\"ws_message_type\": \"event\"}"

	nfExtraMap, err := DecodeNotificationExtra(testExtra)
	if err != nil {
		logger.Errorf(nil, "error=[%+v]", err)

	}

	service := ""
	nfService, ok := (*nfExtraMap)[constants.WsService]
	if ok {
		service = nfService
	}
	logger.Debugf(nil, "service=[%s]", service)

	messageType := ""
	nfExtraType, ok := (*nfExtraMap)[constants.WsMessageType]
	if ok {
		messageType = nfExtraType
	}
	logger.Debugf(nil, "messageType=[%s]", messageType)

}

func TestCheckExtra(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	testExtra := "{\"ws_service1\": \"op\",\"ws_message_type\": \"event\"}"

	err := CheckExtra(nil, testExtra)
	if err != nil {
		logger.Errorf(nil, "error=[%+v]", err)

	}

}
