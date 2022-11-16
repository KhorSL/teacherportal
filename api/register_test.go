package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/khorsl/teacherportal/constants"
	mockdb "github.com/khorsl/teacherportal/db/mock"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestRegisterAPI(t *testing.T) {
	teacher := randomTeacher(t)
	student1 := randomStudent(t)
	student2 := randomStudent(t)

	students := []string{student1.Email, student2.Email}

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"teacher":  teacher.Email,
				"students": students,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.RegisterTxParams{
					Teacher:  teacher.Email,
					Students: students,
				}
				store.EXPECT().RegisterTx(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "InvalidStudentEmailFormat",
			body: gin.H{
				"teacher":  teacher.Email,
				"students": []string{student1.Email, constants.InvalidFormatEmail},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().RegisterTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTeacherEmailFormat",
			body: gin.H{
				"teacher":  constants.InvalidFormatEmail,
				"students": students,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().RegisterTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
