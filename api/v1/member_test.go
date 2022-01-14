package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var routerTest = getRouter()

func getRouter() *gin.Engine {
	intiServerWg := &sync.WaitGroup{}
	intiServerWg.Add(1)
	defer intiServerWg.Wait()
	return router.SetupRouter("v1", intiServerWg)
}

func TestRegister(t *testing.T) {
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/member/register", nil)
	routerTest.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "請傳送正確資料")
}

func TestLogin(t *testing.T) {
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ := http.NewRequest("POST", "/api/v1/member/login", nil)
	routerTest.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "請傳送正確資料")
}
