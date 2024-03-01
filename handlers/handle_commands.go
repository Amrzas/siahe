package handlers

import (
	"gorm.io/gorm"

	"github.com/DearRude/fumTheatreBot/database"
	in "github.com/DearRude/fumTheatreBot/internals"
)

func handleCommands(u in.UpdateMessage) error {
	command := getCommandName(u.Message)
	if command == "" {
		return nil
	}

	StateMap.Set(u.PeerUser.UserID, in.CommandState)

	switch command {
	case "start":
		return startCommand(u)
	case "signup":
		return signupCommand(u)
	case "deleteAccount":
		return deleteAccountCommand(u)
	case "getAccount":
		return getAccountCommand(u)
	default:
		return nil
	}
}

func startCommand(u in.UpdateMessage) error {
	if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageStart(u.PeerUser.UserID)...); err != nil {
		return err
	}
	return nil
}

func signupCommand(u in.UpdateMessage) error {
	var user database.User
	if err := db.Model(&database.User{}).First(&user, u.PeerUser.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Accepted. Ask for first name
			if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAskFirstName()...); err != nil {
				return err
			}
			StateMap.Set(u.PeerUser.UserID, in.SignUpAskFirstName)
			return nil
		} else {
			return err
		}
	}

	// User is found. They can't sign up again
	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageYouAlreadySignedUp(user.FirstName)...)
	return err
}

func deleteAccountCommand(u in.UpdateMessage) error {
	res := db.Delete(&database.User{}, u.PeerUser.UserID)
	if err := res.Error; err != nil {
		return err
	}

	// User does not have an account
	if res.RowsAffected <= 0 {
		if _, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageUserHasNoAccount()...); err != nil {
			return err
		}
		return nil
	}

	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageAccountDeleted()...)
	return err
}

func getAccountCommand(u in.UpdateMessage) error {
	var user database.User
	if err := db.Model(&database.User{}).First(&user, u.PeerUser.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User has no account
			_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessageUserHasNoAccount()...)
			return err
		} else {
			return err
		}
	}

	_, err := sender.Reply(u.Ent, u.Unm).StyledText(u.Ctx, in.MessagePrintUser(user)...)
	return err
}
