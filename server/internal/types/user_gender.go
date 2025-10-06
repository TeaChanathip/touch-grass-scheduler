package types

type UserGender BaseStringEnum

const (
	UserGenderMale   UserGender = "male"
	UserGenderFemale UserGender = "female"
	UserGenderOther  UserGender = "other"
	UserGenderPNTS   UserGender = "prefer_not_to_say"
)
