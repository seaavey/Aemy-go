// Package commands implements a registry for automatic command loading.
package commands

import (
	"aemy/types"
	"strings"
	"sync"
)

// CmdInfo holds information about a registered command
type CmdInfo struct {
	Handler types.CommandHandler
	Cat     string
}

// Registry holds all registered commands
var (
	registry = make(map[string]CmdInfo)
	mutex    = sync.RWMutex{}
	
	// Default category for commands without a specified category
	defaultCat = "Other"
)

// Register registers a command handler with one or more names
func Register(names []string, handler types.CommandHandler, cat string) {
	mutex.Lock()
	defer mutex.Unlock()
	
	// Use default category if none provided
	if cat == "" {
		cat = defaultCat
	}
	
	info := CmdInfo{
		Handler: handler,
		Cat:     cat,
	}
	
	for _, name := range names {
		registry[name] = info
	}
}

// Get retrieves a command handler by name
func Get(name string) (types.CommandHandler, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	
	info, exists := registry[name]
	if !exists {
		return nil, false
	}
	return info.Handler, true
}

// GetInfo retrieves full command information by name
func GetInfo(name string) (CmdInfo, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	
	info, exists := registry[name]
	return info, exists
}

// All returns all registered commands
func All() map[string]CmdInfo {
	mutex.RLock()
	defer mutex.RUnlock()
	
	// Create a copy to avoid concurrent access issues
	result := make(map[string]CmdInfo)
	for k, v := range registry {
		result[k] = v
	}
	return result
}

// ByCategory returns commands grouped by category (with lowercase category names)
func ByCategory() map[string]map[string]CmdInfo {
	mutex.RLock()
	defer mutex.RUnlock()
	
	result := make(map[string]map[string]CmdInfo)
	
	for name, info := range registry {
		// Convert category to lowercase
		category := strings.ToLower(info.Cat)
		
		// If category doesn't exist yet, create it
		if _, exists := result[category]; !exists {
			result[category] = make(map[string]CmdInfo)
		}
		// Add command to its category
		result[category][name] = info
	}
	
	return result
}

// MustRegister is a convenience function that registers a command and panics on error
func MustRegister(names []string, handler types.CommandHandler, cat string) {
	Register(names, handler, cat)
}