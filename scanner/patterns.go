package scanner

import "regexp"

var patterns = map[string]string{
	"AWS Access Key":                  `(?i)AWS_ACCESS_KEY_ID\s*=\s*AKIA[0-9A-Z]{16}`,
	"AWS Secret Key":                  `(?i)AWS_SECRET_ACCESS_KEY\s*=\s*([A-Za-z0-9/+=]{40})`,
	"GitHub Token":                    `(?i)ghp_[0-9a-zA-Z]{36}`,
	"Bearer Token":                    `(?i)bearer[\s]*=[\s]*["']?([A-Za-z0-9-_.]{40})["']?`,
	"Password Config":                 `(?i)(password|pass|pwd)[\s]*[:=][\s]*["']?([^"'\s]+)["']?`,
	"MongoDB URI":                     `(?i)mongodb://[a-zA-Z0-9._%+-]+:[^@]+@[^/]+/[a-zA-Z0-9._%+-]+`,
	"SQL Connection":                  `(?i)(Data Source=|Server=)[^;]+;User ID=[^;]+;Password=[^;]+`,
	"SSH Private Key":                 `(?i)-----BEGIN (RSA|DSA|ECDSA|ED25519) PRIVATE KEY-----[\s\S]+-----END (RSA|DSA|ECDSA|ED25519) PRIVATE KEY-----`,
	"Stripe Secret Key":               `(?i)sk_test_[0-9a-zA-Z]{24}`,
	"API Key":                         `(?i)(api[-_]?key|apikey|key|secret)\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Firebase Secret":                 `(?i)FIREBASE_SECRET\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Mailgun API Key":                 `(?i)key-[0-9a-zA-Z]{32}`,
	"Twilio Auth Token":               `(?i)TWILIO_AUTH_TOKEN\s*=\s*["']?([0-9a-zA-Z]{32})["']?`,
	"Slack Token":                     `(?i)xox[baprs]-[0-9]+-[0-9a-zA-Z]+`,
	"GitLab Personal Access Token":    `(?i)glpat-[0-9a-zA-Z]{20,40}`,
	"Heroku API Key":                  `(?i)HEROKU_API_KEY\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Discord Token":                   `(?i)ND[0-9]{18}\.[0-9]{12}\.[0-9a-zA-Z_-]{27}`,
	"JWT Secret":                      `(?i)(jwt|secret)[\s]*=[\s]*["']?([0-9a-zA-Z_-]+)["']?`,
	"OAuth2 Client Secret":            `(?i)client_secret\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Google API Key":                  `(?i)AIza[0-9A-Za-z-_]{35}`,
	"PayPal Client Secret":            `(?i)["']?A21[0-9a-zA-Z]{80}["']?`,
	"Microsoft Azure Key":             `(?i)azure[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Telegram Bot Token":              `(?i)bot[0-9]+:[0-9a-zA-Z_-]{35}`,
	"Redis Password":                  `(?i)REDIS_PASSWORD\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Stripe Webhook Secret":           `(?i)whsec_[0-9a-zA-Z]{32}`,
	"Salesforce Token":                `(?i)Salesforce[ -]?token\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Oracle DB Password":              `(?i)password\s*=\s*["']?([0-9a-zA-Z!@#$%^&*()-_+=<>?]{8,})["']?`,
	"YAML API Key":                    `(?i)api_key:\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"OpenAI API Key":                  `(?i)sk-[0-9a-zA-Z]{48}`,
	"Dropbox Access Token":            `(?i)DROPBOX_ACCESS_TOKEN\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"MongoDB Atlas Connection String": `(?i)mongodb(?:\+srv)?://[a-zA-Z0-9]+:[^@]+@([^/]+)/[a-zA-Z0-9]+`,
	"Zoom JWT Token":                  `(?i)eyJ(.*)\.[A-Za-z0-9-]+`,
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
	"Facebook Access Token":           `(?i)(fb|facebook)[-_ ]?access[_ ]?token\s*=\s*([A-Za-z0-9]{20,})`,
	"Slack OAuth Token":               `(?i)xoxp-[0-9]+-[0-9]+-[0-9a-zA-Z]+`,
	"HubSpot API Key":                 `(?i)hubspot[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Mailchimp API Key":               `(?i)mailchimp[ -]?api[ -]?key\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"LinkedIn Access Token":           `(?i)linkedin[ -]?access[ -]?token\s*=\s*["']?([0-9a-zA-Z_-]+)["']?`,
	"Hardcoded Admin Password":        `(?i)adminPassword\s*[:=][\s]*([^"\s]+)`,
	"Hardcoded Password":              `(?i)(password\s*[:=]\s*["']?([^"']+)["']?|password\s*=\s*["']?([^"']+)["']?)`,
	"JSON Hardcoded Password":         `(?i)"password"\s*:\s*["']?([^"']+)["']?`,
	"JSON Hardcoded Username":         `(?i)"user"\s*:\s*["']?([^"']+)["']?`,
	"JSON Hardcoded API Key":          `(?i)"apiKey"\s*:\s*["']?([^"']+)["']?`,
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
