// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	ValidationEmailNotifyTemplateEn = `
 <!DOCTYPE html>
<html>
  <head>
    <meta name="viewport" content="width=device-width" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Simple Transactional Email</title>
    <style>
      /* -------------------------------------
			  GLOBAL RESETS
		  ------------------------------------- */
      /*All the styling goes here*/
      img {
        border: none;
        -ms-interpolation-mode: bicubic;
        max-width: 100%;
      }

      body {
        background-color: #ecf1f7;
        font-family: Roboto, PingFang SC, Lantinghei SC, Helvetica Neue,
          Helvetica, Arial, Microsoft YaHei, 微软雅黑, STHeitiSC-Light, simsun,
          宋体, WenQuanYi Zen Hei, WenQuanYi Micro Hei, sans-serif;
        -webkit-font-smoothing: antialiased;
        font-size: 14px;
        line-height: 1.4;
        margin: 0;
        padding: 0;
        -ms-text-size-adjust: 100%;
        -webkit-text-size-adjust: 100%;
        color: #242e42;
      }

      table {
        border-collapse: separate;
        mso-table-lspace: 0pt;
        mso-table-rspace: 0pt;
        width: 100%;
      }
      table td {
        font-size: 14px;
        vertical-align: top;
      }

      /* -------------------------------------
			  BODY & CONTAINER
		  ------------------------------------- */

      .body {
        background-color: #eff0f5;
        width: 100%;
      }

      /* Set a max-width, and make it display as block so it will automatically stretch to that width, but will also shrink down on a phone or something */
      .container {
        display: block;
        margin: 0 auto !important;
        /* makes it centered */
        max-width: 780px;
        padding: 10px;
        padding-top: 80px;
        width: 780px;
      }

      /* This should also be a block element, so that it will fill 100% of the .container */
      .content {
        box-sizing: border-box;
        display: block;
        margin: 0 auto;
        max-width: 780px;
        padding: 10px;
      }

      /* -------------------------------------
			  HEADER, FOOTER, MAIN
		  ------------------------------------- */
      .main {
        background: #ffffff;
        border-radius: 4px;
        border: solid 1px #e3e9ef;
        background-color: #ffffff;
      }

      .wrapper {
        box-sizing: border-box;
        padding: 48px;
      }

      .content-block {
        padding-bottom: 10px;
        padding-top: 10px;
      }

      .footer {
        clear: both;
        margin-top: 14px;
        text-align: center;
        width: 100%;
      }
      .footer td,
      .footer p,
      .footer span,
      .footer a {
        color: #8c96ad;
        font-size: 12px;
        text-align: center;
      }
      .gray {
        color: #8c96ad;
      }
      .logo {
        display: inline-block;
        vertical-align: middle;
      }

      /* -------------------------------------
			  TYPOGRAPHY
		  ------------------------------------- */
      h1,
      h2,
      h3,
      h4 {
        color: #000000;
        font-weight: 400;
        line-height: 1.4;
        margin: 0;
        margin-bottom: 30px;
      }

      h1 {
        font-size: 35px;
        font-weight: 300;
        text-align: center;
        text-transform: capitalize;
      }

      p,
      ul,
      ol {
        font-size: 14px;
        font-weight: normal;
        line-height: 2;
        margin: 0;
        margin-bottom: 15px;
      }
      p li,
      ul li,
      ol li {
        list-style-position: inside;
        margin-left: 5px;
      }

      a {
        color: #329dce;
        font-weight: bold;
        text-decoration: none;
      }

      /* -------------------------------------
			  BUTTONS
		  ------------------------------------- */
      .btn {
        box-sizing: border-box;
        width: 100%;
      }
      .btn > tbody > tr > td {
        padding-bottom: 15px;
      }
      .btn table {
        width: auto;
      }
      .btn table td {
        background-color: #ffffff;
        border-radius: 5px;
        text-align: center;
      }
      .btn a {
        background-color: #ffffff;
        border: solid 1px #3498db;
        border-radius: 5px;
        box-sizing: border-box;
        color: #3498db;
        cursor: pointer;
        display: inline-block;
        font-size: 14px;
        font-weight: bold;
        margin: 0;
        padding: 12px 25px;
        text-decoration: none;
        text-transform: capitalize;
      }

      .btn-primary table td {
        background-color: #3498db;
      }

      .btn-primary a {
        background-color: #3498db;
        border-color: #3498db;
        color: #ffffff;
      }

      /* -------------------------------------
			  OTHER STYLES THAT MIGHT BE USEFUL
		  ------------------------------------- */
      .last {
        margin-bottom: 0;
      }

      .first {
        margin-top: 0;
      }

      .align-center {
        text-align: center;
      }

      .align-right {
        text-align: right;
      }

      .align-left {
        text-align: left;
      }

      .clear {
        clear: both;
      }

      .mt0 {
        margin-top: 0;
      }

      .mb0 {
        margin-bottom: 0;
      }

      .preheader {
        color: transparent;
        display: none;
        height: 0;
        max-height: 0;
        max-width: 0;
        opacity: 0;
        overflow: hidden;
        mso-hide: all;
        visibility: hidden;
        width: 0;
      }

      .powered-by a {
        text-decoration: none;
      }

      hr {
        border: 0;
        border-bottom: 1px solid #eff0f5;
        margin: 50px 0 12px;
      }
      .linkBtn {
        border-radius: 2px;
        box-shadow: 0 1px 3px 0 rgba(73, 33, 173, 0.16),
          0 1px 2px 0 rgba(52, 57, 69, 0.03);
        background-color: #242e42;
        color: #fff;
        line-height: 20px;
        padding: 8px 20px;
      }
      .link {
        font-size: 12px;
        font-style: normal;
        font-stretch: normal;
        line-height: 28px;
        letter-spacing: normal;
      }
      .platform {
        font-size: 14px;
        font-weight: 500;
        font-style: normal;
        font-stretch: normal;
        line-height: 20px;
        letter-spacing: normal;
        color: #343945;
        margin-left: 12px;
      }
      .line1 {
        margin-top: 25px;
        margin-bottom: 16px;
      }
      .line2 {
        font-size: 16px;
        font-weight: bold;
        line-height: 2;
        margin-top: 16px;
        margin-bottom: 20px;
      }
      .line3 {
        margin-bottom: 40px;
        margin-top: 16px;
      }
      .line4 {
        margin-bottom: 0px;
      }
      .line5 {
        margin-top: 0px;
      }
      .line6 {
        margin-bottom: 0px;
      }

      /* -------------------------------------
			  RESPONSIVE AND MOBILE FRIENDLY STYLES
		  ------------------------------------- */
      @media only screen and (max-width: 620px) {
        table[class='body'] h1 {
          font-size: 28px !important;
          margin-bottom: 10px !important;
        }
        table[class='body'] p,
        table[class='body'] ul,
        table[class='body'] ol,
        table[class='body'] td,
        table[class='body'] span,
        table[class='body'] a {
          font-size: 16px !important;
        }
        table[class='body'] .wrapper,
        table[class='body'] .article {
          padding: 10px !important;
        }
        table[class='body'] .content {
          padding: 0 !important;
        }
        table[class='body'] .container {
          padding: 0 !important;
          width: 100% !important;
        }
        table[class='body'] .main {
          border-left-width: 0 !important;
          border-radius: 0 !important;
          border-right-width: 0 !important;
        }
        table[class='body'] .btn table {
          width: 100% !important;
        }
        table[class='body'] .btn a {
          width: 100% !important;
        }
        table[class='body'] .img-responsive {
          height: auto !important;
          max-width: 100% !important;
          width: auto !important;
        }
      }

      /* -------------------------------------
			  PRESERVE THESE STYLES IN THE HEAD
		  ------------------------------------- */
      @media all {
        .ExternalClass {
          width: 100%;
        }
        .ExternalClass,
        .ExternalClass p,
        .ExternalClass span,
        .ExternalClass font,
        .ExternalClass td,
        .ExternalClass div {
          line-height: 100%;
        }
        .apple-link a {
          color: inherit !important;
          font-family: inherit !important;
          font-size: inherit !important;
          font-weight: inherit !important;
          line-height: inherit !important;
          text-decoration: none !important;
        }
        .btn-primary table td:hover {
          background-color: #34495e !important;
        }
        .btn-primary a:hover {
          background-color: #34495e !important;
          border-color: #34495e !important;
        }
      }
    </style>
  </head>
  <body class="">
    <span class="preheader"
      >This is preheader text. Some clients will show this text as a
      preview.</span
    >
    <table
      role="presentation"
      border="0"
      cellpadding="0"
      cellspacing="0"
      class="body"
    >
      <tr>
        <td>&nbsp;</td>
        <td class="container">
          <div class="content">
            <!-- START CENTERED WHITE CONTAINER -->
            <table role="presentation" class="main">
              <!-- START MAIN CONTENT AREA -->
              <tr>
                <td class="wrapper">
                  <table
                    role="presentation"
                    border="0"
                    cellpadding="0"
                    cellspacing="0"
                  >
                    <tr>
                      <td>
                        <p>
                          <img
                            width="150"
                            height="30"
                            src="{{.Icon}}"  
                          />
                        </p>
                        <p class="line1">Hi</p>
                        <p class="line2">
							The email notification service is successfully enabled.                         
                        </p>
                        <hr />
                        <p class="gray">
                          * Don't reply it, this is a system email.
                        </p>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>

            <div class="footer">
              <table
                role="presentation"
                border="0"
                cellpadding="0"
                cellspacing="0"
              >
                <tr>
                  <td class="content-block">                    
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </td>
        <td>&nbsp;</td>
      </tr>
    </table>
  </body>
</html>



`
)
