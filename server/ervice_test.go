package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	user "github.com/Santosh-Sohan/user-service/api"
	"github.com/Santosh-Sohan/user-service/repository"
	"github.com/Santosh-Sohan/user-service/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Santosh-Sohan/user-service/db"
)

// ========== MOCK REPOSITORY FOR UNIT TEST ==========

type mockRepo struct{}

func (m *mockRepo) CreateUser(u *user.User) error                    { return nil }
func (m *mockRepo) UpdateUser(r *user.UpdateUserRequest) error       { return nil }
func (m *mockRepo) UpdateContact(r *user.UpdateContactRequest) error { return nil }
func (m *mockRepo) BlockUser(id string) error                        { return nil }
func (m *mockRepo) UnblockUser(id string) error                      { return nil }
func (m *mockRepo) GetUserByEmailOrPhone(phone, email string) (*user.User, error) {
	return &user.User{Id: "mock-id", FirstName: "Mock"}, nil
}

// ========== UNIT TEST ==========

func TestUnit_CreateUser(t *testing.T) {
	svc := &service.UserServiceServer{Repo: &mockRepo{}}
	req := &user.CreateUserRequest{
		User: &user.User{
			FirstName:   "Unit",
			LastName:    "Test",
			PhoneNumber: "1234567890",
			Email:       "unit@test.com",
		},
	}

	resp, err := svc.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Unit", resp.User.FirstName)
}

// ========== INTEGRATION TEST (Real Cassandra) ==========

func TestIntegration_CreateAndFetchUser(t *testing.T) {
	db.InitCassandra() // assumes you already have this function to connect to DB
	defer db.Session.Close()

	svc := &service.UserServiceServer{Repo: repository.NewUserRepository()}

	newUser := &user.User{
		Id:          uuid.New().String(),
		FirstName:   "Integration",
		LastName:    "Test",
		Gender:      "Other",
		DateOfBirth: "2000-01-01",
		PhoneNumber: "9876543210",
		Email:       "integration@test.com",
		IsBlocked:   false,
	}

	_, err := svc.CreateUser(context.Background(), &user.CreateUserRequest{User: newUser})
	assert.NoError(t, err)

	resp, err := svc.GetUserByEmailOrPhone(context.Background(), &user.EmailOrPhoneRequest{
		PhoneNumber: newUser.PhoneNumber,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Integration", resp.User.FirstName)
}

// ========== BENCHMARK: gRPC vs REST ==========

func Benchmark_gRPC_CreateUser(b *testing.B) {
	db.InitCassandra()
	defer db.Session.Close()

	svc := &service.UserServiceServer{Repo: repository.NewUserRepository()}
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, _ = svc.CreateUser(ctx, &user.CreateUserRequest{
			User: &user.User{
				FirstName:   "gRPC",
				LastName:    "Bench",
				PhoneNumber: uuid.New().String(),
				Email:       uuid.New().String() + "@bench.com",
			},
		})
	}
}

func Benchmark_REST_CreateUser(b *testing.B) {
	go main()                   // start server
	time.Sleep(2 * time.Second) // wait for it to start

	url := "http://localhost:8080/v1/users"

	for i := 0; i < b.N; i++ {
		userPayload := map[string]interface{}{
			"id":            uuid.New().String(),
			"first_name":    "REST",
			"last_name":     "Bench",
			"gender":        "M",
			"date_of_birth": "1999-01-01",
			"phone_number":  uuid.New().String(),
			"email":         uuid.New().String() + "@rest.com",
		}

		body, _ := json.Marshal(userPayload)
		_, _ = http.Post(url, "application/json", bytes.NewReader(body))
	}
}
