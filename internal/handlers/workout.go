package handlers

import (
	"FitClassMaster/internal/auth"
	"FitClassMaster/internal/services"
	"FitClassMaster/internal/templates"
	"net/http"
	"strconv"
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

// AddExerciseRow returns ONE HTML row for the exercise list (HTMX endpoint)
func (h *WorkoutHandler) AddExerciseRow(w http.ResponseWriter, r *http.Request) {
	exercises, _ := h.ExerciseService.GetAll()

	data := map[string]any{
		"Exercises": exercises,
	}
	// Render ONLY the "row" fragment, not the whole page layout
	templates.SmartRender(w, r, "workout_create", "exercise-row", data)
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
	canManage := (role == "trainer" || role == "admin")

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
	canManage := (role == "trainer" || role == "admin")

	data := map[string]any{
		"Title":     plan.Name,
		"Plan":      plan,
		"CanManage": canManage,
	}

	templates.SmartRender(w, r, "workout_plan_details", "", data)
}
