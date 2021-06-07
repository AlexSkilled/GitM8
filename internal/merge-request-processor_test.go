package internal

import (
	"encoding/json"
	"gitlab-tg-bot/internal/model"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestMergeRequestProcessor(t *testing.T) {
	processor := MergeRequestProcessor{}
	requests := make([]model.MergeRequest, 10)
	err := faker.FakeData(&requests)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}
	for _, request := range requests {
		reqBytes, _ := json.Marshal(request)
		msg, skip, err := processor.Process(reqBytes)
		assert.NotEmpty(t, msg, "Ожидается не пустое сообщение")
		assert.False(t, skip, "Ожидается, что результат обработки не должен быть пропущен")
		assert.NoError(t, err, "Ожидается отсутствие ошибки")
	}
}

func Test_IfStateIsUnknown_ShouldSkip(t *testing.T) {
	processor, request := MergeRequestProcessor{}, model.MergeRequest{}
	request.ObjectAttributes.State = "abcdefg"
	reqBytes, _ := json.Marshal(request)
	msg, skip, err := processor.Process(reqBytes)
	assert.Empty(t, msg, "Ожидается пустое сообщение")
	assert.True(t, skip, "Ожидается, что результат обработки должен быть пропущен")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
}
