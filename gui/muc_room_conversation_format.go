package gui

import (
	"time"

	"github.com/coyim/coyim/i18n"
)

// displayTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTimestamp() {
	c.addTextWithTag(i18n.Localf("[%s] ", getTimestamp()), "timestamp")
}

// displayNotificationWhenOccupantJoinedRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantJoinedRoom(nickname string) {
	c.displayTextLineWithTimestamp(i18n.Localf("%s joined the room", nickname), "joinedRoom")
}

// displayNotificationWhenOccupantLeftTheRoom MUST be called from the UI thread
func (c *roomViewConversation) displayNotificationWhenOccupantLeftTheRoom(nickname string) {
	c.displayTextLineWithTimestamp(i18n.Localf("%s left the room", nickname), "leftRoom")
}

// displayNickname MUST be called from the UI thread
func (c *roomViewConversation) displayNickname(nickname string) {
	c.addTextWithTag(i18n.Localf("%s: ", nickname), "nickname")
}

// displayRoomSubject MUST be called from the UI thread
func (c *roomViewConversation) displayRoomSubject(subject string) {
	c.displayTextLineWithTimestamp(subject, "subject")
}

// displayMessage MUST be called from the UI thread
func (c *roomViewConversation) displayMessage(message string) {
	c.addTextWithTag(message, "message")
}

// displayInfoMessage MUST be called from the UI thread
func (c *roomViewConversation) displayInfoMessage(message string) {
	c.addTextWithTag(message, "infoMessage")
}

// displayNewLiveMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewLiveMessage(nickname, message string) {
	c.displayTimestamp()

	c.displayNickname(nickname)
	c.displayMessage(message)

	c.addNewLine()
}

// displayNewInfoMessage MUST be called from the UI thread
func (c *roomViewConversation) displayNewInfoMessage(message string) {
	c.displayTimestamp()
	c.displayInfoMessage(message)
	c.addNewLine()
}

// displayWarningMessage MUST be called from the UI thread
func (c *roomViewConversation) displayWarningMessage(message string) {
	c.displayTextLineWithTimestamp(message, "warning")
}

// displayErrorMessage MUST be called from the UI thread
func (c *roomViewConversation) displayErrorMessage(message string) {
	c.displayTextLineWithTimestamp(message, "error")
}

// displayTextLineWithTimestamp MUST be called from the UI thread
func (c *roomViewConversation) displayTextLineWithTimestamp(text string, tag string) {
	c.displayTimestamp()
	c.addTextWithTag(text, tag)
	c.addNewLine()
}

func getTimestamp() string {
	return time.Now().Format("15:04:05")
}
