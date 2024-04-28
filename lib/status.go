package lib

import "strings"

// StatusKind represents a status of a model
type StatusKind int

// Model statuses
const (
	StatusUnknown  			StatusKind = 0
	StatusOffline  			StatusKind = 1
	StatusOnline   			StatusKind = 2
	StatusNotFound 			StatusKind = 4
	StatusDenied   			StatusKind = 8

	StatusPrivatChat 		StatusKind = 16
	StatusFullPrivatChat   	StatusKind = 32
	StatusGroupPrivatChat 	StatusKind = 64
	StatusVipShow			StatusKind = 128
)

func (s StatusKind) String() string {
	if s == StatusUnknown || s == StatusOffline|StatusOnline|StatusNotFound|StatusDenied|StatusPrivatChat|StatusFullPrivatChat|StatusGroupPrivatChat|StatusVipShow {
		return "unknown"
	}
	var words []string
	if s&StatusOffline != 0 {
		words = append(words, "offline")
	}
	if s&StatusOnline != 0 {
		words = append(words, "online")
	}
	if s&StatusNotFound != 0 {
		words = append(words, "not found")
	}
	if s&StatusDenied != 0 {
		words = append(words, "denied")
	}
	if s&StatusPrivatChat != 0 {
		words = append(words, "private chat")
	}
	if s&StatusFullPrivatChat != 0 {
		words = append(words, "full private chat")
	}
	if s&StatusGroupPrivatChat != 0 {
		words = append(words, "group private chat")
	}
	if s&StatusVipShow != 0 {
		words = append(words, "vip show")
	}
	return strings.Join(words, "|")
}
