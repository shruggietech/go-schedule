package cron

import (
	"testing"
	"time"
)

func BenchmarkCompileCron(b *testing.B) {
	now := time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	for _, expression := range []string{
		"0 9 * * *",
		"*/10 9-17 * * MON,WED,FRI",
		"0 0 29 FEB *",
	} {
		b.Run(expression, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, bad, err := Compile(expression, "America/New_York", now); err != nil || bad.Reason != "" {
					b.Fatalf("compile: refusal=%q err=%v", bad.Reason, err)
				}
			}
		})
	}
}
