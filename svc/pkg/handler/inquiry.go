package handler

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/gin-gonic/gin"
	"kitakyusyu-hackathon/pkg/sendgrid"
	"kitakyusyu-hackathon/pkg/slack"
	"kitakyusyu-hackathon/svc/pkg/gas"
	"kitakyusyu-hackathon/svc/pkg/schema"
	"kitakyusyu-hackathon/svc/pkg/uc"
	"log"
	"time"
)

type InquiryHandler struct {
	slackClient *slack.Slack
	inviteUC    *uc.InviteSlack
	sendgrid    *sendgrid.Sendgrid
	fs          *firestore.Client
	g           *gas.GAS
}

func NewInquiryHandler(fs *firestore.Client) *InquiryHandler {
	s := slack.NewSlack()
	g := gas.NewGAS()
	return &InquiryHandler{
		slackClient: &s,
		inviteUC:    uc.NewInviteSlack(s),
		sendgrid:    sendgrid.NewSendgrid(),
		fs:          fs,
		g:           &g,
	}
}

func (h *InquiryHandler) HandleInquiry() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("inquiry request")
		data := schema.InquiryData{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			log.Printf("failed to bind json, err: %v\n", err)
			return
		}

		if err := data.Validate(); err != nil {
			c.JSON(400, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			log.Printf("validation error: %v\n", err)
			return
		}

		log.Printf("inquiry data: %+v\n", data)

		_, _, err := h.fs.Collection("inquiries").Add(c, data)
		if err != nil {
			log.Printf("failed to add inquiry to firestore, err: %v\n", err)
		}

		var slackChannelURL string
		if data.UseSlack {
			slackChannelURL = h.handleSlack(data)
		} else {
			h.handleMail(data)
		}

		go h.g.PostData(gas.InquiryData{
			Firstname:       data.Firstname,
			Lastname:        data.Lastname,
			CompanyName:     data.CompanyName,
			EmailAddress:    data.EmailAddress,
			Purpose:         data.Purpose,
			InquiryDetails:  data.InquiryDetails,
			UseSlack:        data.UseSlack,
			SlackChannelURL: slackChannelURL,
		})

		log.Printf("inquery process succeeded\n")
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}
}

func (h InquiryHandler) handleSlack(data schema.InquiryData) string {
	var guests []uc.GuestInfo
	if *data.SlackInfo != nil {
		guests = make([]uc.GuestInfo, 0, len(*data.SlackInfo)+1)
		for _, s := range *data.SlackInfo {
			guests = append(guests, uc.GuestInfo{
				Email:     s.Email,
				Firstname: s.Firstname,
				Lastname:  s.Lastname,
			})
		}
	}
	guests = append(guests, uc.GuestInfo{
		Email:     data.EmailAddress,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
	})

	timeStr := time.Now().Format("200601021504")
	inviteInput := uc.InviteSlackInput{
		ChannelName: fmt.Sprintf("%s-%s様", timeStr, data.CompanyName),
		StaffIDs:    []string{"U04936U1UEB"},
		GuestInfo:   guests,
	}
	inviteResult, err := h.inviteUC.Do(inviteInput)
	if err != nil {
		log.Printf("failed to invite slack, err: %v\n", err)
	}
	if inviteResult == nil {
		return ""
	}
	log.Printf("channel created: %s\n", inviteResult.ChannelName)
	log.Printf("channel link: %s\n", inviteResult.ChannelLink)
	return inviteResult.ChannelLink
}

func (h InquiryHandler) handleMail(data schema.InquiryData) {
	h.sendgrid.SendMailNotify(fmt.Sprintf("%s %s", data.Firstname, data.Lastname), data.EmailAddress)
}
