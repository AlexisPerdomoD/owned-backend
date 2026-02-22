package col

import "testing"

func TestSetBasicOperations(t *testing.T) {
	s := NewSet[int]()

	// Add
	if ok := s.Add(1); !ok {
		t.Fatal("expected Add to return true for new element")
	}

	if ok := s.Add(1); ok {
		t.Fatal("expected Add to return false for duplicate element")
	}

	// Contains
	if !s.Contains(1) {
		t.Fatal("expected set to contain element")
	}

	if s.Contains(2) {
		t.Fatal("did not expect set to contain missing element")
	}

	// Remove
	if ok := s.Remove(2); ok {
		t.Fatal("expected Remove to return false for missing element")
	}

	if ok := s.Remove(1); !ok {
		t.Fatal("expected Remove to return true for existing element")
	}

	if s.Contains(1) {
		t.Fatal("expected element to be removed")
	}
}

func TestSetSlice(t *testing.T) {
	s := NewSet[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	items := s.Slice()

	if len(items) != 3 {
		t.Fatalf("expected slice length 3, got %d", len(items))
	}

	seen := make(map[int]bool)
	for _, v := range items {
		seen[v] = true
	}

	for _, v := range []int{1, 2, 3} {
		if !seen[v] {
			t.Fatalf("missing value %d in slice", v)
		}
	}
}
