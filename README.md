# HERMES
hermes is stand alone mongodb migrator, build with go, *HERMES* provide simple feature like automate make schema for your mongodb collection,
you just need to create file using makefile.

```shell
make seeder path=./path name=fileName
```

## COMMAND
1. command [up / down] [DEFAULT: up]
   1. up, migrate all seeder
   2. down, delete seeder
2. path, path destination for seeder file [DEFAULT: ./seed]
3. collection [all / collectionName], name for collection you want to delete [DEFAULT: all]
   1. all, will delete all of your collections
   2. collectionName, example: post,user this is will delete post and user collection
4. dbname, name of database you will used [DEFAULT: test]
5. host, mongodb host [DEFAULT: localhost]
6. port, mongodb port [DEFAULT: 27017]
7. username, mongodb username (OPTIONAL)
8. password, mongodb password (OPTIONAL)

### UP
```shell
go run main.go --command=up --path=./seederPath // migrate all schema inside seederPath
```

### DOWN
```shell
go run main.go --command=down --collection=all // delete all collection
go run main.go --command=down --collection=post,user //delete post and user collection
```
