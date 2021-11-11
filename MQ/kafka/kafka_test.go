package kafka

import "testing"

func TestSyncMessage(t *testing.T) {
	var Address = []string{"127.0.0.1:9092"}
	SyncProducerWithMessage(Address, "test", "testMessage")
	SyncProducerWithMessage(Address, "test", "testMessage1")
	SyncProducerWithMessage(Address, "test", "testMessage2")
	SyncProducerWithMessage(Address, "test", "testMessage3")
	err := SyncProducerWithMessage(Address, "test", "testMessage4")
	if err != nil {
		t.Fatalf("ProducerErr:%v", err)
	}
	err = ConsumerWithMessage(Address, "test", 0, 0)
	if err != nil {
		t.Fatalf("ConsumerErr:%v", err)
	}
}

func TestAsyncMessage(t *testing.T) {
	var Address = []string{"127.0.0.1:9092"}
	err := AsyncProducerWithMessage(Address, "test2", "tt")
	if err != nil {
		t.Fatalf("ProducerErr:%v", err)
	}

	err = ConsumerWithMessage(Address, "test2", 1, 0)
	if err != nil {
		t.Fatalf("ConsumerErr:%v", err)
	}
}
