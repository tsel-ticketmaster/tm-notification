package mailer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tsel-ticketmaster/tm-notification/pkg/mailer/mocks"

	"github.com/tsel-ticketmaster/tm-notification/pkg/mailer"
)

func TestGomailConstructor_SuccessWithoutLogrusInjected(t *testing.T) {
	gomailDialerMock := &mocks.GomailDialer{}
	m := mailer.NewGomailAdapter(nil, "default-sender@mail.com", gomailDialerMock, true)

	assert.NotNil(t, m)

	gomailDialerMock.AssertExpectations(t)
}

func TestGomailAdapterSend_Success(t *testing.T) {
	t.Run("return success with given sender", func(t *testing.T) {
		gomailDialerMock := &mocks.GomailDialer{}
		gomailDialerMock.On("DialAndSend", mock.Anything).Return(nil)

		m := mailer.NewGomailAdapter(logrus.New(), "default-sender@mail.com", gomailDialerMock, true)

		msg := mailer.Message{
			From: "from@mail.com",
			To: []mailer.Recepient{
				{
					Name:    "Testing Testing 1",
					Address: "testing1@mail.com",
				},
			},
		}

		err := m.Send(context.Background(), msg)
		assert.NoError(t, err)

		gomailDialerMock.AssertExpectations(t)
	})

	t.Run("return success with default sender", func(t *testing.T) {
		gomailDialerMock := &mocks.GomailDialer{}
		gomailDialerMock.On("DialAndSend", mock.Anything).Return(nil)

		m := mailer.NewGomailAdapter(logrus.New(), "default-sender@mail.com", gomailDialerMock, true)

		msg := mailer.Message{
			To: []mailer.Recepient{
				{
					Name:    "Testing Testing 1",
					Address: "testing1@mail.com",
				},
			},
			CC: []mailer.Recepient{
				{
					Name:    "CC Testing Testing 1",
					Address: "cctesting1@mail.com",
				},
			},
			Subject: "test subject",
			MessageBody: mailer.MessageBody{
				ContentType: mailer.ContentTypePlaintext,
				Body:        []byte("Hallo test."),
			},
		}

		err := m.Send(context.TODO(), msg)
		assert.NoError(t, err)

		gomailDialerMock.AssertExpectations(t)
	})

	t.Run("return success with multiple recipients and carbon copy", func(t *testing.T) {
		gomailDialerMock := &mocks.GomailDialer{}
		gomailDialerMock.On("DialAndSend", mock.Anything).Return(nil)

		m := mailer.NewGomailAdapter(logrus.New(), "default-sender@mail.com", gomailDialerMock, true)

		msg := mailer.Message{
			From: "from@mail.com",
			To: []mailer.Recepient{
				{
					Name:    "To Testing Testing 1",
					Address: "totesting1@mail.com",
				},
				{
					Name:    "To Testing Testing 2",
					Address: "totesting2@mail.com",
				},
			},
			CC: []mailer.Recepient{
				{
					Name:    "CC Testing Testing 1",
					Address: "cctesting1@mail.com",
				},
				{
					Name:    "CC Testing Testing 2",
					Address: "cctesting2@mail.com",
				},
			},
		}

		err := m.Send(context.TODO(), msg)
		assert.NoError(t, err)

		gomailDialerMock.AssertExpectations(t)
	})
}

func TestGomailAdapterSend_Error_DialUp(t *testing.T) {
	gomailDialerMock := &mocks.GomailDialer{}
	gomailDialerMock.On("DialAndSend", mock.Anything).Return(fmt.Errorf("tcp: Timeout"))

	m := mailer.NewGomailAdapter(logrus.New(), "default-sender@mail.com", gomailDialerMock, true)

	msg := mailer.Message{
		From: "from@mail.com",
		To: []mailer.Recepient{
			{
				Name:    "Testing Testing 1",
				Address: "testing1@mail.com",
			},
		},
		CC: []mailer.Recepient{
			{
				Name:    "CC Testing Testing 1",
				Address: "cctesting1@mail.com",
			},
		},
		Subject: "test subject",
		MessageBody: mailer.MessageBody{
			ContentType: mailer.ContentTypePlaintext,
			Body:        []byte("Hallo test."),
		},
	}

	err := m.Send(context.TODO(), msg)
	assert.Error(t, err)

	gomailDialerMock.AssertExpectations(t)
}

func TestGomailAdapterSend_Error_NoMessage(t *testing.T) {
	gomailDialerMock := &mocks.GomailDialer{}

	m := mailer.NewGomailAdapter(logrus.New(), "default-sender@mail.com", gomailDialerMock, true)

	err := m.Send(context.TODO())

	assert.Equal(t, mailer.ErrNoMessage, err)

	gomailDialerMock.AssertExpectations(t)
}

func TestGomailAdapterSend_Error_NoRecipient(t *testing.T) {
	gomailDialerMock := &mocks.GomailDialer{}

	m := mailer.NewGomailAdapter(logrus.New(), "default-sender@mail.com", gomailDialerMock, true)

	msg := mailer.Message{
		CC: []mailer.Recepient{
			{
				Name:    "CC Testing Testing 1",
				Address: "cctesting1@mail.com",
			},
		},
		Subject: "test subject",
		MessageBody: mailer.MessageBody{
			ContentType: mailer.ContentTypePlaintext,
			Body:        []byte("Hallo test."),
		},
	}

	err := m.Send(context.TODO(), msg)

	assert.Equal(t, mailer.ErrNoRecipient, err)

	gomailDialerMock.AssertExpectations(t)
}
