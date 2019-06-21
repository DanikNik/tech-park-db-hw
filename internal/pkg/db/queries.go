package db

var CreateUserQuery = "INSERT INTO tp_forum.users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4);"

var SelectUsersWithNickOrEmail = "SELECT nickname, fullname, about, email FROM tp_forum.users WHERE nickname = $1 OR email = $2"
var GetUserQuery = "SELECT nickname, email, fullname, about FROM tp_forum.users WHERE nickname = $1;"
var UpdateUserQuery = "UPDATE tp_forum.users SET fullname=$1, email=$2, about=$3;"

var CreateForumQuery = "INSERT INTO tp_forum.forum (slug, title, author) VALUES ($1, $2, $3);"
var GetForumBySlugQuery = "SELECT * FROM tp_forum.forum WHERE slug = $1;"
var GetForumByIdQuery = "SELECT * FROM tp_forum.forum WHERE id = $1;"
var UpdateForumQuery string
