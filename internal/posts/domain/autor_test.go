package domain

import "testing"

func TestNewAutor(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		fullName string
		wantErr  error
	}{
		{
			name:     "success",
			id:       1,
			fullName: "John Doe",
			wantErr:  nil,
		},
		{
			name:     "error - zero ID",
			id:       0,
			fullName: "John Doe",
			wantErr:  ErrAutorIDRequired,
		},
		{
			name:     "error - negative ID",
			id:       -1,
			fullName: "John Doe",
			wantErr:  ErrAutorIDRequired,
		},
		{
			name:     "error - empty fullName",
			id:       1,
			fullName: "",
			wantErr:  ErrAutorFullNameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAutor(tt.id, tt.fullName)
			if err != tt.wantErr {
				t.Errorf("NewAutor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if a.ID() != tt.id {
					t.Errorf("NewAutor().ID() = %v, want %v", a.ID(), tt.id)
				}
				if a.FullName() != tt.fullName {
					t.Errorf("NewAutor().FullName() = %v, want %v", a.FullName(), tt.fullName)
				}
			}
		})
	}
}

func TestAutor_Equals(t *testing.T) {
	tests := []struct {
		name string
		e1   *Autor
		e2   *Autor
		want bool
	}{
		{
			name: "equal - same ID",
			e1:   mustNewAutor(t, 1, "John Doe"),
			e2:   mustNewAutor(t, 1, "Jane Doe"),
			want: true,
		},
		{
			name: "not equal - different ID",
			e1:   mustNewAutor(t, 1, "John Doe"),
			e2:   mustNewAutor(t, 2, "John Doe"),
			want: false,
		},
		{
			name: "not equal - nil other",
			e1:   mustNewAutor(t, 1, "John Doe"),
			e2:   nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e1.Equals(tt.e2); got != tt.want {
				t.Errorf("Autor.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutor_ToDTO(t *testing.T) {
	a := mustNewAutor(t, 1, "John Doe")

	dto := a.ToDTO()

	if dto.ID != 1 {
		t.Errorf("Autor.ToDTO().ID = %v, want %v", dto.ID, 1)
	}
	if dto.FullName != "John Doe" {
		t.Errorf("Autor.ToDTO().FullName = %v, want %v", dto.FullName, "John Doe")
	}
}

func mustNewAutor(t *testing.T, id int64, fullName string) *Autor {
	t.Helper()
	a, err := NewAutor(id, fullName)
	if err != nil {
		t.Fatalf("mustNewAutor() error = %v", err)
	}
	return a
}
