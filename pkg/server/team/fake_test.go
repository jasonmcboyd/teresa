package team

import (
	"testing"

	"github.com/luizalabs/teresa-api/models/storage"
	"github.com/luizalabs/teresa-api/pkg/server/user"
)

func TestFakeOperationsCreate(t *testing.T) {
	fake := NewFakeOperations()

	expectedEmail := "teresa@luizalabs.com"
	expectedName := "teresa"
	expectedURL := "http://teresa.io"

	if err := fake.Create(expectedName, expectedEmail, expectedURL); err != nil {
		t.Fatal("error trying to create a fake team", err)
	}

	fakeTeam := fake.(*FakeOperations).Storage[expectedName]
	if fakeTeam == nil {
		t.Fatal("expected a valid team, got nil")
	}

	if fakeTeam.Name != expectedName {
		t.Errorf("expected %s, got %s", expectedName, fakeTeam.Name)
	}
	if fakeTeam.Email != expectedEmail {
		t.Errorf("expected %s, got %s", expectedEmail, fakeTeam.Email)
	}
	if fakeTeam.URL != expectedURL {
		t.Errorf("expected %s, got %s", expectedURL, fakeTeam.URL)
	}
}

func TestFakeOperationsCreateTeamAlreadyExists(t *testing.T) {
	fake := NewFakeOperations()

	teamName := "teresa"
	fake.(*FakeOperations).Storage[teamName] = &storage.Team{Name: teamName}

	if err := fake.Create(teamName, "", ""); err != ErrTeamAlreadyExists {
		t.Errorf("expected ErrTeamAlreadyExists, got %v", err)
	}
}

func TestFakeOperationsAddUser(t *testing.T) {
	fake := NewFakeOperations()

	expectedUserEmail := "gopher"
	fake.(*FakeOperations).UserOps.(*user.FakeOperations).Storage[expectedUserEmail] = &storage.User{Email: expectedUserEmail}

	expectedTeam := "teresa"
	if err := fake.Create(expectedTeam, "", ""); err != nil {
		t.Fatal("error trying to create a fake team:", err)
	}

	if err := fake.AddUser(expectedTeam, expectedUserEmail); err != nil {
		t.Errorf("error trying on add user to a team: %v", err)
	}
}

func TestFakeOperationsAddUserTeamNotFound(t *testing.T) {
	fake := NewFakeOperations()

	if err := fake.AddUser("teresa", "gopher"); err != ErrNotFound {
		t.Errorf("expected error ErrNotFound, got %v", err)
	}
}

func TestFakeOperationsAddUserUserNotFound(t *testing.T) {
	fake := NewFakeOperations()

	expectedTeam := "teresa"
	if err := fake.Create(expectedTeam, "", ""); err != nil {
		t.Fatal("error trying to create a fake team:", err)
	}

	if err := fake.AddUser(expectedTeam, "gopher"); err != user.ErrNotFound {
		t.Errorf("expected error ErrNotFound, got %v", err)
	}
}

func TestFakeOperationsAddUserUserAlreadyInTeam(t *testing.T) {
	fake := NewFakeOperations()

	expectedUserEmail := "gopher"
	expectedName := "teresa"
	fake.(*FakeOperations).UserOps.(*user.FakeOperations).Storage[expectedUserEmail] = &storage.User{Email: expectedUserEmail}
	fake.(*FakeOperations).Storage[expectedName] = &storage.Team{
		Name:  expectedName,
		Users: []storage.User{storage.User{Email: expectedUserEmail}},
	}

	if err := fake.AddUser(expectedName, expectedUserEmail); err != ErrUserAlreadyInTeam {
		t.Errorf("expected error ErrUserAlreadyInTeam, got %v", err)
	}
}

func TestFakeOperationsList(t *testing.T) {
	var testData = []struct {
		teamName   string
		usersEmail []string
	}{
		{teamName: "Empty"},
		{teamName: "teresa", usersEmail: []string{"gopher", "k8s"}},
	}

	fake := NewFakeOperations()
	for _, tc := range testData {
		fakeTeam := &storage.Team{Name: tc.teamName}
		for _, email := range tc.usersEmail {
			fakeTeam.Users = append(fakeTeam.Users, storage.User{Email: email})
		}
		fake.(*FakeOperations).Storage[tc.teamName] = fakeTeam
	}

	teams, err := fake.List()
	if err != nil {
		t.Fatal("error on list teams:", err)
	}

	if len(teams) != len(testData) {
		t.Errorf("expected %d, got %d", len(testData), len(teams))
	}
}

func TestFakeOperationsListWithoutTeams(t *testing.T) {
	fake := NewFakeOperations()

	teams, err := fake.List()
	if err != nil {
		t.Fatal("error trying to list teams:", err)
	}
	if len(teams) > 0 {
		t.Errorf("expected 0, got %d", len(teams))
	}
}

func TestFakeOperationsListByUser(t *testing.T) {
	expectedUserEmail := "gopher"

	var testData = []struct {
		teamName   string
		usersEmail []string
	}{
		{teamName: "Empty"},
		{teamName: "teresa", usersEmail: []string{expectedUserEmail, "k8s"}},
		{teamName: "gophers", usersEmail: []string{expectedUserEmail, "john"}},
		{teamName: "vimers", usersEmail: []string{"k8s", "john"}},
	}

	fake := NewFakeOperations()
	for _, tc := range testData {
		fakeTeam := &storage.Team{Name: tc.teamName}
		for _, email := range tc.usersEmail {
			fakeTeam.Users = append(fakeTeam.Users, storage.User{Email: email})
		}
		fake.(*FakeOperations).Storage[tc.teamName] = fakeTeam
	}

	teams, err := fake.ListByUser(expectedUserEmail)
	if err != nil {
		t.Fatal("error on list teams:", err)
	}

	if len(teams) != 2 {
		t.Fatalf("expected 2, got %d", len(teams))
	}

	for _, currentTeam := range teams {
		if currentTeam.Name != "gophers" && currentTeam.Name != "teresa" {
			t.Errorf("expecter gophers or teresa, got %s", currentTeam.Name)
		}
	}
}

func TestFakeOperationsListByUserWithoutTeams(t *testing.T) {
	fake := NewFakeOperations()

	teams, err := fake.ListByUser("gopher@luizalabs.com")
	if err != nil {
		t.Fatal("error trying to list teams:", err)
	}
	if len(teams) > 0 {
		t.Errorf("expected 0, got %d", len(teams))
	}
}