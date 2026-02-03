package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/models"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	Service         *services.WorkoutService
	ExerciseService *services.ExerciseService
}

func NewWorkoutHandler(ws *services.WorkoutService, es *services.ExerciseService) *WorkoutHandler {
	return &WorkoutHandler{Service: ws, ExerciseService: es}
}

// CreatePage renders the form to start a new plan
func (h *WorkoutHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	// We need the list of exercises for the dropdown menus
	exercises, _ := h.ExerciseService.GetAll()

	data := map[string]any{
		"Title":     "Create Workout Plan",
		"Exercises": exercises,
	}
	templates.SmartRender(w, r, "workout_create", "", data)
}

// CreatePost handles the form submission
func (h *WorkoutHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get Basic Info
	name := r.FormValue("name")
	desc := r.FormValue("description")
	trainerID, _ := auth.GetUserIDFromSession(r)

	// Parse the Arrays (The Dynamic Rows)
	exerciseIDs := r.Form["exercise_id[]"]
	sets := r.Form["sets[]"]
	reps := r.Form["reps[]"]
	notes := r.Form["notes[]"]

	var inputs []services.ExerciseInput

	// Loop through the arrays and build the input structs
	for i, idStr := range exerciseIDs {
		eid, _ := strconv.ParseUint(idStr, 10, 32)
		s, _ := strconv.Atoi(sets[i])
		r_val, _ := strconv.Atoi(reps[i])

		inputs = append(inputs, services.ExerciseInput{
			ExerciseID: uint(eid),
			Sets:       s,
			Reps:       r_val,
			Notes:      notes[i],
		})
	}

	// Call Service
	err = h.Service.CreatePlan(name, desc, trainerID, inputs)
	if err != nil {
		http.Error(w, "Failed to create plan", http.StatusInternalServerError)
		return
	}

	// Redirect to the list of plans (which we will build next)
	http.Redirect(w, r, "/workout-plans", http.StatusSeeOther)
}

// List handles GET /workout-plans
func (h *WorkoutHandler) List(w http.ResponseWriter, r *http.Request) {
	plans, err := h.Service.ListAll()
	if err != nil {
		http.Error(w, "Failed to fetch plans", http.StatusInternalServerError)
		return
	}

	// Check Role for UI elements (Show "Create" button?)
	role, _ := auth.GetUserRoleFromSession(r)
	canManage := (role == models.RoleTrainer || role == models.RoleAdmin)

	data := map[string]any{
		"Title":     "Workout Plans",
		"Plans":     plans,
		"CanManage": canManage,
	}

	templates.SmartRender(w, r, "workout_plans", "", data)
}

// DetailsPage handles GET /workout-plans/{id}
func (h *WorkoutHandler) DetailsPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch Deep Details (Plan + Exercises + Specifics)
	plan, err := h.Service.GetFullDetails(uint(id))
	if err != nil {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	}

	// Determine Permissions (For "Edit/Delete" buttons)
	role, _ := auth.GetUserRoleFromSession(r)
	canManage := (role == models.RoleTrainer || role == models.RoleAdmin)

	data := map[string]any{
		"Title":     plan.Name,
		"Plan":      plan,
		"CanManage": canManage,
	}

	templates.SmartRender(w, r, "workout_plan_details", "", data)
}

// EditPage renders the form pre-filled with existing data
func (h *WorkoutHandler) EditPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	// Get the plan with details
	plan, err := h.Service.GetFullDetails(uint(id))
	if err != nil {
		http.Error(w, "Plan not found", http.StatusNotFound)
		return
	}

	// Get all exercises for dropdowns
	// Usage: h.ExerciseService.GetAll (Matches your ExerciseService)
	allExercises, _ := h.ExerciseService.GetAll()

	// Prepare row data for the template loop
	var rowData []map[string]any
	for _, we := range plan.WorkoutExercises {
		rowData = append(rowData, map[string]any{
			"Order":      we.Order,
			"SelectedID": we.ExerciseID, // Pre-selects the dropdown
			"Exercises":  allExercises,  // Passes the list to every row
			"Sets":       we.Sets,
			"Reps":       we.Reps,
			"Notes":      we.Notes,
		})
	}

	data := map[string]any{
		"Title":     "Edit Workout Plan",
		"IsEdit":    true, // Toggles the form action to /edit
		"Plan":      plan,
		"Exercises": allExercises, // For the "Add Row" button
		"Rows":      rowData,      // For the existing rows loop
	}

	templates.SmartRender(w, r, "workout_create", "", data)
}

// UpdatePost handles the form submission when editing a plan
func (h *WorkoutHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Capture Basic Info
	name := r.FormValue("name")
	description := r.FormValue("description")

	// Capture Arrays (The exercise rows)
	exIDs := r.Form["exercise_id[]"]
	sets := r.Form["sets[]"]
	reps := r.Form["reps[]"]
	rowNotes := r.Form["notes[]"]

	// Call Service to update
	err = h.Service.UpdatePlan(uint(id), name, description, exIDs, sets, reps, rowNotes)
	if err != nil {
		http.Error(w, "Failed to update plan", http.StatusInternalServerError)
		return
	}

	// Redirect back to list
	http.Redirect(w, r, "/workout-plans", http.StatusSeeOther)
}

// AddExerciseRow handles HTMX requests
func (h *WorkoutHandler) AddExerciseRow(w http.ResponseWriter, r *http.Request) {
	exercises, _ := h.ExerciseService.GetAll()

	data := map[string]any{
		"Exercises": exercises,
		"Sets":      3,
		"Reps":      10,
	}

	templates.SmartRender(w, r, "workout_create", "exercise-row", data)
}

// DeletePost handles the form submission to delete a plan
func (h *WorkoutHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	// Call the service
	if err := h.Service.DeletePlan(uint(id)); err != nil {
		http.Error(w, "Failed to delete plan", http.StatusInternalServerError)
		return
	}

	// Redirect to the list because this page no longer exists
	http.Redirect(w, r, "/workout-plans", http.StatusSeeOther)
}
