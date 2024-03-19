// Copyright (c) Microsoft Corporation. All rights reserved.
package appliance

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/antihax/optional"
	pureclient "github.com/cwedgwood/pureclient/client"
	"github.com/go-logr/logr"
	util "github.com/jezogwza/nc-toolbox-bin/pkg/utils"
)

type PureArray struct {
	apiClient  *pureclient.APIClient
	username   string
	password   string
	endpointIP string
	token      string
	logger     logr.Logger
}

var (
	pureAPIHTTPSPort string       = "443"
	httpClient       *http.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}} // nolint:gosec
)

// Nothing fancy. Pure supports up to 100 characters in length
const randomStringLength = 64

func NewPureArrayWithCredentials(endpointIP string, username string, password string, logger logr.Logger) (*PureArray, error) {

	pureArray := &PureArray{endpointIP: endpointIP,
		username: username,
		password: password,
		token:    "",
		logger:   logger}

	err := pureArray.getPureApiClient()
	if err != nil {
		return nil, err
	}
	return pureArray, nil
}

func NewPureArrayWithAPIToken(endpointIP string, token string, logger logr.Logger) (*PureArray, error) {

	pureArray := &PureArray{endpointIP: endpointIP,
		username: "",
		password: "",
		token:    token,
		logger:   logger}

	err := pureArray.getPureApiClient()
	if err != nil {
		return nil, err
	}
	return pureArray, nil
}

func (pa *PureArray) GetApiClient() *pureclient.APIClient {
	return pa.apiClient
}
func (pa *PureArray) getApiToken() (string, error) {
	// curl -k -X POST https://<endpointIP>/api/1.0/auth/apitoken  -H "Content-Type: application/json"  -d '{"username": "<username>", "password": "<password>"}'

	postBody, err := json.Marshal(map[string]string{
		"username": pa.username,
		"password": pa.password,
	})
	if err != nil {
		pa.logger.Error(err, "unable to json marshal")
		return "", err
	}

	authURL := "https://" + pa.endpointIP + ":" + pureAPIHTTPSPort + "/api/1.0/auth/apitoken"

	apiTokenReq, err := http.NewRequest("POST", authURL, bytes.NewBuffer(postBody))
	if err != nil {
		pa.logger.Error(err, "getApiToken failed to construct request from", "URL", authURL, "body", postBody)
		return "", err
	}

	apiTokenReq.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(apiTokenReq)
	if err != nil || resp == nil {
		pa.logger.Error(err, "getApiToken HTTP request failed") // Can't print request, as it has sensitive info in it.
		return "", err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		pa.logger.Error(err, "getApiToken failed to get response") // Can't print request, as it has sensitive info in it.
		return "", err
	}

	if resp.StatusCode >= 300 {
		err = fmt.Errorf("authentication to Pure arrray failed. status code: %v, status: %v, response body: %v", resp.StatusCode, resp.Status, string(responseBody))
		pa.logger.Error(err, "getApiToken failed") // Can't print request, as it has sensitive info in it.
		return "", err
	}

	responseBodyMap := make(map[string]interface{})
	err = json.Unmarshal(responseBody, &responseBodyMap)
	if err != nil {
		pa.logger.Error(err, "failed to unmarshal json")
		return "", err
	}

	return responseBodyMap["api_token"].(string), nil
}

func (pa *PureArray) getXAuthToken(apiToken string) (string, error) {
	// curl -X POST -k https://127.0.0.1:9999/api/2.8/login -H "api-token: fd99837e-f419-2b54-2ed3-f1fdfdc591ec" -H "Content-Type: application/json"

	tokenURL := "https://" + pa.endpointIP + ":" + pureAPIHTTPSPort + "/api/2.8/login"
	xAuthTokenReq, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		pa.logger.Error(err, "failed to get new http object for login")
		return "", err
	}

	xAuthTokenReq.Header.Add("Content-Type", "application/json")
	xAuthTokenReq.Header.Add("api-token", apiToken)

	resp, err := httpClient.Do(xAuthTokenReq)
	if err != nil {
		pa.logger.Error(err, "failed to get x-auth-token from Pure array") // Can't print request, as it has sensitive info in it.
		return "", err
	}

	resp.Body.Close()

	return resp.Header.Get("X-Auth-Token"), nil
}

func (pa *PureArray) getPureApiClient() error {
	// Two step process
	// 1. Get API token
	// 2. Use the API token to generate x-auth-token to use in REST calls to Pure array

	var err error
	if len(pa.username) != 0 {
		// username is used.
		pa.token, err = pa.getApiToken()
		if err != nil {
			return err
		}
	}

	// By now, pa.token is either retrieved from the array or set during the invokation of
	// NewPureArrayWithAPIToken.
	if len(pa.token) == 0 {
		err = fmt.Errorf("pure API Token is empty")
		pa.logger.Error(err, "API Token is needed to access Pure array")
		return err
	}
	xAuthToken, err := pa.getXAuthToken(pa.token)
	if err != nil {
		return err
	}

	pureClientConfig := pureclient.NewConfiguration()
	pureClientConfig.BasePath = "https://" + pa.endpointIP + ":" + pureAPIHTTPSPort
	pureClientConfig.DefaultHeader["x-auth-token"] = xAuthToken
	pureClientConfig.HTTPClient = httpClient
	pa.apiClient = pureclient.NewAPIClient(pureClientConfig)

	return nil
}

func (pa *PureArray) DeleteUser(username string) error {
	names := optional.NewInterface(username)
	query := &pureclient.AdministratorsApiApi28AdminsDeleteOpts{Names: names}

	httpResp, err := pa.apiClient.AdministratorsApi.Api28AdminsDelete(context.Background(), query)
	if err != nil {
		if httpResp == nil {
			pa.logger.Error(err, "Delete User: http response null", "name", username, "query", query)
			return err
		}

		errBody := err.Error()
		if swaggerErr, ok := err.(pureclient.GenericSwaggerError); ok {
			errBody = string(swaggerErr.Body())
			if strings.Contains(strings.ToLower(errBody), strings.ToLower("Unable to find specified local user")) {
				pa.logger.Info("User not found", "name", username, "query", query)
				return nil
			}

			err = fmt.Errorf(swaggerErr.Error())
		}

		pa.logger.Error(err, "failed to delete user on Pure array", "name", username,
			"http response status: ", httpResp.Status,
			"http response code: ", httpResp.StatusCode, "errBody", errBody)
		return fmt.Errorf("error deleting user %s: %v", username, err)
	}
	pa.logger.Info("Deleting user successful", "name", username)
	return nil
}

// GetUsers gets a list of all the users on the pure array.
func (pa *PureArray) GetUsers() ([]User, error) {
	// TODO: Limit how many to get at once, there should only be three.
	userList, httpResp, err := pa.apiClient.AdministratorsApi.Api28AdminsGet(context.Background(), nil)
	if err != nil {
		if httpResp == nil {
			pa.logger.Error(err, "GetUsers: http response null")
			return nil, err
		}
		errBody := err.Error()
		// If no users are found, Pure will respond with a XXX CONFIRM XXX 400 instead of 404 as well as
		// any status code above 300 will create a GenericSwaggerError, so we need to confirm
		// if this error is real or just an empty result so we can continue
		if swaggerErr, ok := err.(pureclient.GenericSwaggerError); (httpResp.StatusCode == 400) && ok {
			errBody = string(swaggerErr.Body())
			// TODO: What does no users look like?
			if strings.Contains(strings.ToLower(errBody), strings.ToLower("not found")) { // Guessing at error string
				pa.logger.Info("No Users found")
				return nil, nil
			}
			err = fmt.Errorf(swaggerErr.Error())
		}

		pa.logger.Error(err, "failed to read users on Pure array", "http response status: ", httpResp.Status,
			"http response code: ", httpResp.StatusCode, "errBody", errBody)
		return nil, err
	}

	// Likely won't happen as an empty result usually results in an error,
	// but will add in just in case
	if len(userList.Items) == 0 {
		pa.logger.Info("No Users found in Pure")
		return nil, nil
	}
	var pureUser User
	var pureUserList []User
	for _, user := range userList.Items {
		// if apiToken is created an earlier call, Pure would not let you retrieve it unless
		// we login to Pure with the same user's credentials or API token
		pureUser = User{Name: user.Name, Password: user.Password, Role: user.Role}
		pureUserList = append(pureUserList, pureUser)
	}

	pa.logger.Info("Users successfully retrieved.  List returned", "Users", userList.Items)
	return pureUserList, nil
}

func (pa *PureArray) getUser(username string) (*pureclient.Admin, error) {

	names := optional.NewInterface(username)
	query := &pureclient.AdministratorsApiApi28AdminsGetOpts{Names: names}

	userList, httpResp, err := pa.apiClient.AdministratorsApi.Api28AdminsGet(context.Background(), query)
	if err != nil {
		if httpResp == nil {
			pa.logger.Error(err, "getUser: http response null", "name", username)
			return nil, err
		}
		errBody := err.Error()
		// If the user is not found, Pure will respond with a 400 instead of 404 as well as
		// any status code above 300 will create a GenericSwaggerError, so we need to confirm
		// if this error is real or just an empty result so we can continue
		if swaggerErr, ok := err.(pureclient.GenericSwaggerError); (httpResp.StatusCode == 400) && ok {
			errBody = string(swaggerErr.Body())
			if strings.Contains(strings.ToLower(errBody), strings.ToLower("Unable to find specified user")) {
				pa.logger.Info("User not found", "User", username)
				return nil, nil
			}

			err = fmt.Errorf(swaggerErr.Error())
		}

		pa.logger.Error(err, "failed to read user on Pure array", "username", username,
			"http response status: ", httpResp.Status,
			"http response code: ", httpResp.StatusCode, "errBody", errBody)
		return nil, err
	}

	// Likely won't happen as an empty result usually results in an error,
	// but will add in just in case
	if len(userList.Items) == 0 {
		pa.logger.Info("User not found in Pure", "username", username)
		return nil, nil
	}

	pa.logger.Info("User successfully retrieved", "User", userList.Items[0])
	// There should just be one user with this name, if it exists. TODO Check
	return &userList.Items[0], nil
}

func (pa *PureArray) createLocalUser(user *User) error {
	userFetched, err := pa.getUser(user.Name)
	if err != nil {
		return err
	}

	if userFetched != nil {
		// Always delete existing user, so we can regenerate password and API token
		// They cannot be retrieved on lab redeployment but array not reset
		pa.logger.Info("User exists already, deleting it")
		pa.DeleteUser(user.Name)
	}

	names := optional.NewInterface(user.Name)
	query := &pureclient.AdministratorsApiApi28AdminsPostOpts{Names: names}

	if len(user.Password) == 0 {
		randomString := generateRandomString(randomStringLength)
		if len(randomString) == 0 { // ideally this should not happen, but just in case.
			randomString = user.Name
		}

		user.Password = randomString
	}

	body := pureclient.Model28AdminsBody{
		Password: user.Password,
		Role: &pureclient.AllOf28AdminsBodyRole{
			Name: user.Role,
		},
	}

	postResponse, httpResp, err := pa.apiClient.AdministratorsApi.Api28AdminsPost(context.Background(), body, query)
	if err != nil {
		if httpResp == nil {
			pa.logger.Error(err, "createLocalUser: http response null", "name", user.Name)
			return err
		}

		errBody := err.Error()
		if swaggerErr, ok := err.(pureclient.GenericSwaggerError); ok {
			errBody = string(swaggerErr.Body())
			err = fmt.Errorf(swaggerErr.Error())
		}

		pa.logger.Error(err, "failed to delete user on Pure array", "username", user.Name,
			"http response status: ", httpResp.Status,
			"http response code: ", httpResp.StatusCode, "errBody", errBody)
		return fmt.Errorf("error creating user %s: %v", user.Name, err)
	}

	pa.logger.Info("User created successfully. List of users returned for create call", "Users", postResponse.Items)
	return nil
}

func (pa *PureArray) createUserApiToken(username string) (string, error) {
	names := optional.NewInterface(username)
	query := &pureclient.AdministratorsApiApi28AdminsApiTokensPostOpts{Names: names}

	postResponse, httpResp, err := pa.apiClient.AdministratorsApi.Api28AdminsApiTokensPost(context.Background(), query)
	if err != nil {
		if httpResp == nil {
			pa.logger.Error(err, "createUserApiToken: http response null", "name", username)
			return "", err
		}
		errBody := err.Error()
		if swaggerErr, ok := err.(pureclient.GenericSwaggerError); ok {
			errBody = string(swaggerErr.Body())
			err = fmt.Errorf(swaggerErr.Error())

			if strings.Contains(strings.ToLower(errBody), strings.ToLower("API token already created.")) {
				pa.logger.Info("Token already created, return doing nothing", "username", username)
				return "", nil
			}
		}

		pa.logger.Error(err, "failed to create user api token on Pure array", "username", username,
			"http response status: ", httpResp.Status,
			"http response code: ", httpResp.StatusCode, "errBody", errBody)
		return "", fmt.Errorf("error creating api token for user %s: %v", username, err)
	}

	return postResponse.Items[0].ApiToken.Token, nil
}

// This function will change the password of the specified user. The invoker
// needs to be the user or have a role of array admin
func (pa *PureArray) patchLocalUser(userName string, password string, newPassword string) error {
	userFetched, err := pa.getUser(userName)
	if err != nil {
		return err
	}

	if userFetched == nil {
		err = fmt.Errorf("user does not exist")
		pa.logger.Error(err, "patchLocalUser failed", "username", userName)
		return err
	}

	names := optional.NewInterface(userName)
	query := &pureclient.AdministratorsApiApi28AdminsPatchOpts{Names: names}

	body := pureclient.Model28AdminsBody1{
		Password:    newPassword,
		OldPassword: password,
	}

	patchResponse, httpResp, err := pa.apiClient.AdministratorsApi.Api28AdminsPatch(context.Background(), body, query)
	if err != nil {
		if httpResp == nil {
			pa.logger.Error(err, "patchLocalUser: http response null", "name", userName)
			return err
		}

		errBody := err.Error()
		if swaggerErr, ok := err.(pureclient.GenericSwaggerError); ok {
			errBody = string(swaggerErr.Body())
			err = fmt.Errorf(swaggerErr.Error())
		}

		pa.logger.Error(err, "failed to delete user on Pure array", "username", userName,
			"http response status: ", httpResp.Status,
			"http response code: ", httpResp.StatusCode, "errBody", errBody)
		return fmt.Errorf("error changing user password %s: %v", userName, err)
	}

	pa.logger.Info("User password changed successfully, changed user returned", "User", patchResponse.Items[0])
	return nil
}

// CreateUsers creates local users on the array
// It is advised to create one user at a time unless client consuming this API
// keeps track of successful user creation of each user input
// API token for an user once created will not be recreated or retrieved.
func (pa *PureArray) CreateUsers(users []User) ([]User, error) {
	// This should be the first call from controller during the reconcile
	// Not expected to be called more than once
	usersCreated := []User{}

	var err error
	for i := 0; i < len(users); i++ {
		user := users[i]
		pa.logger.Info("Creating", "user", user.Name, "role", user.Role)

		err = pa.createLocalUser(&user)
		if err != nil {
			return usersCreated, fmt.Errorf("error creating user")
		}

		apiToken, err := pa.createUserApiToken(user.Name)
		if err != nil {
			return usersCreated, fmt.Errorf("error creating user token")
		}

		var createdUser User
		// if apiToken is created an earlier call, Pure would not let you retrieve it unless
		// we login to Pure with the same user's credentials or API token
		if len(apiToken) > 0 {
			createdUser = User{Name: user.Name, Password: user.Password, Role: user.Role, ApiToken: apiToken}
		} else {
			createdUser = User{Name: user.Name, Password: user.Password, Role: user.Role}
		}
		usersCreated = append(usersCreated, createdUser)
	}

	return usersCreated, nil
}

// Change the password of a user. The newPassword cannot be empty.
// Currently this function is only supported for Bootstrap User.
func (pa *PureArray) ChangeUserPassword(userName string, password string, newPassword string) error {
	if userName != pa.username || password != pa.password {
		err := fmt.Errorf("password change is not supported for the user")
		pa.logger.Error(err, "ChangeUserPassword", "User", userName)
		return err
	}

	if len(newPassword) == 0 {
		err := fmt.Errorf("newPassword cannot be empty")
		pa.logger.Error(err, "ChangeUserPassword", "User", userName)
		return err
	}

	return pa.patchLocalUser(userName, password, newPassword)
}

func generateRandomString(length int) string {
	const pwLen = 12
	var randomString string
	for i := 0; i < pwLen; i++ { // this is fairly large number
		randomString = util.GenerateRandomString(length)

		// See if generated string contains at least one character.
		// It is a requirement from Pure
		for _, c := range randomString {
			if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
				continue
			} else {
				return randomString
			}
		}
	}

	return ""
}
