package domain

const (
	RoleAdmin = 1
	RoleTeacher = 2
)

var RoleText = map[int]string {
	RoleAdmin:		"Administrator",
	RoleTeacher: 	"Teacher", 
}