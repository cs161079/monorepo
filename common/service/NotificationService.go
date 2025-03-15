package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"google.golang.org/api/fcm/v1"
	"google.golang.org/api/option"
)

type NotificationService interface {
	SendPushNotification(string, models.Notification) error
}

type notificationServiceImpl struct {
}

func NewNotificationService() NotificationService {
	return &notificationServiceImpl{}
}

func (s notificationServiceImpl) SendPushNotification(serviceAccountPath string, notEntity models.Notification) error {
	// Initialize a context and authenticate using the service account file
	ctx := context.Background()
	client, err := fcm.NewService(ctx, option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return fmt.Errorf("failed to create FCM client: %v", err)
	}

	if notEntity.Topic == "" {
		return errors.New("Topic is empty. As can be written notification.")
	}

	// Prepare the FCM message request
	message := &fcm.Message{
		Token: notEntity.Topic,
		Notification: &fcm.Notification{
			Title: notEntity.Title,
			Body:  notEntity.Message,
		},
	}

	// Send the message using the new HTTP v1 API
	respCall := client.Projects.Messages.Send("projects/bus-telematics", &fcm.SendMessageRequest{Message: message})
	resp, err := respCall.Do()

	if err != nil {
		return fmt.Errorf("Failed to send push notification: %v", err)
	}

	// Log the response
	logger.INFO(fmt.Sprintf("Successfully sent message: %s", resp.Name))
	return nil
}
