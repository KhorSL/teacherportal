package api

import (
	"bytes"
	"database/sql"
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

func TestSuspendAPI(t *testing.T) {
	student := randomStudent(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"student": student.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Eq(student.Email)).Times(1).Return(student, nil)
				store.EXPECT().GetNotSuspendedStudentByEmail(gomock.Any(), gomock.Eq(student.Email)).Times(1).Return(student, nil)
				store.EXPECT().SuspendStudentByEmail(gomock.Any(), gomock.Eq(student.Email)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name: "InvalidStudentEmailFormat",
			body: gin.H{
				"student": constants.InvalidFormatEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetNotSuspendedStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().SuspendStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "StudentNotExist",
			body: gin.H{
				"student": constants.DoesNotExistEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Eq(constants.DoesNotExistEmail)).Times(1).Return(db.Student{}, sql.ErrNoRows)
				store.EXPECT().GetNotSuspendedStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().SuspendStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "StudentAlreadySuspended",
			body: gin.H{
				"student": constants.DoesNotExistEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Eq(constants.DoesNotExistEmail)).Times(1).Return(student, nil)
				store.EXPECT().GetNotSuspendedStudentByEmail(gomock.Any(), gomock.Any()).Times(1).Return(db.Student{}, sql.ErrNoRows)
				store.EXPECT().SuspendStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoEmail",
			body: gin.H{
				"student": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetNotSuspendedStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().SuspendStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoBody",
			body: gin.H{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetNotSuspendedStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().SuspendStudentByEmail(gomock.Any(), gomock.Any()).Times(0)
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

			url := "/api/suspend"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
