package data

import "fmt"

const (
	// RoleNone represents XMPP muc 'none' role
	RoleNone = "none"
	// RoleVisitor represents XMPP muc 'visitor' role
	RoleVisitor = "visitor"
	// RoleParticipant represents XMPP muc 'participant' role
	RoleParticipant = "participant"
	// RoleModerator represents XMPP muc 'moderator' role
	RoleModerator = "moderator"
)

// Role represents the specific role that a user has inside a specific room
type Role interface {
	// HasVoice returns true if the user can speak in this room
	HasVoice() bool
	// WithVoice returns the closest role upwards that has voice privilege. For Participants and Moderators, it returns itself, otherwise it returns participant
	WithVoice() Role
	// AsModerator returns the closest role upwards that can act as a moderator
	AsModerator() Role
	// Name returns the string name of the role type
	Name() string
}

// NoneRole is a representation of MUC's "none" role
type NoneRole struct{}

// VisitorRole is a representation of MUC's "visitor" role
type VisitorRole struct{}

// ParticipantRole is a representation of MUC's "participant" role
type ParticipantRole struct{}

// ModeratorRole is a representation of MUC's "moderator" role
type ModeratorRole struct{}

// HasVoice implements Role interface
func (*NoneRole) HasVoice() bool { return false }

// HasVoice implements Role interface
func (*VisitorRole) HasVoice() bool { return false }

// HasVoice implements Role interface
func (*ParticipantRole) HasVoice() bool { return true }

// HasVoice implements Role interface
func (*ModeratorRole) HasVoice() bool { return true }

// WithVoice implements Role interface
func (*NoneRole) WithVoice() Role { return &ParticipantRole{} }

// WithVoice implements Role interface
func (*VisitorRole) WithVoice() Role { return &ParticipantRole{} }

// WithVoice implements Role interface
func (*ParticipantRole) WithVoice() Role { return &ParticipantRole{} }

// WithVoice implements Role interface
func (*ModeratorRole) WithVoice() Role { return &ModeratorRole{} }

// AsModerator implements Role interface
func (*NoneRole) AsModerator() Role { return &ModeratorRole{} }

// AsModerator implements Role interface
func (*VisitorRole) AsModerator() Role { return &ModeratorRole{} }

// AsModerator implements Role interface
func (*ParticipantRole) AsModerator() Role { return &ModeratorRole{} }

// AsModerator implements Role interface
func (*ModeratorRole) AsModerator() Role { return &ModeratorRole{} }

// Name implements Role interface
func (*NoneRole) Name() string { return RoleNone }

// Name implements Role interface
func (*VisitorRole) Name() string { return RoleVisitor }

// Name implements Role interface
func (*ParticipantRole) Name() string { return RoleParticipant }

// Name implements Role interface
func (*ModeratorRole) Name() string { return RoleModerator }

// RoleFromString returns the role object that matches the string given, or an error if the string given doesn't match a known role
func RoleFromString(s string) (Role, error) {
	switch s {
	case RoleNone:
		return &NoneRole{}, nil
	case RoleVisitor:
		return &VisitorRole{}, nil
	case RoleParticipant:
		return &ParticipantRole{}, nil
	case RoleModerator:
		return &ModeratorRole{}, nil
	default:
		return nil, fmt.Errorf("unknown role string: '%s'", s)
	}
}
