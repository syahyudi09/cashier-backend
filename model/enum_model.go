package model

type RoleEnum string

const (
	Admin    RoleEnum   = "Admin"
	User     RoleEnum   = "User"
	Cashier  RoleEnum   = "Cashier"
	Employee RoleEnum = "Employee"
)

type StatusEnum string

const (
	Active    StatusEnum = "Active"
	NonActive StatusEnum = "Non Active"
)
