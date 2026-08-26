package util

import "strings"

// CountryCodeToFlag converts an ISO 3166-1 alpha-2 country code (e.g., "DE", "US") to an emoji flag.
func CountryCodeToFlag(countryCode string) string {
	code := strings.ToUpper(strings.TrimSpace(countryCode))
	if len(code) != 2 {
		return "🌐"
	}
	// Regional indicator symbol letter A is 0x1F1E6 (127462 in decimal)
	// Rune 'A' is 65
	const base = 127397
	r1 := rune(code[0]) + base
	r2 := rune(code[1]) + base
	return string([]rune{r1, r2})
}

// DetectCountry tries to infer country from remark/server name or returns fallback
func DetectCountry(remark string) string {
	// Most subscription remarks already embed a flag emoji (e.g. "🇸🇨 MyServer").
	// Decoding it directly covers every ISO country, not just the keywords below.
	if code := extractFlagCode(remark); code != "" {
		return code
	}

	r := strings.ToUpper(remark)
	countryKeywords := map[string]string{
		"GERMANY": "DE", "DEUTSCHLAND": "DE", "DE-": "DE", " DE ": "DE", "[DE]": "DE", "🇩🇪": "DE",
		"NETHERLAND": "NL", "HOLLAND": "NL", "NL-": "NL", " NL ": "NL", "[NL]": "NL", "🇳🇱": "NL",
		"UNITED STATES": "US", "AMERICA": "US", "USA": "US", "US-": "US", " US ": "US", "[US]": "US", "🇺🇸": "US",
		"UNITED KINGDOM": "GB", "BRITAIN": "GB", "UK-": "GB", " UK ": "GB", "[UK]": "GB", "🇬🇧": "GB",
		"FRANCE": "FR", "FR-": "FR", " FR ": "FR", "[FR]": "FR", "🇫🇷": "FR",
		"TURKEY": "TR", "TURKIYE": "TR", "TR-": "TR", " TR ": "TR", "[TR]": "TR", "🇹🇷": "TR",
		"FINLAND": "FI", "FI-": "FI", " FI ": "FI", "[FI]": "FI", "🇫🇮": "FI",
		"SWEDEN": "SE", "SE-": "SE", " SE ": "SE", "[SE]": "SE", "🇸🇪": "SE",
		"CANADA": "CA", "CA-": "CA", " CA ": "CA", "[CA]": "CA", "🇨🇦": "CA",
		"SINGAPORE": "SG", "SG-": "SG", " SG ": "SG", "[SG]": "SG", "🇸🇬": "SG",
		"JAPAN": "JP", "JP-": "JP", " JP ": "JP", "[JP]": "JP", "🇯🇵": "JP",
		"IRAN": "IR", "IR-": "IR", " IR ": "IR", "[IR]": "IR", "🇮🇷": "IR",
		"UNITED ARAB EMIRATES": "AE", "DUBAI": "AE", "AE-": "AE", "🇦🇪": "AE",
	}

	for kw, code := range countryKeywords {
		if strings.Contains(r, kw) {
			return code
		}
	}

	// Look for 2-letter tokens
	tokens := strings.FieldsFunc(r, func(r rune) bool {
		return r == '-' || r == '_' || r == '|' || r == ' ' || r == '[' || r == ']' || r == '(' || r == ')'
	})
	for _, tok := range tokens {
		if len(tok) == 2 {
			if _, ok := countryKeywords[tok]; ok {
				return tok
			}
		}
	}

	return "UN"
}

// extractFlagCode reverses CountryCodeToFlag's regional-indicator encoding to
// recover the 2-letter code from a flag emoji already present in the text.
func extractFlagCode(remark string) string {
	const base = 127397 // matches the base used in CountryCodeToFlag
	runes := []rune(remark)
	for i := 0; i < len(runes)-1; i++ {
		a, b := runes[i], runes[i+1]
		if a >= 0x1F1E6 && a <= 0x1F1FF && b >= 0x1F1E6 && b <= 0x1F1FF {
			return string([]rune{a - base, b - base})
		}
	}
	return ""
}
