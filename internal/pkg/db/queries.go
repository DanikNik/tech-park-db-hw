package db

var CreateUserQuery = "INSERT INTO forum.users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4);"
var GetUserQuery = "SELECT nickname, email, fullname, about FROM forum.users WHERE nickname = $1;"
var UpdateUserQuery = "UPDATE forum.users SET fullname=$1, email=$2, about=$3;"

var CreateForumQuery = "INSERT INTO forum.forum (slug, title, author) VALUES ($1, $2, $3);"
var GetForumQuery = "SELECT * FROM forum.forum WHERE slug = $1;"
var UpdateForumQuery string
