package push

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"polly"
	"polly/database"
	"polly/log"
	"text/template"

	"github.com/alexjlockwood/gcm"
	"github.com/timehop/apns"
)

const (
	cIOSGateway         = apns.SandboxGateway
	cCertDir            = "certs/"
	cIOSCertFile        = cCertDir + "apns-dev-cert.pem"
	cIOSKeyFile         = cCertDir + "apns-dev-key.key"
	cVoteTemplateText   = "{{.Voter}} voted for {{.Option}}."
	cAndroidServerToken = "AIzaSyCi-zeWU_moOdFtUWggHXMulWGQK72wBuk"

	cPushClientTag    = "PUSHCLIENT"
	cPushServerLogFmt = "Failed to send notification %d: %s"

	cAndroidRetries                = 2
	cNotificationChannelBufferSize = 1

	TYPE_NEW_VOTE = 0
	TYPE_UPVOTE   = 1
	TYPE_NEW_POLL = 2
)

type sVoteTemplateData struct {
	Voter  string
	Option string
}

type NotificationData struct {
	deviceInfos []polly.DeviceInfo
	//Message     string]
	Type   int    `json:"type"`
	User   string `json:"user"`
	Title  string `json:"title"`
	PollID int    `json:"poll_id"`
}

type IPushClient interface {
	StartErrorLogger(*logger.Logger) error
	NotifyForUpvote(*database.Database, *polly.PrivateUser, int) error
}

type sPushClient struct {
	iosClient           apns.Client
	androidClient       gcm.Sender
	logger              log.Logger
	notificationChannel chan *NotificationData
	voteTemplate        template.Template
}

func NewPushCLient() (IPushClient, error) {
	var pushClient = sPushClient{}

	// create ios client
	pollyHome, err := polly.GetPollyHome()
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(pollyHome+cIOSCertFile,
		pollyHome+cIOSKeyFile)
	if err != nil {
		return nil, err
	}

	iosClient := apns.NewClientWithCert(cIOSGateway, cert)
	pushClient.iosClient = iosClient

	template := template.New("VoteTemplate")
	template, err = template.Parse(cVoteTemplateText)
	if err != nil {
		return nil, err
	}

	pushClient.androidClient = gcm.Sender{ApiKey: cAndroidServerToken}
	pushClient.voteTemplate = *template
	pushClient.startNotificationHandling()

	return &pushClient, nil
}

func (pushClient *sPushClient) StartErrorLogger(logger *logger.Logger) error {
	if logger == nil {
		return errors.New("Logger may not be nil.")
	}

	go func() {
		for failures := range pushClient.iosClient.FailedNotifs {
			logger.Log(cPushServerTag, fmt.Sprintf(cPushServerLogFmt,
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
					log.Println("Skipping:", notificationData.deviceInfos[i].
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
					log.Println("Notifying (ios):",
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
		pushClient.logger.Log(cPushServerTag, err.Error(), "::1")
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
		pushClient.logger.Log(cPushServerTag, fmt.Sprintf(
			"Android push error %s", err), "::1")
		return
	}

	// check for failures
	pushClient.logger.Log(cPushServerTag, fmt.Sprint(response), "::1")
}

func (pushClient *sPushClient) NotifyForUpvote(db *database.Database,
	user *polly.PrivateUser, optionID int) error {

	// retrieve option name
	option, err := db.GetOptionById(optionID)
	if err != nil {
		return err
	}

	// retrieve all poll participants
	deviceInfos, err := db.DeviceInfosForPollExcludeCreator(option.PollId,
		user.Id)
	if err != nil {
		return err
	}

	// don't notify for empty polls
	if len(deviceInfos) == 0 {
		return nil
	}

	// prepare notification
	// var buffer bytes.Buffer
	// templateData := voteTemplateData{}
	// templateData.Option = option.Value
	// templateData.Voter = user.DisplayName
	// pushClient.voteTemplate.Execute(&buffer, templateData)
	notificationData := NotificationData{}
	// notificationData.Message = buffer.String()
	notificationData.deviceInfos = deviceInfos
	notificationData.PollID = option.PollId
	notificationData.Type = TYPE_UPVOTE
	notificationData.User = user.DisplayName
	notificationData.Title = option.Value

	// let the notification handler goroutine take care of the rest
	pushClient.notificationChannel <- &notificationData

	return nil
}
