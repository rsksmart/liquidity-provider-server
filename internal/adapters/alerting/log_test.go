package alerting_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/alerting"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSubject    = "Test Alert Subject"
	testBody       = "Test alert body message"
	testRecipient1 = "admin@example.com"
	testRecipient2 = "support@example.com"
)

func TestLogAlertSender_SendAlert_SingleRecipient(t *testing.T) {
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert! - Subject: Test Alert Subject | Recipients: admin@example.com | Body: Test alert body message")

	err := sender.SendAlert(context.Background(), testSubject, testBody, []string{testRecipient1})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_MultipleRecipients(t *testing.T) {
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert! - Subject: Test Alert Subject | Recipients: admin@example.com, support@example.com | Body: Test alert body message")

	err := sender.SendAlert(context.Background(), testSubject, testBody, []string{testRecipient1, testRecipient2})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_EmptyRecipients(t *testing.T) {
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert! - Subject: Test Alert Subject | Recipients:  | Body: Test alert body message")

	err := sender.SendAlert(context.Background(), testSubject, testBody, []string{})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_NilRecipients(t *testing.T) {
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert! - Subject: Test Alert Subject | Recipients:  | Body: Test alert body message")

	err := sender.SendAlert(context.Background(), testSubject, testBody, nil)

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_EmptySubject(t *testing.T) {
	sender := alerting.NewLogAlertSender()

	err := sender.SendAlert(context.Background(), "", testBody, []string{testRecipient1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "alert subject cannot be empty")
}

func TestLogAlertSender_SendAlert_WhitespaceOnlySubject(t *testing.T) {
	sender := alerting.NewLogAlertSender()

	err := sender.SendAlert(context.Background(), "   \t\n   ", testBody, []string{testRecipient1})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "alert subject cannot be empty")
}

func TestLogAlertSender_SendAlert_ValidSubjectEmptyBody(t *testing.T) {
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert! - Subject: Test Alert Subject | Recipients: admin@example.com | Body:")

	err := sender.SendAlert(context.Background(), testSubject, "", []string{testRecipient1})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_SpecialCharacters(t *testing.T) {
	subjectWithSpecialChars := "Alert: System Error! (Critical)"
	bodyWithSpecialChars := "Error occurred: connection failed | status: 500 & timeout"
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert: System Error! (Critical)")

	err := sender.SendAlert(context.Background(), subjectWithSpecialChars, bodyWithSpecialChars, []string{testRecipient1})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_LongContent(t *testing.T) {
	longSubject := "Very long subject that contains a lot of text to test how the logger handles extended content"
	longBody := "This is a very long alert body that contains multiple sentences and a lot of information that might be present in a real alert message. It should be logged completely without truncation."
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, longSubject)

	err := sender.SendAlert(context.Background(), longSubject, longBody, []string{testRecipient1})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestLogAlertSender_SendAlert_JSONBody(t *testing.T) {
	jsonBody := `{"error": "connection_failed", "details": {"host": "localhost", "port": 8080}, "timestamp": "2023-12-01T10:30:00Z"}`
	sender := alerting.NewLogAlertSender()
	assertLogContains := test.AssertLogContains(t, "Alert! - Subject: Test Alert Subject")

	err := sender.SendAlert(context.Background(), testSubject, jsonBody, []string{testRecipient1})

	require.NoError(t, err)
	assert.True(t, assertLogContains())
}

func TestNewLogAlertSender(t *testing.T) {
	sender := alerting.NewLogAlertSender()
	assert.NotNil(t, sender)
}

func TestLogAlertSender_AlwaysSucceeds_WithValidSubject(t *testing.T) {
	sender := alerting.NewLogAlertSender()

	// Test multiple calls to ensure it always returns nil when subject is valid
	for i := 0; i < 5; i++ {
		err := sender.SendAlert(context.Background(), testSubject, testBody, []string{testRecipient1})
		require.NoError(t, err)
	}
}
