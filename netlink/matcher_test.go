package netlink

import "testing"

func TestRules(testing *testing.T) {
	type testcase struct {
		object interface{}
		valid  bool
	}

	t := testingWrapper{testing}

	// Given
	uevent := UEvent{
		Action: ADD,
		KObj:   "/devices/pci0000:00/0000:00:14.0/usb2/2-1/2-1:1.2/0003:04F2:0976.0008/hidraw/hidraw4",
		Env: map[string]string{
			"ACTION":    "add",
			"DEVPATH":   "/devices/pci0000:00/0000:00:14.0/usb2/2-1/2-1:1.2/0003:04F2:0976.0008/hidraw/hidraw4",
			"SUBSYSTEM": "hidraw",
			"MAJOR":     "247",
			"MINOR":     "4",
			"DEVNAME":   "hidraw4",
			"SEQNUM":    "2569",
		},
	}

	add := ADD.MatchingString()
	wrongAction := "can't match"

	// When
	rules := []RuleDefinition{
		{
			Action: nil,
			Env: map[string]string{
				"DEVNAME": "hidraw\\d+",
				"MAJOR":   "\\d+",
			},
		},

		{
			Action: &add,
			Env:    make(map[string]string),
		},

		{
			Action: nil,
			Env: map[string]string{
				"SUBSYSTEM": "can't match",
				"MAJOR":     "247",
			},
		},

		{
			Action: &add,
			Env: map[string]string{
				"SUBSYSTEM": "hidraw",
				"MAJOR":     "\\d+",
			},
		},

		{
			Action: &wrongAction,
			Env: map[string]string{
				"SUBSYSTEM": "hidraw",
				"MAJOR":     "\\d+",
			},
		},
		{
			Action: &add,
			Env: map[string]string{
				"DEVNAME": "hidraw\\d+",
				"MAJOR":   "247",
			},
		},
	}

	// Then
	testcases := []testcase{
		{
			object: &rules[0],
			valid:  true,
		},
		{
			object: &rules[1],
			valid:  true,
		},
		{
			object: &rules[2],
			valid:  false,
		},
		{
			object: &rules[3],
			valid:  true,
		},
		{
			object: &rules[4],
			valid:  false,
		},
		{
			object: &RuleDefinitions{[]RuleDefinition{rules[0], rules[4]}},
			valid:  true,
		},
		{
			object: &RuleDefinitions{[]RuleDefinition{rules[4], rules[0]}},
			valid:  true,
		},
		{
			object: &RuleDefinitions{[]RuleDefinition{rules[2], rules[4]}},
			valid:  false,
		},
		{
			object: &RuleDefinitions{[]RuleDefinition{rules[3], rules[1]}},
			valid:  true,
		},
		{
			object: &rules[5],
			valid:  true,
		},
	}

	for k, tcase := range testcases {
		matcher := tcase.object.(Matcher)

		err := matcher.Compile()
		t.FatalfIf(err != nil, "Testcase n°%d should compile without error, err: %v", k+1, err)

		ok := matcher.Evaluate(uevent)
		t.FatalfIf((ok != tcase.valid) && tcase.valid, "Testcase n°%d should evaluate event", k+1)
		t.FatalfIf((ok != tcase.valid) && !tcase.valid, "Testcase n°%d shouldn't evaluate event", k+1)
	}
}
