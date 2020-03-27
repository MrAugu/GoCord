package gocord

import "fmt"

// Permission holds the bitfield and the permsissions included in it.
type Permission struct {
	Bitfield    int
	Permissions []string
}

// CalculateBitfield returns the Permission interface for that bitfield.
func CalculateBitfield(bitfield int) Permission {
	resultingPermission := Permission{Bitfield: bitfield}
	var permissionsSlice []string = make([]string, 0, 31)
	if (bitfield & 0x00000001) == 0x00000001 {
		fmt.Println(len(permissionsSlice))
		permissionsSlice = append(permissionsSlice, "CREATE_INSTANT_INVITE")
	}
	if (bitfield & 0x00000002) == 0x00000002 {
		permissionsSlice = append(permissionsSlice, "KICK_MEMBERS")
	}
	if (bitfield & 0x00000004) == 0x00000004 {
		permissionsSlice = append(permissionsSlice, "BAN_MEMBERS")
	}
	if (bitfield & 0x00000008) == 0x00000008 {
		permissionsSlice = append(permissionsSlice, "ADMINISTRATOR")
	}
	if (bitfield & 0x00000010) == 0x00000010 {
		permissionsSlice = append(permissionsSlice, "MANAGE_CHANNELS")
	}
	if (bitfield & 0x00000020) == 0x00000020 {
		permissionsSlice = append(permissionsSlice, "MANAGE_GUILD")
	}
	if (bitfield & 0x00000040) == 0x00000040 {
		permissionsSlice = append(permissionsSlice, "ADD_REACTIONS")
	}
	if (bitfield & 0x00000080) == 0x00000080 {
		permissionsSlice = append(permissionsSlice, "VIEW_AUDIT_LOG")
	}
	if (bitfield & 0x00000400) == 0x00000400 {
		permissionsSlice = append(permissionsSlice, "VIEW_CHANNEL")
	}
	if (bitfield & 0x00000800) == 0x00000800 {
		permissionsSlice = append(permissionsSlice, "SEND_MESSAGES")
	}
	if (bitfield & 0x00001000) == 0x00001000 {
		permissionsSlice = append(permissionsSlice, "SEND_TTS_MESSAGES")
	}
	if (bitfield & 0x00002000) == 0x00002000 {
		permissionsSlice = append(permissionsSlice, "MANAGE_MESSAGES")
	}
	if (bitfield & 0x00004000) == 0x00004000 {
		permissionsSlice = append(permissionsSlice, "EMBED_LINKS")
	}
	if (bitfield & 0x00008000) == 0x00008000 {
		permissionsSlice = append(permissionsSlice, "ATTACH_FILES")
	}
	if (bitfield & 0x00010000) == 0x00010000 {
		permissionsSlice = append(permissionsSlice, "READ_MESSAGE_HISTORY")
	}
	if (bitfield & 0x00020000) == 0x00020000 {
		permissionsSlice = append(permissionsSlice, "MENTION_EVERYONE")
	}
	if (bitfield & 0x00040000) == 0x00040000 {
		permissionsSlice = append(permissionsSlice, "USE_EXTERNAL_EMOJIS")
	}
	if (bitfield & 0x00080000) == 0x00080000 {
		permissionsSlice = append(permissionsSlice, "VIEW_GUILD_INSIGHTS")
	}
	if (bitfield & 0x00100000) == 0x00100000 {
		permissionsSlice = append(permissionsSlice, "CONNECT")
	}
	if (bitfield & 0x00200000) == 0x00200000 {
		permissionsSlice = append(permissionsSlice, "SPEAK")
	}
	if (bitfield & 0x00400000) == 0x00400000 {
		permissionsSlice = append(permissionsSlice, "MUTE_MEMBERS")
	}
	if (bitfield & 0x00800000) == 0x00800000 {
		permissionsSlice = append(permissionsSlice, "DEAFEN_MEMBERS")
	}
	if (bitfield & 0x01000000) == 0x01000000 {
		permissionsSlice = append(permissionsSlice, "MOVE_MEMBERS")
	}
	if (bitfield & 0x02000000) == 0x02000000 {
		permissionsSlice = append(permissionsSlice, "USE_VAD")
	}
	if (bitfield & 0x00000100) == 0x00000100 {
		permissionsSlice = append(permissionsSlice, "PRIORITY_SPEAKER")
	}
	if (bitfield & 0x00000200) == 0x00000200 {
		permissionsSlice = append(permissionsSlice, "STREAM")
	}
	if (bitfield & 0x04000000) == 0x04000000 {
		permissionsSlice = append(permissionsSlice, "CHANGE_NICKNAME")
	}
	if (bitfield & 0x08000000) == 0x08000000 {
		permissionsSlice = append(permissionsSlice, "MANAGE_NICKNAMES")
	}
	if (bitfield & 0x10000000) == 0x10000000 {
		permissionsSlice = append(permissionsSlice, "MANAGE_ROLES")
	}
	if (bitfield & 0x20000000) == 0x20000000 {
		permissionsSlice = append(permissionsSlice, "MANAGE_WEBHOOKS")
	}
	if (bitfield & 0x40000000) == 0x40000000 {
		permissionsSlice = append(permissionsSlice, "MANAGE_EMOJIS")
	}
	resultingPermission.Permissions = permissionsSlice
	return resultingPermission
}
