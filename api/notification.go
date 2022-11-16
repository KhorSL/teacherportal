package api

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khorsl/teacherportal/constants"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/khorsl/teacherportal/util"
)

type retreiveForNoticationsRequest struct {
	Teacher      string `json:"teacher" binding:"required,email"`
	Notification string `json:"notification" binding:"required"`
}

type retreiveForNoticationsReseponse struct {
	Recipients []string `json:"recipients"`
}

func newRetreiveForNoticationsReseponse(recipients []string) retreiveForNoticationsReseponse {
	return retreiveForNoticationsReseponse{
		Recipients: recipients,
	}
}

func getMentionedEmails(message string) (emails []string) {
	re := regexp.MustCompile(constants.MentionRegex)

	splitMessage := strings.Split(message, " ")

	for _, value := range splitMessage {
		if re.FindStringIndex(value) != nil {
			emails = append(emails, value[1:])
		}
	}

	return emails
}

func (server *Server) retreiveStudentEmailsForNotifcations(ctx *gin.Context) {
	var req retreiveForNoticationsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	teacher, err := server.store.GetTeacherByEmail(ctx, req.Teacher)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponseCustom("teacher does not exist."))
		return
	}

	emails := getMentionedEmails(req.Notification)

	arg := db.GetStudentsEmailForNotificationParams{
		TeacherID:     teacher.ID,
		StudentEmails: util.StringsToSingleQuoteCommaSep(emails),
	}

	emailsForNotifcation, err := server.store.GetStudentsEmailForNotification(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newRetreiveForNoticationsReseponse(emailsForNotifcation)
	ctx.JSON(http.StatusOK, response)
}
