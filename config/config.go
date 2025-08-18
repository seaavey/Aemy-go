// Package config stores global configuration settings for the WhatsApp bot.
// This includes settings like command prefixes and bot owner identifiers,
// allowing for easy customization without modifying core application logic.
package config

// Prefixes defines a list of strings that are recognized as command prefixes.
var Prefixes = []string{"!", ".", "😂", "🔥", "🐱‍👤"}

// Owners contains a list of WhatsApp user IDs (phone numbers) that have
// administrative privileges. These users can access special commands
// and perform restricted actions.
var Owners = []string{
	"6289513081052",
}

// Self determines whether the bot should process its own messages.
// Setting this to true means the bot will react to commands sent from its own number.
// It is generally recommended to keep this false to prevent infinite loops.
var Self = true

// ReadStatus determines whether the bot should automatically mark incoming messages as "read".
// When set to true, the bot will send a read receipt for every received message,
// letting the sender know that the message has been read.
// It is recommended to set this to false if you want the bot to remain passive
// or to maintain privacy when monitoring messages.
var ReadStatus = true
