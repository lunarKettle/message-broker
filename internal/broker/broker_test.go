package broker

import (
	"sync"
	"testing"
)

func TestBroker_Publish(t *testing.T) {
	t.Parallel()

	const testTopic = "testTopic"
	const testMessage = "testMessage"

	broker := NewMessageBroker()
	broker.Start()

	subChannel := make(chan string)
	subCmd := &Command{
		ClientID: "testSubID",
		Action:   "SUBSCRIBE",
		Queue:    testTopic,
		Response: subChannel,
	}
	broker.ExecuteCommand(subCmd)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		response := <-subChannel
		if response != "OK" {
			t.Error("Expected OK, but got ", response)
		}
		message := <-subChannel
		if message != testMessage {
			t.Errorf("Expected message %s, but got %s", testMessage, message)
		}
	}()

	responseCh := make(chan string)
	publishCmd := &Command{
		ClientID: "testPubID",
		Action:   "PUBLISH",
		Queue:    testTopic,
		Message:  testMessage,
		Response: responseCh,
	}

	broker.ExecuteCommand(publishCmd)
	response := <-responseCh
	if response != "OK" {
		t.Error("Expected OK, got ", response)
	}

	if len(broker.queues[testTopic]) == 0 {
		t.Error("Expected message to be added to the queue")
	}
	wg.Wait()
}

func TestBroker_Subscribe(t *testing.T) {
	t.Parallel()
	const testTopic = "testTopic"

	broker := NewMessageBroker()
	broker.Start()

	subChannel := make(chan string)
	subCmd := &Command{
		ClientID: "testSubID",
		Action:   "SUBSCRIBE",
		Queue:    testTopic,
		Response: subChannel,
	}
	broker.ExecuteCommand(subCmd)

	response := <-subChannel
	if response != "OK" {
		t.Error("Expected OK, but got ", response)
	}

	if len(broker.subs[testTopic]) == 0 {
		t.Error("Expected sub to be added to the subs")
	}
}

func TestBroker_Consume(t *testing.T) {
	t.Parallel()

	const testTopic = "testTopic"
	const testMessage = "testMessage"
	const emptyTopicMessage = "Empty topic"

	broker := NewMessageBroker()
	broker.Start()

	consumerCh1 := make(chan string)
	consumeCmd1 := &Command{
		ClientID: "testConsumerID",
		Action:   "CONSUME",
		Queue:    testTopic,
		Response: consumerCh1,
	}
	broker.ExecuteCommand(consumeCmd1)

	responseBeforePub := <-consumerCh1
	if responseBeforePub != emptyTopicMessage {
		t.Errorf("Expected message %s, but got %s", emptyTopicMessage, responseBeforePub)
	}

	pubCh := make(chan string)
	publishCmd := &Command{
		ClientID: "testPubID",
		Action:   "PUBLISH",
		Queue:    testTopic,
		Message:  testMessage,
		Response: pubCh,
	}
	broker.ExecuteCommand(publishCmd)
	pubResponse := <-pubCh
	if pubResponse != "OK" {
		t.Error("Expected OK, got ", pubResponse)
	}

	consumerCh2 := make(chan string)
	consumeCmd2 := &Command{
		ClientID: "testConsumerID",
		Action:   "CONSUME",
		Queue:    testTopic,
		Response: consumerCh2,
	}
	broker.ExecuteCommand(consumeCmd2)

	responseAfterPub := <-consumerCh2
	if responseAfterPub != testMessage {
		t.Errorf("Expected message %s, but got %s", testMessage, responseAfterPub)
	}
}
