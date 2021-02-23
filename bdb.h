// Copyright 2013 Matt Butcher
// Copyright 2021 Bernd Fix >Y<
// MIT License

extern int go_db_open(DB *, DB_TXN *, char *, char *, DBTYPE, u_int32_t, int);
extern int go_db_close(DB *, u_int32_t);
extern int go_db_get_open_flags(DB *, u_int32_t *);
extern int go_db_remove(DB *, char *);
extern int go_db_rename(DB *, char *, char *);
extern int go_env_open(DB_ENV *, char *, u_int32_t, u_int32_t);
extern int go_env_close(DB_ENV *, u_int32_t);

int go_db_put(DB *, DBT *, DBT *, u_int32_t);
int go_db_get(DB *, DBT *, DBT *);
int go_db_del(DB *, DBT *);
int go_db_cursor(DB *, DBC **);
int go_cursor_get(DBC *, DBT *, DBT *, int);
