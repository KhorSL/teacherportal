package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/khorsl/teacherportal/db/mock"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/khorsl/teacherportal/util"
	"github.com/stretchr/testify/require"
)

func randomStudent(t *testing.T) (student db.Student) {
	student = db.Student{
		ID:       util.RandomInt(1, 1000),
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}
	return
}

func TestCreateStudentAPI(t *testing.T) {
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
				"full_name": student.FullName,
				"email":     student.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateStudentParams{
					FullName: student.FullName,
					Email:    student.Email,
				}
				store.EXPECT().CreateStudent(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil)
				store.EXPECT().GetStudentByEmail(gomock.Any(), gomock.Eq(student.Email)).Times(1).Return(student, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			url := "/api/students"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
