package models

import (
	"strings"
	"time"
)

// Chore data type.
type Chore struct {
	ID            string   `json:"id,omitempty" bson:"_id,omitempty" example:""`
	FamilyID      string   `json:"familyID,omitempty" bson:"familyID,omitempty" example:""`
	Title         string   `json:"title" bson:"title" example:"Mop Floors"`
	Description   string   `json:"description,omitempty" bson:"description,omitempty" example:"Mop floors in kitchen"`
	Status        string   `json:"status,omitempty" bson:"status,omitempty" example:"Unassigned"`
	Notes         []Note   `json:"notes,omitempty" bson:"notes,omitempty" example:""`
	DateCreated   int64    `json:"dateCreated,omitempty" bson:"dateCreated,omitempty" example:"1584588677"`
	DateCompleted int64    `json:"dateCompleted,omitempty" bson:"dateCompleted,omitempty" example:"1584588677"`
	DueDate       int64    `json:"dueDate,omitempty" bson:"dueDate,omitempty" example:"1584588677"`
	CreatorID     string   `json:"creatorID,omitempty" bson:"creatorID,omitempty" example:""`
	AssigneeID    []string `json:"assigneeID,omitempty" bson:"assigneeID,omitempty" example:""`
	VerifiedBy    string   `json:"verifiedBy,omitempty" bson:"verifiedBy,omitempty" example:""`
	Incentive     string   `json:"incentive,omitempty" bson:"incentive,omitempty" example:"1hr of reddit time"`
}

// Validate Chore input fields
func (c Chore) Validate() bool {
	if !c.ValidateTitle() {
		return false
	} else if !c.ValidateDescription() {
		return false
	} else if !c.ValidateIncentive() {
		return false
	} else if !c.ValidateStatus() {
		return false
	} else if !c.ValidateDueDate() {
		return false
	}
	return true
}

// ValidateCreate extends the Validate method to ensure we have
// certain restrictions when creating a chore.
func (c Chore) ValidateCreate() bool {
	if !c.Validate() {
		return false
	} else if c.DateCompleted > 0 {
		return false
	} else if c.Notes != nil {
		return false
	} else if len(c.VerifiedBy) > 0 {
		return false
	}

	return true
}

// ValidateTitle validates title is within bounds.
func (c Chore) ValidateTitle() bool {
	if len(c.Title) < 1 || len(c.Title) > 128 {
		return false
	}
	return true
}

// ValidateDescription validates description is within bounds.
func (c Chore) ValidateDescription() bool {
	if len(c.Description) > 4096 {
		return false
	}
	return true
}

// ValidateStatus validates status is a proper setting.
func (c Chore) ValidateStatus() bool {
	switch strings.ToLower(c.Status) {
	case "unassigned":
		return true
	case "assigned":
		return true
	case "review":
		return true
	case "completed":
		return true
	}
	return false
}

// ValidateIncentive validates incentive is within bounds.
func (c Chore) ValidateIncentive() bool {
	if len(c.Incentive) > 1024 {
		return false
	}
	return true
}

// ValidateDueDate validates duedate is a proper value.
// There is a small window in case there is a time skew.
func (c Chore) ValidateDueDate() bool {
	if c.DueDate == 0 { // If No Due Date Assigned.
		return true
	} else if c.DueDate < (time.Now().Unix() - 86400) { // 24hrs ago.
		return false
	} else if c.DueDate > (time.Now().Unix() + 31557600) { // One year..
		return false
	}
	return true
}
