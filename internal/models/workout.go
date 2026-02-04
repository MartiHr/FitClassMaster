package models

import (
	"gorm.io/gorm"
)

// WorkoutPlan represents a template created by a Trainer.
type WorkoutPlan struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	TrainerID uint `json:"trainer_id"`
	Trainer   User `gorm:"foreignKey:TrainerID" json:"trainer"`

	WorkoutExercises []WorkoutExercise `gorm:"foreignKey:WorkoutPlanID" json:"workout_exercises"`
	DeletedAt        gorm.DeletedAt    `gorm:"index"`
}

// WorkoutExercise is the join table between a WorkoutPlan and an Exercise, including specific targets.
type WorkoutExercise struct {
	gorm.Model
	WorkoutPlanID uint `json:"workout_plan_id"`

	ExerciseID uint     `json:"exercise_id"`
	Exercise   Exercise `gorm:"foreignKey:ExerciseID" json:"exercise"`

	// Specific instructions for this plan.
	Sets     int    `json:"sets"`     // e.g., 3
	Reps     int    `json:"reps"`     // e.g., 12
	Duration int    `json:"duration"` // Duration in seconds (for exercises like planks or cardio)
	Notes    string `json:"notes"`    // Additional instructions (e.g., "Rest 60s between sets")
	Order    int    `json:"order"`    // Sequence order of the exercise within the plan
}

// SetList returns a slice of integers from 1 up to the number of Sets.
// This allows the HTML templates to easily loop over sets.
func (we WorkoutExercise) SetList() []int {
	var sets []int
	for i := 1; i <= we.Sets; i++ {
		sets = append(sets, i)
	}
	return sets
}
