package templates

import (
	"fmt"
	"strings"
	"time"
)

const otpTemplate = `
<!DOCTYPE html>
<html>
<body style="margin:0;padding:24px;background-color:#fcfffd;font-family:Roboto,Arial,sans-serif;color:#202224;">
  
  <table width="100%" cellpadding="0" cellspacing="0" role="presentation">
    <tr>
      <td align="center">

        <table width="600" cellpadding="0" cellspacing="0" role="presentation" style="background:#fcfffd;">

          <!-- Title -->
          <tr>
            <td style="font-size:36px;font-weight:300;line-height:40px;color:#000;">
              Welcome to AIOZ Clawbot!
            </td>
          </tr>

          <!-- Spacer -->
          <tr><td height="32"></td></tr>

          <!-- Description -->
          <tr>
            <td style="font-size:16px;line-height:24px;color:#000;">
              Thank you for signing up. Your confirmation code is below. Enter it in your open browser window and we'll help you get signed up.
            </td>

          <!-- Spacer -->
          <tr><td height="32"></td></tr>

          <!-- OTP -->
          <tr>
            <td>
              <div style="
                display:inline-block;
                padding:18px 20px;
                font-size:32px;
                font-weight:500;
                border:1px solid #dde6de;
                border-radius:8px;
                letter-spacing:4px;
              ">
                {{OTP_CODE}}
              </div>
            </td>
          </tr>

          <!-- Spacer -->
          <tr><td height="32"></td></tr>

          <!-- Note -->
          <tr>
            <td style="font-size:16px;line-height:24px;">
              <strong>
                This code can only be used once and will expire in {{EXPIRE_MINUTES}} minutes.
              </strong>
              <br>
              <span>
                If you didn't request a code, please ignore this email.
              </span>
              <br>
              <strong >
                Never share this code with anyone else.
              </strong>
              <br><br>
              Cheers!<br>
              <strong>AIOZ Clawbot Service</strong>
            </td>
          </tr>

          <!-- Divider -->
          <tr><td height="32"></td></tr>
          <tr>
            <td>
              <hr style="border:none;border-top:1px solid #dde6de;">
            </td>
          </tr>

          <!-- Footer -->
          <tr><td height="32"></td></tr>

          <tr>
            <td style="font-size:14px;">
              Please send any feedback or bug info to 
              <a href="mailto:support@aiozstream.network" style="color:#21975d;">
                support@aiozstream.network
              </a>
              <br><br>
              <a href="https://aiozstream.network/terms-of-service" style="color:#21975d;">Terms of Service</a> |
              <a href="https://aiozstream.network/privacy-policy" style="color:#21975d;">Privacy Policy</a>
            </td>
          </tr>

          <!-- Logo -->
          <tr><td height="24"></td></tr>
          <tr>
            <td>
              <img 
                src="https://content.aioz.network/logo/png/light/logo-aioz_stream_md.png"
                width="173"
                style="display:block;border:0;"
              />
            </td>
          </tr>

          <!-- Copyright -->
          <tr><td height="16"></td></tr>
          <tr>
            <td style="font-size:10px;">
              © 2026 AIOZ Network. All rights reserved.
            </td>
          </tr>

        </table>

      </td>
    </tr>
  </table>

</body>
</html>
`

func BuildOTPEmail(to string, otp string, expire time.Duration) string {
	minutes := int(expire.Minutes())

	html := otpTemplate
	html = strings.Replace(html, "{{OTP_CODE}}", otp, 1)
	html = strings.Replace(html, "{{EXPIRE_MINUTES}}", fmt.Sprintf("%d", minutes), 1)

	return html
}
