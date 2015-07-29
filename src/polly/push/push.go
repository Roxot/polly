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

	cPushClientTag    = "PUSHCLIENT"
	cPushServerLogFmt = "Failed to send notification %d: %s"

	cAndroidRetries                = 2
	cNotificationChannelBufferSize = 1 // TODO these in config files as well i guess

	TYPE_NEW_VOTE = 0 // TODO merge with VOTE types and such in model.go
	TYPE_UPVOTE   = 1
	TYPE_NEW_POLL = 2
	// TODO future types for added users, changed settting, etc. after APIs come
	// available for this
)

type sVoteTemplateData struct {
	Voter  string
	Option string
}

type NotificationData struct { // TODO to model or decapitalize
	deviceInfos []polly.DeviceInfo `json:"-"`
	//Message     string]
	Type   int    `json:"type"`
	User   string `json:"user"`
	Title  string `json:"title"`
	PollID int    `json:"poll_id"`
}

type IPushClient interface {
	StartErrorLogger(log.ILogger) error
	NotifyForVote(db *database.Database, user *polly.PrivateUser,
		optionID, voteType int) error
	NotifyForNewPoll(db *database.Database, user *polly.PrivateUser,
		pollID int, pollTitle string) error
}

type sPushClient struct {
	iosClient           apns.Client
	androidClient       gcm.Sender
	logger              log.ILogger
	notificationChannel chan *NotificationData
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
	var notificationData *NotificationData
	var numDevices int
	pushClient.notificationChannel = make(chan *NotificationData,
		cNotificationChannelBufferSize)

	go func() {
		for {
			notificationData = <-pushClient.notificationChannel
			numDevices = len(notificationData.deviceInfos)

			for i := 0; i < numDevices; i++ {
				if len(notificationData.deviceInfos[i].DeviceGUID) == 0 {
					fmt.Println("Skipping:", notificationData.deviceInfos[i].
						DeviceGUID)
					continue
				}

				if notificationData.deviceInfos[i].DeviceType == polly.
					DEVICE_TYPE_AD {

					fmt.Println("Notifying (and):",
						notificationData.deviceInfos[i].DeviceGUID)
					pushClient.sendAndroidNotification(
						notificationData.deviceInfos[i].DeviceGUID,
						notificationData)
				} else {
					fmt.Println("Notifying (ios):",
						notificationData.deviceInfos[i].DeviceGUID)
					pushClient.sendIosNotification(
						notificationData.deviceInfos[i].DeviceGUID,
						notificationData)
				}
			}
		}
	}()
}

func (pushClient *sPushClient) sendIosNotification(deviceGUID string,
	notificationData *NotificationData) {

	data, err := json.MarshalIndent(notificationData, "", "\t")
	if err != nil {
		pushClient.logger.Log(cPushClientTag, err.Error(), "::1")
		return
	}

	payload := apns.NewPayload()
	payload.APS.Alert.Body = string(data)
	notification := apns.NewNotification()
	notification.Payload = payload
	notification.DeviceToken = deviceGUID
	pushClient.iosClient.Send(notification)
}

func (pushClient *sPushClient) sendAndroidNotification(deviceGUID string,
	notificationData *NotificationData) {

	// construct the notifcation
	data := map[string]interface{}{"poll_id": notificationData.PollID,
		"type": notificationData.Type, "user": notificationData.User,
		"title": notificationData.Title}
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
	user *polly.PrivateUser, optionID, voteType int) error {
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
	notificationData := NotificationData{}
	notificationData.deviceInfos = deviceInfos
	notificationData.PollID = option.PollID
	notificationData.Type = voteType
	notificationData.User = user.DisplayName
	notificationData.Title = option.Value

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationData

	return nil
}

func (pushClient *sPushClient) NotifyForNewPoll(db *database.Database,
	user *polly.PrivateUser, pollID int, pollTitle string) error { // TODO public user?

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
	notificationData := NotificationData{}
	notificationData.deviceInfos = deviceInfos
	notificationData.PollID = pollID
	notificationData.Type = TYPE_NEW_POLL
	notificationData.User = user.DisplayName
	notificationData.Title = pollTitle

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationData

	return nil
}
