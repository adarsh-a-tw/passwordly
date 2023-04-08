package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func PrepareContextAndResponseRecorder(t *testing.T, url string, method string, reqBody any) (ctx *gin.Context, rec *httptest.ResponseRecorder) {
	rec = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(rec)

	var req *http.Request

	if reqBody != nil {
		jsonBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err) // json.Marshal error
		buffer := bytes.NewBuffer(jsonBytes)
		req = httptest.NewRequest(method, url, buffer)
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	ctx.Request = req
	return
}

func DecodeJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, obj any) {
	if err := json.NewDecoder(rec.Body).Decode(obj); err != nil {
		t.Failed()
	}
}