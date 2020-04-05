package usecase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/auth/identityserver"
	"github.com/models"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type identityserverUsecase struct {
	baseUrl				string
	basicAuth 			string
	accountStorage		string
	accessKeyStorage 	string
}


// NewidentityserverUsecase will create new an identityserverUsecase object representation of identityserver.Usecase interface
func NewidentityserverUsecase(baseUrl string,basicAuth string,accountStorage string,accessKeyStorage string) identityserver.Usecase {
	return &identityserverUsecase{
		baseUrl: 			baseUrl,
		basicAuth:			basicAuth,
		accountStorage:		accountStorage,
		accessKeyStorage:	accessKeyStorage,
	}
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

func handleErrors(err error) {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				fmt.Println("Received 409. Container already exists")
				return
			}
		}
		log.Fatal(err)
	}
}
func (m identityserverUsecase) UploadFileToBlob(image string,folder string) (string, error) {
	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := m.accountStorage, m.accessKeyStorage
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a random string for the quick start container
	containerName := "cgo-storage"

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create the container
	fmt.Printf("Creating a container named %s\n", containerName)
	ctx := context.Background() // This example uses a never-expiring context
	//_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	//handleErrors(err)

	// Create a file to test the upload and download.
	fmt.Printf("Creating a dummy file to test the upload and download\n")

	data ,erread:= ioutil.ReadFile(image)
	if erread != nil {
		return "",erread
	}
	fileName := randomString()
	fileName = fileName + ".jpg"
	err = ioutil.WriteFile(fileName, data, 0700)
	handleErrors(err)

	// Here's how to upload a blob.
	folderPath := folder + "/"
	blobURL := containerURL.NewBlockBlobURL(folderPath + fileName)
	file, err := os.Open(fileName)
	handleErrors(err)

	// You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that PutBlob can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// handleErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	fmt.Printf("Uploading the file with blob name: %s\n", fileName)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	handleErrors(err)
	//
	//// List the container that we have created above
	//fmt.Println("Listing the blobs in the container:")
	//for marker := (azblob.Marker{}); marker.NotDone(); {
	//	// Get a result segment starting with the blob indicated by the current Marker.
	//	listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
	//	handleErrors(err)
	//
	//	// ListBlobs returns the start of the next segment; you MUST use this to get
	//	// the next segment (after processing the current result segment).
	//	marker = listBlob.NextMarker
	//
	//	// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
	//	for _, blobInfo := range listBlob.Segment.BlobItems {
	//		fmt.Print("	Blob name: " + blobInfo.Name + "\n")
	//	}
	//}
	//
	//// Here's how to download the blob
	//downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	//
	//// NOTE: automatically retries are performed if the connection fails
	//bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	//
	//// read the body into a buffer
	//downloadedData := bytes.Buffer{}
	//_, err = downloadedData.ReadFrom(bodyStream)
	//handleErrors(err)
	//
	//// The downloaded blob data is in downloadData's buffer. :Let's print it
	//fmt.Printf("Downloaded the blob: " + downloadedData.String())
	//
	//// Cleaning up the quick start by deleting the container and the file created locally
	//fmt.Printf("Press enter key to delete the sample files, example container, and exit the application.\n")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	//fmt.Printf("Cleaning up.\n")
	//containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	file.Close()
	os.Remove(fileName)
	//os.Remove(image)
	return blobURL.String(),err
}
func (m identityserverUsecase) GetUserInfo(token string) (*models.GetUserInfo, error) {

	req, err := http.NewRequest("POST", m.baseUrl + "/connect/userinfo", nil)
	//os.Exit(1)
	req.Header.Set("Authorization", "Bearer " + token)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrUnAuthorize
	}
	user := models.GetUserInfo{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func (m identityserverUsecase) GetToken(username string, password string) (*models.GetToken, error) {

	var param = url.Values{}
	param.Set("grant_type", "password")
	param.Set("username", username)
	param.Set("password", password)
	param.Set("scope", "openid")
	var payload = bytes.NewBufferString(param.Encode())

	req, err := http.NewRequest("POST", m.baseUrl + "/connect/token", payload)
	//os.Exit(1)
	req.Header.Set("Authorization", "Basic " + m.basicAuth)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil  {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrUsernamePassword
	}
	user := models.GetToken{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}


func (m identityserverUsecase) UpdateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser,error) {

	data, _:= json.Marshal(ar)

	req, err := http.NewRequest("POST", m.baseUrl + "/connect/update-user", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type","application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrBadParamInput
	}
	user := models.RegisterAndUpdateUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func (m identityserverUsecase) CreateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser, error) {
	data, _:= json.Marshal(ar)
	req, err := http.NewRequest("POST", m.baseUrl + "/connect/register", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type","application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrBadParamInput
	}
	user := models.RegisterAndUpdateUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}
func (m identityserverUsecase) SendingEmail(r *models.SendingEmail) (*models.SendingEmail, error) {
	data, _:= json.Marshal(r)
	req, err := http.NewRequest("POST", m.baseUrl + "/connect/push-email", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type","application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrBadParamInput
	}
	user := models.SendingEmail{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func (m identityserverUsecase) VerifiedEmail(r *models.VerifiedEmail) (*models.VerifiedEmail, error) {
	data, _:= json.Marshal(r)
	req, err := http.NewRequest("POST", m.baseUrl + "/connect/verified-email", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type","application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrInvalidOTP
	}
	user := models.VerifiedEmail{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
