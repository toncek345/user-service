package storage

import "context"

type MockUser struct {
	InsertUserFn func(ctx context.Context, user *InsertUser) (*UserModel, error)
	DeleteUserFn func(ctx context.Context, id string) error
	UpdateUserFn func(ctx context.Context, user *UpdateUser) (*UserModel, error)
	SearchUserFn func(ctx context.Context, filters *Filters, offset, limit int64) ([]*UserModel, error)
}

func (m *MockUser) InsertUser(ctx context.Context, user *InsertUser) (*UserModel, error) {
	return m.InsertUserFn(ctx, user)
}
func (m *MockUser) DeleteUser(ctx context.Context, id string) error {
	return m.DeleteUserFn(ctx, id)
}
func (m *MockUser) UpdateUser(ctx context.Context, user *UpdateUser) (*UserModel, error) {
	return m.UpdateUserFn(ctx, user)
}
func (m *MockUser) SearchUser(ctx context.Context, filters *Filters, offset, limit int64) ([]*UserModel, error) {
	return m.SearchUserFn(ctx, filters, offset, limit)
}
