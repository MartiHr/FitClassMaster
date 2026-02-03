package models

import (
	"gorm.io/gorm"
)

// WorkoutPlan represents a template created by a Trainer
type WorkoutPlan struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	TrainerID uint `json:"trainer_id"`
	Trainer   User `gorm:"foreignKey:TrainerID" json:"trainer"`

	WorkoutExercises []WorkoutExercise `gorm:"foreignKey:WorkoutPlanID" json:"workout_exercises"`
}

// WorkoutExercise is the link between a Plan and an Exercise
type WorkoutExercise struct {
	gorm.Model
	WorkoutPlanID uint `json:"workout_plan_id"`

	ExerciseID uint     `json:"exercise_id"`
	Exercise   Exercise `gorm:"foreignKey:ExerciseID" json:"exercise"`

	// Specific instructions for this plan
	Sets     int    `json:"sets"`     // e.g., 3
	Reps     int    `json:"reps"`     // e.g., 12
	Duration int    `json:"duration"` // In seconds (for planks/cardio)
	Notes    string `json:"notes"`    // e.g., "Rest 60s between sets"
	Order    int    `json:"order"`    // To keep exercises in sequence (1, 2, 3...)
}

// SetList returns a slice of integers
// This allows the template to loop over sets without a custom 'seq' function.
func (we WorkoutExercise) SetList() []int {
	var sets []int
	for i := 1; i <= we.Sets; i++ {
		sets = append(sets, i)
	}
	return sets
}
