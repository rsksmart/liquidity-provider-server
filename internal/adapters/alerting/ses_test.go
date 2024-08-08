package alerting_test

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	fromAddress = "from@email"
	toAddress   = "to@email"
	subject     = "any subject"
	body        = "any body"
)

func TestSesAlertSender_SendAlert(t *testing.T) {
	client := &mocks.SesClientMock{}
	client.On("SendEmail", test.AnyCtx, mock.MatchedBy(func(input *ses.SendEmailInput) bool {
		return assert.Equal(t, []string{toAddress}, input.Destination.ToAddresses) &&
			assert.Equal(t, body, *input.Message.Body.Text.Data) &&
			assert.Equal(t, subject, *input.Message.Subject.Data) &&
			assert.Equal(t, fromAddress, *input.Source)
	})).Return(&ses.SendEmailOutput{MessageId: aws.String("msgId")}, nil)

	sender := alerting.NewSesAlertSender(client, fromAddress)
	err := sender.SendAlert(context.Background(), subject, body, []string{toAddress})
	require.NoError(t, err)
	client.AssertExpectations(t)
}

func TestSesAlertSender_SendAlert_ErrorHandling(t *testing.T) {
	client := &mocks.SesClientMock{}

	client.On("SendEmail", test.AnyCtx, mock.Anything).Return(nil, assert.AnError)
	sender := alerting.NewSesAlertSender(client, fromAddress)
	err := sender.SendAlert(context.Background(), subject, body, []string{toAddress})
	require.Error(t, err)
	client.AssertExpectations(t)
}
