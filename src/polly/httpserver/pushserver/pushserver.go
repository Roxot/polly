package pushserver

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"polly"
	"polly/database"
	"polly/logger"
	"polly/push/android"
	"text/template"

	"github.com/timehop/apns"
)

const (
	cIosGateway         = apns.SandboxGateway
	cCertDir            = "certs/"
	cIosCertFile        = cCertDir + "apns-dev-cert.pem"
	cIosKeyFile         = cCertDir + "apns-dev-key.key"
	cVoteTemplateText   = "{{.Voter}} voted for {{.Option}}."
	cAndroidServerToken = "AIzaSyCi-zeWU_moOdFtUWggHXMulWGQK72wBuk"

	cPushServerTag    = "PUSHSERVER"
	cPushServerLogFmt = "Failed to send notification %d: %s"

	cNotificationChannelBufferSize = 1
)

type voteTemplateData struct {
	Voter  string
	Option string
}

type notificationData struct {
	Message     string
	DeviceInfos []polly.DeviceInfo
}

type PushServer struct {
	iosClient           *apns.Client
	androidClient       *android.Client
	logger              *logger.Logger
	notificationChannel chan *notificationData
	voteTemplate        *template.Template
}

func New() (*PushServer, error) {
	var pushSrv = PushServer{}

	// create ios client
	pollyHome, err := polly.GetPollyHome()
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(pollyHome+cIosCertFile,
		pollyHome+cIosKeyFile)
	if err != nil {
		return nil, err
	}

	iosClient := apns.NewClientWithCert(cIosGateway, cert)
	pushSrv.iosClient = &iosClient

	template := template.New("VoteTemplate")
	template, err = template.Parse(cVoteTemplateText)
	if err != nil {
		return nil, err
	}

	pushSrv.androidClient = android.New(cAndroidServerToken)
	pushSrv.voteTemplate = template
	pushSrv.startNotificationHandling()

	return &pushSrv, nil
}

func (pushSrv *PushServer) StartErrorLogger(logger *logger.Logger) error {
	if logger == nil {
		return errors.New("Logger may not be nil.")
	}

	go func() {
		for failures := range pushSrv.iosClient.FailedNotifs {
			logger.Log(cPushServerTag, fmt.Sprintf(cPushServerLogFmt,
				failures.Notif.ID, failures.Err.Error()), "::1")
		}
	}()

	pushSrv.logger = logger

	return nil
}

func (pushSrv *PushServer) startNotificationHandling() {
	var notData *notificationData
	var numDevices int
	pushSrv.notificationChannel = make(chan *notificationData,
		cNotificationChannelBufferSize)

	go func() {
		for {
			notData = <-pushSrv.notificationChannel
			numDevices = len(notData.DeviceInfos)

			for i := 0; i < numDevices; i++ {
				if len(notData.DeviceInfos[i].DeviceGUID) == 0 {
					log.Println("Skipping:", notData.DeviceInfos[i].DeviceGUID)
					continue
				}

				if notData.DeviceInfos[i].DeviceType == polly.DEVICE_TYPE_AD {
					fmt.Println("Notifying (and):",
						notData.DeviceInfos[i].DeviceGUID)
					pushSrv.sendAndroidNotification(
						notData.DeviceInfos[i].DeviceGUID, notData.Message)
				} else {
					log.Println("Notifying (ios):",
						notData.DeviceInfos[i].DeviceGUID)
					pushSrv.sendIosNotification(
						notData.DeviceInfos[i].DeviceGUID, notData.Message)
				}
			}
		}
	}()
}

func (pushSrv *PushServer) sendIosNotification(deviceGUID, msg string) {
	payload := apns.NewPayload()
	payload.APS.Alert.Body = msg
	notif := apns.NewNotification()
	notif.Payload = payload
	notif.DeviceToken = deviceGUID
	pushSrv.iosClient.Send(notif)
}

func (pushSrv *PushServer) sendAndroidNotification(deviceGUID, msgText string) {
	msg := android.NewMessage(deviceGUID)
	msg.SetPayload("message", msgText)
	resp, err := pushSrv.androidClient.Send(msg)
	if err != nil && pushSrv.logger != nil {
		pushSrv.logger.Log(cPushServerTag, err.Error(), "::1")
	}

	errorIndices := resp.ErrorIndexes()
	if len(errorIndices) != 0 {
		pushSrv.logger.Log(cPushServerTag, "error indices not nil", "::1")
	}

	refreshIndices := resp.RefreshIndexes()
	if len(refreshIndices) != 0 {
		pushSrv.logger.Log(cPushServerTag, "refresh indices not nil", "::1")
	}
}

func (pushSrv *PushServer) NotifyForUpvote(db *database.Database,
	usr *polly.PrivateUser, optionID int) error {

	// retrieve option name
	option, err := db.OptionById(optionID)
	if err != nil {
		return err
	}

	// retrieve all poll participants
	dvcInfos, err := db.DeviceInfosForPollExcludeCreator(option.PollId, usr.Id)
	if err != nil {
		return err
	}

	// don't notify for empty polls
	if len(dvcInfos) == 0 {
		return nil
	}

	// prepare notification
	var buffer bytes.Buffer
	templateData := voteTemplateData{}
	templateData.Option = option.Value
	templateData.Voter = usr.DisplayName
	pushSrv.voteTemplate.Execute(&buffer, templateData)
	notifData := notificationData{}
	notifData.Message = buffer.String()
	notifData.DeviceInfos = dvcInfos

	// let the notification handler goroutine take care of the rest
	pushSrv.notificationChannel <- &notifData

	return nil
}
