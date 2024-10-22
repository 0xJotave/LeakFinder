package scanner

import "regexp"

var patterns = map[string]string{
	"AWS Access Key":    `AKIA[0-9A-Z]{16}`,
	"AWS Secret Key":    `AWS_SECRET_ACCESS_KEY\s*=\s*([A-Za-z0-9/+=]{40})`,
	"GitHub Token":      `ghp_[0-9a-zA-Z]{36}`,
	"Bearer Token":      `[A-Za-z0-9-_.]{40}`,
	"Password Config":   `password\s*=\s*["']?([^\s"']+)["']?`,
	"MongoDB URI":       `mongodb://[a-zA-Z0-9._%+-]+:[^@]+@[^/]+/[a-zA-Z0-9._%+-]+`,
	"SQL Connection":    `(Data Source=|Server=)[^;]+;User ID=[^;]+;Password=[^;]+`,
	"SSH Private Key":   `-----BEGIN (RSA|DSA|ECDSA|ED25519) PRIVATE KEY-----[\\s\\S]+-----END \\1 PRIVATE KEY-----`, // Escapado corretamente
	"Stripe Secret Key": `sk_test_[0-9a-zA-Z]{24}`,
	"API Key":           `(?i)(api[-_]?key|apikey|key|secret)`,
	"Senha":             `(?i)(senha|password|pwd|pass)`,
	"Token":             `(?i)(token|auth|bearer)`,
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
