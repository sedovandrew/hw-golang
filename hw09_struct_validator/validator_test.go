package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	House struct {
		Owners    []User `validate:"nested"`
		Country   string `validate:"regexp:^[A-Z][a-z]*$"`
		City      string `validate:"regexp:^[A-Z][a-z]*$"`
		Street    string `validate:"regexp:^([-0-9A-Z][-0-9a-z]*)( [-0-9A-Za-z]+)*$"`
		Number    int    `validate:"min:1"`
		Apartment int
	}

	User struct {
		ID     string `json:"id" validate:"len:36|regexp:^[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}$"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	WrongRegexp struct {
		Field string `validate:"regexp:five[zero"`
	}

	UnknownValidate struct {
		Name string `validate:"unknown:3"`
	}

	NotPublic struct {
		notPublicField string `validate:"len:1"`
	}

	ValidationValueNotSet struct {
		NotValue int `validate:"min:5|max"`
	}
)

//nolint:funlen
func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// case_0: Good example
		{
			in:          App{"5 len"},
			expectedErr: nil,
		},

		// case_1: Length error
		{
			in: App{"not 5 len"},
			expectedErr: ValidationErrors{
				ValidationError{
					"Version",
					ErrLengthValidation,
				},
			},
		},

		// case_2: Good example
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},

		// case_3: Not allowed
		{
			in: Response{
				Code: 403,
				Body: "Forbidden",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					"Code",
					ErrInValidation,
				},
			},
		},

		// case_4: All field satisfy the checks.
		{
			in: User{
				ID:     "503f659e-c113-45d6-ad73-183ca5601fef",
				Name:   "Alex",
				Age:    20,
				Email:  "alex@example.com",
				Role:   "stuff",
				Phones: []string{"79261234567", "79031234567"},
				meta:   []byte{},
			},
			expectedErr: nil,
		},

		// case_5: Good test with recursion
		{
			in: House{
				Owners: []User{
					{
						ID:     "ddae45b6-a0f6-4b9d-b97c-b4b54bc9f5e7",
						Name:   "Fred",
						Age:    30,
						Email:  "fred@example.com",
						Role:   "admin",
						Phones: []string{"79151234567"},
					},
				},
				Country: "Russia",
				City:    "Moscow",
				Street:  "1st Cavalry Army",
				Number:  1,
			},
			expectedErr: nil,
		},

		// case_6: Negative test with recursion
		{
			in: House{
				Owners: []User{
					{
						ID:     "123",
						Name:   "Eva",
						Age:    17,
						Email:  "Eva_email.com",
						Role:   "student",
						Phones: []string{"eva123", "12345678901", "eva"},
					},
					{
						ID:     "1a80203e-b925-4b7d-af4e-76e38ac1e676",
						Name:   "Adam",
						Age:    81,
						Email:  "adam@example.com",
						Role:   "worker",
						Phones: []string{"7903123adam"},
					},
				},
				Country: "Engl@nd",
				City:    "london",
				Street:  "baker St",
				Number:  0,
			},
			expectedErr: ValidationErrors{
				// House Country
				ValidationError{
					"Country",
					ErrRegexpValidation,
				},
				// House City
				ValidationError{
					"City",
					ErrRegexpValidation,
				},
				// House Street
				ValidationError{
					"Street",
					ErrRegexpValidation,
				},
				// House Number
				ValidationError{
					"Number",
					ErrMinValidation,
				},
				// Eva ID len
				ValidationError{
					"ID",
					ErrLengthValidation,
				},
				// Eva ID regexp
				ValidationError{
					"ID",
					ErrRegexpValidation,
				},
				// Eva Age
				ValidationError{
					"Age",
					ErrMinValidation,
				},
				// Eva Email
				ValidationError{
					"Email",
					ErrRegexpValidation,
				},
				// Eva Role
				ValidationError{
					"Role",
					ErrInValidation,
				},
				// Adam Age
				ValidationError{
					"Age",
					ErrMaxValidation,
				},
				// Adam Role
				ValidationError{
					"Role",
					ErrInValidation,
				},
				// Eva first phone
				ValidationError{
					"Phones",
					ErrLengthValidation,
				},
				// Eva third phone
				ValidationError{
					"Phones",
					ErrLengthValidation,
				},
			},
		},

		// case_7: Wrong validate
		{
			in: WrongRegexp{
				Field: "one,two,three",
			},
			expectedErr: ErrRegexpCheck,
		},

		// case_8: Unknown validate
		{
			in: UnknownValidate{
				Name: "Yes",
			},
			expectedErr: ErrUnknownCheck,
		},

		// case_9: A non-public field should not be validated.
		{
			in: NotPublic{
				notPublicField: "123",
			},
			expectedErr: nil,
		},

		// case_10: Validation parsing error
		{
			in: ValidationValueNotSet{
				NotValue: 9,
			},
			expectedErr: ErrParseCheck,
		},

		// case_11: Not structure.
		{
			in:          555,
			expectedErr: ErrNonStructCheck,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.Nil(t, err)
			} else {
				var errs ValidationErrors
				if errors.As(tt.expectedErr, &errs) {
					var vErr ValidationErrors
					require.ErrorAs(t, err, &vErr)
					for i, e := range vErr {
						require.Lessf(t, i, len(errs), "The nubmer of errors returned is greater than required")
						require.Equal(t, e.Field, errs[i].Field)
						require.ErrorIs(t, e.Err, errs[i].Err)
					}
					require.Equalf(t, len(errs), len(vErr), "The nubmer of errors returned is less than required")
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			}
		})
	}
}
