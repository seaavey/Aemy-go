package utils

import "regexp"

// Match returns a boolean indicating whether the pattern matches the string.
var URLRegex = regexp.MustCompile(`https?://[^\s]+`)
var TiktokRegex = regexp.MustCompile(`https:\/\/([a-z0-9]+\.)?tiktok\.com\/[^\s]+`)

