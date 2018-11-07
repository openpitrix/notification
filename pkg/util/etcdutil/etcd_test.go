package etcdutil

import (
	"fmt"
	"log"
	"openpitrix.io/openpitrix/pkg/etcd"
	"testing"
)

func TestConnect(t *testing.T) {
	//e := new(Etcd)
	endpoints := []string{"192.168.0.7:2379"}
	//endpoints:=[]string{"192.168.0.3:2379"}
	prefix := "test"
	e, err := etcd.Connect(endpoints, prefix)
	log.Println(e)
	if err != nil {
		t.Fatal(err)
	}

}

func TestNewQueue(t *testing.T) {
	//endpoints:=[]string{"192.168.0.7:2379,192.168.0.8:2379,192.168.0.6:2379"}
	endpoints := []string{"192.168.0.7:2379"}
	prefix := "test"
	e, err := etcd.Connect(endpoints, prefix)
	log.Println(e)
	if err != nil {
		t.Fatal(err)
	}

	q := e.NewQueue("notification")
	q.Enqueue("ssss")
}

func TestEnqueue(t *testing.T) {
	endpoints := []string{"192.168.0.7:2379"}
	prefix := "test"
	e, err := etcd.Connect(endpoints, prefix)
	if err != nil {
		t.Fatal(err)
	}
	queue := e.NewQueue("notification")
	go func() {
		for i := 0; i < 100; i++ {
			err := queue.Enqueue(fmt.Sprintf("%d", i))
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Push message to queue, worker number [%d]", i)
		}

	}()
	for i := 0; i < 100; i++ {
		n, err := queue.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Got message [%s] from queue, worker number [%d]", n, i)
	}
}


func TestEnqueue2(t *testing.T) {
	endpoints := []string{"192.168.0.7:2379"}
	prefix := "nf_"
	e, err := etcd.Connect(endpoints, prefix)
	if err != nil {
		t.Fatal(err)
	}
	queue := e.NewQueue("nf_")
	//go func() {
	//	for i := 0; i < 100; i++ {
	//		err := queue.Enqueue(fmt.Sprintf("%d", i))
	//		if err != nil {
	//			t.Fatal(err)
	//		}
	//		t.Logf("Push message to queue, worker number [%d]", i)
	//	}
	//
	//}()
	for i := 0; i < 100; i++ {
		n, err := queue.Dequeue()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Got message [%s] from queue, worker number [%d]", n, i)
	}
}
