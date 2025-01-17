package db

const (
	addRoleQuery = `
		INSERT INTO role_assignments (user_id, role, resource_id) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (user_id, resource_id) 
		DO UPDATE SET role = $2`
	deleteRoleQuery = `DELETE FROM role_assignments WHERE user_id = $1 AND resource_id = $2`
	selectRoleQuery = `SELECT role FROM role_assignments WHERE user_id = $1 AND resource_id = $2`
)
