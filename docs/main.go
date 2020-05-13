package main

import (
	"fmt"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     "875572922963749",
		ClientSecret: "1a8edb7b4bd80e7945ac56412f440fdf",
		RedirectURL:  "http://localhost:9090/oauth2callback",
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}
	oauthStateString = "thisshouldberandom"
)

const htmlIndex = `<html><body>
Logged in with <a href="/login">facebook</a>
</body></html>
`

func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func handleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res := string(response)
	log.Printf("parseResponseBody: %s\n", res)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func main() {
	//http.HandleFunc("/", handleMain)
	//http.HandleFunc("/login", handleFacebookLogin)
	//http.HandleFunc("/oauth2callback", handleFacebookCallback)
	//fmt.Print("Started running on http://localhost:9090\n")
	//log.Fatal(http.ListenAndServe(":9090", nil))
	pdfg, err :=  wkhtml.NewPDFGenerator()
	if err != nil{
		//return
	}
	htmlStr := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="../css/style.css">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">
    <title>CGO - Ticket Experience PDF</title>
</head>
<body>
    <!-- Header -->
    <nav class="navbar navbar-light bg-light justify-content-between" id="navPayment">
        <div class="container">
            <a class="navbar-brand ticketPDF-nav py-4">
                <img src="../img/cGO Fix(1)-02.png" class="mr-4" alt="Kitten">
                E-Ticket
            </a>
        
        </div>
      </nav>

      <div id="down-payment-waiting">
        <div >
            <div class="container">
    
                <div class="mt-5">
                    <div class="row" style="    padding-right: 15px;
                    padding-left: 15px;">
                        <div class="col-md-7 p-4" style="box-shadow: 0px 5px 8px 4px #f1f1f1;">
                            <div class="mt-1 mr-2">
                                <span style="    padding: .3rem .8rem;
                            background: #dadada;
                            border-radius: 1rem;">Sailing</span>
                            </div>
                            <div class="mt-4 pb-4" style="border-bottom: 1px solid #e0e0e0;">
                            <div class="mt-3">
                                <p>11 Feb 2020</p>
                                <h4 class="info" style="font-size: 15px;">Sailing with Rising Fast Boat with Tropical Yacht Charters 2 Days 1 Night</h4>
                            </div>
                            <div class="d-flex mt-3">
                                <img height="24px" src="../img/pin-outline 3.png" alt="">
                                <h6 class="align-self-center mb-0 ml-2 text-muted" style="font-size: 15px;">Bali, Indonesia</h6>
                            </div>
                            </div>
                            <div class="row mt-3">
                                <div class="col-md-6">
                                    <h6 style="font-weight: 600 !important;">Meeting Point</h6>
                                    <div class="mt-2">
                                        <div class="d-flex justify-content-between">
                                            <p>Place</p>
                                            <p>Nusa Penida Harbour</p>
                                        </div>
                                        <div class="d-flex justify-content-between">
                                            <p>Time</p>
                                            <p>06.00 AM</p>
                                        </div>
                                    </div>
                                </div>

                                <div class="col-12">
                                    <div class="alert alert-info" role="alert">
                                        <div class="d-flex justify-content-between">
                                            <div class="d-flex">
                                                <img src="../img/Ellipse 485.png" alt="" style="width: 32px;  height: 32px; object-fit: cover;">
                                                <p class="ml-2 align-self-center mb-0">Tropical Yacht Charter</p>
                                            </div>
                                            <p class="align-self-center mb-0">Contact:   +62 856  219 2264</p>
                                        </div>
                                      </div>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-1"></div>
                        <div class="text-center col-md-4 p-4 barcode-ticketPDF align-self-center" style="box-shadow: 0px 5px 8px 4px #f1f1f1;">
                            <img src="../img/websiteQRCode_noFrame 1.png" alt="" style="width: 200px;">
                            <p>Order ID</p>
                            <h4>1234567</h4>
                        </div>
                    </div>
                    
                </div>

                <table class="table mt-5">
                    <thead class="thead-light">
                      <tr>
                        <th scope="col">No</th>
                        <th scope="col">Guest</th>
                        <th scope="col">Type</th>
                        <th scope="col">ID Type</th>
                        <th scope="col">ID Number</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <th scope="row">1</th>
                        <td>Emma Watson</td>
                        <td>Adult</td>
                        <td>ID Card</td>
                        <td>12345678</td>
                      </tr>
                      <tr>
                        <th scope="row">2</th>
                        <td>Helmi Ismail</td>
                        <td>Adult</td>
                        <td>ID Card</td>
                        <td>12345678</td>
                      </tr>
                      <tr>
                        <th scope="row">3</th>
                        <td>Arief Askar</td>
                        <td>Adult</td>
                        <td>ID Card</td>
                        <td>12345678</td>
                      </tr>
                    </tbody>
                  </table>
            </div>
        </div>
    </div>

    

    <footer class="footer-distributed">
    
        <div class="container">
            <div class="row">
                <div class="col-4">
                    <div class="d-flex">
                        <img src="../img/jam_ticket.png" alt="" width="53" height="53" class="mr-3">
                    <p>Show e-ticket to check-in at your departure place </p>
                    </div>
                </div>

                <div class="col-4">
                    <div class="d-flex">
                        <img src="../img/fa-regular_address-card.png" width="55" height="59" alt="" class="mr-3">
                    <p>Bring your official identity document as used in your booking</p>
                    </div>
                </div>

                <div class="col-4">
                    <div class="d-flex">
                        <img src="../img/Group 1618.png" alt="" width="53" height="47" class="mr-3">
                        <p>Please arrive at the harbour 60 minutes before departure</p>
                    </div>
                </div>
            </div>
        
    
        </footer>
    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
</body>
</html>`

	pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(htmlStr)))


	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	//Your Pdf Name
	err = pdfg.WriteFile("./Your_pdfname.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}