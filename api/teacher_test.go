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

func randomTeacher(t *testing.T) (teacher db.Teacher) {
	teacher = db.Teacher{
		ID:       util.RandomInt(1, 1000),
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}
	return
}

func TestCreateTeacherAPI(t *testing.T) {
	teacher := randomTeacher(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"full_name": teacher.FullName,
				"email":     teacher.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateTeacherParams{
					FullName: teacher.FullName,
					Email:    teacher.Email,
				}
				store.EXPECT().CreateTeacher(gomock.Any(), gomock.Eq(arg)).Times(1).Return(nil)
				store.EXPECT().GetTeacherByEmail(gomock.Any(), gomock.Eq(teacher.Email)).Times(1).Return(teacher, nil)
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

			url := "/api/teachers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
