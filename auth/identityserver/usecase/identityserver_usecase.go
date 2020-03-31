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
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func (m identityserverUsecase) UploadFileToBlob(image string,files *os.File) (string, error) {
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
	ctx := context.Background()
	// Create a file to test the upload and download.
	fmt.Printf("Creating a dummy file to test the upload and download\n")
	data := []byte(image)
	fileName := image
	err = ioutil.WriteFile(fileName, data, 0700)
	if err != nil {
		return "#",models.ErrBadParamInput
	}

	// Here's how to upload a blob.
	blobURL := containerURL.NewBlockBlobURL(fileName)
	//file, err := os.Open(fileName)
	//if err != nil {
	//	return "#",models.ErrBadParamInput
	//}

	dir, err := os.Getwd()
	fileLocation := filepath.Join(dir, "files", fileName)
	fmt.Println(fileLocation)
	// You can use the low-level Upload (PutBlob) API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that Upload can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/rest/api/storageservices/put-blob
	// To upload more than 256MB, use StageBlock (PutBlock) and CommitBlockList (PutBlockList) functions.
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.Upload(ctx, file, azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// handleErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// This function calls StageBlock/CommitBlockList for files larger 256 MBs, and calls Upload for any file smaller
	fmt.Printf("Uploading the file with blob name: %s\n", fileName)
	_, err = azblob.UploadFileToBlockBlob(ctx, files, blobURL, azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "image/jpg",   //  Add any needed headers here
		},
	})
	if err != nil {
		return "#",models.ErrBadParamInput
	}
	return image,nil
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

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
