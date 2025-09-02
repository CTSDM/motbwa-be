package ws

// import (
// 	"bytes"
// 	"encoding/json"
// 	"sync"
// 	"testing"
// 	"time"
//
// 	"github.com/gorilla/websocket"
// )
//
// type mockConn struct {
// 	messages []Event
// 	countOut int
// 	dropConn []bool
// 	closed   bool
// 	errors   []error
// }
//
// func (mc *mockConn) ReadMessage() (messageType int, p []byte, err error) {
// 	ticker := time.NewTicker(time.Millisecond * 1)
// 	for {
// 		<-ticker.C
// 		if !mc.closed {
// 			mc.countOut += 1
// 			mc.closed = mc.dropConn[mc.countOut]
// 			data, err := json.Marshal(mc.messages[mc.countOut-1])
// 			if err != nil {
// 				return -1, []byte{}, err
// 			}
// 			return websocket.TextMessage, data, nil
// 		} else {
// 			return -1, []byte{}, nil
// 		}
// 	}
// }
//
// func (mc *mockConn) WriteMessage(messageType int, data []byte) error {
// 	return nil
// }
//
// func (mc *mockConn) Close() error {
// 	mc.closed = true
// 	return nil
// }
//
// func TestReadMessages(t *testing.T) {
// 	type casesTest struct {
// 		input    []Event
// 		want     []Event
// 		dropConn []bool
// 		err      error
// 	}
//
// 	t.Run("Should only receive messages when the websocket connection is open", func(t *testing.T) {
// 		cases := []casesTest{
// 			{
// 				input: []Event{
// 					{Type: "send_message", Payload: []byte("hello")},
// 					{Type: "send_message", Payload: []byte("there!")},
// 				},
// 				want: []Event{
// 					{Type: "send_message", Payload: []byte("hello")},
// 					{},
// 				},
// 				dropConn: []bool{false, true},
// 				err:      nil,
// 			},
// 			{
// 				input:    []Event{{Type: "send_message", Payload: []byte("hello")}},
// 				want:     []Event{{}},
// 				dropConn: []bool{true},
// 				err:      nil,
// 			},
// 			{
// 				input: []Event{
// 					{Type: "send_message", Payload: []byte("hi!")},
// 					{Type: "send_message", Payload: []byte("long time no see")},
// 					{Type: "send_message", Payload: []byte("hello?")},
// 				},
// 				want: []Event{
// 					{Type: "send_message", Payload: []byte("hi!")},
// 					{Type: "send_message", Payload: []byte("long time no see")},
// 					{},
// 				},
// 				dropConn: []bool{false, false, true},
// 				err:      nil,
// 			},
// 		}
//
// 		for _, c := range cases {
// 			// setup the mock connection
// 			var messages []Event
// 			var dropConn []bool
// 			for i := range c.input {
// 				messages = append(messages, c.input[i])
// 				dropConn = append(dropConn, c.dropConn[i])
// 			}
// 			newConn := mockConn{
// 				countOut: 0,
// 				dropConn: dropConn,
// 				messages: messages,
// 				closed:   dropConn[0],
// 			}
//
// 			client := NewClient(&newConn)
// 			// call the func in a go routine
// 			go func() {
// 				err := client.ReadMessages()
// 				if err != nil {
// 					t.Errorf("Got err: %v, want nil", err)
// 				}
// 			}()
// 			// loop until all data is collected
// 			var wg sync.WaitGroup
// 			for i := range c.input {
// 				wg.Add(1)
// 				got := <-client.data
// 				if !bytes.Equal(got.Payload, c.want[i].Payload) {
// 					t.Errorf("For run %v, Got %q, want %q", i, string(got.Payload), string(c.want[i].Payload))
// 				}
// 				if got.Type != c.want[i].Type {
// 					t.Errorf("For run %v, Got %q, want %q", i, got.Type, c.want[i].Type)
// 				}
// 				wg.Done()
// 			}
// 			wg.Wait()
// 		}
// 	})
// }
