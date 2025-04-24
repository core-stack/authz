package email

import "html/template"

type TemplateType string

const (
	TemplateActiveAccount       TemplateType = "active_account"
	TemplateResetPassword       TemplateType = "reset_password"
	TemplateNotifyResetPassword TemplateType = "notify_reset_password"
	TemplateChangePassword      TemplateType = "change_password"
	TemplateDeleteAccount       TemplateType = "delete_account"
)

type TemplateItem struct {
	Template *template.Template
	Subject  string
}
type TemplateRegistry map[TemplateType]TemplateItem

func DefaultTemplates() TemplateRegistry {
	return TemplateRegistry{
		TemplateActiveAccount: TemplateItem{
			Template: template.Must(template.ParseFiles("templates/email/active_account.html")),
			Subject:  "Activate your account",
		},
		TemplateResetPassword: TemplateItem{
			Template: template.Must(template.ParseFiles("templates/email/reset_password.html")),
			Subject:  "Reset your password",
		},
		TemplateNotifyResetPassword: TemplateItem{
			Template: template.Must(template.ParseFiles("templates/email/notify_reset_password.html")),
			Subject:  "Reset your password",
		},
		TemplateChangePassword: TemplateItem{
			Template: template.Must(template.ParseFiles("templates/email/change_password.html")),
			Subject:  "Change your password",
		},
		TemplateDeleteAccount: TemplateItem{
			Template: template.Must(template.ParseFiles("templates/email/delete_account.html")),
			Subject:  "Delete your account",
		},
	}
}
