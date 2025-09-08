export type Shot = {
	shot_id: string
	brew_time: string // ISO string from backend
	machine_id: string
	user_id: string
	software_bundle: string
	coffee_type: string
	recipe_id: string
	grind_size_actual: number
	grind_size_target: number
	dose_grams: number
	dose_target_grams: number
	brew_time_seconds: number
	peak_pressure_bar: number
	last_status: string
}
