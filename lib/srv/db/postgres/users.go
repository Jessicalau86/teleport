/*
Copyright 2023 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package postgres

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gravitational/trace"
	"github.com/jackc/pgx/v4"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/srv/db/common"
)

// ActivateUser creates or enables the database user.
func (e *Engine) ActivateUser(ctx context.Context, sessionCtx *common.Session) error {
	if sessionCtx.Database.GetAdminUser().Name == "" {
		return trace.BadParameter("Teleport does not have admin user configured for this database")
	}

	conn, err := e.pgxConnect(ctx, sessionCtx.WithUser(sessionCtx.Database.GetAdminUser().Name))
	if err != nil {
		return trace.Wrap(err)
	}
	defer conn.Close(ctx)

	// We could call this once when the database is being initialized but
	// doing it here has a nice "self-healing" property in case the Teleport
	// bookkeeping group or stored procedures get deleted or changed offband.
	err = e.initAutoUsers(ctx, sessionCtx, conn)
	if err != nil {
		return trace.Wrap(err)
	}

	roles, err := prepareRoles(sessionCtx)
	if err != nil {
		return trace.Wrap(err)
	}

	e.Log.Infof("Activating PostgreSQL user %q with roles %v.", sessionCtx.DatabaseUser, roles)

	_, err = conn.Exec(ctx, activateQuery, sessionCtx.DatabaseUser, roles)
	if err != nil {
		e.Log.Debugf("Call teleport_activate_user failed: %v", err)
		return trace.Wrap(convertActivateError(sessionCtx, err))
	}
	return nil

}

// DeactivateUser disables the database user.
func (e *Engine) DeactivateUser(ctx context.Context, sessionCtx *common.Session) error {
	if sessionCtx.Database.GetAdminUser().Name == "" {
		return trace.BadParameter("Teleport does not have admin user configured for this database")
	}

	conn, err := e.pgxConnect(ctx, sessionCtx.WithUser(sessionCtx.Database.GetAdminUser().Name))
	if err != nil {
		return trace.Wrap(err)
	}
	defer conn.Close(ctx)

	e.Log.Infof("Deactivating PostgreSQL user %q.", sessionCtx.DatabaseUser)

	_, err = conn.Exec(ctx, deactivateQuery, sessionCtx.DatabaseUser)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// DeleteUser deletes the database user.
func (e *Engine) DeleteUser(ctx context.Context, sessionCtx *common.Session) error {
	// TODO support DeleteUser for Redshift
	if sessionCtx.Database.IsRedshift() {
		e.Log.Debug("DeleteUser is not supported for Redshift yet, it was disabled instead.")
		return trace.Wrap(e.DeactivateUser(ctx, sessionCtx))
	}

	if sessionCtx.Database.GetAdminUser().Name == "" {
		return trace.BadParameter("Teleport does not have admin user configured for this database")
	}

	conn, err := e.pgxConnect(ctx, sessionCtx.WithUser(sessionCtx.Database.GetAdminUser().Name))
	if err != nil {
		return trace.Wrap(err)
	}
	defer conn.Close(ctx)

	e.Log.Infof("Deleting PostgreSQL user %q.", sessionCtx.DatabaseUser)

	var state string
	err = conn.QueryRow(ctx, deleteQuery, sessionCtx.DatabaseUser).Scan(&state)
	if err != nil {
		return trace.Wrap(err)
	}

	switch state {
	case common.SQLStateUserDropped:
		e.Log.Debug("User %q deleted successfully.", sessionCtx.DatabaseUser)
	case common.SQLStateUserDeactivated:
		e.Log.Infof("Unable to delete user %q, it was disabled instead.", sessionCtx.DatabaseUser)
	default:
		e.Log.Warnf("Unable to determine user %q deletion state.", sessionCtx.DatabaseUser)
	}

	return nil
}

// initAutoUsers installs procedures for activating and deactivating users and
// creates the bookkeeping role for auto-provisioned users.
func (e *Engine) initAutoUsers(ctx context.Context, sessionCtx *common.Session, conn *pgx.Conn) error {
	// Create a role/group which all auto-created users will be a part of.
	_, err := conn.Exec(ctx, fmt.Sprintf("create role %q", teleportAutoUserRole))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return trace.Wrap(err)
		}
		e.Log.Debugf("PostgreSQL role %q already exists.", teleportAutoUserRole)
	} else {
		e.Log.Debugf("Created PostgreSQL role %q.", teleportAutoUserRole)
	}

	// Install stored procedures for creating and disabling database users.
	for name, sql := range pickProcedures(sessionCtx) {
		_, err := conn.Exec(ctx, sql)
		if err != nil {
			return trace.Wrap(err)
		}
		e.Log.Debugf("Installed PostgreSQL stored procedure %q.", name)
	}
	return nil
}

// pgxConnect connects to the database using pgx driver which is higher-level
// than pgconn and is easier to use for executing queries.
func (e *Engine) pgxConnect(ctx context.Context, sessionCtx *common.Session) (*pgx.Conn, error) {
	config, err := e.getConnectConfig(ctx, sessionCtx)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	pgxConf, err := pgx.ParseConfig("")
	if err != nil {
		return nil, trace.Wrap(err)
	}
	pgxConf.Config = *config
	return pgx.ConnectConfig(ctx, pgxConf)
}

func prepareRoles(sessionCtx *common.Session) (any, error) {
	switch sessionCtx.Database.GetType() {
	case types.DatabaseTypeRDS:
		return append(sessionCtx.DatabaseRoles, "rds_iam"), nil

	case types.DatabaseTypeRedshift:
		// Redshift does not support array. Encode roles in JSON (type text).
		rolesJSON, err := json.Marshal(sessionCtx.DatabaseRoles)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		return string(rolesJSON), nil

	default:
		return sessionCtx.DatabaseRoles, nil
	}
}

func convertActivateError(sessionCtx *common.Session, err error) error {
	switch {
	case strings.Contains(err.Error(), "already exists"):
		return trace.AlreadyExists("user %q already exists in this PostgreSQL database and is not managed by Teleport", sessionCtx.DatabaseUser)

	case strings.Contains(err.Error(), "TP002: User has active connections and roles have changed"):
		return trace.CompareFailed("roles for user %q has changed. Please quit all active connections and try again.", sessionCtx.DatabaseUser)

	default:
		return trace.Wrap(err)
	}
}

func pickProcedures(sessionCtx *common.Session) map[string]string {
	if sessionCtx.Database.IsRedshift() {
		return redshiftProcs
	}
	return procs
}

const (
	// activateProcName is the name of the stored procedure Teleport will use
	// to automatically provision/activate database users.
	activateProcName = "teleport_activate_user"
	// deactivateProcName is the name of the stored procedure Teleport will use
	// to automatically deactivate database users after session ends.
	deactivateProcName = "teleport_deactivate_user"
	// deleteProcName is the name of the stored procedure Teleport will use to
	// automatically delete database users after session ends.
	deleteProcName = "teleport_delete_user"

	// teleportAutoUserRole is the name of a PostgreSQL role that all Teleport
	// managed users will be a part of.
	teleportAutoUserRole = "teleport-auto-user"
)

var (
	//go:embed sql/activate-user.sql
	activateProc string
	// activateQuery is the query for calling user activation procedure.
	activateQuery = fmt.Sprintf(`call %v($1, $2)`, activateProcName)

	//go:embed sql/deactivate-user.sql
	deactivateProc string
	// deactivateQuery is the query for calling user deactivation procedure.
	deactivateQuery = fmt.Sprintf(`call %v($1)`, deactivateProcName)

	//go:embed sql/delete-user.sql
	deleteProc string
	// deleteQuery is the query for calling user deletion procedure.
	deleteQuery = fmt.Sprintf(`call %v($1)`, deleteProcName)

	//go:embed sql/redshift-activate-user.sql
	redshiftActivateProc string
	//go:embed sql/redshift-deactivate-user.sql
	redshiftDeactivateProc string

	procs = map[string]string{
		activateProcName:   activateProc,
		deactivateProcName: deactivateProc,
		deleteProcName:     deleteProc,
	}
	redshiftProcs = map[string]string{
		activateProcName:   redshiftActivateProc,
		deactivateProcName: redshiftDeactivateProc,
	}
)
