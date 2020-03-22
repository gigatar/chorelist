package models

import (
	"testing"
	"time"
)

func TestChoreValidate(t *testing.T) {
	testCases := []struct {
		name           string
		in             Chore
		expectedStatus bool
	}{
		{
			name: "Good test - unassigned",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: true,
		},
		{
			name: "Good test - assigned",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "Assigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: true,
		},
		{
			name: "Good test - review",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "Review",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: true,
		},
		{
			name: "Good test - completed",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "Completed",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: true,
		},
		{
			name: "Bad test - status",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "Bad",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: false,
		},
		{
			name: "Bad test - Title Too Short",
			in: Chore{
				Title:       "",
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: false,
		},
		{
			name: "Bad test - Title Too Long",
			in: Chore{
				Title:       randSeq(129),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: false,
		},
		{
			name: "Good test - No Description",
			in: Chore{
				Title:     randSeq(5),
				Incentive: randSeq(5),
				Status:    "UnAssigned",
				DueDate:   time.Now().Unix(),
			},
			expectedStatus: true,
		},
		{
			name: "Bad test - Description Too Long",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(4097),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: false,
		},
		{
			name: "Good test - No Incentive",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: true,
		},
		{
			name: "Bad test - Incentive Too Long",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(1025),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix(),
			},
			expectedStatus: false,
		},
		{
			name: "Good test - No Due Date",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
			},
			expectedStatus: true,
		},
		{
			name: "Bad test - Due Date Too Far In Past",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix() - 86401,
			},
			expectedStatus: false,
		},
		{
			name: "Bad test - Due Date Too Far In Future",
			in: Chore{
				Title:       randSeq(5),
				Description: randSeq(5),
				Incentive:   randSeq(5),
				Status:      "UnAssigned",
				DueDate:     time.Now().Unix() + 31557601,
			},
			expectedStatus: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if !test.in.Validate() == test.expectedStatus {
				t.Error("Failed to properly validate chore")
				t.Fail()
			}
		})
	}
}
