package api

type User struct {
	ApiTokenName string `json:"apiTokenName"`
	Id           string `json:"id"`
	SubjectType  string `json:"subjectType"`
	Username     string `json:"username"`
	FullName     string `json:"fullName"`
	AvatarUrl    string `json:"avatarUrl"`
}
