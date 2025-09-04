package entconst

var (
	ErrInvalidUserType = NewValidationError("invalid user type")
)

type UserType string

type DelegateStatus string

const (
	SystemAdmin UserType = "SystemAdmin" // システム管理者
	GuestUser   UserType = "GuestUser"   // ゲストユーザー
	MemberUser  UserType = "MemberUser"  //会員
)

func (u UserType) String() string {
	return string(u)
}

func (u UserType) StringJa() string {
	switch u {
	case SystemAdmin:
		return "システム管理者"
	case GuestUser:
		return "ゲストユーザー"
	case MemberUser:
		return "会員"

	default:
		return ""
	}
}

func UserTypeList() []UserType {
	return []UserType{
		SystemAdmin,
		GuestUser,
		MemberUser,
	}
}

func UserTypeFromString(userType string) (*UserType, error) {
	ret := UserType(userType)
	switch ret {
	case SystemAdmin, MemberUser, GuestUser:
		return &ret, nil
	default:
		return nil, ErrInvalidUserType
	}
}

func ConvertUserTypeToExportName(userType UserType) string {
	switch userType {
	case SystemAdmin:
		return "システム管理者"
	case GuestUser:
		return "代議員"
	case MemberUser:
		return "会員"
	default:
		return ""
	}
}
