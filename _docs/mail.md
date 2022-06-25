# Send email via Gmail SMTP server with Golang

# Demo:
You can try my forgot password workflow: [forgot password](https://api.tdo.works/swagger/index.html#/auth/post_api_v1_auth_forgot_password)

# Installation

1. Enable IMAP setting
https://support.google.com/mail/answer/7126229?hl=en

2. Enable 2FA on your Gmail
https://myaccount.google.com/signinoptions/two-step-verification/enroll-welcome

3. Create an app-specific password
https://myaccount.google.com/apppasswords
After this, put the password to the settings:[config_local.yml](../components/appctx/config_local.yml).

```bash
# mail
mail:
  sender: "yourgmail@gmail.com"
  host: "smtp.gmail.com"
  port: 587
  username: "yourgmail@gmail.com"
  password: "yourpassword" <-- put the password here
```

4. Checkout [Mailer component](../components/mailer/mailer.go)

5. Run test:
```bash
mail="yourgmail@gmail.com" go run ./packages/mail/main.go
```
