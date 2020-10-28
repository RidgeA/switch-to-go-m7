package crawler

import (
	"reflect"
	"sync"
	"testing"
)

func Test_set_Add(t *testing.T) {
	type fields struct {
		s map[string]struct{}
	}
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMap map[string]struct{}
	}{
		{
			name: "should add element",
			fields: fields{
				s: map[string]struct{}{},
			},
			args: args{
				v: "item",
			},
			wantMap: map[string]struct{}{"item": {}},
		},
		{
			name: "should not add an element if it exists",
			fields: fields{
				s: map[string]struct{}{"item": {}},
			},
			args: args{
				v: "item",
			},
			wantMap: map[string]struct{}{"item": {}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &set{
				RWMutex: sync.RWMutex{},
				s:       tt.fields.s,
			}

			s.Add(tt.args.v)

			if !reflect.DeepEqual(s.s, tt.wantMap) {
				t.Errorf("Want internal map: %v, got: %v", tt.wantMap, s.s)
			}
		})
	}
}

func Test_set_Has(t *testing.T) {
	type fields struct {
		s map[string]struct{}
	}
	type args struct {
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "should return true if element exists",
			fields: fields{
				s: map[string]struct{}{"item": {}},
			},
			args: args{
				v: "item",
			},
			want: true,
		},
		{
			name: "should return false if element exists",
			fields: fields{
				s: map[string]struct{}{},
			},
			args: args{
				v: "item",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &set{
				RWMutex: sync.RWMutex{},
				s:       tt.fields.s,
			}
			if got := s.Has(tt.args.v); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Del(t *testing.T) {
	type fields struct {
		s map[string]struct{}
	}
	type args struct {
		v string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMap map[string]struct{}
	}{
		{
			name: "should delete element",
			fields: fields{
				s: map[string]struct{}{"item": {}},
			},
			args: args{
				v: "item",
			},
			wantMap: map[string]struct{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &set{
				RWMutex: sync.RWMutex{},
				s:       tt.fields.s,
			}

			s.Del(tt.args.v)

			if !reflect.DeepEqual(s.s, tt.wantMap) {
				t.Errorf("Want internal map: %v, got: %v", tt.wantMap, s.s)
			}
		})
	}
}

func Test_set_Empty(t *testing.T) {
	type fields struct {
		s map[string]struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should return true for empty set",
			fields: fields{
				s: map[string]struct{}{},
			},
			want: true,
		},
		{
			name: "should return false for non-empty set",
			fields: fields{
				s: map[string]struct{}{"item":{}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &set{
				RWMutex: sync.RWMutex{},
				s:       tt.fields.s,
			}
			if got := s.Empty(); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}
