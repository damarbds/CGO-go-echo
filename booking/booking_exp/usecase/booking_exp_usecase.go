package usecase

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/product/reviews"

	"github.com/third-party/xendit"

	"github.com/product/experience_add_ons"
	"github.com/service/exp_payment"

	"github.com/service/experience"
	"github.com/transactions/transaction"

	"github.com/third-party/paypal"

	"github.com/third-party/midtrans"

	"github.com/auth/identityserver"
	"github.com/auth/merchant"
	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/skip2/go-qrcode"
	"golang.org/x/net/context"
)

type bookingExpUsecase struct {
	reviewRepo                reviews.Repository
	adOnsRepo                 experience_add_ons.Repository
	experiencePaymentTypeRepo exp_payment.Repository
	bookingExpRepo            booking_exp.Repository
	userUsecase               user.Usecase
	merchantUsecase           merchant.Usecase
	isUsecase                 identityserver.Usecase
	expRepo                   experience.Repository
	transactionRepo           transaction.Repository
	contextTimeout            time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewbookingExpUsecase(reviewRepo reviews.Repository, adOnsRepo experience_add_ons.Repository, ept exp_payment.Repository, a booking_exp.Repository, u user.Usecase, m merchant.Usecase, is identityserver.Usecase, er experience.Repository, tr transaction.Repository, timeout time.Duration) booking_exp.Usecase {
	return &bookingExpUsecase{
		reviewRepo:                reviewRepo,
		adOnsRepo:                 adOnsRepo,
		experiencePaymentTypeRepo: ept,
		bookingExpRepo:            a,
		userUsecase:               u,
		merchantUsecase:           m,
		isUsecase:                 is,
		expRepo:                   er,
		transactionRepo:           tr,
		contextTimeout:            timeout,
	}
}

const (
	templateWaitingApprovalDP string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
    <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Waiting Approval Down Payment</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Please wait for your booking confirmation</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Thank you for choosing cGO Indonesia. 
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We have received your booking <b>{{.title}}</b> with trip date on <b>{{.tripDate}}</b>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            The guide will confirm for availability confirmation within 1x24 hr for the trip you have booked. 
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            You have chosen to make a payment with Down Payment. Please note that your booking is reserved, but to get your official E-ticket from us, you must pay the remaining payment within determined time. Your guide will contact you regarding payment instructions.
                        </td>
                    </tr>
                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                            Down Payment
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">IDR {{.payment}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Remaining Payment
                                            </td>
                                            <td align="right" style="color: #35405A; ">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">IDR {{.remainingPayment}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Payment Deadline
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.paymentDeadline}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 20px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                           <b>Important:</b> we advise not to make any travel arrangements or pay for the <br> remaining payment before you receive guide’s confirmation.

                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 20px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateWaitingApprovalFP string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
    <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Waiting Approval Full Payment</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Please wait for your booking confirmation</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Thank you for choosing cGO Indonesia. 
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We have received your booking <b>{{.title}}</b> with trip date on <b>{{.tripDate}}</b>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            The guide will confirm for availability confirmation within <b>1x24 hr</b> for the trip you have booked. 
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Please note that your booking is reserved and you will get your official E-ticket which can be used for check in after we get your guide’s availability <br> confirmation.

                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 20px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                           <b>Important:</b> we advise not to make any travel arrangements before you receive guide’s confirmation.

                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateBookingCancelled string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Booking Cancelled</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your booking was cancelled</b>
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 40px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We regret to inform you that your trip <b>{{.title}}</b> with trip date on <b>{{.tripDate}} </b> was cancelled. This cancellation occurs because 
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 1.5rem 20px;">
                            <table border="0" cellpadding="0" cellspacing="50" width="100%" bgcolor="#F2F2F2" style="    border-radius: .5rem;">
                                <tr>
                                    <td style="text-align: center; font-family: 'Nunito Sans', sans-serif;
                                    font-style: normal;
                                    font-weight: normal;
                                    font-size: 15px;
                                    line-height: 24px;">
                                        Sailing ban from Indonesian government
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            If you wish to apply for a refund, please submit your bank account information and you will receive your refund within <font color="red">3 working days</font>.
If you wish your payment to be transmitted to credits, please click transmit to credits button.
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 2rem;
                            border-radius: 2rem; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">SUBMIT MY BANK ACCOUNT</a>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 25px 0 40px 0; text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 3rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">TRANSMIT TO CREDITS</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 25px 0 25px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            As a valued customer, your satisfaction is one of our concerns and we apologize for any inconvenience this cancellation caused. We suggest you to book another trip or choose different trip dates.
                        </td>
                    </tr>

                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketFP string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket FP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your booking has been successfully confirmed. Please find your E-ticket <br> attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: rgb(76, 76, 76); font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Meeting Point
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.meetingPoint}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Time
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.time}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style=" background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;"> 
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketFPWithoutMeetingPoint string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket FP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your booking has been successfully confirmed. Please find your E-ticket <br> attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Time
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.time}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style=" background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;"> 
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketFPWithoutTime string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket FP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your booking has been successfully confirmed. Please find your E-ticket <br> attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: rgb(76, 76, 76); font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Meeting Point
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.meetingPoint}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style=" background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;"> 
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketFPWithoutMeetingPointAndTime string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket FP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your booking has been successfully confirmed. Please find your E-ticket <br> attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style=" background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;"> 
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketDP string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket DP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your remaining payment has been successfully received. Please find your E-ticket attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: rgb(76, 76, 76); font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Meeting Point
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.meetingPoint}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Time
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.time}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketDPWithoutMeetingPoint string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket DP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your remaining payment has been successfully received. Please find your E-ticket attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Time
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.time}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketDPWithoutTime string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket DP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your remaining payment has been successfully received. Please find your E-ticket attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: rgb(76, 76, 76); font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Meeting Point
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.meetingPoint}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketDPWithoutMeetingPointAndTime string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket DP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your remaining payment has been successfully received. Please find your E-ticket attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;">
                                                <b style="font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                font-size: 15px;
                                                line-height: 24px;">{{.title}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 20px 30px 20px 30px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketTransportation string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
  <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Ticket TRANSPORTATION</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your E-ticket is here</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your booking has been successfully confirmed. Please find your E-ticket <br> attached.
                        </td>
                    </tr>

                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 0px 20px;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="padding: 15px 0;">
                                                <img src="{{.merchantPicture}}" alt="" width="53" height="24" style="object-fit: cover;" />
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="color: black;" width="15%">
                                                <b style="font-size: 17px; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                line-height: 24px;">{{.sourceTime}}</b>
                                            </td>
                                            <td style="color: #8E8E8E;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;" width="15%">
                                                {{.duration}}
                                            </td>
                                            <td style="color: black;">
                                                <b style="font-size: 17px;font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                line-height: 24px;">{{.desTime}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 0px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: black; padding: 15px 0 5px; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;" width="25%">
                                                {{.source}}
                                            </td>
                                            <td width="15%">
                                                <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/arrow-back.png" alt="">
                                            </td>
                                            <td style="color: black; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                 {{.dest}}
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="color:#7E7E7E; font-weight:600 !important;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.class}}</td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Dates
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.tripDate}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Guest
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.guestCount}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 4rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-size: 15px;
                            line-height: 24px;">ADD TO CALENDAR</a>
                        </td>
                    </tr>
                    
                    <tr>
                        <td style="padding: 30px 0 20px 0;font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We wish you a pleasant trip ahead.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0;font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateWaitingRemainingDP string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
    <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Waiting Remaining DP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Please pay for your remaining payment </b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We are happy to inform you that your booking <b>{{.title}}</b> with trip <br> date on <b> {{.tripDate}}</b> has been confirmed with your guide.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            You have chosen to make a payment with Down Payment. Please note that your booking is reserved but to get your official E-ticket from us, you must pay the remaining payment within determined time. After you pay the remaining payment, we will send you an E-ticket that can be used to check in.

                        </td>
                    </tr>
                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                            Down Payment
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">IDR {{.payment}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Remaining Payment
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">IDR {{.remainingPayment}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Payment Deadline
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.paymentDeadline}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 30px 0 5px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            <b>How to pay your remaining payment</b>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 5px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your guide will contact you regarding payment instructions. Please wait for them to contact you.  
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; ">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>

                    <tr>
                        <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
                <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                            <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;">
                                Please have your Order ID {{.orderId}} handy when contacting us.
    
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;" >
                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                    <tr>
                                        <td width="35%">
                                            <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">For Question</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                        <td>
                                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                <tr>
                                                    <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: normal;
                                                    line-height: 24px;">More Information</td>
                                                </tr>
                                                <tr>
                                                    <td >
                                                        <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                    </td>
                                                </tr>
                                            </table>
                                        </td>
                                    </tr>
                                    
                                </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">Download cGO app</td>
                        </tr>
                        <tr>
                            <td style="padding: 0px 20px 0px 20px;">
                                <table border="0" cellpadding="0" cellspacing="0">
                                    <tr>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                     <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                     <td>
                                      <a href="http://www.twitter.com/">
                                       <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                      </a>
                                     </td>
                                    </tr>
                                   </table>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: normal;
                            line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                        </tr>
                        </table>
                 </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketExperiencePDF string = `<html>
    <head>
       <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
       <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
       <style>
       @page {
           /*
            * Size can be a length (not a percentage) for width and height
            * or a standard page size such as: 
            * a4, a5, a3, b3, b4, letter, legal, ledger.
            * A standard page size can be followed by either 'portrait' or 'landscape'.
            *
            * In theory, you can use different page sizes in one document, but this renderer
            * currently uses the first page width as the width of the body. That means it
            * is only practical to use different page heights in the one document.
            * See danfickle/openhtmltopdf#176 or #119 for more information.
            */
           size: A4 portrait !important;
           
           /*
            * Margin box for each page. Accepts one-to-four values, similar
            * to normal margin property.
            */
           margin: 0px 0px 0px 0px !important;
           padding: 0px 0px 0px 0px !important;
       }
       html{
           margin: 0px 0px 0px 0px !important; 
       }
       body{
           margin: 0px 0px 0px 0px !important; 
       }
       </style>
   </head>
   <body style="margin: 0; padding: 0;">
   <table align="center" border="0" cellpadding="0" cellspacing="0" width="100%">
       <tr style="background: url('https://cgostorage.blob.core.windows.net/cgo-storage/img/img/backgroundColorCGO.jpeg'); background-size: cover;">
           <td style="padding: 15px 50px 15px 50px;">
               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                   <tr>
                    <td width="10%">
                     <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                    </td>
                    <td>
                     <font color="#ffffff" style="font-size: 20px; font-family: 'Rubik', sans-serif;font-weight: normal !important;">E-Ticket</font>
                    </td>
                   </tr>
                  </table>
           </td>
       </tr>
       <tr>
        <td bgcolor="#ffffff" style="padding: 50px 50px 15px 50px;">
           <table border="0" cellpadding="0" cellspacing="0" width="100%">
   
               <tr >
                <td bgcolor="#ffffff">
                   <table border="0" cellpadding="0" cellspacing="0" width="100%">
                       <tr>
                           <td style="padding: 20px;border-radius: .8rem; border: 1px solid #D1D1D1;vertical-align: initial;width: 55%;">
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           {{range .expType}}<a style="padding: 4px 1rem;
                                           background: #e8e5e5;
                                           border-radius: 1rem;
                                           font-size: 10px;
                                           font-family: 'Nunito Sans', sans-serif;
                                           font-style: normal;
                                           font-weight: 600;">{{.}}</a>{{end}}                                                
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 15px 0 10px 0;
                                       font-size: 11px;
                                       font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600;
                                       color: #35405A">
                                           {{.tripDate}}
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 0px 0 10px 0; font-size: 12px; font-family: 'Rubik', sans-serif;font-weight: normal !important; color: #35405A;">
                                           {{.title}}
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="    border-bottom: 1px solid #efeaea !important;
                                       padding-bottom: 1rem;">
                                           <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <td width="24">
                                                   <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/pin-outline_3.png" alt="" width="14" height="14">
                                               </td>
                                               <td style="color: #8E8E8E;font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: 700;">
                                                   {{.city}}, {{.country}}
                                               </td>
                                               <td style="font-size: 0; line-height: 0;" width="120">
                                                   &nbsp;
                                                   </td>
                                           </table>
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 10px 0px 6px 0px;">
                                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <tr>
                                                   <td style="color:#7E7E7E; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; padding-top: 10px;">Meeting Point</td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 0 0 .6rem 0;">
                                           <table >
                                               <tr>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: normal; padding-right: 1rem;">
                                                       Place
                                                   </td>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A;">
                                                       {{.meetingPoint}}
                                                   </td>
                                               </tr>
                                               <tr>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: normal; padding-right: 1rem; padding-top: 6px;">
                                                       Time
                                                   </td>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A; padding-top: 6px;">
                                                       {{.time}}
                                                   </td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
   
                                   <tr>
                                       <td  bgcolor="#E1FAFF" style="border: 1px solid #56CCF2; border-radius: .3rem; padding: 4px 7px;">
                                           <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <tr>
                                                   <td>
                                                       <img src="{{.merchantPicture}}" style="width: 32px; height: 32px: object-fit: cover;" alt="">
                                                   </td>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A">
                                                       by {{.merchantName}}
                                                   </td>
                                                   <td align="right" style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A">
                                                       Contact:   {{.merchantPhone}}
                                                   </td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
                               </table>
                           </td>
                           <td style="font-size: 0; line-height: 0;" width="5%">
                           &nbsp;
                           </td>
                           <td width="150" style="padding: 10px 20px 0px 20px; border-radius: .8rem; border: 1px solid #D1D1D1; width: 40%">
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td style="padding: 15px 0;text-align: center;">
                                           <img src="{{.qrCode}}" alt="" width="154" height="154" style="object-fit: cover;" />
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="color: black;text-align: center; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: normal;">
                                           Order ID
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="color: #3C7DF0;text-align: center; padding-bottom: 20px; font-size: 25px;font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 700;">
                                           {{.orderId}}
                                       </td>
                                   </tr>
                               </table>
                           </td>
                       </tr>
                      </table>
                </td>
               </tr>
               
               <tr>
                   <td style="padding: 50px 0 20px 0;">
                       <table  border="0" cellpadding="4" cellspacing="0" width="100%">
                           <tr bgcolor="#e6e6e6">
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">No</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">Guest </th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">Type</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">ID Type</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">ID Number</th>
                           </tr>
                           {{range .guestDesc}}<tr>
                               {{range rangeStruct .}}<td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                font-style: normal;
                                font-weight: 600; color: #35405A;" >{{.}}</td>{{end}}
                              </tr>{{end}}
                         </table>
                   </td>
               </tr>
              </table>
        </td>
       </tr>
       
      </table>
      <div style="width: 100%; position: fixed;bottom: 0">
          <table style="width: 100%">
           <tr>
                <td bgcolor="#EFF3FF" style="padding: 20px 30px 40px 30px;">
                   <table border="0" cellpadding="0" cellspacing="0" width="100%">
                       <tr>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/jam_ticket.png" alt="" width="35" height="35">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Show e-ticket to check-in at <br//> your departure place </td>
                                   </tr>
                               </table>
                               
                           </td>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/fa-regular_address-card.png" alt="" width="35" height="29">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal; padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Bring your official identity <br/>document as used in your <br/>booking </td>
                                   </tr>
                               </table>
                               
                           </td>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1618.png" alt="" width="33" height="27">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal;  padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Please arrive at the harbour 60 <br/> minutes before departure </td>
                                   </tr>
                               </table>
                               
                           </td>
                       </tr>
                   </table>
                </td>
               </tr>
          </table>
      </div>
        
   </body>
   </html>`

	templateTicketTransportationPDF string = `<html>
 <head>
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
	<style>
	@page {
		/*
		 * Size can be a length (not a percentage) for width and height
		 * or a standard page size such as: 
		 * a4, a5, a3, b3, b4, letter, legal, ledger.
		 * A standard page size can be followed by either 'portrait' or 'landscape'.
		 *
		 * In theory, you can use different page sizes in one document, but this renderer
		 * currently uses the first page width as the width of the body. That means it
		 * is only practical to use different page heights in the one document.
		 * See danfickle/openhtmltopdf#176 or #119 for more information.
		 */
		size: A4 portrait !important;
		
		/*
		 * Margin box for each page. Accepts one-to-four values, similar
		 * to normal margin property.
		 */
		margin: 0px 0px 0px 0px !important;
		padding: 0px 0px 0px 0px !important;
	}
	html{
		margin: 0px 0px 0px 0px !important; 
	}
	body{
		margin: 0px 0px 0px 0px !important; 
	}
	</style>
</head>

<body>
    <table align="center" border="0" cellpadding="0" cellspacing="0" width="100%">
		<tr style="background: url('https://cgostorage.blob.core.windows.net/cgo-storage/img/img/backgroundColorCGO.jpeg'); background-size: cover;">
			<td style="padding: 15px 50px 15px 50px;">
				<table border="0" cellpadding="0" cellspacing="0" width="100%">
					<tr>
					 <td width="10%">
					  <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
					 </td>
					 <td>
					  <font color="#ffffff" style="font-size: 20px; font-family: 'Rubik', sans-serif;font-weight: normal !important;">E-Ticket</font>
					 </td>
					</tr>
				   </table>
			</td>
		</tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 50px 50px 15px 50px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr >
                     <td bgcolor="#ffffff">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 20px;border-radius: .8rem; border: 1px solid #D1D1D1;vertical-align: initial;width: 55%;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td>
                                                <a style="padding: 4px 1rem;
                                                background: #e8e5e5;
                                                border-radius: 1rem;
                                                font-size: 10px;
                                                font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 600;">Transportation</a>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="padding: 15px 0;">
                                                <img src="{{.merchantPicture}}" alt="" style="object-fit: cover; width: 53px;" />
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="padding: 0 0 10px 0;
                                            font-size: 11px;
                                            font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: 600;
                                            color: #35405A">
                                                {{.tripDate}}
                                            </td>
                                        </tr>
                                        <tr>
                                            <td>
                                                <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                                    <td style="color: black;">
                                                        <b style="font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                                        font-style: normal;
                                                        font-weight: 700;">{{.sourceTime}}</b>
                                                    </td>
                                                    <td style="color: #8E8E8E;font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                    font-style: normal;
                                                    font-weight: 700;">
                                                        {{.duration}}
                                                    </td>
                                                    <td style="color: black;">
                                                        <b style="font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                                        font-style: normal;
                                                        font-weight: 700;">{{.desTime}}</b>
                                                    </td>
                                                    <td style="font-size: 0; line-height: 0;" width="120">
                                                        &nbsp;
                                                        </td>
                                                </table>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="padding: 10px 0px 10px 0px;">
                                                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                                    <tr>
                                                        <td style="font-family: 'Rubik', sans-serif;font-weight: normal !important; font-size: 13px;">
                                                            {{.source}}
                                                        </td>
                                                        <td>
                                                            <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/arrow-back.png" alt="">
                                                        </td>
                                                        <td style="color: black; text-align: right;font-family: 'Rubik', sans-serif;font-weight: normal !important; font-size: 13px;">
                                                            {{.dest}}
                                                        </td>
                                                        <td style="font-size: 0; line-height: 0;" width="76">
                                                            &nbsp;
                                                            </td>
                                                    </tr>
                                                    <tr>
                                                        <td style="color:#7E7E7E; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                                        font-style: normal;
                                                        font-weight: 600; padding-top: 10px;">{{.class}}</td>
                                                    </tr>
                                                </table>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                                <td style="font-size: 0; line-height: 0;" width="5%">
								&nbsp;
								</td>
                                <td width="150" style="padding: 10px 20px 0px 20px; border-radius: .8rem; border: 1px solid #D1D1D1; width: 40%">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="padding: 15px 0;text-align: center;">
                                                <img src="{{.qrCode}}" alt="" width="154" height="154" style="object-fit: cover;" />
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="color: black;text-align: center; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;">
                                                Order ID
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="color: #3C7DF0;text-align: center; padding-bottom: 20px; font-size: 25px;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: 700;">
                                                {{.orderId}}
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>
                    
                    <tr>
                        <td style="padding: 50px 0 20px 0;">
                            <table  border="0" cellpadding="4" cellspacing="0" width="100%">
                                <tr bgcolor="#e6e6e6">
                                  <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                  font-style: normal;
                                  font-weight: 600; color: #35405A;">No</th>
                                  <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                  font-style: normal;
                                  font-weight: 600; color: #35405A;">Guest </th>
                                  <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                  font-style: normal;
                                  font-weight: 600; color: #35405A;">Type</th>
                                  <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                  font-style: normal;
                                  font-weight: 600; color: #35405A;">ID Type</th>
                                  <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                  font-style: normal;
                                  font-weight: 600; color: #35405A;">ID Number</th>
                                </tr>
                                {{range .guestDesc}}<tr>
                                    {{range rangeStruct .}}<td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                     font-style: normal;
                                     font-weight: 600; color: #35405A;" >{{.}}</td>{{end}}
                                   </tr>{{end}}
                              </table>
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
	<div style="width: 100%; position: fixed;bottom: 0">
	   <table style="width: 100%">
		<tr>
			 <td bgcolor="#EFF3FF" style="padding: 20px 30px 40px 30px;">
				<table border="0" cellpadding="0" cellspacing="0" width="100%">
					<tr>
						<td>
							<table border="0" cellpadding="0" cellspacing="0" width="100%">
								<tr>
									<td>
										<img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/jam_ticket.png" alt="" width="35" height="35">
									</td>
									<td style="font-size: 11px; line-height: normal; font-family: 'Nunito Sans', sans-serif;
									font-style: normal;
									font-weight: 600; color: #35405A;">Show e-ticket to check-in at <br//> your departure place </td>
								</tr>
							</table>
							
						</td>
						<td>
							<table border="0" cellpadding="0" cellspacing="0" width="100%">
								<tr>
									<td>
										<img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/fa-regular_address-card.png" alt="" width="35" height="29">
									</td>
									<td style="font-size: 11px; line-height: normal; padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
									font-style: normal;
									font-weight: 600; color: #35405A;">Bring your official identity <br/>document as used in your <br/>booking </td>
								</tr>
							</table>
							
						</td>
						<td>
							<table border="0" cellpadding="0" cellspacing="0" width="100%">
								<tr>
									<td>
										<img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1618.png" alt="" width="33" height="27">
									</td>
									<td style="font-size: 11px; line-height: normal;  padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
									font-style: normal;
									font-weight: 600; color: #35405A;">Please arrive at the harbour 60 <br/> minutes before departure </td>
								</tr>
							</table>
							
						</td>
					</tr>
				</table>
			 </td>
			</tr>
	   </table>
   </div>
   </body>
</html>`

	templateBookingApprovalDP string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Booking Approved Down Payment</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
</head>
</html>
<body style="margin: 0; padding: 0;">
    <table bgcolor="#F2F2F2" border="0" cellpadding="0" cellspacing="0" width="100%">
     <tr>
      <td>
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="628">
            <tr>
                <td style="padding: 15px 30px 15px 30px; background:linear-gradient(90deg, rgba(35,62,152,1) 0%, rgba(35,62,152,1) 35%, rgba(53,116,222,1) 100%);">
                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                        <tr>
                         <td>
                          <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                         </td>
                         <td align="right" style="color: white; font-family: 'Nunito Sans', sans-serif;
                         font-weight: 700 !important;
                         font-size: 17px;">
                            Order ID: {{.orderId}}
                         </td>
                        </tr>
                       </table>
                </td>
            </tr>
            <tr>
             <td bgcolor="#ffffff" style="padding: 40px 30px 40px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                     <td>
                        <b style="font-size: 20px; font-family: 'Rubik', sans-serif;
                        color: #35405A;font-weight: normal !important;">Your booking has been confirmed</b>
                     </td>
                    </tr>
                    <tr>
                     <td style="padding: 30px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                     font-style: normal;
                     font-weight: normal;
                     font-size: 15px;
                     line-height: 24px;">
                        Dear {{.user}},
                     </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            We are happy to inform you that your booking <b>{{.title}}</b> with trip <br> date on <b> {{.tripDate}}</b> has been confirmed with your guide.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            You have chosen to make a payment with Down Payment. Please note that your booking is reserved but to get your official E-ticket from us, you must pay the remaining payment within determined time.
                        </td>
                    </tr>
                    <tr >
                     <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                            Down Payment
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                            <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">IDR {{.payment}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Remaining Payment
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">IDR {{.remainingPayment}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                            <tr>
                                <td style="padding: 10px 20px 10px 20px; ">
                                    <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                        <tr>
                                            <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                Payment Deadline
                                            </td>
                                            <td align="right" style="color: #35405A;">
                                                <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-size: 15px;
                                                line-height: 24px;">{{.paymentDeadline}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>

                    <tr>
                        <td style="padding: 20px 0 5px 0; ">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">How to pay your remaining payment</b>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 5px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Your guide will contact you regarding payment instructions. Please wait for them to contact you. 
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0;">
                            <b style="font-family: 'Nunito Sans', sans-serif;
                            font-style: normal;
                            font-weight: 700;
                            font-size: 15px;
                            line-height: 24px;">Your guide contact</b>
                        </td>
                    </tr>
                    <tr >
                        <td bgcolor="#E1FAFF" style="border-radius: .8rem;">
                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                               <tr>
                                   <td style="padding: 10px 20px 10px 20px;     border-bottom: 1px solid #E0E0E0;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                               <b style="font-weight: 800 !important;font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-size: 15px;
                                               line-height: 24px;">{{.userGuide}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                               <tr>
                                   <td style="padding: 10px 20px 10px 20px;">
                                       <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                           <tr>
                                               <td style="color: #35405A; font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: normal;
                                               font-size: 15px;
                                               line-height: 24px;">
                                                Guide Contact
                                               </td>
                                               <td align="right" style="color: #35405A;">
                                                   <b style="font-weight: 800 !important; font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-size: 15px;
                                                   line-height: 24px;">{{.guideContact}}</b>
                                               </td>
                                           </tr>
                                       </table>
                                   </td>
                               </tr>
                              </table>
                        </td>
                       </tr>

                    <tr>
                        <td style="padding: 20px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Do not hesitate to contact us if you have any questions or if you need additional information.
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            Kind regards,
                        </td>
                    </tr>
                    <tr>
                        <td style="font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            cGO Indonesia
                        </td>
                    </tr>
                   </table>
             </td>
            </tr>
            <tr>
             <td bgcolor="#E1FAFF" style="padding: 20px 30px 10px 30px;">
                <table border="0" cellpadding="0" cellspacing="0" width="100%">
                    <tr>
                        <td style="padding: 10px 20px 10px 20px; font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;">
                            Please have your Order ID {{.orderId}} handy when contacting us.

                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 0px 20px 10px 20px;" >
                            <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                <tr>
                                    <td width="35%">
                                        <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                            <tr>
                                                <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: normal;
                                                line-height: 24px;">For Question</td>
                                            </tr>
                                            <tr>
                                                <td >
                                                    <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1637.png" alt="">
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                    <td>
                                        <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                            <tr>
                                                <td style="padding: 10px 20px 10px 6px; color: #7A7A7A;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: normal;
                                                line-height: 24px;">More Information</td>
                                            </tr>
                                            <tr>
                                                <td >
                                                    <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1638.png" alt="">
                                                </td>
                                            </tr>
                                        </table>
                                    </td>
                                </tr>
                                
                            </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 0px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        line-height: 24px;">Download cGO app</td>
                    </tr>
                    <tr>
                        <td style="padding: 0px 20px 0px 20px;">
                            <table border="0" cellpadding="0" cellspacing="0">
                                <tr>
                                 <td>
                                  <a href="http://www.twitter.com/">
                                   <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/en_badge_web_generic_1.png" alt="Twitter" style="display: block;" border="0" />
                                  </a>
                                 </td>
                                 <td style="font-size: 0; line-height: 0;" width="20">&nbsp;</td>
                                 <td>
                                  <a href="http://www.twitter.com/">
                                   <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/6341429_preview_1.png" alt="Twitter" style="display: block;" border="0" />
                                  </a>
                                 </td>
                                </tr>
                               </table>
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 20px 10px 20px;font-size: 12px; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        line-height: 24px;">© 2019-2020, PT DTech Solusi Bisnis</td>
                    </tr>
                    </table>
             </td>
            </tr>
           </table>
      </td>
     </tr>
    </table>
   </body>`

	templateTicketExperiencePDFWithoutTime string = `<html>
    <head>
       <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
       <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
       <style>
       @page {
           /*
            * Size can be a length (not a percentage) for width and height
            * or a standard page size such as: 
            * a4, a5, a3, b3, b4, letter, legal, ledger.
            * A standard page size can be followed by either 'portrait' or 'landscape'.
            *
            * In theory, you can use different page sizes in one document, but this renderer
            * currently uses the first page width as the width of the body. That means it
            * is only practical to use different page heights in the one document.
            * See danfickle/openhtmltopdf#176 or #119 for more information.
            */
           size: A4 portrait !important;
           
           /*
            * Margin box for each page. Accepts one-to-four values, similar
            * to normal margin property.
            */
           margin: 0px 0px 0px 0px !important;
           padding: 0px 0px 0px 0px !important;
       }
       html{
           margin: 0px 0px 0px 0px !important; 
       }
       body{
           margin: 0px 0px 0px 0px !important; 
       }
       </style>
   </head>
   <body style="margin: 0; padding: 0;">
   <table align="center" border="0" cellpadding="0" cellspacing="0" width="100%">
       <tr style="background: url('https://cgostorage.blob.core.windows.net/cgo-storage/img/img/backgroundColorCGO.jpeg'); background-size: cover;">
           <td style="padding: 15px 50px 15px 50px;">
               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                   <tr>
                    <td width="10%">
                     <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                    </td>
                    <td>
                     <font color="#ffffff" style="font-size: 20px; font-family: 'Rubik', sans-serif;font-weight: normal !important;">E-Ticket</font>
                    </td>
                   </tr>
                  </table>
           </td>
       </tr>
       <tr>
        <td bgcolor="#ffffff" style="padding: 50px 50px 15px 50px;">
           <table border="0" cellpadding="0" cellspacing="0" width="100%">
   
               <tr >
                <td bgcolor="#ffffff">
                   <table border="0" cellpadding="0" cellspacing="0" width="100%" height="296px;">
                       <tr>
                           <td style="padding: 20px;border-radius: .8rem; border: 1px solid #D1D1D1;vertical-align: initial;width: 55%;">
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           {{range .expType}}<a style="padding: 4px 1rem;
                                           background: #e8e5e5;
                                           border-radius: 1rem;
                                           font-size: 10px;
                                           font-family: 'Nunito Sans', sans-serif;
                                           font-style: normal;
                                           font-weight: 600;">{{.}}</a>{{end}}                                                
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 15px 0 10px 0;
                                       font-size: 11px;
                                       font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600;
                                       color: #35405A">
                                           {{.tripDate}}
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 0px 0 10px 0; font-size: 12px; font-family: 'Rubik', sans-serif;font-weight: normal !important; color: #35405A;">
                                           {{.title}}
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="    border-bottom: 1px solid #efeaea !important;
                                       padding-bottom: 1rem;">
                                           <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <td width="24">
                                                   <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/pin-outline_3.png" alt="" width="14" height="14">
                                               </td>
                                               <td style="color: #8E8E8E;font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: 700;">
                                                   {{.city}}, {{.country}}
                                               </td>
                                               <td style="font-size: 0; line-height: 0;" width="120">
                                                   &nbsp;
                                                   </td>
                                           </table>
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 10px 0px 6px 0px;">
                                           <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <tr>
                                                   <td style="color:#7E7E7E; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; padding-top: 10px;">Meeting Point</td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 0 0 .6rem 0;">
                                           <table >
                                               <tr>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: normal; padding-right: 1rem;">
                                                       Place
                                                   </td>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A;">
                                                       {{.meetingPoint}}
                                                   </td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
   
                                   <tr>
                                       <td  bgcolor="#E1FAFF" style="border: 1px solid #56CCF2; border-radius: .3rem; padding: 4px 7px;">
                                           <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <tr>
                                                   <td>
                                                       <img src="{{.merchantPicture}}" style="width: 32px; height: 32px: object-fit: cover;" alt="">
                                                   </td>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A">
                                                       by {{.merchantName}}
                                                   </td>
                                                   <td align="right" style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A">
                                                       Contact:   {{.merchantPhone}}
                                                   </td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
                               </table>
                           </td>
                           <td style="font-size: 0; line-height: 0;" width="5%">
                           &nbsp;
                           </td>
                           <td width="150" style="padding: 10px 20px 0px 20px; border-radius: .8rem; border: 1px solid #D1D1D1; width: 40%">
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td style="padding: 15px 0;text-align: center;">
                                           <img src="{{.qrCode}}" alt="" width="154" height="154" style="object-fit: cover;" />
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="color: black;text-align: center; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: normal;">
                                           Order ID
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="color: #3C7DF0;text-align: center; padding-bottom: 20px; font-size: 25px;font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 700;">
                                           {{.orderId}}
                                       </td>
                                   </tr>
                               </table>
                           </td>
                       </tr>
                      </table>
                </td>
               </tr>
               
               <tr>
                   <td style="padding: 50px 0 20px 0;">
                       <table  border="0" cellpadding="4" cellspacing="0" width="100%">
                           <tr bgcolor="#e6e6e6">
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">No</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">Guest </th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">Type</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">ID Type</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">ID Number</th>
                           </tr>
                           {{range .guestDesc}}<tr>
                               {{range rangeStruct .}}<td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                font-style: normal;
                                font-weight: 600; color: #35405A;" >{{.}}</td>{{end}}
                              </tr>{{end}}
                         </table>
                   </td>
               </tr>
              </table>
        </td>
       </tr>
       
      </table>
      <div style="width: 100%; position: fixed;bottom: 0">
          <table style="width: 100%">
           <tr>
                <td bgcolor="#EFF3FF" style="padding: 20px 30px 40px 30px;">
                   <table border="0" cellpadding="0" cellspacing="0" width="100%">
                       <tr>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/jam_ticket.png" alt="" width="35" height="35">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Show e-ticket to check-in at <br//> your departure place </td>
                                   </tr>
                               </table>
                               
                           </td>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/fa-regular_address-card.png" alt="" width="35" height="29">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal; padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Bring your official identity <br/>document as used in your <br/>booking </td>
                                   </tr>
                               </table>
                               
                           </td>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1618.png" alt="" width="33" height="27">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal;  padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Please arrive at the harbour 60 <br/> minutes before departure </td>
                                   </tr>
                               </table>
                               
                           </td>
                       </tr>
                   </table>
                </td>
               </tr>
          </table>
      </div>
        
   </body>
   </html>`

	templateTicketExperiencePDFWithoutMeetingPointAndTime string = `<html>
    <head>
       <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
       <link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
       <style>
       @page {
           /*
            * Size can be a length (not a percentage) for width and height
            * or a standard page size such as: 
            * a4, a5, a3, b3, b4, letter, legal, ledger.
            * A standard page size can be followed by either 'portrait' or 'landscape'.
            *
            * In theory, you can use different page sizes in one document, but this renderer
            * currently uses the first page width as the width of the body. That means it
            * is only practical to use different page heights in the one document.
            * See danfickle/openhtmltopdf#176 or #119 for more information.
            */
           size: A4 portrait !important;
           
           /*
            * Margin box for each page. Accepts one-to-four values, similar
            * to normal margin property.
            */
           margin: 0px 0px 0px 0px !important;
           padding: 0px 0px 0px 0px !important;
       }
       html{
           margin: 0px 0px 0px 0px !important; 
       }
       body{
           margin: 0px 0px 0px 0px !important; 
       }
       </style>
   </head>
   <body style="margin: 0; padding: 0;">
   <table align="center" border="0" cellpadding="0" cellspacing="0" width="100%">
       <tr style="background: url('https://cgostorage.blob.core.windows.net/cgo-storage/img/img/backgroundColorCGO.jpeg'); background-size: cover;">
           <td style="padding: 15px 50px 15px 50px;">
               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                   <tr>
                    <td width="10%">
                     <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/cGO_Fix(1)-02.png" alt="">
                    </td>
                    <td>
                     <font color="#ffffff" style="font-size: 20px; font-family: 'Rubik', sans-serif;font-weight: normal !important;">E-Ticket</font>
                    </td>
                   </tr>
                  </table>
           </td>
       </tr>
       <tr>
        <td bgcolor="#ffffff" style="padding: 50px 50px 15px 50px;">
           <table border="0" cellpadding="0" cellspacing="0" width="100%">
   
               <tr >
                <td bgcolor="#ffffff">
                   <table border="0" cellpadding="0" cellspacing="0" width="100%" height="296px;">
                       <tr>
                           <td style="padding: 20px;border-radius: .8rem; border: 1px solid #D1D1D1;vertical-align: initial;width: 55%;">
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           {{range .expType}}<a style="padding: 4px 1rem;
                                           background: #e8e5e5;
                                           border-radius: 1rem;
                                           font-size: 10px;
                                           font-family: 'Nunito Sans', sans-serif;
                                           font-style: normal;
                                           font-weight: 600;">{{.}}</a>{{end}}                                                
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 15px 0 10px 0;
                                       font-size: 11px;
                                       font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600;
                                       color: #35405A">
                                           {{.tripDate}}
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="padding: 0px 0 10px 0; font-size: 12px; font-family: 'Rubik', sans-serif;font-weight: normal !important; color: #35405A;">
                                           {{.title}}
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="    border-bottom: 1px solid #efeaea !important;
                                       padding-bottom: 1rem;">
                                           <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <td width="24">
                                                   <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/pin-outline_3.png" alt="" width="14" height="14">
                                               </td>
                                               <td style="color: #8E8E8E;font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                               font-style: normal;
                                               font-weight: 700;">
                                                   {{.city}}, {{.country}}
                                               </td>
                                               <td style="font-size: 0; line-height: 0;" width="120">
                                                   &nbsp;
                                                   </td>
                                           </table>
                                       </td>
                                   </tr>                                  
                                   <tr>
                                       <td  bgcolor="#E1FAFF" style="border: 1px solid #56CCF2; border-radius: .3rem; padding: 4px 7px;">
                                           <table  border="0" cellpadding="0" cellspacing="0" width="100%">
                                               <tr>
                                                   <td>
                                                       <img src="{{.merchantPicture}}" style="width: 32px; height: 32px: object-fit: cover;" alt="">
                                                   </td>
                                                   <td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A">
                                                       by {{.merchantName}}
                                                   </td>
                                                   <td align="right" style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                                   font-style: normal;
                                                   font-weight: 600; color: #35405A">
                                                       Contact:   {{.merchantPhone}}
                                                   </td>
                                               </tr>
                                           </table>
                                       </td>
                                   </tr>
                               </table>
                           </td>
                           <td style="font-size: 0; line-height: 0;" width="5%">
                           &nbsp;
                           </td>
                           <td width="150" style="padding: 10px 20px 0px 20px; border-radius: .8rem; border: 1px solid #D1D1D1; width: 40%">
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td style="padding: 15px 0;text-align: center;">
                                           <img src="{{.qrCode}}" alt="" width="154" height="154" style="object-fit: cover;" />
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="color: black;text-align: center; font-size: 13px;font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: normal;">
                                           Order ID
                                       </td>
                                   </tr>
                                   <tr>
                                       <td style="color: #3C7DF0;text-align: center; padding-bottom: 20px; font-size: 25px;font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 700;">
                                           {{.orderId}}
                                       </td>
                                   </tr>
                               </table>
                           </td>
                       </tr>
                      </table>
                </td>
               </tr>
               
               <tr>
                   <td style="padding: 50px 0 20px 0;">
                       <table  border="0" cellpadding="4" cellspacing="0" width="100%">
                           <tr bgcolor="#e6e6e6">
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">No</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">Guest </th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">Type</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">ID Type</th>
                             <th style="text-align: left; font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                             font-style: normal;
                             font-weight: 600; color: #35405A;">ID Number</th>
                           </tr>
                           {{range .guestDesc}}<tr>
                               {{range rangeStruct .}}<td style="font-size: 11px;font-family: 'Nunito Sans', sans-serif;
                                font-style: normal;
                                font-weight: 600; color: #35405A;" >{{.}}</td>{{end}}
                              </tr>{{end}}
                         </table>
                   </td>
               </tr>
              </table>
        </td>
       </tr>
       
      </table>
      <div style="width: 100%; position: fixed;bottom: 0">
          <table style="width: 100%">
           <tr>
                <td bgcolor="#EFF3FF" style="padding: 20px 30px 40px 30px;">
                   <table border="0" cellpadding="0" cellspacing="0" width="100%">
                       <tr>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/jam_ticket.png" alt="" width="35" height="35">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Show e-ticket to check-in at <br//> your departure place </td>
                                   </tr>
                               </table>
                               
                           </td>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/fa-regular_address-card.png" alt="" width="35" height="29">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal; padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Bring your official identity <br/>document as used in your <br/>booking </td>
                                   </tr>
                               </table>
                               
                           </td>
                           <td>
                               <table border="0" cellpadding="0" cellspacing="0" width="100%">
                                   <tr>
                                       <td>
                                           <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Group_1618.png" alt="" width="33" height="27">
                                       </td>
                                       <td style="font-size: 11px; line-height: normal;  padding-left: 10px; font-family: 'Nunito Sans', sans-serif;
                                       font-style: normal;
                                       font-weight: 600; color: #35405A;">Please arrive at the harbour 60 <br/> minutes before departure </td>
                                   </tr>
                               </table>
                               
                           </td>
                       </tr>
                   </table>
                </td>
               </tr>
          </table>
      </div>
        
   </body>
   </html>`
)

var templateFuncs = template.FuncMap{"rangeStruct": rangeStructer}

func (b bookingExpUsecase) RemainingPaymentNotification(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()
	list, err := b.transactionRepo.GetTransactionDownPaymentByDate(ctx)
	if err != nil {
		return err
	}
	for _, element := range list {
		var bookedBy []models.BookedByObj
		if element.BookedBy != "" {
			if errUnmarshal := json.Unmarshal([]byte(element.BookedBy), &bookedBy); errUnmarshal != nil {
				return err
			}
		}
		user := bookedBy[0].Title + `.` + bookedBy[0].FullName
		tripDate := element.BookingDate.Format("02 January 2006")
		duration := 0
		if element.ExpDuration != 0 && element.ExpDuration != 1 {
			duration = element.ExpDuration - 1
			tripDate = tripDate + ` - ` + element.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
		}
		paymentDeadline := element.BookingDate
		if element.ExpPaymentDeadlineType != nil && element.ExpPaymentDeadlineAmount != nil {
			if *element.ExpPaymentDeadlineType == "Days" {
				paymentDeadline = paymentDeadline.AddDate(0, 0, -*element.ExpPaymentDeadlineAmount)
			} else if *element.ExpPaymentDeadlineType == "Week" {
				paymentDeadline = paymentDeadline.AddDate(0, 0, -*element.ExpPaymentDeadlineAmount*7)
			} else if *element.ExpPaymentDeadlineType == "Month" {
				paymentDeadline = paymentDeadline.AddDate(0, -*element.ExpPaymentDeadlineAmount, 0)
			}
		}
		var tmpl = template.Must(template.New("main-template").Parse(templateWaitingRemainingDP))
		remainingPayment := element.Price - element.TotalPrice
		var data = map[string]interface{}{
			"title":            element.ExpTitle,
			"user":             user,
			"payment":          element.TotalPrice,
			"remainingPayment": remainingPayment,
			"paymentDeadline":  paymentDeadline.Format("02 January 2006"),
			"orderId":          element.OrderId,
			"tripDate":         tripDate,
			"userGuide":        element.MerchantName,
			"guideContact":     element.MerchantPhone,
		}
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, data)
		if err != nil {
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//maxTime := time.Now().AddDate(0, 0, 1)
		msg := tpl.String()
		pushEmail := &models.SendingEmail{
			Subject:  "Waiting Remaining DP",
			Message:  msg,
			From:     "CGO Indonesia",
			To:       bookedBy[0].Email,
			FileName: "",
		}

		_, err = b.isUsecase.SendingEmail(pushEmail)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b bookingExpUsecase) XenPayment(ctx context.Context, amount float64, tokenId, authId, orderId, paymentType string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	var result map[string]interface{}

	xendit.XenditSetup()

	booking, err := b.bookingExpRepo.GetByID(ctx, orderId)
	if err != nil {
		return nil, err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, errUnmarshal
		}
	}

	re, err := regexp.Compile("[0-9]+")
	if err != nil {
		return nil, err
	}
	fullname := re.ReplaceAllString(bookedBy[0].FullName, "")
	splitName := strings.Split(fullname, " ")
	var name string
	if len(splitName) < 2 {
		name = bookedBy[0].Title + " " + strings.Title(fullname)
	} else {
		name = fullname
	}

	if paymentType == "BRI" {
		va := &xendit.VirtualAccount{
			Client:     xendit.XenClient.VirtualAccount,
			ExternalID: orderId,
			BankCode:   paymentType,
			Name:       name,
			ExpireDate: booking.ExpiredDatePayment,
		}
		resVA, err := va.CreateFixedVA(ctx)
		if err != nil {
			return result, err
		}

		var bookingCode string
		if booking.ExpId != nil {
			bookingCode = booking.Id
		} else {
			bookingCode = booking.OrderId
		}
		if err := b.transactionRepo.UpdateAfterPayment(ctx, 0, resVA.AccountNumber, "", bookingCode); err != nil {
			return nil, err
		}

		result = structToMap(resVA)
	}

	if paymentType == "cc" || (authId != "" && tokenId != "") {
		cc := &xendit.CreditCard{
			Client:     xendit.XenClient.Card,
			TokenID:    tokenId,
			AuthID:     authId,
			ExternalID: orderId,
			Amount:     amount,
			IsCapture:  true,
		}
		resCC, err := cc.CreateCharge(ctx)
		if err != nil {
			return result, err
		}

		if err := b.SetAfterCCPayment(ctx, resCC.ExternalID, resCC.MaskedCardNumber, resCC.Status); err != nil {
			return result, err
		}

		result = structToMap(resCC)
	}

	return result, nil
}

func (b bookingExpUsecase) SetAfterCCPayment(ctx context.Context, externalId, accountNumber, status string) error {
	booking, err := b.bookingExpRepo.GetByID(ctx, externalId)
	if err != nil {
		return err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return errUnmarshal
		}
	}

	var transactionStatus int
	if status == "CAPTURED" {
		if booking.ExpId != nil {
			exp, err := b.expRepo.GetByID(ctx, *booking.ExpId)
			if err != nil {
				return err
			}
			bookingDetail, err := b.GetDetailBookingID(ctx, booking.Id, "")
			if err != nil {
				return err
			}
			if exp.ExpBookingType == "No Instant Booking" {
				transactionStatus = 1
				if bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
					user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
					tripDate := bookingDetail.BookingDate.Format("02 January 2006")
					duration := 0
					if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
						duration = bookingDetail.Experience[0].ExpDuration - 1
						tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
					}
					paymentDeadline := bookingDetail.BookingDate
					if bookingDetail.Experience[0].ExpPaymentDeadlineType != nil && bookingDetail.Experience[0].ExpPaymentDeadlineAmount != nil {
						if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Days" {
							paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount)
						} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Week" {
							paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount*7)
						} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Month" {
							paymentDeadline = paymentDeadline.AddDate(0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount, 0)
						}
					}
					var tmpl = template.Must(template.New("main-template").Parse(templateWaitingApprovalDP))
					var data = map[string]interface{}{
						"title":            exp.ExpTitle,
						"user":             user,
						"payment":          bookingDetail.TotalPrice,
						"remainingPayment": bookingDetail.ExperiencePaymentType.RemainingPayment,
						"paymentDeadline":  paymentDeadline.Format("02 January 2006"),
						"orderId":          bookingDetail.OrderId,
						"tripDate":         tripDate,
					}
					var tpl bytes.Buffer
					err = tmpl.Execute(&tpl, data)
					if err != nil {
						//http.Error(w, err.Error(), http.StatusInternalServerError)
					}

					//maxTime := time.Now().AddDate(0, 0, 1)
					msg := tpl.String()
					pushEmail := &models.SendingEmail{
						Subject:  "Waiting Approval For Merchant",
						Message:  msg,
						From:     "CGO Indonesia",
						To:       bookedBy[0].Email,
						FileName: "",
					}
					if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
						return nil
					}
				} else {
					user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
					tripDate := bookingDetail.BookingDate.Format("02 January 2006")
					duration := 0
					if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
						duration = bookingDetail.Experience[0].ExpDuration - 1
						tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
					}
					var tmpl = template.Must(template.New("main-template").Parse(templateWaitingApprovalFP))
					var data = map[string]interface{}{
						"title":    exp.ExpTitle,
						"user":     user,
						"tripDate": tripDate,
						"orderId":  bookingDetail.OrderId,
					}
					var tpl bytes.Buffer
					err = tmpl.Execute(&tpl, data)
					if err != nil {
						//http.Error(w, err.Error(), http.StatusInternalServerError)
					}

					//maxTime := time.Now().AddDate(0, 0, 1)
					msg := tpl.String()
					pushEmail := &models.SendingEmail{
						Subject:  "Waiting Approval For Merchant",
						Message:  msg,
						From:     "CGO Indonesia",
						To:       bookedBy[0].Email,
						FileName: "",
					}

					if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
						return nil
					}
				}

			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
				transactionStatus = 5
				user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
				tripDate := bookingDetail.BookingDate.Format("02 January 2006")
				duration := 0
				if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
					duration = bookingDetail.Experience[0].ExpDuration - 1
					tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
				}
				paymentDeadline := bookingDetail.BookingDate
				if bookingDetail.Experience[0].ExpPaymentDeadlineType != nil && bookingDetail.Experience[0].ExpPaymentDeadlineAmount != nil {
					if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Days" {
						paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount)
					} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Week" {
						paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount*7)
					} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Month" {
						paymentDeadline = paymentDeadline.AddDate(0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount, 0)
					}
				}
				var tmpl = template.Must(template.New("main-template").Parse(templateBookingApprovalDP))
				var data = map[string]interface{}{
					"title":            bookingDetail.Experience[0].ExpTitle,
					"user":             user,
					"payment":          bookingDetail.TotalPrice,
					"remainingPayment": bookingDetail.ExperiencePaymentType.RemainingPayment,
					"paymentDeadline":  paymentDeadline.Format("02 January 2006"),
					"orderId":          bookingDetail.OrderId,
					"tripDate":         tripDate,
					"userGuide":        bookingDetail.Experience[0].MerchantName,
					"guideContact":     bookingDetail.Experience[0].MerchantPhone,
				}
				var tpl bytes.Buffer
				err = tmpl.Execute(&tpl, data)
				if err != nil {
					//http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				//ticketPDF Bind HTML
				var htmlPDFTicket bytes.Buffer

				var guestDesc []models.GuestDescObjForHTML
				for i, element := range bookingDetail.GuestDesc {
					guest := models.GuestDescObjForHTML{
						No:       i + 1,
						FullName: element.FullName,
						Type:     element.Type,
						IdType:   element.IdType,
						IdNumber: element.IdNumber,
					}
					guestDesc = append(guestDesc, guest)
				}

				dataMapping := map[string]interface{}{
					"guestDesc":       guestDesc,
					"expType":         bookingDetail.Experience[0].ExpType,
					"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
					"title":           bookingDetail.Experience[0].ExpTitle,
					"city":            bookingDetail.Experience[0].City,
					"country":         bookingDetail.Experience[0].CountryName,
					"meetingPoint":    bookingDetail.Experience[0].ExpPickupPlace,
					"time":            bookingDetail.Experience[0].ExpPickupTime,
					"merchantName":    bookingDetail.Experience[0].MerchantName,
					"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
					"orderId":         bookingDetail.OrderId,
					"qrCode":          bookingDetail.TicketQRCode,
					"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
				}
				// We create the template and register out template function
				t := template.New("t").Funcs(templateFuncs)
				t, err := t.Parse(templateTicketExperiencePDF)
				if err != nil {
					panic(err)
				}

				err = t.Execute(&htmlPDFTicket, dataMapping)
				if err != nil {
					panic(err)
				}

				msg := tpl.String()
				// pdf := htmlPDFTicket.String()
				pushEmail := &models.SendingEmail{
					Subject: "Ticket DP",
					Message: msg,
					From:    "CGO Indonesia",
					To:      bookedBy[0].Email,
					// FileName:          "E-Ticket.pdf",
					// AttachmentFileUrl: pdf,
				}
				if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil
				}

			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Full Payment" {
				transactionStatus = 2
				user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
				tripDate := bookingDetail.BookingDate.Format("02 January 2006")
				duration := 0
				if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
					duration = bookingDetail.Experience[0].ExpDuration - 1
					tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
				}
				guestCount := len(bookingDetail.GuestDesc)

				var tmpl *template.Template
				var data map[string]interface{}
				var t *template.Template
				var dataMapping map[string]interface{}
				if bookingDetail.Experience[0].ExpPickupPlace == "" && (bookingDetail.Experience[0].ExpPickupTime == "" || bookingDetail.Experience[0].ExpPickupTime == "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutMeetingPointAndTime))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDFWithoutMeetingPointAndTime)
					if err != nil {
						panic(err)
					}

					t = temp

					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}

				} else if bookingDetail.Experience[0].ExpPickupPlace != "" && (bookingDetail.Experience[0].ExpPickupTime == "" || bookingDetail.Experience[0].ExpPickupTime == "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutTime))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"meetingPoint": bookingDetail.Experience[0].ExpPickupPlace,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDFWithoutTime)
					if err != nil {
						panic(err)
					}

					t = temp

					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"meetingPoint":    bookingDetail.Experience[0].ExpPickupPlace,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}

				} else if bookingDetail.Experience[0].ExpPickupPlace == "" && (bookingDetail.Experience[0].ExpPickupTime != "" && bookingDetail.Experience[0].ExpPickupTime != "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutMeetingPoint))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"time":         bookingDetail.Experience[0].ExpPickupTime,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDFWithoutMeetingPointAndTime)
					if err != nil {
						panic(err)
					}

					t = temp

					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}

				} else {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFP))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"meetingPoint": bookingDetail.Experience[0].ExpPickupPlace,
						"time":         bookingDetail.Experience[0].ExpPickupTime,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDF)
					if err != nil {
						panic(err)
					}

					t = temp

					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"meetingPoint":    bookingDetail.Experience[0].ExpPickupPlace,
						"time":            bookingDetail.Experience[0].ExpPickupTime,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}
				}
				var tpl bytes.Buffer
				err = tmpl.Execute(&tpl, data)
				if err != nil {
					//http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				//ticketPDF Bind HTML
				var htmlPDFTicket bytes.Buffer

				err = t.Execute(&htmlPDFTicket, dataMapping)
				if err != nil {
					panic(err)
				}

				msg := tpl.String()
				pdf := htmlPDFTicket.String()
				pushEmail := &models.SendingEmail{
					Subject:           "Ticket FP",
					Message:           msg,
					From:              "CGO Indonesia",
					To:                bookedBy[0].Email,
					FileName:          "E-Ticket.pdf",
					AttachmentFileUrl: pdf,
				}

				if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil
				}

			}
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, accountNumber, "", booking.Id); err != nil {
				return err
			}
		} else {
			bookingDetail, err := b.GetDetailTransportBookingID(ctx, booking.Id, "")
			if err != nil {
				return err
			}
			user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
			tripDate := bookingDetail.BookingDate.Format("02 January 2006")
			guestCount := len(bookingDetail.GuestDesc)

			layoutFormat := "15:04:05"
			departureTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].DepartureTime)
			arrivalTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].ArrivalTime)

			var tmpl = template.Must(template.New("main-template").Parse(templateTicketTransportation))
			var data = map[string]interface{}{
				"title":           bookingDetail.Transportation[0].TransTitle,
				"user":            user,
				"tripDate":        tripDate,
				"guestCount":      strconv.Itoa(guestCount) + " Guest(s)",
				"sourceTime":      departureTime.Format("15:04"),
				"desTime":         arrivalTime.Format("15:04"),
				"duration":        bookingDetail.Transportation[0].TripDuration,
				"source":          bookingDetail.Transportation[0].HarborSourceName,
				"dest":            bookingDetail.Transportation[0].HarborDestName,
				"class":           bookingDetail.Transportation[0].TransClass,
				"orderId":         bookingDetail.OrderId,
				"merchantPicture": bookingDetail.Transportation[0].MerchantPicture,
			}
			var tpl bytes.Buffer
			err = tmpl.Execute(&tpl, data)
			if err != nil {
				//http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			//ticketPDF Bind HTML
			var htmlPDFTicket bytes.Buffer

			var guestDesc []models.GuestDescObjForHTML
			for i, element := range bookingDetail.GuestDesc {
				guest := models.GuestDescObjForHTML{
					No:       i + 1,
					FullName: element.FullName,
					Type:     element.Type,
					IdType:   element.IdType,
					IdNumber: element.IdNumber,
				}
				guestDesc = append(guestDesc, guest)
			}

			dataMapping := map[string]interface{}{
				"guestDesc":       guestDesc,
				"tripDate":        tripDate,
				"sourceTime":      departureTime.Format("15:04"),
				"desTime":         arrivalTime.Format("15:04"),
				"duration":        bookingDetail.Transportation[0].TripDuration,
				"source":          bookingDetail.Transportation[0].HarborSourceName,
				"dest":            bookingDetail.Transportation[0].HarborDestName,
				"class":           bookingDetail.Transportation[0].TransClass,
				"qrCode":          bookingDetail.TicketQRCode,
				"merchantPicture": bookingDetail.Transportation[0].MerchantPicture,
				"orderId":         bookingDetail.OrderId,
			}
			// We create the template and register out template function
			t := template.New("t").Funcs(templateFuncs)
			t, err = t.Parse(templateTicketTransportationPDF)
			if err != nil {
				panic(err)
			}

			err = t.Execute(&htmlPDFTicket, dataMapping)
			if err != nil {
				panic(err)
			}

			msg := tpl.String()
			pdf := htmlPDFTicket.String()
			pushEmail := &models.SendingEmail{
				Subject:           "Ticket FP Transportation",
				Message:           msg,
				From:              "CGO Indonesia",
				To:                bookedBy[0].Email,
				FileName:          "E-Ticket.pdf",
				AttachmentFileUrl: pdf,
			}
			if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil
			}

			transactionStatus = 2
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, accountNumber, "", booking.OrderId); err != nil {
				return err
			}
		}

	} else if status == "FAILED" {
		var bookingCode string
		if booking.ExpId != nil {
			bookingCode = booking.Id
		} else {
			bookingCode = booking.OrderId
		}
		transactionStatus = 3

		//bookingDetail, err := b.GetDetailBookingID(ctx, booking.Id, "")
		//if err != nil {
		//	return err
		//}
		//tripDate := bookingDetail.BookingDate.Format("02 January 2006")
		//tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0,0,bookingDetail.Experience[0].ExpDuration).Format("02 January 2006")
		//var tmpl = template.Must(template.New("main-template").Parse(templateBookingCancelled))
		//var data = map[string]interface{}{
		//	"title": bookingDetail.Experience[0].ExpTitle,
		//	"orderId" : bookingDetail.OrderId,
		//	"tripDate" : tripDate,
		//}
		//var tpl bytes.Buffer
		//err = tmpl.Execute(&tpl, data)
		//if err != nil {
		//	//http.Error(w, err.Error(), http.StatusInternalServerError)
		//}
		//
		////maxTime := time.Now().AddDate(0, 0, 1)
		//msg := tpl.String()
		//pushEmail := &models.SendingEmail{
		//	Subject:  "Booking Cancelled",
		//	Message:  msg,
		//	From:     "CGO Indonesia",
		//	To:       bookedBy[0].Email,
		//	FileName: "",
		//}
		//if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
		//	return nil
		//}

		if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, accountNumber, "", bookingCode); err != nil {
			return err
		}
	}

	return nil
}

func (b bookingExpUsecase) GetByGuestCount(ctx context.Context, expId string, date string, guest int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()
	getExperience, err := b.expRepo.GetByID(ctx, expId)
	if err != nil {
		return false, err
	}
	getBooking, err := b.transactionRepo.GetCountByExpId(ctx, date, expId)
	if err != nil {
		return false, err
	}
	guestDesc := make([]models.GuestDescObj, 0)
	if getBooking != nil && *getBooking != "" {
		if errUnmarshal := json.Unmarshal([]byte(*getBooking), &guestDesc); errUnmarshal != nil {
			return false, models.ErrInternalServerError
		}
	}
	var result = false
	currentAmountBooking := len(guestDesc)
	remainingSeat := getExperience.ExpMaxGuest - currentAmountBooking
	if guest > remainingSeat {
		result = true
	}
	return result, nil
}
func (b bookingExpUsecase) Verify(ctx context.Context, orderId, bookingCode string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	var result map[string]interface{}

	cfg := paypal.PaypalConfig{
		OAuthUrl: paypal.PaypalOauthUrl,
		OrderUrl: paypal.PaypalOrderUrl,
	}

	res, err := paypal.PaypalSetup(cfg, orderId)
	if err != nil {
		return nil, err
	}

	if orderId != res.ID {
		return nil, errors.New("Incorrect Paypal Order ID")
	}

	booking, err := b.bookingExpRepo.GetByID(ctx, bookingCode)
	if err != nil {
		return nil, err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, errUnmarshal
		}
	}

	var transactionStatus int
	if res.Status == "COMPLETED" {
		if booking.ExpId != nil {
			exp, err := b.expRepo.GetByID(ctx, *booking.ExpId)
			if err != nil {
				return nil, err
			}
			bookingDetail, err := b.GetDetailBookingID(ctx, booking.Id, "")
			if err != nil {
				return nil, err
			}
			if exp.ExpBookingType == "No Instant Booking" {
				transactionStatus = 1
				if bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
					user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
					tripDate := bookingDetail.BookingDate.Format("02 January 2006")
					duration := 0
					if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
						duration = bookingDetail.Experience[0].ExpDuration - 1
						tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
					}
					paymentDeadline := bookingDetail.BookingDate
					if bookingDetail.Experience[0].ExpPaymentDeadlineType != nil && bookingDetail.Experience[0].ExpPaymentDeadlineAmount != nil {
						if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Days" {
							paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount)
						} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Week" {
							paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount*7)
						} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Month" {
							paymentDeadline = paymentDeadline.AddDate(0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount, 0)
						}
					}
					var tmpl = template.Must(template.New("main-template").Parse(templateWaitingApprovalDP))
					var data = map[string]interface{}{
						"title":            exp.ExpTitle,
						"user":             user,
						"payment":          bookingDetail.TotalPrice,
						"remainingPayment": bookingDetail.ExperiencePaymentType.RemainingPayment,
						"paymentDeadline":  paymentDeadline.Format("02 January 2006"),
						"orderId":          bookingDetail.OrderId,
						"tripDate":         tripDate,
					}
					var tpl bytes.Buffer
					err = tmpl.Execute(&tpl, data)
					if err != nil {
						//http.Error(w, err.Error(), http.StatusInternalServerError)
					}

					//maxTime := time.Now().AddDate(0, 0, 1)
					msg := tpl.String()
					pushEmail := &models.SendingEmail{
						Subject:  "Waiting Approval For Merchant",
						Message:  msg,
						From:     "CGO Indonesia",
						To:       bookedBy[0].Email,
						FileName: "",
					}
					if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
						return nil, nil
					}
				} else {
					user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
					tripDate := bookingDetail.BookingDate.Format("02 January 2006")
					duration := 0
					if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
						duration = bookingDetail.Experience[0].ExpDuration - 1
						tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
					}
					var tmpl = template.Must(template.New("main-template").Parse(templateWaitingApprovalFP))
					var data = map[string]interface{}{
						"title":    exp.ExpTitle,
						"user":     user,
						"tripDate": tripDate,
						"orderId":  bookingDetail.OrderId,
					}
					var tpl bytes.Buffer
					err = tmpl.Execute(&tpl, data)
					if err != nil {
						//http.Error(w, err.Error(), http.StatusInternalServerError)
					}

					//maxTime := time.Now().AddDate(0, 0, 1)
					msg := tpl.String()
					pushEmail := &models.SendingEmail{
						Subject:  "Waiting Approval For Merchant",
						Message:  msg,
						From:     "CGO Indonesia",
						To:       bookedBy[0].Email,
						FileName: "",
					}

					if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
						return nil, nil
					}
				}

			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
				transactionStatus = 5
				user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
				tripDate := bookingDetail.BookingDate.Format("02 January 2006")
				duration := 0
				if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
					duration = bookingDetail.Experience[0].ExpDuration - 1
					tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
				}
				paymentDeadline := bookingDetail.BookingDate
				if bookingDetail.Experience[0].ExpPaymentDeadlineType != nil && bookingDetail.Experience[0].ExpPaymentDeadlineAmount != nil {
					if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Days" {
						paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount)
					} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Week" {
						paymentDeadline = paymentDeadline.AddDate(0, 0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount*7)
					} else if *bookingDetail.Experience[0].ExpPaymentDeadlineType == "Month" {
						paymentDeadline = paymentDeadline.AddDate(0, -*bookingDetail.Experience[0].ExpPaymentDeadlineAmount, 0)
					}
				}
				var tmpl = template.Must(template.New("main-template").Parse(templateBookingApprovalDP))
				var data = map[string]interface{}{
					"title":            bookingDetail.Experience[0].ExpTitle,
					"user":             user,
					"payment":          bookingDetail.TotalPrice,
					"remainingPayment": bookingDetail.ExperiencePaymentType.RemainingPayment,
					"paymentDeadline":  paymentDeadline.Format("02 January 2006"),
					"orderId":          bookingDetail.OrderId,
					"tripDate":         tripDate,
					"userGuide":        bookingDetail.Experience[0].MerchantName,
					"guideContact":     bookingDetail.Experience[0].MerchantPhone,
				}
				var tpl bytes.Buffer
				err = tmpl.Execute(&tpl, data)
				if err != nil {
					//http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				//ticketPDF Bind HTML
				var htmlPDFTicket bytes.Buffer

				var guestDesc []models.GuestDescObjForHTML
				for i, element := range bookingDetail.GuestDesc {
					guest := models.GuestDescObjForHTML{
						No:       i + 1,
						FullName: element.FullName,
						Type:     element.Type,
						IdType:   element.IdType,
						IdNumber: element.IdNumber,
					}
					guestDesc = append(guestDesc, guest)
				}

				dataMapping := map[string]interface{}{
					"guestDesc":       guestDesc,
					"expType":         bookingDetail.Experience[0].ExpType,
					"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
					"title":           bookingDetail.Experience[0].ExpTitle,
					"city":            bookingDetail.Experience[0].City,
					"country":         bookingDetail.Experience[0].CountryName,
					"meetingPoint":    bookingDetail.Experience[0].ExpPickupPlace,
					"time":            bookingDetail.Experience[0].ExpPickupTime,
					"merchantName":    bookingDetail.Experience[0].MerchantName,
					"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
					"orderId":         bookingDetail.OrderId,
					"qrCode":          bookingDetail.TicketQRCode,
					"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
				}
				// We create the template and register out template function
				t := template.New("t").Funcs(templateFuncs)
				t, err := t.Parse(templateTicketExperiencePDF)
				if err != nil {
					panic(err)
				}

				err = t.Execute(&htmlPDFTicket, dataMapping)
				if err != nil {
					panic(err)
				}

				msg := tpl.String()
				// pdf := htmlPDFTicket.String()
				pushEmail := &models.SendingEmail{
					Subject: "Ticket DP",
					Message: msg,
					From:    "CGO Indonesia",
					To:      bookedBy[0].Email,
					// FileName:          "E-Ticket.pdf",
					// AttachmentFileUrl: pdf,
				}
				if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil, nil
				}

			} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Full Payment" {
				transactionStatus = 2
				user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
				tripDate := bookingDetail.BookingDate.Format("02 January 2006")
				duration := 0
				if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
					duration = bookingDetail.Experience[0].ExpDuration - 1
					tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
				}
				guestCount := len(bookingDetail.GuestDesc)

				var tmpl *template.Template
				var data map[string]interface{}
				var t *template.Template
				var dataMapping map[string]interface{}
				if bookingDetail.Experience[0].ExpPickupPlace == "" && (bookingDetail.Experience[0].ExpPickupTime == "" || bookingDetail.Experience[0].ExpPickupTime == "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutMeetingPointAndTime))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDFWithoutMeetingPointAndTime)
					if err != nil {
						panic(err)
					}

					t = temp

					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}

				} else if bookingDetail.Experience[0].ExpPickupPlace != "" && (bookingDetail.Experience[0].ExpPickupTime == "" || bookingDetail.Experience[0].ExpPickupTime == "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutTime))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"meetingPoint": bookingDetail.Experience[0].ExpPickupPlace,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDFWithoutTime)
					if err != nil {
						panic(err)
					}
					t = temp
					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"meetingPoint":    bookingDetail.Experience[0].ExpPickupPlace,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}

				} else if bookingDetail.Experience[0].ExpPickupPlace == "" && (bookingDetail.Experience[0].ExpPickupTime != "" && bookingDetail.Experience[0].ExpPickupTime != "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutMeetingPoint))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"time":         bookingDetail.Experience[0].ExpPickupTime,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDFWithoutMeetingPointAndTime)
					if err != nil {
						panic(err)
					}
					t = temp
					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}

				} else {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFP))
					data = map[string]interface{}{
						"title":        exp.ExpTitle,
						"user":         user,
						"tripDate":     tripDate,
						"orderId":      bookingDetail.OrderId,
						"meetingPoint": bookingDetail.Experience[0].ExpPickupPlace,
						"time":         bookingDetail.Experience[0].ExpPickupTime,
						"userGuide":    bookingDetail.Experience[0].MerchantName,
						"guideContact": bookingDetail.Experience[0].MerchantPhone,
						"guestCount":   strconv.Itoa(guestCount) + " Guest(s)",
					}

					//for html pdf
					var guestDesc []models.GuestDescObjForHTML
					for i, element := range bookingDetail.GuestDesc {
						guest := models.GuestDescObjForHTML{
							No:       i + 1,
							FullName: element.FullName,
							Type:     element.Type,
							IdType:   element.IdType,
							IdNumber: element.IdNumber,
						}
						guestDesc = append(guestDesc, guest)
					}

					// We create the template and register out template function
					temp := template.New("t").Funcs(templateFuncs)
					temp, err := temp.Parse(templateTicketExperiencePDF)
					if err != nil {
						panic(err)
					}
					t = temp
					dataMapping = map[string]interface{}{
						"guestDesc":       guestDesc,
						"expType":         bookingDetail.Experience[0].ExpType,
						"tripDate":        bookingDetail.BookingDate.Format("02 January 2006"),
						"title":           bookingDetail.Experience[0].ExpTitle,
						"city":            bookingDetail.Experience[0].City,
						"country":         bookingDetail.Experience[0].CountryName,
						"meetingPoint":    bookingDetail.Experience[0].ExpPickupPlace,
						"time":            bookingDetail.Experience[0].ExpPickupTime,
						"merchantName":    bookingDetail.Experience[0].MerchantName,
						"merchantPhone":   bookingDetail.Experience[0].MerchantPhone,
						"orderId":         bookingDetail.OrderId,
						"qrCode":          bookingDetail.TicketQRCode,
						"merchantPicture": bookingDetail.Experience[0].MerchantPicture,
					}
				}
				var tpl bytes.Buffer
				err = tmpl.Execute(&tpl, data)
				if err != nil {
					//http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				//ticketPDF Bind HTML
				var htmlPDFTicket bytes.Buffer

				err = t.Execute(&htmlPDFTicket, dataMapping)
				if err != nil {
					panic(err)
				}

				msg := tpl.String()
				pdf := htmlPDFTicket.String()
				pushEmail := &models.SendingEmail{
					Subject:           "Ticket FP",
					Message:           msg,
					From:              "CGO Indonesia",
					To:                bookedBy[0].Email,
					FileName:          "E-Ticket.pdf",
					AttachmentFileUrl: pdf,
				}

				if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil, nil
				}
			}
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, "", "", booking.Id); err != nil {
				return nil, err
			}
		} else {

			bookingDetail, err := b.GetDetailTransportBookingID(ctx, booking.Id, "")
			if err != nil {
				return nil, err
			}
			user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
			tripDate := bookingDetail.BookingDate.Format("02 January 2006")
			guestCount := len(bookingDetail.GuestDesc)
			layoutFormat := "15:04:05"
			departureTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].DepartureTime)
			arrivalTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].ArrivalTime)

			var tmpl = template.Must(template.New("main-template").Parse(templateTicketTransportation))
			var data = map[string]interface{}{
				"title":           bookingDetail.Transportation[0].TransTitle,
				"user":            user,
				"tripDate":        tripDate,
				"guestCount":      strconv.Itoa(guestCount) + " Guest(s)",
				"sourceTime":      departureTime.Format("15:04"),
				"desTime":         arrivalTime.Format("15:04"),
				"duration":        bookingDetail.Transportation[0].TripDuration,
				"source":          bookingDetail.Transportation[0].HarborSourceName,
				"dest":            bookingDetail.Transportation[0].HarborDestName,
				"class":           bookingDetail.Transportation[0].TransClass,
				"orderId":         bookingDetail.OrderId,
				"merchantPicture": bookingDetail.Transportation[0].MerchantPicture,
			}
			var tpl bytes.Buffer
			err = tmpl.Execute(&tpl, data)
			if err != nil {
				//http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			//ticketPDF Bind HTML
			var htmlPDFTicket bytes.Buffer

			var guestDesc []models.GuestDescObjForHTML
			for i, element := range bookingDetail.GuestDesc {
				guest := models.GuestDescObjForHTML{
					No:       i + 1,
					FullName: element.FullName,
					Type:     element.Type,
					IdType:   element.IdType,
					IdNumber: element.IdNumber,
				}
				guestDesc = append(guestDesc, guest)
			}

			dataMapping := map[string]interface{}{
				"guestDesc":       guestDesc,
				"tripDate":        tripDate,
				"sourceTime":      departureTime.Format("15:04"),
				"desTime":         arrivalTime.Format("15:04"),
				"duration":        bookingDetail.Transportation[0].TripDuration,
				"source":          bookingDetail.Transportation[0].HarborSourceName,
				"dest":            bookingDetail.Transportation[0].HarborDestName,
				"class":           bookingDetail.Transportation[0].TransClass,
				"qrCode":          bookingDetail.TicketQRCode,
				"merchantPicture": bookingDetail.Transportation[0].MerchantPicture,
				"orderId":         bookingDetail.OrderId,
			}
			// We create the template and register out template function
			t := template.New("t").Funcs(templateFuncs)
			t, err = t.Parse(templateTicketTransportationPDF)
			if err != nil {
				panic(err)
			}

			err = t.Execute(&htmlPDFTicket, dataMapping)
			if err != nil {
				panic(err)
			}

			msg := tpl.String()
			pdf := htmlPDFTicket.String()
			pushEmail := &models.SendingEmail{
				Subject:           "Ticket FP Transportation",
				Message:           msg,
				From:              "CGO Indonesia",
				To:                bookedBy[0].Email,
				FileName:          "E-Ticket.pdf",
				AttachmentFileUrl: pdf,
			}
			if _, err := b.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil, nil
			}

			transactionStatus = 2
			if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, "", "", booking.OrderId); err != nil {
				return nil, err
			}
		}

	}

	var bookCode string
	if booking.ExpId != nil {
		bookCode = booking.Id
	} else {
		bookCode = booking.OrderId
	}
	if res.Status == "VOIDED" {
		transactionStatus = 3
		if err := b.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, "", "", bookCode); err != nil {
			return nil, err
		}
	}

	data, _ := json.Marshal(res)
	json.Unmarshal(data, &result)

	return result, nil
}

func (b bookingExpUsecase) GetDetailTransportBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	details, err := b.bookingExpRepo.GetDetailTransportBookingID(ctx, bookingId, bookingCode)
	if err != nil {
		return nil, err
	}

	transport := make([]models.BookingTransportationDetail, len(details))
	for i, detail := range details {
		var tripDuration string
		if detail.DepartureTime != nil && detail.ArrivalTime != nil {
			departureTime, _ := time.Parse("15:04:00", *detail.DepartureTime)
			arrivalTime, _ := time.Parse("15:04:00", *detail.ArrivalTime)

			tripHour := arrivalTime.Hour() - departureTime.Hour()
			tripMinute := arrivalTime.Minute() - departureTime.Minute()
			tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
		}
		transport[i] = models.BookingTransportationDetail{
			TransID:          *detail.TransId,
			TransName:        *detail.TransName,
			TransTitle:       *detail.TransTitle,
			TransStatus:      *detail.TransStatus,
			TransClass:       *detail.TransClass,
			DepartureDate:    *detail.DepartureDate,
			DepartureTime:    *detail.DepartureTime,
			ArrivalTime:      *detail.ArrivalTime,
			TripDuration:     tripDuration,
			HarborSourceName: *detail.HarborSourceName,
			HarborDestName:   *detail.HarborDestName,
			MerchantName:     detail.MerchantName.String,
			MerchantPhone:    detail.MerchantPhone.String,
			MerchantPicture:  detail.MerchantPicture.String,
		}
	}

	var bookedBy []models.BookedByObj
	var guestDesc []models.GuestDescObj
	var accountBank models.AccountDesc
	if details[0].BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(details[0].BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if details[0].GuestDesc != "" {
		if errUnmarshal := json.Unmarshal([]byte(details[0].GuestDesc), &guestDesc); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if details[0].AccountBank != nil {
		if errUnmarshal := json.Unmarshal([]byte(*details[0].AccountBank), &accountBank); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	var currency string
	if details[0].Currency == 1 {
		currency = "USD"
	} else {
		currency = "IDR"
	}
	transport[0].TotalGuest = len(guestDesc)
	if len(transport) > 1 {
		transport[1].TotalGuest = len(guestDesc)
	}
	vaNumber := ""
	if details[0].VaNumber != nil {
		vaNumber = *details[0].VaNumber
	}

	results := &models.BookingExpDetailDto{
		Id:                     details[0].Id,
		GuestDesc:              guestDesc,
		BookedBy:               bookedBy,
		BookedByEmail:          details[0].BookedByEmail,
		BookingDate:            details[0].BookingDate,
		ExpiredDatePayment:     details[0].ExpiredDatePayment,
		CreatedDateTransaction: details[0].CreatedDateTransaction,
		UserId:                 details[0].UserId,
		Status:                 details[0].Status,
		TransactionStatus:      details[0].TransactionStatus,
		OrderId:                details[0].OrderId,
		TicketQRCode:           details[0].TicketQRCode,
		ExperienceAddOnId:      details[0].ExperienceAddOnId,
		TotalPrice:             details[0].TotalPrice,
		Currency:               currency,
		PaymentType:            details[0].PaymentType,
		AccountNumber:          vaNumber,
		AccountHolder:          accountBank.AccHolder,
		BankIcon:               details[0].Icon,
		ExperiencePaymentId:    details[0].ExperiencePaymentId,
		Transportation:         transport,
		MidtransUrl:            details[0].PaymentUrl,
		ExChangeRates:          details[0].ExChangeRates,
		ExChangeCurrency:       details[0].ExChangeCurrency,
	}

	return results, nil
}

func (b bookingExpUsecase) SendCharge(ctx context.Context, bookingId, paymentType string) (map[string]interface{}, error) {
	var data map[string]interface{}

	midtrans.SetupMidtrans()
	client := &http.Client{}

	booking, err := b.bookingExpRepo.GetByID(ctx, bookingId)
	if err != nil {
		return nil, err
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, err
		}
	}

	fullName := bookedBy[0].FullName
	email := bookedBy[0].Email

	var phoneNumber string
	if phoneStr, ok := bookedBy[0].PhoneNumber.(string); ok {
		phoneNumber = phoneStr
	} else if phoneInt, ok := bookedBy[0].PhoneNumber.(int); ok {
		phoneNumber = strconv.Itoa(phoneInt)
	}

	name := strings.Split(fullName, " ")
	var first, last string
	if len(name) < 2 {
		first = fullName
		last = fullName
	} else {
		first = name[0]
		last = name[1]
	}

	var charge midtrans.MidtransCharge
	charge.CustomerDetail = midtrans.CustomerDetail{
		FirstName: first,
		LastName:  last,
		Phone:     phoneNumber,
		Email:     email,
	}

	charge.TransactionDetails.GrossAmount = math.Round(booking.TotalPrice)
	charge.TransactionDetails.OrderID = booking.OrderId

	charge.EnablePayment = []string{paymentType}
	charge.OptionColorTheme = midtrans.OptionColorTheme{
		Primary:     "#c51f1f",
		PrimaryDark: "#1a4794",
		Secondary:   "#1fce38",
	}
	j, _ := json.Marshal(charge)
	fmt.Println(string(j))
	AUTH_STRING := b64.StdEncoding.EncodeToString([]byte(midtrans.Midclient.ServerKey + ":"))
	req, _ := http.NewRequest("POST", midtrans.TransactionEndpoint, bytes.NewBuffer(j))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+AUTH_STRING)

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return data, err
		}
		return data, nil
	} else {
		err := errors.New("MIDTRANS ERROR : " + resp.Status)
		return data, err
	}
}

func (b bookingExpUsecase) CountThisMonth(ctx context.Context) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	count, err := b.bookingExpRepo.CountThisMonth(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (b bookingExpUsecase) GetGrowthByMerchantID(ctx context.Context, token string) ([]*models.BookingGrowthDto, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	currentMerchant, err := b.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	growth, err := b.bookingExpRepo.GetGrowthByMerchantID(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}

	results := make([]*models.BookingGrowthDto, len(growth))
	for i, g := range growth {
		results[i] = &models.BookingGrowthDto{
			Date:  g.Date.Format("2006-01-02"),
			Count: g.Count,
		}
	}

	return results, nil
}

func (b bookingExpUsecase) GetByUserID(ctx context.Context, status string, token string, page, limit, offset int) (*models.MyBookingWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	currentUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, err
	}
	bookingIds, err := b.bookingExpRepo.GetBookingIdByUserID(ctx, status, currentUser.Id, limit, offset)
	if err != nil {
		return nil, err
	}

	expList, err := b.bookingExpRepo.GetBookingExpByUserID(ctx, bookingIds)
	if err != nil {
		return nil, err
	}

	myBooking := make([]*models.MyBooking, len(expList))
	for i, b := range expList {
		var guestDesc []models.GuestDescObj
		if b.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(b.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, err
			}
		}
		expType := make([]string, 0)
		if b.ExpType != nil {
			if errUnmarshal := json.Unmarshal([]byte(*b.ExpType), &expType); errUnmarshal != nil {
				return nil, err
			}
		}
		var expGuest models.TotalGuestTransportation
		if len(guestDesc) > 0 {
			for _, guest := range guestDesc {
				if guest.Type == "Adult" {
					expGuest.Adult = expGuest.Adult + 1
				} else if guest.Type == "Children" {
					expGuest.Children = expGuest.Children + 1
				}
			}
		}
		myBooking[i] = &models.MyBooking{
			OrderId:     b.OrderId,
			ExpType:     expType,
			ExpId:       *b.ExpId,
			ExpTitle:    *b.ExpTitle,
			BookingDate: b.BookingDate,
			ExpDuration: *b.ExpDuration,
			TotalGuest:  len(guestDesc),
			ExpGuest:    expGuest,
			City:        b.City,
			Province:    b.Province,
			Country:     b.Country,
		}
	}

	transList, err := b.bookingExpRepo.GetBookingTransByUserID(ctx, bookingIds)
	if err != nil {
		return nil, err
	}
	for _, b := range transList {
		var guestDesc []models.GuestDescObj
		if b.GuestDesc != "" {
			if errUnmarshal := json.Unmarshal([]byte(b.GuestDesc), &guestDesc); errUnmarshal != nil {
				return nil, err
			}
		}
		var transGuest models.TotalGuestTransportation
		if len(guestDesc) > 0 {
			for _, guest := range guestDesc {
				if guest.Type == "Adult" {
					transGuest.Adult = transGuest.Adult + 1
				} else if guest.Type == "Children" {
					transGuest.Children = transGuest.Children + 1
				}
			}
		}
		var tripDuration string
		if b.DepartureTime != nil && b.ArrivalTime != nil {
			departureTime, _ := time.Parse("15:04:00", *b.DepartureTime)
			arrivalTime, _ := time.Parse("15:04:00", *b.ArrivalTime)

			tripHour := arrivalTime.Hour() - departureTime.Hour()
			tripMinute := arrivalTime.Minute() - departureTime.Minute()
			tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
		}
		booking := models.MyBooking{
			OrderId:            b.OrderId,
			ExpId:              "",
			ExpTitle:           "",
			TransId:            *b.TransId,
			TransName:          *b.TransName,
			TransFrom:          *b.HarborDestName,
			TransTo:            *b.HarborSourceName,
			TransDepartureTime: b.DepartureTime,
			TransArrivalTime:   b.ArrivalTime,
			TripDuration:       tripDuration,
			TransClass:         *b.TransClass,
			TransGuest:         transGuest,
			BookingDate:        b.BookingDate,
			ExpDuration:        0,
			TotalGuest:         len(guestDesc),
			City:               b.City,
			Province:           b.Province,
			Country:            b.Country,
		}

		myBooking = append(myBooking, &booking)
	}

	totalRecords, _ := b.bookingExpRepo.GetBookingCountByUserID(ctx, status, currentUser.Id)

	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(myBooking),
	}

	response := &models.MyBookingWithPagination{
		Data: myBooking,
		Meta: meta,
	}
	return response, nil
}

func (b bookingExpUsecase) GetDetailBookingID(c context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	getDetailBooking, err := b.bookingExpRepo.GetDetailBookingID(ctx, bookingId, bookingCode)
	if err != nil {
		return nil, err
	}
	var bookedBy []models.BookedByObj
	var guestDesc []models.GuestDescObj
	var accountBank models.AccountDesc
	var expType []string
	if getDetailBooking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.BookedBy), &bookedBy); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if getDetailBooking.GuestDesc != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailBooking.GuestDesc), &guestDesc); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if getDetailBooking.ExpType != nil {
		if errUnmarshal := json.Unmarshal([]byte(*getDetailBooking.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	if getDetailBooking.AccountBank != nil {
		if errUnmarshal := json.Unmarshal([]byte(*getDetailBooking.AccountBank), &accountBank); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	var currency string
	if getDetailBooking.Currency == 1 {
		currency = "USD"
	} else {
		currency = "IDR"
	}
	var experiencePaymentType *models.ExperiencePaymentTypeDto
	if getDetailBooking.ExperiencePaymentId != "" {
		query, err := b.experiencePaymentTypeRepo.GetByExpID(ctx, *getDetailBooking.ExpId)
		if err != nil {

		}
		for _, element := range query {
			if element.Id == getDetailBooking.ExperiencePaymentId {
				paymentType := models.ExperiencePaymentTypeDto{
					Id:   element.ExpPaymentTypeId,
					Name: element.ExpPaymentTypeName,
					Desc: element.ExpPaymentTypeDesc,
				}
				if paymentType.Name == "Down Payment" {
					remainingPayment := element.Price - *getDetailBooking.TotalPrice
					paymentType.RemainingPayment = remainingPayment
				} else {
					paymentType.RemainingPayment = 0
				}
				experiencePaymentType = &paymentType
			}
		}
	}
	expAddOns := make([]models.ExperienceAddOnObj, 0)
	expAddOnsQuery, errorQuery := b.adOnsRepo.GetByExpId(ctx, *getDetailBooking.ExpId)
	if errorQuery != nil {
		return nil, err
	}
	if expAddOnsQuery != nil {
		for _, element := range expAddOnsQuery {
			var currency string
			if element.Currency == 1 {
				currency = "USD"
			} else {
				currency = "IDR"
			}
			addOns := models.ExperienceAddOnObj{
				Id:       element.Id,
				Name:     element.Name,
				Desc:     element.Desc,
				Currency: currency,
				Amount:   element.Amount,
			}
			expAddOns = append(expAddOns, addOns)
		}
	}

	vaNumber := ""
	if getDetailBooking.VaNumber != nil {
		vaNumber = *getDetailBooking.VaNumber
	}

	expDetail := make([]models.BookingExpDetail, 1)
	expDetail[0] = models.BookingExpDetail{
		ExpId:                    *getDetailBooking.ExpId,
		ExpTitle:                 *getDetailBooking.ExpTitle,
		ExpType:                  expType,
		ExpPickupPlace:           *getDetailBooking.ExpPickupPlace,
		ExpPickupTime:            *getDetailBooking.ExpPickupTime,
		MerchantName:             getDetailBooking.MerchantName.String,
		MerchantPhone:            getDetailBooking.MerchantPhone.String,
		MerchantPicture:          getDetailBooking.MerchantPicture.String,
		TotalGuest:               len(guestDesc),
		City:                     getDetailBooking.City,
		ProvinceName:             getDetailBooking.Province,
		ExpDuration:              *getDetailBooking.ExpDuration,
		HarborsName:              *getDetailBooking.HarborsName,
		ExperienceAddOn:          expAddOns,
		CountryName:              getDetailBooking.Country,
		ExpPaymentDeadlineAmount: getDetailBooking.ExpPaymentDeadlineAmount,
		ExpPaymentDeadlineType:   getDetailBooking.ExpPaymentDeadlineType,
	}
	if getDetailBooking.UserId == nil {
		getDetailBooking.UserId = new(string)
	}
	if getDetailBooking.ExpId == nil {
		getDetailBooking.ExpId = new(string)
	}
	reviews, _ := b.reviewRepo.GetByExpId(ctx, *getDetailBooking.ExpId, "", 0, 1, 0, *getDetailBooking.UserId)
	if err != nil {
		return nil, err
	}
	var isReview = false
	bookingExp := models.BookingExpDetailDto{
		Id:                     getDetailBooking.Id,
		OrderId:                getDetailBooking.OrderId,
		GuestDesc:              guestDesc,
		BookedBy:               bookedBy,
		BookedByEmail:          getDetailBooking.BookedByEmail,
		BookingDate:            getDetailBooking.BookingDate,
		ExpiredDatePayment:     getDetailBooking.ExpiredDatePayment,
		CreatedDateTransaction: getDetailBooking.CreatedDateTransaction,
		UserId:                 getDetailBooking.UserId,
		Status:                 getDetailBooking.Status,
		TransactionStatus:      getDetailBooking.TransactionStatus,
		//TicketCode:        getDetailBooking.TicketCode,
		TicketQRCode:          getDetailBooking.TicketQRCode,
		ExperienceAddOnId:     getDetailBooking.ExperienceAddOnId,
		TotalPrice:            getDetailBooking.TotalPrice,
		Currency:              currency,
		PaymentType:           getDetailBooking.PaymentType,
		AccountNumber:         vaNumber,
		AccountHolder:         accountBank.AccHolder,
		BankIcon:              getDetailBooking.Icon,
		ExperiencePaymentId:   getDetailBooking.ExperiencePaymentId,
		Experience:            expDetail,
		ExperiencePaymentType: experiencePaymentType,
		IsReview:              isReview,
		MidtransUrl:           getDetailBooking.PaymentUrl,
		ExChangeRates:         getDetailBooking.ExChangeRates,
		ExChangeCurrency:      getDetailBooking.ExChangeCurrency,
	}
	if len(reviews) != 0 {
		desc := models.ReviewDtoObject{}
		if reviews[0].Desc != "" {
			if errUnmarshal := json.Unmarshal([]byte(reviews[0].Desc), &desc); errUnmarshal != nil {
				return nil, models.ErrInternalServerError
			}
		}
		bookingExp.ReviewDesc = &desc.Desc
		bookingExp.IsReview = true
		bookingExp.GuideReview = reviews[0].GuideReview
		bookingExp.ActivitiesReview = reviews[0].ActivitiesReview
		bookingExp.ServiceReview = reviews[0].ServiceReview
		bookingExp.CleanlinessReview = reviews[0].CleanlinessReview
		bookingExp.ValueReview = reviews[0].ValueReview
	}
	return &bookingExp, nil

}

func (b bookingExpUsecase) Insert(c context.Context, booking *models.NewBookingExpCommand, transReturnId, scheduleReturnId, token string) ([]*models.NewBookingExpCommand, error, error) {

	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()

	if booking.ExpId == "" && booking.TransId == nil && booking.ScheduleId != nil {
		return nil, models.ValidationExpId, nil
	}
	if booking.BookingDate == "" {
		return nil, models.ValidationBookedDate, nil
	}
	if booking.Status == "" {
		return nil, models.ValidationStatus, nil
	}
	if booking.BookedBy == "" {
		return nil, models.ValidationBookedBy, nil
	}
	layoutFormat := "2006-01-02 15:04:05"
	bookingDate, errDate := time.Parse(layoutFormat, booking.BookingDate)
	if errDate != nil {
		return nil, errDate, nil
	}
	orderId, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}

	// re-generate if duplicate order id
	if b.bookingExpRepo.CheckBookingCode(ctx, orderId) {
		orderId, err = generateRandomString(12)
		if err != nil {
			return nil, models.ErrInternalServerError, nil
		}
	}

	ticketCode, err := generateRandomString(12)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}
	var createdBy string
	if token != "" {
		currentUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
		if err != nil {
			return nil, err, nil
		}
		createdBy = currentUser.UserEmail
	} else {
		createdBy = booking.BookedByEmail
	}
	booking.OrderId = orderId
	booking.TicketCode = ticketCode
	fileNameQrCode, err := generateQRCode(orderId)
	if err != nil {
		return nil, models.ErrInternalServerError, nil
	}
	imagePath, _ := b.isUsecase.UploadFileToBlob(*fileNameQrCode, "TicketBookingQRCode")

	errRemove := os.Remove(*fileNameQrCode)
	if errRemove != nil {
		return nil, models.ErrInternalServerError, nil
	}
	booking.TicketQRCode = imagePath

	reqBooking := make([]*models.BookingExp, 0)

	bookingExp := models.BookingExp{
		Id:                "",
		CreatedBy:         createdBy,
		CreatedDate:       time.Now(),
		ModifiedBy:        nil,
		ModifiedDate:      nil,
		DeletedBy:         nil,
		DeletedDate:       nil,
		IsDeleted:         0,
		IsActive:          1,
		ExpId:             &booking.ExpId,
		OrderId:           orderId,
		GuestDesc:         booking.GuestDesc,
		BookedBy:          booking.BookedBy,
		BookedByEmail:     booking.BookedByEmail,
		BookingDate:       bookingDate,
		UserId:            booking.UserId,
		Status:            0,
		TicketCode:        ticketCode,
		TicketQRCode:      imagePath,
		ExperienceAddOnId: booking.ExperienceAddOnId,
		TransId:           booking.TransId,
		ScheduleId:        booking.ScheduleId,
	}
	if *bookingExp.ExperienceAddOnId == "" {
		bookingExp.ExperienceAddOnId = nil
	}
	if *bookingExp.UserId == "" {
		bookingExp.UserId = nil
	}
	if *bookingExp.TransId == "" {
		bookingExp.TransId = nil
	}
	if *bookingExp.ExpId == "" {
		bookingExp.ExpId = nil
	}
	if *bookingExp.ScheduleId == "" {
		bookingExp.ScheduleId = nil
	}

	reqBooking = append(reqBooking, &bookingExp)

	if transReturnId != "" && scheduleReturnId != "" {
		bookingReturn := models.BookingExp{
			Id:                "",
			CreatedBy:         createdBy,
			CreatedDate:       time.Now(),
			ModifiedBy:        nil,
			ModifiedDate:      nil,
			DeletedBy:         nil,
			DeletedDate:       nil,
			IsDeleted:         0,
			IsActive:          1,
			ExpId:             &booking.ExpId,
			OrderId:           orderId,
			GuestDesc:         booking.GuestDesc,
			BookedBy:          booking.BookedBy,
			BookedByEmail:     booking.BookedByEmail,
			BookingDate:       bookingDate,
			UserId:            booking.UserId,
			Status:            0,
			TicketCode:        ticketCode,
			TicketQRCode:      imagePath,
			ExperienceAddOnId: booking.ExperienceAddOnId,
			TransId:           &transReturnId,
			ScheduleId:        &scheduleReturnId,
		}
		if *bookingReturn.ExperienceAddOnId == "" {
			bookingReturn.ExperienceAddOnId = nil
		}
		if *bookingReturn.UserId == "" {
			bookingReturn.UserId = nil
		}
		if *bookingReturn.TransId == "" {
			bookingReturn.TransId = nil
		}
		if *bookingReturn.ExpId == "" {
			bookingReturn.ExpId = nil
		}
		if *bookingExp.ScheduleId == "" {
			bookingExp.ScheduleId = nil
		}

		reqBooking = append(reqBooking, &bookingReturn)
	}

	resBooking := make([]*models.NewBookingExpCommand, len(reqBooking))
	for i, req := range reqBooking {
		res, err := b.bookingExpRepo.Insert(ctx, req)
		if err != nil {
			return nil, err, nil
		}
		reqBooking[i].Id = res.Id
		resBooking[i] = &models.NewBookingExpCommand{
			Id:                res.Id,
			ExpId:             booking.ExpId,
			GuestDesc:         res.GuestDesc,
			BookedBy:          res.BookedBy,
			BookedByEmail:     res.BookedByEmail,
			BookingDate:       res.BookingDate.String(),
			UserId:            res.UserId,
			Status:            strconv.Itoa(res.Status),
			OrderId:           res.OrderId,
			TicketCode:        res.TicketCode,
			TicketQRCode:      res.TicketQRCode,
			ExperienceAddOnId: res.ExperienceAddOnId,
			TransId:           res.TransId,
			ScheduleId:        res.ScheduleId,
		}
	}

	return resBooking, nil, nil
}

func (b bookingExpUsecase) GetHistoryBookingByUserId(c context.Context, token string, monthType string, page, limit, offset int) (*models.BookingHistoryDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(c, b.contextTimeout)
	defer cancel()
	var currentUserId string
	if token != "" {
		validateUser, err := b.userUsecase.ValidateTokenUser(ctx, token)
		if err != nil {
			return nil, err
		}
		currentUserId = validateUser.Id
	}
	var guestDesc []models.GuestDescObj
	var result []*models.BookingHistoryDto
	if monthType == "past-30-days" {
		bookingIds, err := b.bookingExpRepo.QuerySelectIdHistoryByUserId(ctx, currentUserId, "", limit, offset)
		if err != nil {
			return nil, err
		}
		query, err := b.bookingExpRepo.QueryHistoryPer30DaysExpByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		historyDto := models.BookingHistoryDto{
			Category: "past-30-days",
			Items:    nil,
		}
		for _, element := range query {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}

			var expGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						expGuest.Adult = expGuest.Adult + 1
					} else if guest.Type == "Children" {
						expGuest.Children = expGuest.Children + 1
					}
				}
			}
			totalGuest := len(guestDesc)
			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if element.StatusTransaction == 0 || element.StatusTransaction == 3 {
					status = "Payment Expired"
				} else if element.StatusTransaction == 1 || element.StatusTransaction == 4 || element.StatusTransaction == 5 {
					status = "Cancelled"
				} else if element.StatusTransaction == 2 {
					status = "Success"
				}
			} else {
				if element.StatusTransaction == 0 && time.Now().Add(7*time.Hour).After(element.ExpiredDatePayment.Add(7*time.Hour)) {
					status = "Payment Expired"
				} else if element.StatusTransaction == 3 || element.StatusTransaction == 4 {
					status = "Cancelled"
				}
			}

			if element.UserId == nil {
				element.UserId = new(string)
			}
			checkReview, _ := b.reviewRepo.GetByExpId(ctx, element.ExpId, "", 0, 1, 0, *element.UserId)
			if err != nil {
				return nil, err
			}
			var isReview = false
			if len(checkReview) != 0 {
				isReview = true
			}

			itemDto := models.ItemsHistoryDto{
				OrderId:        element.OrderId,
				ExpId:          element.ExpId,
				ExpTitle:       element.ExpTitle,
				ExpType:        expType,
				ExpBookingDate: element.BookingDate,
				ExpDuration:    element.ExpDuration,
				TotalGuest:     totalGuest,
				ExpGuest:       expGuest,
				City:           element.CityName,
				Province:       element.ProvinceName,
				Country:        element.CountryName,
				Status:         status,
				IsReview:       isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}

		queryTrans, err := b.bookingExpRepo.QueryHistoryPer30DaysTransByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		for _, element := range queryTrans {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			var transGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						transGuest.Adult = transGuest.Adult + 1
					} else if guest.Type == "Children" {
						transGuest.Children = transGuest.Children + 1
					}
				}
			}
			//totalGuest := len(guestDesc)
			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if *element.TransactionStatus == 0 || *element.TransactionStatus == 3 {
					status = "Payment Expired"
				} else if *element.TransactionStatus == 1 || *element.TransactionStatus == 4 || *element.TransactionStatus == 5 {
					status = "Cancelled"
				} else if *element.TransactionStatus == 2 {
					status = "Success"
				}
			} else {
				if *element.TransactionStatus == 0 && time.Now().Add(7*time.Hour).After(element.ExpiredDatePayment.Add(7*time.Hour)) {
					status = "Payment Expired"
				} else if *element.TransactionStatus == 3 || *element.TransactionStatus == 4 {
					status = "Cancelled"
				}
			}
			var tripDuration string
			if element.DepartureTime != nil && element.ArrivalTime != nil {
				departureTime, _ := time.Parse("15:04:00", *element.DepartureTime)
				arrivalTime, _ := time.Parse("15:04:00", *element.ArrivalTime)

				tripHour := arrivalTime.Hour() - departureTime.Hour()
				tripMinute := arrivalTime.Minute() - departureTime.Minute()
				tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
			}
			//if element.UserId == nil {
			//	*element.UserId = ""
			//}
			//if element.ExpId == nil {
			//	*element.ExpId = ""
			//}
			//checkReview ,_:= b.reviewRepo.GetByExpId(ctx , *element.ExpId,"",0,0,1,*element.UserId)
			//if err != nil {
			//	return nil,err
			//}
			var isReview = false
			//if len(checkReview) != 0 {
			//	isReview = true
			//}
			itemDto := models.ItemsHistoryDto{
				OrderId:            element.OrderId,
				ExpId:              "",
				ExpTitle:           "",
				ExpType:            nil,
				TransId:            *element.TransId,
				TransName:          *element.TransName,
				TransFrom:          *element.HarborDestName,
				TransTo:            *element.HarborSourceName,
				TransDepartureTime: element.DepartureTime,
				TransArrivalTime:   element.ArrivalTime,
				TripDuration:       tripDuration,
				TransClass:         *element.TransClass,
				TransGuest:         transGuest,
				ExpBookingDate:     element.BookingDate,
				ExpDuration:        0,
				TotalGuest:         0,
				City:               element.City,
				Province:           element.Province,
				Country:            element.Country,
				Status:             status,
				IsReview:           isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}
		result = append(result, &historyDto)
	} else {
		bookingIds, err := b.bookingExpRepo.QuerySelectIdHistoryByUserId(ctx, currentUserId, monthType, limit, offset)
		if err != nil {
			return nil, err
		}
		queryExp, err := b.bookingExpRepo.QueryHistoryPerMonthExpByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		monthType = monthType + "-" + "01" + " 00:00:00"
		layoutFormat := "2006-01-02 15:04:05"
		dt, _ := time.Parse(layoutFormat, monthType)
		dtstr2 := dt.Format("Jan '06")
		historyDto := models.BookingHistoryDto{
			Category: dtstr2,
			Items:    nil,
		}
		for _, element := range queryExp {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			var expGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						expGuest.Adult = expGuest.Adult + 1
					} else if guest.Type == "Children" {
						expGuest.Children = expGuest.Children + 1
					}
				}
			}
			totalGuest := len(guestDesc)

			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if element.StatusTransaction == 0 || element.StatusTransaction == 3 {
					status = "Payment Expired"
				} else if element.StatusTransaction == 1 || element.StatusTransaction == 4 || element.StatusTransaction == 5 {
					status = "Cancelled"
				} else if element.StatusTransaction == 2 {
					status = "Success"
				}
			} else {
				if element.StatusTransaction == 0 && time.Now().Add(7*time.Hour).After(element.ExpiredDatePayment.Add(7*time.Hour)) {
					status = "Payment Expired"
				} else if element.StatusTransaction == 3 || element.StatusTransaction == 4 {
					status = "Cancelled"
				}
			}
			if element.UserId == nil {
				element.UserId = new(string)
			}
			checkReview, err := b.reviewRepo.GetByExpId(ctx, element.ExpId, "", 0, 1, 0, *element.UserId)
			if err != nil {
				return nil, err
			}
			var isReview = false
			if len(checkReview) != 0 {
				isReview = true
			}
			itemDto := models.ItemsHistoryDto{
				OrderId:        element.OrderId,
				ExpId:          element.ExpId,
				ExpTitle:       element.ExpTitle,
				ExpType:        expType,
				ExpBookingDate: element.BookingDate,
				ExpDuration:    element.ExpDuration,
				TotalGuest:     totalGuest,
				ExpGuest:       expGuest,
				City:           element.CityName,
				Province:       element.ProvinceName,
				Country:        element.CountryName,
				Status:         status,
				IsReview:       isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}

		queryTrans, err := b.bookingExpRepo.QueryHistoryPerMonthTransByUserId(ctx, bookingIds)
		if err != nil {
			return nil, err
		}
		for _, element := range queryTrans {
			var expType []string
			if element.ExpType != nil {
				if errUnmarshal := json.Unmarshal([]byte(*element.ExpType), &expType); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			if element.GuestDesc != "" {
				if errUnmarshal := json.Unmarshal([]byte(element.GuestDesc), &guestDesc); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
			}
			var transGuest models.TotalGuestTransportation
			if len(guestDesc) > 0 {
				for _, guest := range guestDesc {
					if guest.Type == "Adult" {
						transGuest.Adult = transGuest.Adult + 1
					} else if guest.Type == "Children" {
						transGuest.Children = transGuest.Children + 1
					}
				}
			}
			//totalGuest := len(guestDesc)
			var status string
			if element.BookingDate.Before(time.Now()) == true {
				if *element.TransactionStatus == 0 || *element.TransactionStatus == 3 {
					status = "Payment Expired"
				} else if *element.TransactionStatus == 1 || *element.TransactionStatus == 4 || *element.TransactionStatus == 5 {
					status = "Cancelled"
				} else if *element.TransactionStatus == 2 {
					status = "Success"
				}
			} else {
				if *element.TransactionStatus == 0 && time.Now().Add(7*time.Hour).After(element.ExpiredDatePayment.Add(7*time.Hour)) {
					status = "Payment Expired"
				} else if *element.TransactionStatus == 3 || *element.TransactionStatus == 4 {
					status = "Cancelled"
				}
			}
			var tripDuration string
			if element.DepartureTime != nil && element.ArrivalTime != nil {
				departureTime, _ := time.Parse("15:04:00", *element.DepartureTime)
				arrivalTime, _ := time.Parse("15:04:00", *element.ArrivalTime)

				tripHour := arrivalTime.Hour() - departureTime.Hour()
				tripMinute := arrivalTime.Minute() - departureTime.Minute()
				tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
			}
			//if element.UserId == nil {
			//	*element.UserId = ""
			//}
			//if element.ExpId == nil {
			//	*element.ExpId = ""
			//}
			//checkReview ,_:= b.reviewRepo.GetByExpId(ctx , *element.ExpId,"",0,0,1,*element.UserId)
			//if err != nil {
			//	return nil,err
			//}
			var isReview = false
			//if len(checkReview) != 0 {
			//	isReview = true
			//}
			itemDto := models.ItemsHistoryDto{
				OrderId:            element.OrderId,
				ExpId:              "",
				ExpTitle:           "",
				ExpType:            nil,
				TransId:            *element.TransId,
				TransName:          *element.TransName,
				TransFrom:          *element.HarborDestName,
				TransTo:            *element.HarborSourceName,
				TransDepartureTime: element.DepartureTime,
				TransArrivalTime:   element.ArrivalTime,
				TripDuration:       tripDuration,
				TransClass:         *element.TransClass,
				TransGuest:         transGuest,
				ExpBookingDate:     element.BookingDate,
				ExpDuration:        0,
				TotalGuest:         0,
				City:               element.City,
				Province:           element.Province,
				Country:            element.Country,
				Status:             status,
				IsReview:           isReview,
			}
			historyDto.Items = append(historyDto.Items, itemDto)
		}
		result = append(result, &historyDto)
	}
	var totalRecords int
	if monthType == "past-30-days" {
		totalRecords, _ = b.bookingExpRepo.QueryCountHistoryByUserId(ctx, currentUserId, "")
	} else {
		totalRecords, _ = b.bookingExpRepo.QueryCountHistoryByUserId(ctx, currentUserId, monthType)
	}
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(result[0].Items),
	}

	response := &models.BookingHistoryDtoWithPagination{
		Data: result,
		Meta: meta,
	}
	return response, nil
}

func rangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}

func generateQRCode(content string) (*string, error) {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	name, err := generateRandomString(5)
	if err != nil {
		return nil, err
	}

	fileName := name + ".png"
	err = ioutil.WriteFile(fileName, png, 0700)
	copy, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	copy.Close()
	return &fileName, nil

	//err := qrcode.WriteFile("https://example.org", qrcode.Medium, 256, "qr.png")

}
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
