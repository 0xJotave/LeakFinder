package scanner

import "regexp"

var patterns = map[string]string{
	"AWS Access Key":                  `(?i)(AWS_ACCESS_KEY_ID|aws_access_key_id|aws:AccessKeyId)[\s]*[:=][\s]*['"]?(AKIA[0-9A-Z]{16})['"]?`,
	"AWS Secret Key":                  `(?i)(AWS_SECRET_ACCESS_KEY|aws_secret_access_key|aws:SecretAccessKey)[\s]*[:=][\s]*['"]?([A-Za-z0-9/+=]{40})['"]?`,
	"GitHub Token":                    `(?i)ghp_[0-9a-zA-Z]{36}`,
	"Bearer Token":                    `(?i)bearer[\s]*=[\s]*['"]?([A-Za-z0-9-_.]{40})['"]?`,
	"Password Config":                 `(?i)(password|pass|pwd)[\s]*[:=][\s]*['"]?([^'"\s]+)['"]?`,
	"MongoDB URI":                     `(?i)mongodb://[a-zA-Z0-9._%+-]+:[^@]+@[^/]+/[a-zA-Z0-9._%+-]+`,
	"SQL Connection":                  `(?i)(Data Source=|Server=)[^;]+;User ID=[^;]+;Password=[^;]+`,
	"SSH Private Key":                 `(?i)-----BEGIN (RSA|DSA|ECDSA|ED25519) PRIVATE KEY-----[\s\S]+-----END (RSA|DSA|ECDSA|ED25519) PRIVATE KEY-----`,
	"Stripe Secret Key":               `(?i)sk_test_[0-9a-zA-Z]{24}`,
	"API Key":                         `(?i)(api[-_]?key|apikey|key)[\s]*=\s*['"]?([0-9a-zA-Z-_.]+)['"]?`,
	"Firebase Secret":                 `(?i)FIREBASE_SECRET\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Mailgun API Key":                 `(?i)key-[0-9a-zA-Z]{32}`,
	"Twilio Auth Token":               `(?i)TWILIO_AUTH_TOKEN\s*=\s*['"]?([0-9a-zA-Z]{32})['"]?`,
	"Slack Token":                     `(?i)xox[baprs]-[0-9]+-[0-9a-zA-Z]+`,
	"GitLab Personal Access Token":    `(?i)glpat-[0-9a-zA-Z]{20,40}`,
	"Heroku API Key":                  `(?i)HEROKU_API_KEY\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Discord Token":                   `(?i)ND[0-9]{18}\.[0-9]{12}\.[0-9a-zA-Z_-]{27}`,
	"JWT Secret":                      `(?i)(jwt|secret)[\s]*=\s*['"]?([^'"]+)['"]?`,
	"OAuth2 Client Secret":            `(?i)(client[-_]?secret|secret)[\s]*=\s*['"]?([^'"]+)['"]?`,
	"OpenAI API Key":                  `(?i)sk-[0-9a-zA-Z]{48}`,
	"Dropbox Access Token":            `(?i)DROPBOX_ACCESS_TOKEN\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"MongoDB Atlas Connection String": `(?i)mongodb(?:\+srv)?://[a-zA-Z0-9]+:[^@]+@([^/]+)/[a-zA-Z0-9]+`,
	"Zoom JWT Token":                  `(?i)eyJ(.*)\.[A-Za-z0-9-]+`,
	"PayPal Webhook ID":               `(?i)PAYPAL_WEBHOOK_ID\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Okta API Token":                  `(?i)okta[ -]?api[ -]?token\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Zoom API Key":                    `(?i)zoom[ -]?api[ -]?key\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Algolia API Key":                 `(?i)algolia[ -]?api[ -]?key\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Telegram API Key":                `(?i)telegram[ -]?api[ -]?key\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
	"Bitbucket Access Token":          `(?i)BITBUCKET_ACCESS_TOKEN\s*=\s*['"]?([0-9a-zA-Z_-]+)['"]?`,
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
