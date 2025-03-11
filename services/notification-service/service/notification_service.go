package service

import (
    "fmt"
    "os"
    "time"

    "github.com/ankan8/swapsync/backend/services/notification-service/models"
    "github.com/ankan8/swapsync/backend/services/notification-service/repository"
    "github.com/google/uuid"

    // SendGrid
    sendgrid "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendNotification handles creating or sending a notification.
// Real implementation might send an email, SMS, or push message.
func SendNotification(userID, message, channel string) (bool, string, error) {
    if userID == "" || message == "" || channel == "" {
        return false, "", fmt.Errorf("invalid notification data")
    }

    // 1) Create a notification record
    notifID := uuid.NewString()
    notif := &models.Notification{
        NotificationID: notifID,
        UserID:         userID,
        Message:        message,
        Channel:        channel,
        Timestamp:      time.Now().Format(time.RFC3339),
    }
    err := repository.InsertNotification(notif)
    if err != nil {
        return false, "", fmt.Errorf("failed to insert notification: %v", err)
    }

    // 2) If channel == "EMAIL", send a real email via SendGrid
    if channel == "EMAIL" {
        // Option A: Directly pass user’s email from the calling service (skip DB lookups)
        // userEmail := "hardcoded@example.com"

        // Option B: Fetch user’s email from DB if you store user data
        userEmail, err := repository.GetUserEmail(userID)
        if err != nil {
            // If user not found or DB error, return an error
            return false, notifID, fmt.Errorf("failed to retrieve user email: %v", err)
        }

        from := mail.NewEmail("NexTrade", "chessankan@gmail.com") // Must be verified in SendGrid
        subject := "Your Trade Executed"
        to := mail.NewEmail("User", userEmail)
        plainTextContent := message
        htmlContent := fmt.Sprintf("<strong>%s</strong>", message)

        mailMessage := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

        client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
        response, sendErr := client.Send(mailMessage)
        if sendErr != nil {
            return false, notifID, fmt.Errorf("error sending email via SendGrid: %v", sendErr)
        }
        fmt.Printf("Email sent to %s, status code=%d\n", userEmail, response.StatusCode)

    } else {
        // For channels other than EMAIL (e.g. "SMS", "PUSH"), just log or implement other logic
        fmt.Printf("Notification (channel=%s) stored for user=%s: %s\n", channel, userID, message)
    }

    // 3) Return success
    return true, notifID, nil
}
