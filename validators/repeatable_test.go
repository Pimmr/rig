package validators

import (
	"testing"
)

func TestToRepeatable(t *testing.T) {
	t.Run("IntRange", func(t *testing.T) {
		g := ToRepeatable(IntRange(2, 16))

		err := g(0)
		if err == nil {
			t.Errorf("ToRepeatable(IntRange(2, 16))(0): expected error, got nil")
		}

		err = g(int8(4))
		if err != nil {
			t.Errorf("ToRepeatable(IntRange(2, 16))(int8(4)): unexpected error %v", err)
		}

		err = g("foo")
		if err == nil {
			t.Errorf("ToRepeatable(IntRange(2, 16))(\"foo\"): expected error, got nil")
		}
	})

	t.Run("Not a function", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected a panic, got none")
				return
			}
		}()

		ToRepeatable(42)
	})

	t.Run("nil", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected a panic, got none")
				return
			}
		}()

		ToRepeatable(nil)
	})

	t.Run("Too many arguments", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected a panic, got none")
				return
			}
		}()

		ToRepeatable(func(int, int) error {
			return nil
		})
	})

	t.Run("Too many values returned", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected a panic, got none")
				return
			}
		}()

		ToRepeatable(func(int) (int, error) {
			return 0, nil
		})
	})

	t.Run("Not returning an error", func(t *testing.T) {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("expected a panic, got none")
				return
			}
		}()

		ToRepeatable(func(int) int {
			return 0
		})
	})
}
