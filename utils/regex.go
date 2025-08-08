package utils

import "regexp"

var TiktokRegex = regexp.MustCompile(`https:\/\/([a-z0-9]+\.)?tiktok\.com\/[^\s]+`)

