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
	"text/template"

	"github.com/alexjlockwood/gcm"
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

	cAndroidRetries                = 2
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
	androidClient       *gcm.Sender
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

	pushSrv.androidClient = &gcm.Sender{ApiKey: cAndroidServerToken}
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

	// construct the notifcation
	data := map[string]interface{}{"poll_id": 0, "message": msgText}
	regIDs := []string{deviceGUID}
	msg := gcm.NewMessage(data, regIDs...)

	// send the notification to the GCM server
	response, err := pushSrv.androidClient.Send(msg, cAndroidRetries)
	if err != nil && pushSrv.logger != nil {
		pushSrv.logger.Log(cPushServerTag, fmt.Sprintf("Android push error %s",
			err), "::1")
		return
	}

	// check for failures
	pushSrv.logger.Log(cPushServerTag, fmt.Sprint(response), "::1")
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
