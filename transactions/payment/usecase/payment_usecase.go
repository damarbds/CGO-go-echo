package usecase

import (
	"bytes"
	"context"
	"html/template"
	"strconv"
	"time"

	"github.com/auth/identityserver"

	"github.com/misc/notif"
	"github.com/transactions/transaction"

	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/transactions/payment"
)

type paymentUsecase struct {
	isUsecase        identityserver.Usecase
	transactionRepo  transaction.Repository
	notificationRepo notif.Repository
	paymentRepo      payment.Repository
	userUsercase     user.Usecase
	bookingRepo      booking_exp.Repository
	userRepo         user.Repository
	contextTimeout   time.Duration
	bookingUsecase   booking_exp.Usecase
}

// NewPaymentUsecase will create new an paymentUsecase object representation of payment.Usecase interface
func NewPaymentUsecase(bookingUsecase booking_exp.Usecase, isUsecase identityserver.Usecase, t transaction.Repository, n notif.Repository, p payment.Repository, u user.Usecase, b booking_exp.Repository, ur user.Repository, timeout time.Duration) payment.Usecase {
	return &paymentUsecase{
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
)

func (p paymentUsecase) Insert(ctx context.Context, payment *models.Transaction, token string, points float64) (string, error) {
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

func (p paymentUsecase) ConfirmPayment(ctx context.Context, confirmIn *models.ConfirmPaymentIn) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	err := p.paymentRepo.ConfirmPayment(ctx, confirmIn)
	if err != nil {
		return err
	}
	getTransaction, err := p.transactionRepo.GetById(ctx, confirmIn.TransactionID)
	if err != nil {
		return err
	}
	notif := models.Notification{
		Id:           "",
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
		bookingDetail, err := p.bookingUsecase.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "")
		if bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
			user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
			tripDate := bookingDetail.BookingDate.Format("02 January 2006")
			tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, bookingDetail.Experience[0].ExpDuration).Format("02 January 2006")
			var tmpl = template.Must(template.New("main-template").Parse(templateBookingApprovalDP))
			var data = map[string]interface{}{
				"title":            bookingDetail.Experience[0].ExpTitle,
				"user":             user,
				"payment":          bookingDetail.TotalPrice,
				"remainingPayment": bookingDetail.ExperiencePaymentType.RemainingPayment,
				"paymentDeadline":  bookingDetail.BookingDate.Format("02 January 2006"),
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

			//maxTime := time.Now().AddDate(0, 0, 1)
			msg := tpl.String()
			pushEmail := &models.SendingEmail{
				Subject:  "Booking Approved DP By Merchant",
				Message:  msg,
				From:     "CGO Indonesia",
				To:       bookingDetail.BookedBy[0].Email,
				FileName: "",
			}
			if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil
			}

		} else {
			user := bookingDetail.BookedBy[0].Title + `.` + bookingDetail.BookedBy[0].FullName
			tripDate := bookingDetail.BookingDate.Format("02 January 2006")
			tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, bookingDetail.Experience[0].ExpDuration).Format("02 January 2006")
			guestCount := len(bookingDetail.GuestDesc)

			var tmpl = template.Must(template.New("main-template").Parse(templateTicketFP))
			var data = map[string]interface{}{
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
			var tpl bytes.Buffer
			err = tmpl.Execute(&tpl, data)
			if err != nil {
				//http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			//maxTime := time.Now().AddDate(0, 0, 1)
			msg := tpl.String()
			pushEmail := &models.SendingEmail{
				Subject:  "Ticket FP",
				Message:  msg,
				From:     "CGO Indonesia",
				To:       bookingDetail.BookedBy[0].Email,
				FileName: "",
			}

			if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
				return nil
			}
		}

	} else if confirmIn.TransactionStatus == 3 && confirmIn.BookingStatus == 1 {
		//cancelled
		bookingDetail, err := p.bookingUsecase.GetDetailBookingID(ctx, *getTransaction.BookingExpId, "")
		if err != nil {
			return err
		}
		tripDate := bookingDetail.BookingDate.Format("02 January 2006")
		tripDate = tripDate + ` - ` + bookingDetail.BookingDate.AddDate(0, 0, bookingDetail.Experience[0].ExpDuration).Format("02 January 2006")
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
			Subject:  "Booking Rejected",
			Message:  msg,
			From:     "CGO Indonesia",
			To:       getTransaction.CreatedBy,
			FileName: "",
		}
		if _, err := p.isUsecase.SendingEmail(pushEmail); err != nil {
			return nil
		}
	}

	return nil
}
