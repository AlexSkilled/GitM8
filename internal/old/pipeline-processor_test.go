package old

import (
	"encoding/json"
	"gitlab-tg-bot/internal/old/model"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestPipelineProcessor(t *testing.T) {
	processor := PipelineProcessor{}
	requests := make([]model.Pipeline, 10)
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

func Test_IfStatusIsUnknown_ShouldSkip(t *testing.T) {
	processor, request := PipelineProcessor{}, model.Pipeline{}
	request.ObjectAttributes.Status = "abcdefg"
	reqBytes, _ := json.Marshal(request)
	msg, skip, err := processor.Process(reqBytes)
	assert.Empty(t, msg, "Ожидается пустое сообщение")
	assert.True(t, skip, "Ожидается, что результат обработки должен быть пропущен")
	assert.NoError(t, err, "Ожидается отсутствие ошибки")
}