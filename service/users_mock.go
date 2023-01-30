package service

import "context"

type UsersMock struct {
	AddUserFn    func(ctx context.Context, user *AddUser) (*User, error)
	DeleteUserFn func(ctx context.Context, id string) error
	UpdateUserFn func(ctx context.Context, user *UpdateUser) (*User, error)
	SearchUserFn func(ctx context.Context, page, page_size int64, country string) ([]*User, error)
}

func (m *UsersMock) AddUser(ctx context.Context, user *AddUser) (*User, error) {
	return m.AddUserFn(ctx, user)
}

func (m *UsersMock) DeleteUser(ctx context.Context, id string) error {
	return m.DeleteUserFn(ctx, id)
}

func (m *UsersMock) UpdateUser(ctx context.Context, user *UpdateUser) (*User, error) {
	return m.UpdateUserFn(ctx, user)
}
func (m *UsersMock) SearchUser(ctx context.Context, page, page_size int64, country string) ([]*User, error) {
	return m.SearchUserFn(ctx, page, page_size, country)
}
