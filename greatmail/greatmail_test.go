package greatmail

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type httpProviderMock struct {
	DoFn func(*http.Request) (*http.Response, error)
}

func (m httpProviderMock) Do(request *http.Request) (*http.Response, error) {
	return m.DoFn(request)
}

func TestSendEmail(t *testing.T) {

	t.Run("should return an error when fail to send email", func(t *testing.T) {
		client := Client{
			http: httpProviderMock{
				DoFn: func(*http.Request) (*http.Response, error) {
					return &http.Response{}, nil
				},
			},
		}

		err := client.SendEmail(Email{})
		if err != nil {
			t.Error("invalid error received, expeted a nil error")
		}
	})

	t.Run("should send an email with a correct body", func(t *testing.T) {
		var receivedRequest *http.Request

		client := Client{
			http: httpProviderMock{
				DoFn: func(request *http.Request) (*http.Response, error) {
					receivedRequest = request
					return &http.Response{StatusCode: 200}, nil
				},
			},
		}
		msg, _ := json.Marshal(createFakeEmail())
		expectedBody := ioutil.NopCloser(bytes.NewReader(msg))
		client.SendEmail(createFakeEmail())
		assert.Equal(t, expectedBody, receivedRequest.Body, "wrong email sent")
	})

	t.Run("should fail in the second sent", func(t *testing.T) {

		count := -1
		errs := []error{errors.New("failed"), nil}

		client := Client{
			http: httpProviderMock{
				DoFn: func(request *http.Request) (*http.Response, error) {
					count++
					return &http.Response{}, errs[count]
				},
			},
		}

		expectedErrs := []error{errors.New("failed"), nil}

		for _, err := range expectedErrs {
			outputErr := client.SendEmail(createFakeEmail())
			assert.Equal(t, err, outputErr, "wrong error returned")
		}
	})
}

func createFakeEmail() Email {
	return Email{
		Message: "message",
		Subject: "subject",
		From:    "someone",
		To:      []string{"foo", "bar"},
	}
}

func TestSendMail(t *testing.T) {

	t.Run("should send an email correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockProvider := NewMockHTTPProvider(ctrl)
		sendTest := Client{http: mockProvider}

		msg, _ := json.Marshal(createFakeEmail())
		request, _ := http.NewRequest(
			"POST",
			"https://api.greatmail.com/send",
			ioutil.NopCloser(bytes.NewReader(msg)),
		)

		mockProvider.EXPECT().Do(request).Return(&http.Response{StatusCode: 200}, nil).Times(1)
		sendTest.SendEmail(createFakeEmail())
	})

}

// func TestMailSends(t *testing.T) {

// 	t.Run("should send an email correctly", func(t *testing.T) {
// 		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
// 			expedtedBody, _ := json.Marshal(createFakeEmail())
// 			body, _ := ioutil.ReadAll(req.Body)

// 			assert.Equal(t, expedtedBody, body, "wrong email received")
// 			rw.Write([]byte(`ok`))
// 		}))
// 		defer server.Close()

// 		client := Client{
// 			http:    &http.Client{},
// 			baseURL: server.URL,
// 		}

// 		client.SendEmail(createFakeEmail())
// 	})
// }
