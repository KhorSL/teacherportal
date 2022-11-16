package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/khorsl/teacherportal/constants"
	mockdb "github.com/khorsl/teacherportal/db/mock"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/stretchr/testify/require"
)

func TestCommonStudentAPI(t *testing.T) {
	teacher1 := randomTeacher(t)
	teacher2 := randomTeacher(t)

	common1 := randomStudent(t)
	common2 := randomStudent(t)

	testCases := []struct {
		name          string
		queryParams   string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			queryParams: fmt.Sprintf("teacher=%s&teacher=%s", teacher1.Email, teacher2.Email),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetCommonStudentsEmailParams{
					Email: fmt.Sprintf("'%s','%s'", teacher1.Email, teacher2.Email),
					Count: int64(2),
				}
				store.EXPECT().GetCommonStudentsEmail(gomock.Any(), gomock.Eq(arg)).Times(1).
					Return([]string{common1.Email, common2.Email}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireCommonStudentBodyMatch(t, recorder.Body, []string{common1.Email, common2.Email})
			},
		},
		{
			name:        "TeacherEmailInvalidFormat",
			queryParams: fmt.Sprintf("teacher=%s&teacher=%s", teacher1.Email, constants.InvalidFormatEmail),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommonStudentsEmail(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "NoEmailProvided",
			queryParams: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetCommonStudentsEmail(gomock.Any(), gomock.Any()).Times(0)
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

			url := fmt.Sprintf("/api/commonstudents?%s", tc.queryParams)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireCommonStudentBodyMatch(t *testing.T, body *bytes.Buffer, expected []string) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse commonStudentsResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)
	require.Equal(t, expected, gotResponse.Students)
}
