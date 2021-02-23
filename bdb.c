// Copyright 2013 Matt Butcher
// Copyright 2021 Bernd Fix >Y<
// MIT License
// This file contains a number of wrapper functions that make it
// possible for Go to call the method-style functions on BerkeleyDB
// structs.

#include <db.h>

int go_db_open(DB *dbp, DB_TXN *txnid, char *filename, char *dbname, DBTYPE type, u_int32_t flags, int mode) {
        return dbp->open(dbp, txnid, filename, dbname, type, flags, mode);
}

int go_db_close(DB *dbp, u_int32_t flags) {
        if (dbp == NULL)
                return 0;
        return dbp->close(dbp, flags);
}

int go_db_get_open_flags(DB *dbp, u_int32_t *open_flags) {
        return dbp->get_open_flags(dbp, open_flags);
}

int go_db_remove(DB *dbp, char *filename) {
        return dbp->remove(dbp, filename, NULL, 0);
}

int go_db_rename(DB *dbp, char *oldname, char *newname) {
        return dbp->rename(dbp, oldname, NULL, newname, 0);
}

int go_env_open(DB_ENV *env, char *dirname, u_int32_t flags, u_int32_t mode) {
	return env->open(env, dirname, flags, mode);
}

int go_env_close(DB_ENV *env, u_int32_t flags) {
	return env->close(env, flags);
}

int go_db_put(DB *dbp, DBT *key, DBT *value, u_int32_t flags) {
	return dbp->put(dbp, NULL, key, value, flags);
}

int go_db_get(DB *dbp, DBT *key, DBT *value) {
	return dbp->get(dbp, NULL, key, value, 0);
}

int go_db_cursor(DB *dbp, DBC **dbcp) {
        return dbp->cursor(dbp, NULL, dbcp, 0);
}

int go_cursor_get(DBC *dbcp, DBT *key, DBT *value, int mode) {
        return dbcp->c_get(dbcp, key, value, mode);
}

int go_db_del(DB *dbp, DBT *key) {
	return dbp->del(dbp, NULL, key, 0);
}

