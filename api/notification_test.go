package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func TestNotificationAPI(t *testing.T) {
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
				"teacher":      teacher.Email,
				"notification": fmt.Sprintf("Hello world @%s @%s beep", student1.Email, student2.Email),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherByEmail(gomock.Any(), gomock.Eq(teacher.Email)).Times(1).Return(teacher, nil)

				arg := db.GetStudentsEmailForNotificationParams{
					TeacherID:     teacher.ID,
					StudentEmails: fmt.Sprintf("'%s','%s'", student1.Email, student2.Email),
				}
				store.EXPECT().GetStudentsEmailForNotification(gomock.Any(), gomock.Eq(arg)).Times(1).
					Return(students, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireNotificationBodyMatch(t, recorder.Body, students)
			},
		},
		{
			name: "TeacherDoesNotExist",
			body: gin.H{
				"teacher":      constants.DoesNotExistEmail,
				"notification": fmt.Sprintf("Hello world @%s @%s beep", student1.Email, student2.Email),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherByEmail(gomock.Any(), gomock.Eq(constants.DoesNotExistEmail)).Times(1).Return(db.Teacher{}, sql.ErrNoRows)
				store.EXPECT().GetStudentsEmailForNotification(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"teacher":      constants.InvalidFormatEmail,
				"notification": fmt.Sprintf("Hello world @%s @%s beep", student1.Email, student2.Email),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetStudentsEmailForNotification(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NoStudents",
			body: gin.H{
				"teacher":      teacher.Email,
				"notification": "Hello world beep",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetTeacherByEmail(gomock.Any(), gomock.Eq(teacher.Email)).Times(1).Return(teacher, nil)

				arg := db.GetStudentsEmailForNotificationParams{
					TeacherID:     teacher.ID,
					StudentEmails: "''",
				}
				store.EXPECT().GetStudentsEmailForNotification(gomock.Any(), gomock.Eq(arg)).Times(1).Return([]string{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireNotificationBodyMatch(t, recorder.Body, []string{})
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

			url := "/api/retrievefornotifications"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireNotificationBodyMatch(t *testing.T, body *bytes.Buffer, expected []string) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse retreiveForNoticationsReseponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, expected, gotResponse.Recipients)
}
