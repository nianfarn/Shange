package main

type user struct {
	name, password, email string
}

func createUser(u *user) user {

	return user{}
}

func getUser(id int) user {
	return user{}
}

func deleteUser(id int) {
}

func recoverPassword(id int) error {
	return nil
}
