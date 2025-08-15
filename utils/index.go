package utils

import "time"

func Ucapan() string {
	jam := time.Now().Hour()

	switch {
	case jam >= 5 && jam < 11:
		return "Selamat Pagi 🌅"
	case jam >= 11 && jam < 15:
		return "Selamat Siang 🌞"
	case jam >= 15 && jam < 18:
		return "Selamat Sore 🌇"
	default:
		return "Selamat Malam 🌙"
	}
}
