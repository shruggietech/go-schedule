package cron

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseCompositeStandardFields(t *testing.T) {
	tests := []struct {
		expr       string
		minute     []int
		hour       []int
		dom        []int
		month      []int
		dow        []int
		minuteStep int
	}{
		{"0 9,17 * * *", []int{0}, []int{9, 17}, rangeOf(1, 31), rangeOf(1, 12), rangeOf(0, 7), 0},
		{"30 8-17 * * 1-5", []int{30}, rangeOf(8, 17), rangeOf(1, 31), rangeOf(1, 12), []int{1, 2, 3, 4, 5}, 0},
		{"*/7 * * * *", []int{0, 7, 14, 21, 28, 35, 42, 49, 56}, rangeOf(0, 23), rangeOf(1, 31), rangeOf(1, 12), rangeOf(0, 7), 7},
		{"10-20/2 * * * *", []int{10, 12, 14, 16, 18, 20}, rangeOf(0, 23), rangeOf(1, 31), rangeOf(1, 12), rangeOf(0, 7), 0},
		{"0 0 1,15 JAN,MAR *", []int{0}, []int{0}, []int{1, 15}, []int{1, 3}, rangeOf(0, 7), 0},
		{"0 12 * * MON,WED,FRI,7", []int{0}, []int{12}, rangeOf(1, 31), rangeOf(1, 12), []int{0, 1, 3, 5}, 0},
		{"1,1,0-2 * * * *", []int{0, 1, 2}, rangeOf(0, 23), rangeOf(1, 31), rangeOf(1, 12), rangeOf(0, 7), 0},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			res, err := Parse(tt.expr)
			if err != nil || !res.OK {
				t.Fatalf("Parse: ok=%v err=%v refusal=%q", res.OK, err, res.Bad.Reason)
			}
			for name, pair := range map[string][2][]int{
				"minute": {res.Spec.Minute.Values, tt.minute},
				"hour":   {res.Spec.Hour.Values, tt.hour},
				"dom":    {res.Spec.DOM.Values, tt.dom},
				"month":  {res.Spec.Month.Values, tt.month},
				"dow":    {res.Spec.DOW.Values, tt.dow},
			} {
				if !reflect.DeepEqual(pair[0], pair[1]) {
					t.Errorf("%s values=%v, want %v", name, pair[0], pair[1])
				}
			}
			if res.Spec.Minute.Step != tt.minuteStep {
				t.Errorf("minute step=%d, want %d", res.Spec.Minute.Step, tt.minuteStep)
			}
		})
	}
}

func TestParseAcceptsCalendarWildcardSteps(t *testing.T) {
	for _, expr := range []string{"0 9 */2 * *", "0 9 * */3 *", "0 9 * * */2"} {
		res, err := Parse(expr)
		if err != nil || !res.OK {
			t.Errorf("Parse(%q): ok=%v err=%v refusal=%q", expr, res.OK, err, res.Bad.Reason)
		}
	}
}

func TestParseStillRefusesCronDayUnion(t *testing.T) {
	for _, expr := range []string{"0 0 13 * 5", "0 9 */2 * MON"} {
		res, err := Parse(expr)
		if err != nil {
			t.Fatal(err)
		}
		if res.OK || !strings.Contains(res.Bad.Reason, "either") {
			t.Fatalf("Parse(%q)=%+v, want DOM/DOW refusal", expr, res)
		}
	}
}
