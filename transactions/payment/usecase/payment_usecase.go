package usecase

import (
	"bytes"
	"context"
	"github.com/auth/merchant"
	guuid "github.com/google/uuid"
	"github.com/service/experience"
	"github.com/service/transportation"
	"html/template"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/auth/identityserver"

	"github.com/misc/notif"
	"github.com/transactions/transaction"

	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/transactions/payment"
)

type paymentUsecase struct {
	merchantUsecase  merchant.Usecase
	isUsecase        identityserver.Usecase
	transactionRepo  transaction.Repository
	notificationRepo notif.Repository
	paymentRepo      payment.Repository
	userUsercase     user.Usecase
	bookingRepo      booking_exp.Repository
	userRepo         user.Repository
	contextTimeout   time.Duration
	bookingUsecase   booking_exp.Usecase
	expRepo 		experience.Repository
	transportationRepo transportation.Repository
}


// NewPaymentUsecase will create new an paymentUsecase object representation of payment.Usecase interface
func NewPaymentUsecase(	merchantUsecase  merchant.Usecase	,transportationRepo transportation.Repository,expRepo experience.Repository,bookingUsecase booking_exp.Usecase, isUsecase identityserver.Usecase, t transaction.Repository, n notif.Repository, p payment.Repository, u user.Usecase, b booking_exp.Repository, ur user.Repository, timeout time.Duration) payment.Usecase {
	return &paymentUsecase{
		merchantUsecase:merchantUsecase,
		transportationRepo:transportationRepo,
		expRepo:expRepo,
		bookingUsecase:   bookingUsecase,
		isUsecase:        isUsecase,
		transactionRepo:  t,
		notificationRepo: n,
		paymentRepo:      p,
		userUsercase:     u,
		bookingRepo:      b,
		userRepo:         ur,
		contextTimeout:   timeout,
	}
}

const (
	templateBookingRejected string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
 <head>
  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
  <link href="https://fonts.googleapis.com/css2?family=Nunito+Sans:wght@400;600;700;800&display=swap" rel="stylesheet" type="text/css">
	<link href="https://fonts.googleapis.com/css2?family=Rubik:wght@500&display=swap" rel="stylesheet" type="text/css">
    <title>Booking Rejected (DP/FP)</title>
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
                            We regret to inform you that your trip <b>{{.title}}</b> with trip date on <b>{{.tripDate}}</b> was cancelled. This cancellation occurs because technical preparations needed for the trip are not available.
                        </td>
                    </tr>
                    <tr>
                        <td style="padding: 10px 0 20px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            If you wish to apply for a refund, please submit your bank account information <br> and you will receive your refund within <font color="red">3 working days</font>.
                            If you wish your <br> payment to be transmitted to credits, please click transmit to credits button.
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 45px 0 20px 0;     text-align: center;">
                            <a href="#" style="    background: linear-gradient(145deg, rgba(55,123,232,1) 0%, rgba(35,62,152,1) 42%, rgba(35,62,152,1) 100%);
                            color: white;
                            text-decoration: none;
                            font-weight: 600;
                            padding: 12px 2rem;
                            border-radius: 2rem;
                            font-family: 'Nunito Sans', sans-serif;
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
                        <td style="padding: 25px 0 10px 0; font-family: 'Nunito Sans', sans-serif;
                        font-style: normal;
                        font-weight: normal;
                        font-size: 15px;
                        line-height: 24px;">
                            As a valued customer, your satisfaction is one of our concerns and we apologize for any inconvenience this cancellation caused. We suggest you to book another trip or choose different trip dates.
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
                                                <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/Unknown.png" alt="" width="53" height="24" style="object-fit: cover;" />
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

	templateTicketTransportationWithReturn string = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
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
                                                line-height: 24px;">{{.sourceTimeDeparture}}</b>
                                            </td>
                                            <td style="color: #8E8E8E;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;" width="15%">
                                                {{.durationDeparture}}
                                            </td>
                                            <td style="color: black;">
                                                <b style="font-size: 17px;font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                line-height: 24px;">{{.desTimeDeparture}}</b>
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
                                                {{.sourceDeparture}}
                                            </td>
                                            <td width="15%">
                                                <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/arrow-back.png" alt="">
                                            </td>
                                            <td style="color: black; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                 {{.destDeparture}}
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="color:#7E7E7E; font-weight:600 !important;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.classDeparture}}</td>
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
                                            line-height: 24px;">{{.tripDateDeparture}}</b>
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
                                                line-height: 24px;">{{.guestCountDeparture}}</b>
                                            </td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                           </table>
                     </td>
                    </tr>
					<tr>
						<td bgcolor="#FFFFF" width="200px">&nbsp;</td>
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
                                                line-height: 24px;">{{.sourceTimeReturn}}</b>
                                            </td>
                                            <td style="color: #8E8E8E;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;" width="15%">
                                                {{.durationReturn}}
                                            </td>
                                            <td style="color: black;">
                                                <b style="font-size: 17px;font-family: 'Nunito Sans', sans-serif;
                                                font-style: normal;
                                                font-weight: 700;
                                                line-height: 24px;">{{.desTimeReturn}}</b>
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
                                                {{.sourceReturn}}
                                            </td>
                                            <td width="15%">
                                                <img src="https://cgostorage.blob.core.windows.net/cgo-storage/img/img/arrow-back.png" alt="">
                                            </td>
                                            <td style="color: black; font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-weight: normal;
                                            font-size: 15px;
                                            line-height: 24px;">
                                                 {{.destReturn}}
                                            </td>
                                        </tr>
                                        <tr>
                                            <td style="color:#7E7E7E; font-weight:600 !important;font-family: 'Nunito Sans', sans-serif;
                                            font-style: normal;
                                            font-size: 15px;
                                            line-height: 24px;">{{.classReturn}}</td>
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
                                            line-height: 24px;">{{.tripDateReturn}}</b>
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
                                                line-height: 24px;">{{.guestCountReturn}}</b>
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
)

var templateFuncs = template.FuncMap{"rangeStruct": rangeStructer}

func (p paymentUsecase) ConfirmPaymentBoarding(ctx context.Context, orderId string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()
	_, err := p.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}
	var bookingId string
	checkOrderId,_ := p.bookingRepo.GetDetailBookingID(ctx,orderId,orderId)
	if checkOrderId == nil{
		checkOrderIdTranspotation,_ := p.bookingRepo.GetDetailTransportBookingID(ctx,orderId,orderId,nil)
		if len(checkOrderIdTranspotation) == 0{
			return nil,models.ErrNotFound
		}else if *checkOrderIdTranspotation[0].TransactionStatus != 2{
			return nil,models.CheckBoarding
		}else {
			bookingId = checkOrderIdTranspotation[0].Id
			if err := p.transactionRepo.UpdateAfterPayment(ctx, 6, "", "", orderId); err != nil {
				return nil,err
			}
			result := &models.ResponseDelete{
				Id:      orderId,
				Message: "Success Updating Status",
			}
			return result,nil
		}

	}else if *checkOrderId.TransactionStatus != 2{
		return nil,models.CheckBoarding
	}else {
		bookingId = checkOrderId.Id
		if err := p.transactionRepo.UpdateAfterPayment(ctx, 6, "", "", bookingId); err != nil {
			return nil,err
		}
		result := &models.ResponseDelete{
			Id:      orderId,
			Message: "Success Updating Status",
		}
		return result,nil
	}


}

func (p paymentUsecase) Insert(ctx context.Context, payment *models.Transaction, token string, points float64,autoComplete bool) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	var userId string

	if payment.PaymentMethodId == "" {
		return "", models.PaymentMethodIdRequired
	}

	if payment.Currency == "" {
		payment.Currency = "IDR"
	}
	bookingCode := payment.OrderId
	if payment.BookingExpId != nil {
		bookingCode = payment.BookingExpId
	}
	createdBy, err := p.bookingRepo.GetEmailByID(ctx, *bookingCode)
	if err != nil {
		return "", err
	}
	if token != "" {
		currentUser, err := p.userUsercase.ValidateTokenUser(ctx, token)
		if err != nil {
			return "", err
		}
		createdBy = currentUser.UserEmail
		userId = currentUser.Id
	}

	newData := &models.Transaction{
		Id:                  "",
		CreatedBy:           createdBy,
		CreatedDate:         time.Now(),
		ModifiedBy:          nil,
		ModifiedDate:        nil,
		DeletedBy:           nil,
		DeletedDate:         nil,
		IsDeleted:           0,
		IsActive:            1,
		BookingType:         payment.BookingType,
		BookingExpId:        payment.BookingExpId,
		PromoId:             payment.PromoId,
		PaymentMethodId:     payment.PaymentMethodId,
		ExperiencePaymentId: payment.ExperiencePaymentId,
		Status:              payment.Status,
		TotalPrice:          payment.TotalPrice,
		Currency:            payment.Currency,
		OrderId:             payment.OrderId,
		ExChangeRates:       payment.ExChangeRates,
		ExChangeCurrency:    payment.ExChangeCurrency,
		Points:payment.Points,
		OriginalPrice:payment.OriginalPrice,
	}
	if autoComplete == true {
		exp, err := p.expRepo.GetExperienceByBookingId(ctx, *newData.BookingExpId,*newData.ExperiencePaymentId)
		if err != models.ErrNotFound{
			if exp.ExpBookingType == "No Instant Booking" {
				newData.Status = 1
			}else if exp.ExpBookingType == "Instant Booking" && exp.ExpPaymentTypeName == "Down Payment" {
				newData.Status = 5
			}else if exp.ExpBookingType == "Instant Booking" && exp.ExpPaymentTypeName == "Full Payment" {
				newData.Status = 2
			}
		}else {
			_ ,err := p.transportationRepo.GetTransportationByBookingId(ctx,*newData.BookingExpId)
			if err == models.ErrNotFound{

			}else {
				newData.Status = 2
			}
		}
	}
	res, err := p.paymentRepo.Insert(ctx, newData)
	if err != nil {
		return "", models.ErrInternalServerError
	}

	expiredPayment := res.CreatedDate.Add(2 * time.Hour)
	err = p.bookingRepo.UpdateStatus(ctx, *bookingCode, expiredPayment)
	if err != nil {
		return "", err
	}

	if points != 0 {
		err = p.userRepo.UpdatePointByID(ctx, points, userId)
		if err != nil {
			return "", err
		}
	}

	return res.Id, nil
}

func (p paymentUsecase) ConfirmPaymentByDate(ctx context.Context, confirmIn *models.ConfirmTransactionPayment) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	err := p.paymentRepo.ChangeStatusTransByDate(ctx, confirmIn)
	if err != nil {
		return err
	}
	if confirmIn.ExpId != "" {
		listTransaction , _ := p.transactionRepo.GetByBookingDate(ctx,confirmIn.BookingDate,"",confirmIn.ExpId)
		for _,getTransaction := range listTransaction{
			if err != nil {
				return err
			}
			notif := models.Notification{
				Id:           guuid.New().String(),
				CreatedBy:    getTransaction.CreatedBy,
				CreatedDate:  time.Now(),
				ModifiedBy:   nil,
				ModifiedDate: nil,
				DeletedBy:    nil,
				DeletedDate:  nil,
				IsDeleted:    0,
				IsActive:     0,
				MerchantId:   getTransaction.MerchantId,
				Type:         0,
				Title:        " New Order Receive: Order ID " + getTransaction.OrderIdBook,
				Desc:         "You got a booking for " + getTransaction.ExpTitle + " , booked by " + getTransaction.CreatedBy,
			}
			pushNotifErr := p.notificationRepo.Insert(ctx, notif)
			if pushNotifErr != nil {
				return nil
			}
			if confirmIn.TransactionStatus == 2 && confirmIn.BookingStatus == 1 {
				//confirm
				bookingDetail, err := p.bookingUsecase.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "","")
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
					var tmpl = template.Must(template.New("main-template").Parse(templateBookingApprovalDP))
					var data = map[string]interface{}{
						"title":            bookingDetail.Experience[0].ExpTitle,
						"user":             user,
						"payment":          message.NewPrinter(language.German).Sprint(*bookingDetail.TotalPrice),
						"remainingPayment": message.NewPrinter(language.German).Sprint(bookingDetail.ExperiencePaymentType.RemainingPayment),
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
						"city":            bookingDetail.Experience[0].HarborsName,
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
					pdf := htmlPDFTicket.String()
					var attachment []*models.Attachment
					eTicket := models.Attachment{
						AttachmentFileUrl: pdf,
						FileName:          "E-Ticket.pdf",
					}
					attachment = append(attachment, &eTicket)
					pushEmail := &models.SendingEmail{
						Subject:    "Your booking has been confirmed",
						Message:    msg,
						From:       "CGO Indonesia",
						To:         bookingDetail.BookedBy[0].Email,
						Attachment: attachment,
					}
					if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
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
					guestCount := len(bookingDetail.GuestDesc)

					var tmpl *template.Template
					var data map[string]interface{}
					var t *template.Template
					var dataMapping map[string]interface{}
					if bookingDetail.Experience[0].ExpPickupPlace == "" && (bookingDetail.Experience[0].ExpPickupTime == "" || bookingDetail.Experience[0].ExpPickupTime == "00:00:00") {
						tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutMeetingPointAndTime))
						data = map[string]interface{}{
							"title":        bookingDetail.Experience[0].ExpTitle,
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
							"city":            bookingDetail.Experience[0].HarborsName,
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
							"title":        bookingDetail.Experience[0].ExpTitle,
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
							"city":            bookingDetail.Experience[0].HarborsName,
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
							"title":        bookingDetail.Experience[0].ExpTitle,
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
							"city":            bookingDetail.Experience[0].HarborsName,
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
							"title":        bookingDetail.Experience[0].ExpTitle,
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
							"city":            bookingDetail.Experience[0].HarborsName,
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
					var attachment []*models.Attachment
					eTicket := models.Attachment{
						AttachmentFileUrl: pdf,
						FileName:          "E-Ticket.pdf",
					}
					attachment = append(attachment, &eTicket)
					pushEmail := &models.SendingEmail{
						Subject:    "Experience E-Ticket",
						Message:    msg,
						From:       "CGO Indonesia",
						To:         bookingDetail.BookedBy[0].Email,
						Attachment: attachment,
					}

					if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
						return nil
					}
				}

			} else if confirmIn.TransactionStatus == 3 && confirmIn.BookingStatus == 1 {
				//cancelled
				bookingDetail, err := p.bookingUsecase.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "","")
				if err != nil {
					return err
				}
				tripDate := bookingDetail.BookingDate.Format("02 January 2006")
				duration := 0
				if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
					duration = bookingDetail.Experience[0].ExpDuration - 1
					tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
				}
				var tmpl = template.Must(template.New("main-template").Parse(templateBookingRejected))
				var data = map[string]interface{}{
					"title":    bookingDetail.Experience[0].ExpTitle,
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
					Subject:    "Cancelled Booking",
					Message:    msg,
					From:       "CGO Indonesia",
					To:         getTransaction.CreatedBy,
					Attachment: nil,
				}
				if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil
				}
			}

		}
	}else if confirmIn.TransId != ""{
		listTransaction , _ := p.transactionRepo.GetByBookingDate(ctx,confirmIn.BookingDate,confirmIn.TransId,"")
		for _,getTransaction := range listTransaction{
			notif := models.Notification{
				Id:           guuid.New().String(),
				CreatedBy:    getTransaction.CreatedBy,
				CreatedDate:  time.Now(),
				ModifiedBy:   nil,
				ModifiedDate: nil,
				DeletedBy:    nil,
				DeletedDate:  nil,
				IsDeleted:    0,
				IsActive:     0,
				MerchantId:   getTransaction.MerchantId,
				Type:         0,
				Title:        " New Order Receive: Order ID " + getTransaction.OrderIdBook,
				Desc:         "You got a booking for " + getTransaction.ExpTitle + " , booked by " + getTransaction.CreatedBy,
			}
			pushNotifErr := p.notificationRepo.Insert(ctx, notif)
			if pushNotifErr != nil {
				return nil
			}
			bookingDetail, err := p.bookingUsecase.GetDetailTransportBookingID(ctx, *getTransaction.OrderId, *getTransaction.OrderId, nil,"")
			if err != nil {
				return err
			}
			user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
			tripDate := bookingDetail.BookingDate.Format("02 January 2006")
			guestCount := len(bookingDetail.GuestDesc)

			layoutFormat := "15:04:05"
			departureTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].DepartureTime)
			arrivalTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].ArrivalTime)

			if bookingDetail.Transportation[0].ReturnTransId != nil && len(bookingDetail.Transportation) > 1 {

				bookingDetailReturn, err := p.bookingUsecase.GetDetailTransportBookingID(ctx, bookingDetail.OrderId, bookingDetail.OrderId, bookingDetail.Transportation[0].ReturnTransId,"")
				if err != nil {
					return err
				}
				tripDateReturn := bookingDetailReturn.BookingDate.Format("02 January 2006")

				departureTimeReturn, _ := time.Parse(layoutFormat, bookingDetailReturn.Transportation[0].DepartureTime)
				arrivalTimeReturn, _ := time.Parse(layoutFormat, bookingDetailReturn.Transportation[0].ArrivalTime)

				tmpl := template.Must(template.New("main-template").Parse(templateTicketTransportationWithReturn))
				data := map[string]interface{}{
					"title":               bookingDetail.Transportation[0].TransTitle,
					"user":                user,
					"tripDateDeparture":   tripDate,
					"guestCountDeparture": strconv.Itoa(guestCount) + " Guest(s)",
					"sourceTimeDeparture": departureTime.Format("15:04"),
					"desTimeDeparture":    arrivalTime.Format("15:04"),
					"durationDeparture":   bookingDetail.Transportation[0].TripDuration,
					"sourceDeparture":     bookingDetail.Transportation[0].HarborSourceName,
					"destDeparture":       bookingDetail.Transportation[0].HarborDestName,
					"classDeparture":      bookingDetail.Transportation[0].TransClass,
					"orderId":             bookingDetail.OrderId,
					"merchantPicture":     bookingDetail.Transportation[0].MerchantPicture,
					"tripDateReturn":      tripDateReturn,
					"guestCountReturn":    strconv.Itoa(guestCount) + " Guest(s)",
					"sourceTimeReturn":    departureTimeReturn.Format("15:04"),
					"desTimeReturn":       arrivalTimeReturn.Format("15:04"),
					"durationReturn":      bookingDetailReturn.Transportation[0].TripDuration,
					"sourceReturn":        bookingDetailReturn.Transportation[0].HarborSourceName,
					"destReturn":          bookingDetailReturn.Transportation[0].HarborDestName,
					"classReturn":         bookingDetailReturn.Transportation[0].TransClass,
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

				//ticketPDF Bind HTML is Return
				var htmlPDFTicketReturn bytes.Buffer

				dataMappingReturn := map[string]interface{}{
					"guestDesc":       guestDesc,
					"tripDate":        tripDateReturn,
					"sourceTime":      departureTimeReturn.Format("15:04"),
					"desTime":         arrivalTimeReturn.Format("15:04"),
					"duration":        bookingDetailReturn.Transportation[0].TripDuration,
					"source":          bookingDetailReturn.Transportation[0].HarborSourceName,
					"dest":            bookingDetailReturn.Transportation[0].HarborDestName,
					"class":           bookingDetailReturn.Transportation[0].TransClass,
					"qrCode":          bookingDetailReturn.TicketQRCode,
					"merchantPicture": bookingDetailReturn.Transportation[0].MerchantPicture,
					"orderId":         bookingDetailReturn.OrderId,
				}
				// We create the template and register out template function
				tReturn := template.New("t").Funcs(templateFuncs)
				tReturn, err = tReturn.Parse(templateTicketTransportationPDF)
				if err != nil {
					panic(err)
				}

				err = tReturn.Execute(&htmlPDFTicketReturn, dataMappingReturn)
				if err != nil {
					panic(err)
				}

				msg := tpl.String()
				pdf := htmlPDFTicket.String()
				pdfReturn := htmlPDFTicketReturn.String()
				var attachment []*models.Attachment
				eTicket := models.Attachment{
					AttachmentFileUrl: pdf,
					FileName:          "E-Ticket.pdf",
				}
				attachment = append(attachment, &eTicket)
				eTicketReturn := models.Attachment{
					AttachmentFileUrl: pdfReturn,
					FileName:          "E-Ticket-Return.pdf",
				}
				attachment = append(attachment, &eTicketReturn)
				pushEmail := &models.SendingEmail{
					Subject:    "Transportation E-Ticket",
					Message:    msg,
					From:       "CGO Indonesia",
					To:         bookingDetail.BookedBy[0].Email,
					Attachment: attachment,
				}
				if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil
				}

			} else {
				tmpl := template.Must(template.New("main-template").Parse(templateTicketTransportation))
				data := map[string]interface{}{
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
				var attachment []*models.Attachment
				eTicket := models.Attachment{
					AttachmentFileUrl: pdf,
					FileName:          "E-Ticket.pdf",
				}
				attachment = append(attachment, &eTicket)
				pushEmail := &models.SendingEmail{
					Subject:    "Transportation E-Ticket",
					Message:    msg,
					From:       "CGO Indonesia",
					To:        bookingDetail.BookedBy[0].Email,
					Attachment: attachment,
				}
				if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
					return err
				}

			}

		}
	}
	return nil
}
func (p paymentUsecase) ConfirmPayment(ctx context.Context, confirmIn *models.ConfirmPaymentIn) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()
	if confirmIn.BookingStatus == 0 {
		transaction ,_:= p.transactionRepo.GetById(ctx,confirmIn.TransactionID)
		if transaction != nil{
			confirmIn.BookingStatus = *transaction.BookingStatus
		}
	}
	err := p.paymentRepo.ConfirmPayment(ctx, confirmIn)
	if err != nil {
		return err
	}
	getTransaction, err := p.transactionRepo.GetById(ctx, confirmIn.TransactionID)
	if err != nil {
		return err
	}
	notif := models.Notification{
		Id:           guuid.New().String(),
		CreatedBy:    getTransaction.CreatedBy,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		MerchantId:   getTransaction.MerchantId,
		Type:         0,
		Title:        " New Order Receive: Order ID " + getTransaction.OrderIdBook,
		Desc:         "You got a booking for " + getTransaction.ExpTitle + " , booked by " + getTransaction.CreatedBy,
	}
	pushNotifErr := p.notificationRepo.Insert(ctx, notif)
	if pushNotifErr != nil {
		return nil
	}
	if getTransaction.ExpId != nil{
		if confirmIn.TransactionStatus == 2 && confirmIn.BookingStatus == 1 {
			//confirm
			bookingDetail, err := p.bookingUsecase.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "","")
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
				var tmpl = template.Must(template.New("main-template").Parse(templateBookingApprovalDP))
				var data = map[string]interface{}{
					"title":            bookingDetail.Experience[0].ExpTitle,
					"user":             user,
					"payment":          message.NewPrinter(language.German).Sprint(*bookingDetail.TotalPrice),
					"remainingPayment": message.NewPrinter(language.German).Sprint(bookingDetail.ExperiencePaymentType.RemainingPayment),
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
					"city":            bookingDetail.Experience[0].HarborsName,
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
				pdf := htmlPDFTicket.String()
				var attachment []*models.Attachment
				eTicket := models.Attachment{
					AttachmentFileUrl: pdf,
					FileName:          "E-Ticket.pdf",
				}
				attachment = append(attachment, &eTicket)
				pushEmail := &models.SendingEmail{
					Subject:    "Your booking has been confirmed",
					Message:    msg,
					From:       "CGO Indonesia",
					To:         bookingDetail.BookedBy[0].Email,
					Attachment: attachment,
				}
				if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
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
				guestCount := len(bookingDetail.GuestDesc)

				var tmpl *template.Template
				var data map[string]interface{}
				var t *template.Template
				var dataMapping map[string]interface{}
				if bookingDetail.Experience[0].ExpPickupPlace == "" && (bookingDetail.Experience[0].ExpPickupTime == "" || bookingDetail.Experience[0].ExpPickupTime == "00:00:00") {
					tmpl = template.Must(template.New("main-template").Parse(templateTicketFPWithoutMeetingPointAndTime))
					data = map[string]interface{}{
						"title":        bookingDetail.Experience[0].ExpTitle,
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
						"city":            bookingDetail.Experience[0].HarborsName,
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
						"title":        bookingDetail.Experience[0].ExpTitle,
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
						"city":            bookingDetail.Experience[0].HarborsName,
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
						"title":        bookingDetail.Experience[0].ExpTitle,
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
						"city":            bookingDetail.Experience[0].HarborsName,
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
						"title":        bookingDetail.Experience[0].ExpTitle,
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
						"city":            bookingDetail.Experience[0].HarborsName,
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
				var attachment []*models.Attachment
				eTicket := models.Attachment{
					AttachmentFileUrl: pdf,
					FileName:          "E-Ticket.pdf",
				}
				attachment = append(attachment, &eTicket)
				pushEmail := &models.SendingEmail{
					Subject:    "Experience E-Ticket",
					Message:    msg,
					From:       "CGO Indonesia",
					To:         bookingDetail.BookedBy[0].Email,
					Attachment: attachment,
				}

				if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
					return nil
				}
			}

		} else if confirmIn.TransactionStatus == 3 && confirmIn.BookingStatus == 1 {
			//cancelled
			bookingDetail, err := p.bookingUsecase.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "","")
			if err != nil {
				return err
			}
			tripDate := bookingDetail.BookingDate.Format("02 January 2006")
			duration := 0
			if bookingDetail.Experience[0].ExpDuration != 0 && bookingDetail.Experience[0].ExpDuration != 1 {
				duration = bookingDetail.Experience[0].ExpDuration - 1
				tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, duration).Format("02 January 2006")
			}
			var tmpl = template.Must(template.New("main-template").Parse(templateBookingRejected))
			var data = map[string]interface{}{
				"title":    bookingDetail.Experience[0].ExpTitle,
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
				Subject:    "Cancelled Booking",
				Message:    msg,
				From:       "CGO Indonesia",
				To:         getTransaction.CreatedBy,
				Attachment: nil,
			}
			if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil
			}
		}
	}else if getTransaction.TransId != nil{
		bookingDetail, err := p.bookingUsecase.GetDetailTransportBookingID(ctx, *getTransaction.OrderId, *getTransaction.OrderId, nil,"")
		if err != nil {
			return err
		}
		user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
		tripDate := bookingDetail.BookingDate.Format("02 January 2006")
		guestCount := len(bookingDetail.GuestDesc)

		layoutFormat := "15:04:05"
		departureTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].DepartureTime)
		arrivalTime, _ := time.Parse(layoutFormat, bookingDetail.Transportation[0].ArrivalTime)

		if bookingDetail.Transportation[0].ReturnTransId != nil && len(bookingDetail.Transportation) > 1 {

			bookingDetailReturn, err := p.bookingUsecase.GetDetailTransportBookingID(ctx, bookingDetail.OrderId, bookingDetail.OrderId, bookingDetail.Transportation[0].ReturnTransId,"")
			if err != nil {
				return err
			}
			tripDateReturn := bookingDetailReturn.BookingDate.Format("02 January 2006")

			departureTimeReturn, _ := time.Parse(layoutFormat, bookingDetailReturn.Transportation[0].DepartureTime)
			arrivalTimeReturn, _ := time.Parse(layoutFormat, bookingDetailReturn.Transportation[0].ArrivalTime)

			tmpl := template.Must(template.New("main-template").Parse(templateTicketTransportationWithReturn))
			data := map[string]interface{}{
				"title":               bookingDetail.Transportation[0].TransTitle,
				"user":                user,
				"tripDateDeparture":   tripDate,
				"guestCountDeparture": strconv.Itoa(guestCount) + " Guest(s)",
				"sourceTimeDeparture": departureTime.Format("15:04"),
				"desTimeDeparture":    arrivalTime.Format("15:04"),
				"durationDeparture":   bookingDetail.Transportation[0].TripDuration,
				"sourceDeparture":     bookingDetail.Transportation[0].HarborSourceName,
				"destDeparture":       bookingDetail.Transportation[0].HarborDestName,
				"classDeparture":      bookingDetail.Transportation[0].TransClass,
				"orderId":             bookingDetail.OrderId,
				"merchantPicture":     bookingDetail.Transportation[0].MerchantPicture,
				"tripDateReturn":      tripDateReturn,
				"guestCountReturn":    strconv.Itoa(guestCount) + " Guest(s)",
				"sourceTimeReturn":    departureTimeReturn.Format("15:04"),
				"desTimeReturn":       arrivalTimeReturn.Format("15:04"),
				"durationReturn":      bookingDetailReturn.Transportation[0].TripDuration,
				"sourceReturn":        bookingDetailReturn.Transportation[0].HarborSourceName,
				"destReturn":          bookingDetailReturn.Transportation[0].HarborDestName,
				"classReturn":         bookingDetailReturn.Transportation[0].TransClass,
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

			//ticketPDF Bind HTML is Return
			var htmlPDFTicketReturn bytes.Buffer

			dataMappingReturn := map[string]interface{}{
				"guestDesc":       guestDesc,
				"tripDate":        tripDateReturn,
				"sourceTime":      departureTimeReturn.Format("15:04"),
				"desTime":         arrivalTimeReturn.Format("15:04"),
				"duration":        bookingDetailReturn.Transportation[0].TripDuration,
				"source":          bookingDetailReturn.Transportation[0].HarborSourceName,
				"dest":            bookingDetailReturn.Transportation[0].HarborDestName,
				"class":           bookingDetailReturn.Transportation[0].TransClass,
				"qrCode":          bookingDetailReturn.TicketQRCode,
				"merchantPicture": bookingDetailReturn.Transportation[0].MerchantPicture,
				"orderId":         bookingDetailReturn.OrderId,
			}
			// We create the template and register out template function
			tReturn := template.New("t").Funcs(templateFuncs)
			tReturn, err = tReturn.Parse(templateTicketTransportationPDF)
			if err != nil {
				panic(err)
			}

			err = tReturn.Execute(&htmlPDFTicketReturn, dataMappingReturn)
			if err != nil {
				panic(err)
			}

			msg := tpl.String()
			pdf := htmlPDFTicket.String()
			pdfReturn := htmlPDFTicketReturn.String()
			var attachment []*models.Attachment
			eTicket := models.Attachment{
				AttachmentFileUrl: pdf,
				FileName:          "E-Ticket.pdf",
			}
			attachment = append(attachment, &eTicket)
			eTicketReturn := models.Attachment{
				AttachmentFileUrl: pdfReturn,
				FileName:          "E-Ticket-Return.pdf",
			}
			attachment = append(attachment, &eTicketReturn)
			pushEmail := &models.SendingEmail{
				Subject:    "Transportation E-Ticket",
				Message:    msg,
				From:       "CGO Indonesia",
				To:         bookingDetail.BookedBy[0].Email,
				Attachment: attachment,
			}
			if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil
			}

		} else {
			tmpl := template.Must(template.New("main-template").Parse(templateTicketTransportation))
			data := map[string]interface{}{
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
			var attachment []*models.Attachment
			eTicket := models.Attachment{
				AttachmentFileUrl: pdf,
				FileName:          "E-Ticket.pdf",
			}
			attachment = append(attachment, &eTicket)
			pushEmail := &models.SendingEmail{
				Subject:    "Transportation E-Ticket",
				Message:    msg,
				From:       "CGO Indonesia",
				To:        bookingDetail.BookedBy[0].Email,
				Attachment: attachment,
			}
			if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
				return err
			}

		}
	}

	return nil
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
