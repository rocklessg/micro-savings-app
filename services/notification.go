package services

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SendEmailNotification sends an email using SendGrid
func SendEmail(toEmail, subject, content string) error {
	from := mail.NewEmail("Your App Name", os.Getenv("SENDGRID_FROM_EMAIL"))
	to := mail.NewEmail("", toEmail)
	message := mail.NewSingleEmail(from, subject, to, content, content)
	
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("error sending email, status code: %v", response.StatusCode)
	}
	
	return nil
}

// SendSMSNotification sends an SMS using Twilio
func SendSMS(toPhone, message string) error {
	client := twilio.NewRestClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(toPhone)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	if resp.ErrorCode != nil {
		return fmt.Errorf("error sending SMS: %v", *resp.ErrorMessage)
	}

	return nil
}