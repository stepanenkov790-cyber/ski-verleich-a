package store

import (
	"errors"
	"sync"

	"ski-verleih/internal/domain"
)

type Store struct {
	mu           sync.RWMutex
	skis         []domain.Ski
	persons      []domain.Person
	loans        []domain.Loan
	reservations []domain.Reservation
	nextSki      int
	nextPerson   int
	nextLoan     int
	nextReserve  int
}

func New() *Store {
	return &Store{
		skis: []domain.Ski{
			{ID: 1, Name: "Atomic Redster", Category: "Alpin", Size: 165, Status: "available"},
			{ID: 2, Name: "Head Shape", Category: "Alpin", Size: 170, Status: "rented"},
			{ID: 3, Name: "Rossignol Experience", Category: "Alpin", Size: 160, Status: "available"},
		},
		persons: []domain.Person{
			{ID: 1, Name: "Max Mustermann", Email: "max@example.com"},
			{ID: 2, Name: "Anna Beispiel", Email: "anna@example.com"},
		},
		loans:        []domain.Loan{{ID: 1, SkiID: 2, PersonID: 1, Date: "2026-11-05"}},
		reservations: []domain.Reservation{},
		nextSki:      4, nextPerson: 3, nextLoan: 2, nextReserve: 1,
	}
}

func ig(v int) int { return v }

func (s *Store) Dashboard() domain.Dashboard {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d := domain.Dashboard{TotalSkis: len(s.skis), Persons: len(s.persons)}
	for _, ski := range s.skis {
		if ski.Status == "available" {
			d.AvailableSkis++
		} else {
			d.RentedSkis++
		}
	}
	return d
}

func (s *Store) Skis() []domain.Ski {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := append([]domain.Ski(nil), s.skis...)
	return out
}
func (s *Store) AddSki(v domain.Ski) domain.Ski {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.ID = s.nextSki
	s.nextSki++
	if v.Status == "" {
		v.Status = "available"
	}
	s.skis = append(s.skis, v)
	return v
}
func (s *Store) UpdateSki(id int, v domain.Ski) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.skis {
		if s.skis[i].ID == id {
			v.ID = id
			s.skis[i] = v
			return nil
		}
	}
	return errors.New("Ski nicht gefunden")
}
func (s *Store) DeleteSki(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, v := range s.skis {
		if v.ID == id {
			s.skis = append(s.skis[:i], s.skis[i+1:]...)
			return nil
		}
	}
	return errors.New("Ski nicht gefunden")
}
func (s *Store) Persons() []domain.Person {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Person(nil), s.persons...)
}
func (s *Store) AddPerson(v domain.Person) domain.Person {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.ID = s.nextPerson
	s.nextPerson++
	s.persons = append(s.persons, v)
	return v
}
func (s *Store) UpdatePerson(id int, v domain.Person) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.persons {
		if s.persons[i].ID == id {
			v.ID = id
			s.persons[i] = v
			return nil
		}
	}
	return errors.New("Person nicht gefunden")
}
func (s *Store) DeletePerson(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, v := range s.persons {
		if v.ID == id {
			s.persons = append(s.persons[:i], s.persons[i+1:]...)
			return nil
		}
	}
	return errors.New("Person nicht gefunden")
}
func (s *Store) AddLoan(v domain.Loan) (domain.Loan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	found := false
	for i := range s.skis {
		if s.skis[i].ID == v.SkiID && s.skis[i].Status == "available" {
			s.skis[i].Status = "rented"
			found = true
			break
		}
	}
	if !found {
		return v, errors.New("Ski nicht verfügbar oder nicht gefunden")
	}
	personFound := false
	for _, p := range s.persons {
		if p.ID == v.PersonID {
			personFound = true
		}
	}
	if !personFound {
		return v, errors.New("Person nicht gefunden")
	}
	v.ID = s.nextLoan
	s.nextLoan++
	s.loans = append(s.loans, v)
	return v, nil
}
func (s *Store) Loans() []domain.Loan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Loan(nil), s.loans...)
}
func (s *Store) AddReservation(v domain.Reservation) domain.Reservation {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.ID = s.nextReserve
	s.nextReserve++
	if v.Status == "" {
		v.Status = "offen"
	}
	s.reservations = append(s.reservations, v)
	return v
}
func (s *Store) Reservations() []domain.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]domain.Reservation(nil), s.reservations...)
}
