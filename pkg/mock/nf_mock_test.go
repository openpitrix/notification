/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package mock

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	notification "notificationService/pkg/nf/pb"
	nfmock "notificationService/pkg/nf/mock/mockgen"
)

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestSayHello(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockNotificationClient := nfmock.NewMockNotificationClient(ctrl)
	req := &notification.HelloRequest{Name: "unit_test"}
	mockNotificationClient.EXPECT().SayHello(
		gomock.Any(),
		&rpcMsg{msg: req},
	).Return(&notification.HelloReply{Message: "Mocked Interface"}, nil)
	testSayHello(t, mockNotificationClient)
}

func testSayHello(t *testing.T, client notification.NotificationClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SayHello(ctx, &notification.HelloRequest{Name: "unit_test"})
	t.Log("Test function : SayHello")
	if err != nil || r.Message != "Mocked Interface" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.Message)
	t.Log("Test Pass : ", r.Message)
}


