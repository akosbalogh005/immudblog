package restapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"immudblog/config"
	"immudblog/immudb"
	immudb_mock "immudblog/immudb/mocks"
	"immudblog/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	immuDBMock *immudb_mock.ImmuDBIF
)

func setup(t *testing.T) {
	immuDBMock = &immudb_mock.ImmuDBIF{}
	//immuDBMock.AssertExpectations(t)
	immudb.ImmuDB = immuDBMock
}

func getAPIResult(t *testing.T, method string, path string, auth bool, ret interface{}, toSend interface{}) (w *httptest.ResponseRecorder, req *http.Request) {
	config.ServerFlags.AuthUsers = "user:user:write"
	router := SetupRouter()
	w = httptest.NewRecorder()
	var err error
	if toSend != nil {
		jsonData, err := json.Marshal(toSend)
		assert.NoError(t, err)
		req, err = http.NewRequest(method, path, bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if auth {
		base64.StdEncoding.EncodeToString(([]byte)("user:user"))
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(([]byte)("user:user")))
	}
	assert.NoError(t, err)
	router.ServeHTTP(w, req)
	if ret != nil {
		err = json.Unmarshal(w.Body.Bytes(), ret)
		assert.NoError(t, err)
	}
	return
}

func TestNotAuth(t *testing.T) {
	w, _ := getAPIResult(t, "GET", "/api/v1/logs/count", false, nil, nil)
	assert.Equal(t, 401, w.Code)
}

func TestGetCountDBError(t *testing.T) {
	setup(t)
	var ret uint64 = 12
	immuDBMock.On("CountLogs").Return(ret, fmt.Errorf("Error")).Once()
	retJson := model.APIResponse{}
	w, _ := getAPIResult(t, "GET", "/api/v1/logs/count", true, &retJson, nil)
	assert.Equal(t, 500, w.Code)
	assert.Equal(t, 500, retJson.Code)
}

func TestGetCountOK(t *testing.T) {
	setup(t)
	var ret uint64 = 12
	immuDBMock.On("CountLogs").Return(ret, nil).Once()
	retJson := model.GetLogsCountResponse{}
	w, _ := getAPIResult(t, "GET", "/api/v1/logs/count", true, &retJson, nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, ret, retJson.Count)
}

func TestGetLogs(t *testing.T) {
	setup(t)
	var ret []model.Log
	ret = append(ret, model.Log{ID: 111})
	immuDBMock.On("GetLogs", uint64(20), "").Return(ret, nil).Once()
	var retJson []model.Log
	w, _ := getAPIResult(t, "GET", "/api/v1/logs?count=20", true, &retJson, nil)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 1, len(retJson))
	assert.Equal(t, int64(111), retJson[0].ID)
}

func TestAddLogs(t *testing.T) {
	setup(t)
	var ret []model.Log
	ret = append(ret, model.Log{ID: 111, Hostname: "testhostname"})
	immuDBMock.On("AddLogs", ret).Return(nil).Once()
	retJson := model.APIResponse{}
	w, _ := getAPIResult(t, "POST", "/api/v1/logs", true, &retJson, ret)
	assert.Equal(t, 200, w.Code)
}
