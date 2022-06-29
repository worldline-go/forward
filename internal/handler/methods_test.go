package handler

import (
	"reflect"
	"testing"
)

func TestFilterMethods_Parse(t *testing.T) {
	type want struct {
		Allow map[string]struct{}
		Deny  map[string]struct{}
	}
	type match struct {
		Method string
		Check  bool
	}
	type args struct {
		methods []string
	}
	tests := []struct {
		name   string
		want   want
		matchs []match
		args   args
		str    string
	}{
		{
			name: "just deny",
			args: args{
				methods: []string{"-POST", "-PUT", "-DELETE"},
			},
			want: want{
				Deny: map[string]struct{}{
					"POST":   {},
					"PUT":    {},
					"DELETE": {},
				},
				Allow: map[string]struct{}{
					"*": {},
				},
			},
			matchs: []match{
				{
					Method: "POST",
					Check:  false,
				},
				{
					Method: "PUT",
					Check:  false,
				},
				{
					Method: "DELETE",
					Check:  false,
				},
				{
					Method: "GET",
					Check:  true,
				},
				{
					Method: "XYZ",
					Check:  true,
				},
			},
			str: "allow: *; deny: DELETE,POST,PUT",
		},
		{
			name: "allow",
			args: args{
				methods: []string{"GET"},
			},
			want: want{
				Allow: map[string]struct{}{
					"GET": {},
				},
			},
			matchs: []match{
				{
					Method: "POST",
					Check:  false,
				},
				{
					Method: "GET",
					Check:  true,
				},
				{
					Method: "XYZ",
					Check:  false,
				},
			},
			str: "allow: GET; deny: ",
		},
		{
			name: "mix",
			args: args{
				methods: []string{"GET", "-DELETE"},
			},
			want: want{
				Allow: map[string]struct{}{
					"GET": {},
				},
				Deny: map[string]struct{}{
					"DELETE": {},
				},
			},
			matchs: []match{
				{
					Method: "DELETE",
					Check:  false,
				},
				{
					Method: "POST",
					Check:  false,
				},
				{
					Method: "GET",
					Check:  true,
				},
				{
					Method: "XYZ",
					Check:  false,
				},
			},
			str: "allow: GET; deny: DELETE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFilterMethods()
			f.Parse(tt.args.methods)

			if !reflect.DeepEqual(f.Allow, tt.want.Allow) {
				t.Errorf("FilterMethods.Parse() = %v, want %v", f.Allow, tt.want.Allow)
			}

			for _, m := range tt.matchs {
				if got := f.Match(m.Method); got != m.Check {
					t.Errorf("FilterMethods.Match() = %v, want %v", got, m.Check)
				}
			}

			if f.String() != tt.str {
				t.Errorf("FilterMethods.String() = %v, want %v", f.String(), tt.str)
			}
		})
	}
}
