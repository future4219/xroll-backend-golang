package entconst

var (
	ErrInvalidUserType = NewValidationError("invalid user type")
)

type UserType string

type DelegateStatus string

const (
	SystemAdmin   UserType = "SystemAdmin"   // システム管理者
	DelegateUser  UserType = "DelegateUser"  // 代議員
	MemberUser    UserType = "MemberUser"    // 会員
	NonMemberUser UserType = "NonMemberUser" // 非会員
)

func (u UserType) String() string {
	return string(u)
}

func (u UserType) StringJa() string {
	switch u {
	case SystemAdmin:
		return "システム管理者"
	case DelegateUser:
		return "代議員"
	case MemberUser:
		return "会員"
	case NonMemberUser:
		return "非会員"
	default:
		return ""
	}
}

func UserTypeList() []UserType {
	return []UserType{
		SystemAdmin,
		DelegateUser,
		MemberUser,
		NonMemberUser,
	}
}

func UserTypeFromString(userType string) (*UserType, error) {
	ret := UserType(userType)
	switch ret {
	case SystemAdmin, MemberUser, NonMemberUser, DelegateUser:
		return &ret, nil
	default:
		return nil, ErrInvalidUserType
	}
}

func ConvertUserTypeToExportName(userType UserType) string {
	switch userType {
	case SystemAdmin:
		return "システム管理者"
	case DelegateUser:
		return "代議員"
	case MemberUser:
		return "会員"
	case NonMemberUser:
		return "非会員"
	default:
		return ""
	}
}
