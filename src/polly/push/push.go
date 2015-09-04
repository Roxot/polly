package push

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"polly"
	"polly/database"
	"polly/log"

	"polly/internal/github.com/alexjlockwood/gcm"
	"polly/internal/github.com/timehop/apns"
)

const (
	cIOSGateway         = apns.SandboxGateway // TODO all this should be in a config file
	cCertDir            = "certs/"
	cIOSCertFile        = cCertDir + "apns-dev-cert.pem"
	cIOSKeyFile         = cCertDir + "apns-dev-key.key"
	cAndroidServerToken = "AIzaSyCi-zeWU_moOdFtUWggHXMulWGQK72wBuk"

	cPushClientTag         = "PUSHCLIENT"
	cPushServerLogFmt      = "Failed to send notification %d: %s"
	cIOSSilentNotification = 1

	cAndroidRetries                = 2
	cNotificationChannelBufferSize = 1 // TODO these in config files as well i guess
)

type IPushClient interface {
	StartErrorLogger(log.ILogger) error
	NotifyForVote(db *database.Database, user *polly.PrivateUser,
		optionID int64, voteType int) error
	NotifyForNewPoll(db *database.Database, user *polly.PrivateUser,
		pollID int64, pollTitle string) error
}

type sPushClient struct {
	iosClient           apns.Client
	androidClient       gcm.Sender
	logger              log.ILogger
	notificationChannel chan *polly.NotificationMessage
}

func NewClient() (IPushClient, error) {
	var pushClient = sPushClient{}

	// get the polly home directory
	pollyHome, err := polly.GetPollyHome()
	if err != nil {
		return nil, err
	}

	// load ios cert
	cert, err := tls.LoadX509KeyPair(pollyHome+cIOSCertFile,
		pollyHome+cIOSKeyFile)
	if err != nil {
		return nil, err
	}

	// create ios client
	iosClient := apns.NewClientWithCert(cIOSGateway, cert)
	pushClient.iosClient = iosClient
	pushClient.androidClient = gcm.Sender{ApiKey: cAndroidServerToken}
	pushClient.startNotificationHandling()

	return &pushClient, nil
}

func (pushClient *sPushClient) StartErrorLogger(logger log.ILogger) error {
	if logger == nil {
		return errors.New("Logger may not be nil.")
	}

	go func() {
		for failures := range pushClient.iosClient.FailedNotifs {
			logger.Log(cPushClientTag, fmt.Sprintf(cPushServerLogFmt,
				failures.Notif.ID, failures.Err.Error()), "::1")
		}
	}()

	pushClient.logger = logger

	return nil
}

func (pushClient *sPushClient) startNotificationHandling() {
	var notificationMsg *polly.NotificationMessage
	var numDevices int
	pushClient.notificationChannel = make(chan *polly.NotificationMessage,
		cNotificationChannelBufferSize)

	go func() {
		for {
			notificationMsg = <-pushClient.notificationChannel
			numDevices = len(notificationMsg.DeviceInfos)

			for i := 0; i < numDevices; i++ {
				if len(notificationMsg.DeviceInfos[i].DeviceGUID) == 0 {
					fmt.Println("Skipping:", notificationMsg.DeviceInfos[i].
						DeviceGUID)
					continue
				}

				if notificationMsg.DeviceInfos[i].DeviceType == polly.
					DEVICE_TYPE_ANDROID {

					fmt.Println("Notifying (and):",
						notificationMsg.DeviceInfos[i].DeviceGUID)
					pushClient.sendAndroidNotification(
						notificationMsg.DeviceInfos[i].DeviceGUID,
						notificationMsg)
				} else {
					fmt.Println("Notifying (ios):",
						notificationMsg.DeviceInfos[i].DeviceGUID)
					pushClient.sendIosNotification(
						notificationMsg.DeviceInfos[i].DeviceGUID,
						notificationMsg)
				}
			}
		}
	}()
}

func (pushClient *sPushClient) sendIosNotification(deviceGUID string,
	notificationMsg *polly.NotificationMessage) {

	data, err := json.MarshalIndent(notificationMsg, "", "\t")
	if err != nil {
		pushClient.logger.Log(cPushClientTag, err.Error(), "::1")
		return
	}

	payload := apns.NewPayload()
	payload.APS.ContentAvailable = cIOSSilentNotification
	payload.SetCustomValue(polly.NOTIFICATION_INFO_FIELD, string(data))
	notification := apns.NewNotification()
	notification.Payload = payload
	notification.DeviceToken = deviceGUID
	pushClient.iosClient.Send(notification)
}

func (pushClient *sPushClient) sendAndroidNotification(deviceGUID string,
	notificationMsg *polly.NotificationMessage) {

	// construct the notifcation
	data := map[string]interface{}{"poll_id": notificationMsg.PollID,
		"type": notificationMsg.Type, "user": notificationMsg.User,
		"title": notificationMsg.Title}
	regIDs := []string{deviceGUID}
	msg := gcm.NewMessage(data, regIDs...)

	// send the notification to the GCM server
	response, err := pushClient.androidClient.Send(msg, cAndroidRetries)
	if err != nil && pushClient.logger != nil {
		pushClient.logger.Log(cPushClientTag, fmt.Sprintf(
			"Android push error %s", err), "::1")
		return
	}

	// check for failures
	pushClient.logger.Log(cPushClientTag, fmt.Sprint(response), "::1")
}

func (pushClient *sPushClient) NotifyForVote(db *database.Database,
	user *polly.PrivateUser, optionID int64, voteType int) error {
	// TODO user->voter, PrivateUser->PublicUser

	// TODO assert votetype

	// retrieve option name
	option, err := db.GetOptionByID(optionID)
	if err != nil {
		return err
	}

	// retrieve all poll participants
	deviceInfos, err := db.GetDeviceInfosForPollExcludeCreator(option.PollID,
		user.ID)
	if err != nil {
		return err
	}

	// don't notify for empty polls
	if len(deviceInfos) == 0 {
		return nil
	}

	// prepare notification
	notificationMsg := polly.NotificationMessage{}
	notificationMsg.DeviceInfos = deviceInfos
	notificationMsg.PollID = option.PollID
	notificationMsg.Type = voteType
	notificationMsg.User = user.DisplayName
	notificationMsg.Title = option.Value

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg

	return nil
}

func (pushClient *sPushClient) NotifyForNewPoll(db *database.Database,
	user *polly.PrivateUser, pollID int64, pollTitle string) error { // TODO public user?

	// retrieve all poll participants
	deviceInfos, err := db.GetDeviceInfosForPollExcludeCreator(pollID,
		user.ID)
	if err != nil {
		return err
	}

	// don't notify for empty polls
	if len(deviceInfos) == 0 {
		return nil
	}

	// prepare notification
	notificationMsg := polly.NotificationMessage{}
	notificationMsg.DeviceInfos = deviceInfos
	notificationMsg.PollID = pollID
	notificationMsg.Type = polly.EVENT_TYPE_NEW_POLL
	notificationMsg.User = user.DisplayName
	notificationMsg.Title = pollTitle

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg

	return nil
}
