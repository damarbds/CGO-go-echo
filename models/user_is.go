package models

type GetToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
type RolesPermissionIs struct {
	Id 		string 		`json:"Id"`
	RoleName string 	`json:"RoleName"`
	RoleType int 		`json:"RoleType"`
	Description string 	`json:"Description"`
	Permissions []PermissionIs	`json:"Permissions"`
}
type PermissionIs struct {
	Id 		int 		`json:"Id"`
	ActivityCode string 	`json:"ActivityCode"`
	ActivityName string	`json:"ActivityName"`
	Description string	`json:"Description"`
	Status 		string	`json:"Status"`
}
type RegisterAndUpdateUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	GivenName     string `json:"givenname"`
	FamilyName    string `json:"familyname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailverified"`
	Website       string `json:"website"`
	Address       string `json:"address"`
	OTP           string `json:"oTP"`
	UserType      int    `json:"userType"`
	PhoneNumber   string `json:"phoneNumber"`
	UserRoles	 []string	`json:"userRoles"`
}

type Response struct {
	Email string `json:"email"`
}

type SendingEmail struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
	AttachmentFileUrl string `json:"attachmentFileUrl"`
	FileName string `json:"FileName"`
	From    string `json:"from"`
	To      string `json:"to"`
}
type SendingSMS struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Text        string `json:"text"`
	Encoding    string `json:"encoding"`
}

type RequestOTP struct {
	OTP              string `json:"OTP"`
	ExpiredDate      string `json:"ExpiredDate"`
	ExpiredInMSecond int    `json:"ExpiredInMSecond"`
}
type VerifiedEmail struct {
	Email   string `json:"email"`
	CodeOTP string `json:"codeOTP"`
}
type GetUserInfo struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	GivenName     string `json:"givenname"`
	FamilyName    string `json:"familyname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailverified"`
	Website       string `json:"website"`
	Address       string `json:"address"`
}
type GetUserDetail struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Password 	  string `json:"password"`
	Name          string `json:"name"`
	GivenName     string `json:"givenname"`
	FamilyName    string `json:"familyname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"emailverified"`
	Website       string `json:"website"`
	Address       string `json:"address"`
	PhoneNumber 	string	`json:"phone_number"`
	Roles 		[]*RolesPermissionIs `json:"Roles"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Scope    string `json:"scope"`
}

type RequestOTPNumber struct {
	PhoneNumber string `json:"phone_number"`
}
type RequestOTPTmpNumber struct {
	PhoneNumber string `json:"phone_number"`
	Email 		string	`json:"email"`
}