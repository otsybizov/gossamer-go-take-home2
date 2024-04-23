package network_test

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/gossamer-go-interview/network"
	"github.com/stretchr/testify/assert"
)

func generateMessage(n int) *network.Message {
	return &network.Message{
		ID:     fmt.Sprintf("someID%d", n),
		PeerID: fmt.Sprintf("somePeerID%d", n),
		Data:   []byte{0, 1, 1},
	}
}

func TestMessageTracker_Add(t *testing.T) {
	t.Run("add, get, then all messages", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)
	})

	t.Run("add, get, then all messages, delete some", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)

		for i := 0; i < length-2; i++ {
			err := mt.Delete(generateMessage(i).ID)
			assert.NoError(t, err)
		}

		msgs = mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(3),
			generateMessage(4),
		}, msgs)

	})

	t.Run("not full, with duplicates", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(length - 2))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})

	t.Run("not full, with duplicates from other peers", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			msg := generateMessage(length - 2)
			msg.PeerID = "somePeerID0"
			err := mt.Add(msg)
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})
}

func TestMessageTracker_Cleanup(t *testing.T) {
	t.Run("overflow and cleanup", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})

	t.Run("overflow and cleanup with duplicate", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		for i := length; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})
}

func TestMessageTracker_Delete(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		err := mt.Delete("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
	})
}

func TestMessageTracker_Message(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		msg, err := mt.Message("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
		assert.Nil(t, msg)
	})
}

// Benchmarks

func createMessagesAndTracker(length int, fillMessageTracker bool, b *testing.B) ([]*network.Message, network.MessageTracker) {
	msgs := make([]*network.Message, 0, length)
	for i := 0; i < length; i++ {
		msgs = append(msgs, generateMessage(i))
	}

	mt := network.NewMessageTracker(length)
	if fillMessageTracker {
		for i := 0; i < length; i++ {
			err := mt.Add(msgs[i])
			assert.NoError(b, err)
		}
	}

	return msgs, mt
}

func incrementIndex(index *int, length int) {
	*index++
	if *index == length {
		*index = 0
	}
}

func BenchmarkMessageTracker(b *testing.B) {
	var err error
	length := 100000

	b.Run("add", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, false, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = mt.Add(msgs[index])
			assert.NoError(b, err)

			incrementIndex(&index, length)
		}
	})

	b.Run("add and get one", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, false, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = mt.Add(msgs[index])
			assert.NoError(b, err)

			msg, err := mt.Message(msgs[index].ID)
			assert.NoError(b, err)
			assert.NotNil(b, msg)

			incrementIndex(&index, length)
		}
	})

	b.Run("add and get all", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, false, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = mt.Add(msgs[index])
			assert.NoError(b, err)

			_ = mt.Messages()

			incrementIndex(&index, length)
		}
	})

	b.Run("remove and add", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, true, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = mt.Delete(msgs[index].ID)
			assert.NoError(b, err)

			err = mt.Add(msgs[index])
			assert.NoError(b, err)

			incrementIndex(&index, length)
		}
	})

	b.Run("remove, add and get one", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, true, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = mt.Delete(msgs[index].ID)
			assert.NoError(b, err)

			err = mt.Add(msgs[index])
			assert.NoError(b, err)

			msg, err := mt.Message(msgs[index].ID)
			assert.NoError(b, err)
			assert.NotNil(b, msg)

			incrementIndex(&index, length)
		}
	})

	b.Run("remove, add and get all", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, true, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err = mt.Delete(msgs[index].ID)
			assert.NoError(b, err)

			err = mt.Add(msgs[index])
			assert.NoError(b, err)

			_ = mt.Messages()

			incrementIndex(&index, length)
		}
	})

	b.Run("remove ignoring errors", func(b *testing.B) {
		msgs, mt := createMessagesAndTracker(length, true, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = mt.Delete(msgs[index].ID)

			incrementIndex(&index, length)
		}
	})

	b.Run("get one message", func(b *testing.B) {
		var msg *network.Message
		msgs, mt := createMessagesAndTracker(length, true, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			msg, err = mt.Message(msgs[index].ID)
			assert.NoError(b, err)
			assert.NotNil(b, msg)

			incrementIndex(&index, length)
		}
	})

	b.Run("get all messages", func(b *testing.B) {
		_, mt := createMessagesAndTracker(length, true, b)
		index := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			msgs := mt.Messages()
			assert.Equal(b, length, len(msgs))

			incrementIndex(&index, length)
		}
	})
}
