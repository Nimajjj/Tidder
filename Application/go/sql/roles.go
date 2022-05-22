package mySQL

import (
	"strconv"
	"strings"

	Util "github.com/Nimajjj/Tidder/go/utility"
)

func (sqlServ SqlServer) CreateRole(idSubject int, data string) {
	roles := strings.Split(data, ";")

	for _, role := range roles {
		access := strings.Split(role, ",")
		if len(access) != 8 {
			continue
		}
		roleName := access[1]

		query := "INSERT INTO subject_access (create_post, pin_post, manage_post, ban_user, manage_role, manage_subtidder) VALUES ("
		for i, acc := range access {
			if i > 1 {
				query += acc
				if i != len(access)-1 {
					query += ", "
				}
			}
		}
		query += ")"
		sqlServ.executeQuery(query) // create new access
		var id int
		q := "SELECT id_subject_access FROM tidder.subject_access ORDER BY id_subject_access DESC LIMIT 1"
		if err := sqlServ.db.QueryRow(q).Scan(&id); err != nil {
			Util.Error(err)
			return
		}

		query = "INSERT INTO subject_roles (name, id_subject, id_subject_access) VALUES ("
		query += "'" + roleName + "', " + strconv.Itoa(idSubject) + ", " + strconv.Itoa(id) + ")"
		sqlServ.executeQuery(query)
	}
}

func (sqlServ SqlServer) UpdateRole(idSubject int, data string) {
	roles := strings.Split(data, ";")

	for _, role := range roles {
		access := strings.Split(role, ",")
		if len(access) != 8 {
			continue
		}
		roleId := access[0]
		roleName := access[1]

		if roleId == "-1" {
			query := "UPDATE subjects SET can_create_post=" + access[2] + " WHERE id_subject=" + strconv.Itoa(idSubject)
			sqlServ.executeQuery(query)
			continue
		}

		query := "UPDATE subject_roles SET name='" + roleName + "'" + " WHERE id_subject_role=" + roleId + " AND id_subject=" + strconv.Itoa(idSubject)
		sqlServ.executeQuery(query)

		var accessId int
		query = "SELECT id_subject_access FROM subject_roles WHERE id_subject_role=" + roleId
		if err := sqlServ.db.QueryRow(query).Scan(&accessId); err != nil {
			Util.Error(err)
			return
		}

		query = "UPDATE subject_access SET "
		for i, acc := range access {
			if i > 1 {
				switch i {
				case 2:
					query += " create_post="
				case 3:
					query += " pin_post="
				case 4:
					query += " manage_post="
				case 5:
					query += " ban_user="
				case 6:
					query += " manage_role="
				case 7:
					query += " manage_subtidder="
				}
				query += acc
				if i != len(access)-1 {
					query += ", "
				}
			}
		}
		query += " WHERE id_subject_access=" + strconv.Itoa(accessId)
		sqlServ.executeQuery(query)
	}
}

func (sqlServ SqlServer) DeleteRole(idSubject int, data string) {
	// CHECK IF ROLE IS USED !!!
	roles := strings.Split(data, ";")

	for _, roleId := range roles {
		var accessId int
		query := "SELECT id_subject_access FROM subject_roles WHERE id_subject_role=" + roleId
		if err := sqlServ.db.QueryRow(query).Scan(&accessId); err != nil {
			Util.Error(err)
			return
		}

		query = "DELETE FROM subject_roles WHERE id_subject_role=" + roleId
		sqlServ.executeQuery(query)

		query = "DELETE FROM subject_access WHERE id_subject_access=" + strconv.Itoa(accessId)
		sqlServ.executeQuery(query)
	}
}

func (sqlServ SqlServer) GenerateRoleAccess(idSubject int) []RoleAccess {
	res := []RoleAccess{}
	roles := []SubjectRoles{}

	query := "SELECT id_subject_role, name, id_subject_access FROM subject_roles WHERE id_subject=" + strconv.Itoa(idSubject)
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}

	for rows.Next() {
		var idRole int
		var name string
		var idAccess int
		if err2 := rows.Scan(
			&idRole,
			&name,
			&idAccess,
		); err2 != nil {
			Util.Error(err2)
			return res
		}
		roles = append(roles, SubjectRoles{idRole, name, idSubject, idAccess})
	}

	for _, role := range roles {
		query = "SELECT * FROM tidder.subject_access WHERE id_subject_access=" + strconv.Itoa(role.IdSubjectAccess)
		var id int
		var createPost int
		var pin int
		var managePost int
		var banUser int
		var manageRole int
		var giveRole int
		if err := sqlServ.db.QueryRow(query).Scan(&id, &createPost, &pin, &managePost, &banUser, &manageRole, &giveRole); err != nil {
			Util.Error(err)
			return res
		}
		res = append(res, RoleAccess{role, SubjectAccess{id, createPost, pin, managePost, banUser, manageRole, giveRole}})
	}

	return res
}

func (sqlServ SqlServer) ChangeRoleAtribution(idSubject int, data string) {
	atr := strings.Split(data, ";")

	for _, atrib := range atr {
		split := strings.Split(atrib, ",")
		if len(split) != 2 {
			continue
		}
		userId := split[0]
		roleId := split[1]

		if roleId == "-1" {
			query := "DELETE FROM has_subject_role WHERE id_account=" + userId + " AND id_subject=" + strconv.Itoa(idSubject)
			sqlServ.executeQuery(query)
			continue
		}

		if sqlServ.RowExists("has_subject_role", "id_account="+userId+" AND id_subject="+strconv.Itoa(idSubject)) {
			query := "UPDATE has_subject_role SET id_subject_role=" + roleId + " WHERE id_account=" + userId + " AND id_subject=" + strconv.Itoa(idSubject)
			sqlServ.executeQuery(query)
			continue
		}

		query := "INSERT INTO has_subject_role (id_account, id_subject, id_subject_role) VALUES (" + userId + ", " + strconv.Itoa(idSubject) + ", " + roleId + ")"
		sqlServ.executeQuery(query)
	}
}

func (sqlServ SqlServer) RowExists(table string, conditions string) bool {
	query := "SELECT * FROM " + table + " WHERE " + conditions
	rows, err := sqlServ.db.Query(query)
	if err != nil {
		Util.Error(err)
	}
	if rows.Next() {
		return true
	}
	return false
}

func (sqlServ SqlServer) GenerateUserRoleAccess(idSubject int, idUser int) RoleAccess {
	res := RoleAccess{}
	role := SubjectRoles{}
	access := SubjectAccess{}

	if !sqlServ.RowExists("has_subject_role", "id_account="+strconv.Itoa(idUser)+" AND id_subject="+strconv.Itoa(idSubject)) {
		role.Id = -1
		res.Role = role
		return res
	}

	query := "SELECT id_subject_role FROM has_subject_role WHERE id_account=" + strconv.Itoa(idUser) + " AND id_subject=" + strconv.Itoa(idSubject)
	var idRole int
	if err := sqlServ.db.QueryRow(query).Scan(&idRole); err != nil {
		Util.Error(err)
	}

	query = "SELECT * FROM subject_roles WHERE id_subject_role=" + strconv.Itoa(idRole)
	if err := sqlServ.db.QueryRow(query).Scan(&role.Id, &role.Name, &role.IdSubject, &role.IdSubjectAccess); err != nil {
		Util.Error(err)
	}

	query = "SELECT * FROM subject_access WHERE id_subject_access=" + strconv.Itoa(role.IdSubjectAccess)
	if err := sqlServ.db.QueryRow(query).Scan(&access.Id, &access.CreatePost, &access.Pin, &access.ManagePost, &access.BanUser, &access.ManageRole, &access.ManageSub); err != nil {
		Util.Error(err)
	}

	res.Role = role
	res.Access = access

	return res
}
