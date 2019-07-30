package models

import (
	"testing"

	"openpitrix.io/logger"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/util/jsonutil"
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
	config.GetInstance()
	testExtra := "{\"ws_service1\": \"op\",\"ws_message_type\": \"event\"}"

	err := CheckExtra(nil, testExtra)
	if err != nil {
		logger.Errorf(nil, "error=[%+v]", err)

	}

}

func TestUseMsgStringToPb(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	dataStr := `{"user_id":"huojiao","service":"ks","message_type":"event","MessageDetail":{"ws_message_id":"msg-XXy43Kkkl95V","ws_user_id":"huojiao","ws_service":"ks","ws_message_type":"event","ws_message":"test_content_normal"}}`
	userMsg := new(UserMessage)
	err := jsonutil.Decode([]byte(dataStr), userMsg)
	if err != nil {
		logger.Errorf(nil, "Decode [%s] into UserMessage failed: %+v", dataStr, err)
	}
	pbUserMsg := UserMessageToPb(userMsg)

	logger.Infof(nil, "pbUserMsg=%+v", pbUserMsg)

}
