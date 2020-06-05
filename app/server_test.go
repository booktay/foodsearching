package main

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestInput1IDOutDocument(t *testing.T) {
	router := gin.Default()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews/0", nil)
	router.ServeHTTP(w, req)

	template = `{
		"Message": "ReviewID is not found"
	}`

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, template, w.Body.String())
}

func TestInput2IDOutDocument(t *testing.T) {
	router := gin.Default()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews/6204", nil)
	router.ServeHTTP(w, req)

	template = `{
		"Message": "ReviewID is not found"
	}`

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, template, w.Body.String())
}

func TestInputIDInDocument(t *testing.T) {
	router := gin.Default()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews/10", nil)
	router.ServeHTTP(w, req)

	template = `{
		"_id": "10",
		"_index": "reviews",
		"_score": 7.81682,
		"_source": {
			"created": 1591378632031689000,
			"modified": 1591385439459338000,
			"reviewid": "10",
			"reviewtext": "ดีใจที่จุดพักรถมอร์เตอร์เวย์ขาเข้า มีร้านอาหารจริงๆเสียที ขับรถมาเหนื่อยๆ อยากทานอาหารในร้านอาหารที่มีนั่งสบาย บริการแบบร้านอาหารจริงๆ ต้องที่นี่ค่ะ\nเมนูขึ้นชื่อของร้านมัลลิการ์ ก็ต้องเย็นตาโฟ ทั้งแบบใจเสาะ (เผ็ดปานกลาง) และรสเจ็บ (เผ็ดจริงจัง) แต่กว่าจะรู้สึกเผ็ดผ่านไปครึ่งชาม 555 รสเข้มข้น ไม่ต้องปรุงเลย เส้นใหญ่เหนียวนุ่ม ทานกับเต้าหู้ทอด หอม กรอบนอกนุ่มในอร่อยมากค่ะ"
		},
		"_type": "_doc"
	}`

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, template, w.Body.String())
}