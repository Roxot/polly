package http

import (
    "encoding/json"
    "net/http"
    "github.com/roxot/polly"

    "github.com/julienschmidt/httprouter"
    "github.com/satori/go.uuid"
)

const (
    cRegisterTag                  = "POST/REGISTER"
    cServiceProviderHeader        = "X-Auth-Service-Provider"
    cClientAuthorizationHeader    = "X-Verify-Credentials-Authorization"
    cOAuthEchoAuthorizationHeader = "Authorization" // TODO correct names?
)

type sTwitterResponse struct {
    IDString    string
    ID          int64
    CreatedAt   string
    PhoneNumber string
    AccessToken sTwitterAccessToken
}

type sTwitterAccessToken struct {
    Secret string
    Token  string
}

// TODO security checks headers
func (server *sServer) Register(writer http.ResponseWriter,
    request *http.Request, _ httprouter.Params) {
    var decoder *json.Decoder

    // retrieve headers for the OAuth Echo
    provider := request.Header.Get(cServiceProviderHeader)
    authorization := request.Header.Get(cClientAuthorizationHeader)
    if len(provider) == 0 || len(authorization) == 0 {
        server.respondWithError(ERR_AUT_NO_OAUTH_HEADERS, nil, cRegisterTag,
            writer, request)
        return
    }

    // create an HTTP Client and the OAuth Echo request
    client := http.Client{}
    oauthEcho, err := http.NewRequest("GET", provider, nil)
    if err != nil {
        server.respondWithError(ERR_INT_CREATE_HTTP, err, cRegisterTag, writer,
            request)
        return
    }

    // add the authorization header to the OAuth Echo request
    oauthEcho.Header.Add(cOAuthEchoAuthorizationHeader, authorization)

    // verify the provided credentials with Twitter's server
    response, err := client.Do(oauthEcho)
    if err != nil {
        server.respondWithError(ERR_INT_DO_HTTP, err, cRegisterTag, writer,
            request)
        return
    } else if response.StatusCode != http.StatusOK {
        server.respondWithError(ERR_AUT_BAD_OAUTH_RESPONSE, nil, cRegisterTag,
            writer, request)
        return
    }

    // decode Twitter's response
    var twitterResponse sTwitterResponse
    decoder = json.NewDecoder(response.Body)
    err = decoder.Decode(&twitterResponse)
    if err != nil {
        server.respondWithError(ERR_INT_DEMARSHALL, err, cRegisterTag, writer,
            request)
        return
    }

    // decode the sent user info
    var user polly.PrivateUser
    decoder = json.NewDecoder(request.Body)
    err = decoder.Decode(&user)
    if err != nil {
        server.respondWithError(ERR_BAD_JSON, err, cRegisterTag, writer,
            request)
        return
    }

    // check that the device type is a correct one
    if !isValidDeviceType(user.DeviceType) {
        server.respondWithError(ERR_BAD_DEVICE_TYPE, nil, cRegisterTag, writer,
            request)
        return
    }

    existentUser, err := server.db.GetUserByID(twitterResponse.ID)
    if err == nil {

        // we're dealing with an already existing user, generate a new token
        // and update his or her device type
        existentUser.DeviceType = user.DeviceType
        existentUser.Token = uuid.NewV4().String()
        err = server.db.UpdateToken(existentUser.ID, existentUser.Token)
        if err != nil {
            server.respondWithError(ERR_INT_DB_UPDATE, err, cRegisterTag,
                writer, request)
            return
        }

        // create the response body
        responseBody, err := json.MarshalIndent(existentUser, "", "\t")
        if err != nil {
            server.respondWithError(ERR_INT_MARSHALL, err, cRegisterTag,
                writer, request)
            return
        }

        // send the user a 200 OK with his new user info
        err = server.respondWithJSONBody(writer, responseBody)
        if err != nil {
            server.respondWithError(ERR_INT_WRITE, err, cRegisterTag, writer,
                request)
            return
        }

    } else {

        // we're dealing with a new user
        user.Token = uuid.NewV4().String()
        user.ID = twitterResponse.ID
        err = server.db.AddUser(&user)
        if err != nil {
            server.respondWithError(ERR_INT_DB_ADD, err, cRegisterTag, writer,
                request)
            return
        }

        // create the response body
        responseBody, err := json.MarshalIndent(user, "", "\t")
        if err != nil {
            server.respondWithError(ERR_INT_MARSHALL, err, cRegisterTag, writer,
                request)
            return
        }

        // send the user a 200 OK with his user info
        err = server.respondWithJSONBody(writer, responseBody)
        if err != nil {
            server.respondWithError(ERR_INT_WRITE, err, cRegisterTag, writer,
                request)
            return
        }
    }

}
