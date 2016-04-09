package push

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/roxot/polly"
	"github.com/roxot/polly/database"
	"github.com/roxot/polly/log"

	"github.com/timehop/apns"
	"github.com/yogyrahmawan/gcm"
)

const (
	cIOSGateway         = apns.SandboxGateway // TODO all this should be in a config file
	cCertDir            = "cert/"
	cIOSCertFile        = cCertDir + "apns-dev-cert.pem"
	cIOSKeyFile         = cCertDir + "apns-dev-key.key"
	cAndroidServerToken = "AIzaSyCi-zeWU_moOdFtUWggHXMulWGQK72wBuk"

	cPushClientTag         = "PUSHCLIENT"
	cPushServerLogFmt      = "Failed to send notification %s: %s"
	cIOSSilentNotification = 1

	cAndroidRetries                = 2
	cNotificationChannelBufferSize = 1 // TODO these in config files as well i guess
)

type IPushClient interface {
	StartErrorLogger(log.ILogger) error
	NotifyForVote(db *database.Database, user *polly.PrivateUser,
		optionTitle string, pollID int64, voteType int) error
	NotifyForNewPoll(db *database.Database, user *polly.PrivateUser,
		pollID int64, pollTitle string) error
	NotifyForClosedEvent(db *database.Database, pollID int64,
		title string) error
	NotifyForUndoneVote(db *database.Database, user *polly.PrivateUser,
		optionTitle string, pollID int64) error
	NotifyForParticipantLeft(db *database.Database, user *polly.PrivateUser,
		pollID int64, pollTitle string) error
	NotifyForNewParticipant(db *database.Database, creator *polly.PrivateUser,
		pollID int64, pollTitle string, newUser *polly.PrivateUser) error
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
	msg.Priority = gcm.HighPriority

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

func (pushClient *sPushClient) NotifyForClosedEvent(db *database.Database,
	pollID int64, title string) error {

	// retrieve all poll participants
	deviceInfos, err := db.GetDeviceInfosForPoll(pollID)
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
	notificationMsg.Type = polly.EVENT_TYPE_POLL_CLOSED
	notificationMsg.Title = title

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg

	return nil
}

func (pushClient *sPushClient) NotifyForVote(db *database.Database,
	user *polly.PrivateUser, optionTitle string, pollID int64, voteType int) error {
	// TODO user->voter, PrivateUser->PublicUser

	// TODO assert votetype

	// retrieve all poll participants
	deviceInfos, err := db.GetDeviceInfosForPollExcludeCreator(pollID, user.ID)
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
	notificationMsg.Type = voteType
	notificationMsg.User = user.DisplayName
	notificationMsg.UserID = user.ID
	notificationMsg.Title = optionTitle

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg

	return nil
}

func (pushClient *sPushClient) NotifyForUndoneVote(db *database.Database,
	user *polly.PrivateUser, optionTitle string, pollID int64) error {
	// TODO user->voter, PrivateUser->PublicUser

	// TODO assert votetype

	// retrieve all poll participants
	deviceInfos, err := db.GetDeviceInfosForPollExcludeCreator(pollID, user.ID)
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
	notificationMsg.Type = polly.EVENT_TYPE_UNDONE_VOTE
	notificationMsg.User = user.DisplayName
	notificationMsg.UserID = user.ID
	notificationMsg.Title = optionTitle

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
	notificationMsg.UserID = user.ID
	notificationMsg.Title = pollTitle

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg

	return nil
}

func (pushClient *sPushClient) NotifyForParticipantLeft(db *database.Database,
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
	notificationMsg.Type = polly.EVENT_TYPE_PARTICIPANT_LEFT
	notificationMsg.User = user.DisplayName
	notificationMsg.UserID = user.ID
	notificationMsg.Title = pollTitle

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg

	return nil
}

func (pushClient *sPushClient) NotifyForNewParticipant(db *database.Database,
	creator *polly.PrivateUser, pollID int64, pollTitle string,
	newUser *polly.PrivateUser) error {

	// retrieve all existing poll participants device infos
	deviceInfos, err := db.GetDeviceInfosForPollExcludeCreatorAndUser(pollID,
		creator.ID, newUser.ID)
	if err != nil {
		return err
	}

	// retrieve the new user's device info
	newUserDeviceInfo, err := db.GetDeviceInfoForUser(newUser.ID)
	if err != nil {
		return err
	}

	notificationMsg1 := polly.NotificationMessage{}
	notificationMsg1.DeviceInfos = deviceInfos
	notificationMsg1.PollID = pollID
	notificationMsg1.Type = polly.EVENT_TYPE_NEW_PARTICIPANT
	notificationMsg1.User = newUser.DisplayName
	notificationMsg1.UserID = newUser.ID
	notificationMsg1.Title = pollTitle

	notificationMsg2 := polly.NotificationMessage{}
	notificationMsg2.DeviceInfos = []polly.DeviceInfo{*newUserDeviceInfo}
	notificationMsg2.PollID = pollID
	notificationMsg2.Type = polly.EVENT_TYPE_ADDED_TO_POLL
	notificationMsg2.User = creator.DisplayName
	notificationMsg1.UserID = creator.ID
	notificationMsg2.Title = pollTitle

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationMsg1
	pushClient.notificationChannel <- &notificationMsg2

	return nil
}
