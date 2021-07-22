package old

import (
	"bytes"
	"encoding/json"
	config "gitlab-tg-bot/conf"
	"gitlab-tg-bot/internal/old/model"
	moq2 "gitlab-tg-bot/internal/old/moq"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/mock"
)

func Test_ShouldHandleMergeRequest(t *testing.T) {
	conf := moq2.Config{}
	conf.On("GetBool", config.NoAuth).Return(true)
	notifier := moq2.Notifier{}
	notifier.On("Notify", mock.Anything).Return(nil)
	handler := NewHandler(conf, notifier)
	srv := httptest.NewServer(handler)
	payload := model.MergeRequest{}
	err := faker.FakeData(&payload)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	plBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", srv.URL, ioutil.NopCloser(bytes.NewReader(plBytes)))
	req.Header.Add(EventHeaderKey, MergeRequestHeader)
	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Ожидается status 200 OK")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
}

func Test_IfMRStateIsUnknown_ShouldNotSendNotification(t *testing.T) {
	conf := moq2.Config{}
	conf.On("GetBool", config.NoAuth).Return(true)
	notifier := moq2.Notifier{}
	notifier.On("Notify", mock.Anything).Return(nil)
	handler := NewHandler(conf, notifier)
	srv := httptest.NewServer(handler)
	payload := model.MergeRequest{}
	payload.ObjectAttributes.State = "unknown"
	plBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", srv.URL, ioutil.NopCloser(bytes.NewReader(plBytes)))
	req.Header.Add(EventHeaderKey, MergeRequestHeader)
	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Ожидается status 200 OK")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
	notifier.AssertNotCalled(t, "Notify")
}

func Test_ShouldHandlePipeline(t *testing.T) {
	conf := moq2.Config{}
	conf.On("GetBool", config.NoAuth).Return(true)
	notifier := moq2.Notifier{}
	notifier.On("Notify", mock.Anything).Return(nil)
	handler := NewHandler(conf, notifier)
	srv := httptest.NewServer(handler)
	payload := model.Pipeline{}
	err := faker.FakeData(&payload)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	plBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", srv.URL, ioutil.NopCloser(bytes.NewReader(plBytes)))
	req.Header.Add(EventHeaderKey, PipelineHeader)
	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Ожидается status 200 OK")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
}

func Test_IfPipelineStatusIsUnknown_ShouldNotSendNotification(t *testing.T) {
	conf := moq2.Config{}
	conf.On("GetBool", config.NoAuth).Return(true)
	notifier := moq2.Notifier{}
	notifier.On("Notify", mock.Anything).Return(nil)
	handler := NewHandler(conf, notifier)
	srv := httptest.NewServer(handler)
	payload := model.Pipeline{}
	payload.ObjectAttributes.Status = "unknown"
	plBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", srv.URL, ioutil.NopCloser(bytes.NewReader(plBytes)))
	req.Header.Add(EventHeaderKey, PipelineHeader)
	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Ожидается status 200 OK")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
	notifier.AssertNotCalled(t, "Notify")
}

func Test_IfTokenIsNotValid_ShouldReturnStatus403(t *testing.T) {
	conf := moq2.Config{}
	conf.On("GetBool", config.NoAuth).Return(false)
	conf.On("GetString", config.SecretKey).Return("my_secret_key")
	notifier := moq2.Notifier{}
	notifier.On("Notify", mock.Anything).Return(nil)
	handler := NewHandler(conf, notifier)
	srv := httptest.NewServer(handler)
	req, _ := http.NewRequest("POST", srv.URL, bytes.NewReader([]byte("{}")))
	req.Header.Add(TokenHeaderKey, "not a valid token")
	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusForbidden, "Ожидается status 403 Forbidden")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
	notifier.AssertNotCalled(t, "Notify")
}

func Test_IfTokenIsValid_ShouldReturnStatus200(t *testing.T) {
	conf := moq2.Config{}
	conf.On("GetBool", config.NoAuth).Return(false)
	conf.On("GetString", config.SecretKey).Return("my_secret_key")
	notifier := moq2.Notifier{}
	notifier.On("Notify", mock.Anything).Return(nil)
	handler := NewHandler(conf, notifier)
	srv := httptest.NewServer(handler)
	req, _ := http.NewRequest("POST", srv.URL, bytes.NewReader([]byte("{}")))
	req.Header.Add(TokenHeaderKey, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.p10yaiAwCVW5_sASd6qj8QeWzaQBeMSeUaV_w22JzSM")
	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Ожидается status 200 OK")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
	notifier.AssertNotCalled(t, "Notify")
}
