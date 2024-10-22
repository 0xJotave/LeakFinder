package scanner

import "regexp"

var patterns = map[string]string{
	"AWS Access Key":                  `AKIA[0-9A-Z]{16}`,
	"AWS Secret Key":                  `(?<=AWS_SECRET_ACCESS_KEY\s*=\s*)[A-Za-z0-9/+=]{40}`,
	"GitHub Token":                    `ghp_[0-9a-zA-Z]{36}`,
	"GitHub Personal Access Token":    `ghp_[0-9a-zA-Z]{36}`,
	"Bearer Token":                    `[A-Za-z0-9-_.]{40}`,
	"Password Config":                 `password\s*=\s*["']?([^\s"']+)["']?`,
	"MongoDB URI":                     `mongodb://[a-zA-Z0-9._%+-]+:[^@]+@[^/]+/[a-zA-Z0-9._%+-]+`,
	"SQL Connection":                  `(Data Source=|Server=)[^;]+;User ID=[^;]+;Password=[^;]+`,
	"SSH Private Key":                 `-----BEGIN (RSA|DSA|ECDSA|ED25519) PRIVATE KEY-----[\s\S]+-----END \1 PRIVATE KEY-----`,
	"Stripe Secret Key":               `sk_test_[0-9a-zA-Z]{24}`,
	"API Key":                         `(?i)(api[-_]?key|apikey|key|secret)[\s]*=[\s]*["']?([0-9a-zA-Z_-]+)["']?`,
	"Firebase Secret":                 `FIREBASE_SECRET\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Mailgun API Key":                 `key-[0-9a-zA-Z]{32}`,
	"Twilio Auth Token":               `(?i)TWILIO_AUTH_TOKEN\s*=\s*["']?([0-9a-zA-Z]{32})["']?`,
	"Slack Token":                     `xox[baprs]-[0-9]+-[0-9a-zA-Z]+`,
	"Token":                           `(?i)(token|auth|bearer)[\s]*=[\s]*["']?([0-9a-zA-Z_-]+)["']?`,
	"GitLab Personal Access Token":    `glpat-[0-9a-zA-Z]{20,40}`,
	"Heroku API Key":                  `HEROKU_API_KEY\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Discord Token":                   `(?i)ND[0-9]{18}\.[0-9]{12}\.[0-9a-zA-Z_-]{27}`,
	"JWT Secret":                      `(?i)(jwt|secret)[\s]*=[\s]*["']?([0-9a-zA-Z_-]+)["']?`,
	"OAuth2 Client Secret":            `(?i)client_secret[\s]*=[\s]*["']?([0-9a-zA-Z_-]+)["']?`,
	"Google API Key":                  `AIza[0-9A-Za-z-_]{35}`,
	"PayPal Client Secret":            `["']?A21[0-9A-Za-z]{80}["']?`,
	"Microsoft Azure Key":             `(?i)azure[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Telegram Bot Token":              `(?i)bot[0-9]+:[0-9a-zA-Z_-]{35}`,
	"Redis Password":                  `(?i)REDIS_PASSWORD\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Stripe Webhook Secret":           `whsec_[0-9a-zA-Z]{32}`,
	"Salesforce Token":                `(?i)Salesforce[ -]?token\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Oracle DB Password":              `(?i)password\s*=\s*["']?([0-9a-zA-Z!@#$%^&*()-_+=<>?]{8,})["']?`,
	"YAML API Key":                    `(?i)api_key:\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"OpenAI API Key":                  `sk-[0-9a-zA-Z]{48}`,
	"Dropbox Access Token":            `(?i)DROPBOX_ACCESS_TOKEN\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"MongoDB Atlas Connection String": `mongodb(?:\+srv)?://[a-zA-Z0-9]+:[^@]+@([^/]+)/[a-zA-Z0-9]+`,
	"Zoom JWT Token":                  `eyJ(.*)\.[A-Za-z0-9-]+`,
	"PayPal Webhook ID":               `(?i)PAYPAL_WEBHOOK_ID\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Okta API Token":                  `(?i)okta[ -]?api[ -]?token\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Zoom API Key":                    `(?i)zoom[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Algolia API Key":                 `(?i)algolia[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Telegram API Key":                `(?i)telegram[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Secret Key":                      `(?i)secret_key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Line Bot Token":                  `(?i)line[ -]?bot[ -]?token\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Segment Write Key":               `(?i)write_key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Jira API Token":                  `(?i)JIRA_API_TOKEN\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Trello API Key":                  `(?i)trello[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Facebook Access Token":           `EA[0-9A-Za-z]+`,
	"Slack OAuth Token":               `xoxp-[0-9]+-[0-9]+-[0-9a-zA-Z]+`,
	"HubSpot API Key":                 `(?i)hubspot[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Mailchimp API Key":               `(?i)mailchimp[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"LinkedIn Access Token":           `(?i)linkedin[ -]?access[ -]?token\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Hardcoded API Key":               `(?i)(api[-_]?key|apikey)[\s]*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Hardcoded Username":              `(?i)(username|user)\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded Access Token":          `(?i)access[-_]?token\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded Client Secret":         `(?i)client[-_]?secret\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded Refresh Token":         `(?i)refresh[-_]?token\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded Database Password":     `(?i)(db|database)[\s]*password\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded FTP Password":          `(?i)ftp[-_]?password\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded MySQL Password":        `(?i)(mysql|sql)[\s]*password\s*=\s*["']?([^"']+)["']?`,
	"Hardcoded API Secret":            `(?i)api[-_]?secret\s*=\s*["']?([^"']+)["']?`,
}

func GetPatterns() (map[string]*regexp.Regexp, error) {
	compiledPatterns := make(map[string]*regexp.Regexp)

	for name, pattern := range patterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		compiledPatterns[name] = compiled
	}

	return compiledPatterns, nil
}
