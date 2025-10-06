package types

type UserRole BaseStringEnum

const (
	UserRoleStudent  UserRole = "student"
	UserRoleTeacher  UserRole = "teacher"
	UserRoleGuardian UserRole = "guardian"
	UserRoleAdmin    UserRole = "admin"
)
