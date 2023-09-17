package main

import (
	"testing"
)

func Test_getNumVerses(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "withBonusWorksCorrectly",
			args: args{
				url: "https://www.churchofjesuschrist.org/study/manual/hymns/how-firm-a-foundation?lang=eng",
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNumVerses(tt.args.url); got != tt.want {
				t.Errorf("getNumVerses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getExtraVerses(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "withBonusWorksCorrectly",
			args: args{
				url: "https://www.churchofjesuschrist.org/study/manual/hymns/how-firm-a-foundation?lang=eng",
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExtraVerses(tt.args.url); got != tt.want {
				t.Errorf("getExtraVerses() = %v, want %v", got, tt.want)
			}
		})
	}
}
