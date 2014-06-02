package mysql

var fieldTypes = map[string]string{}

var mysqlStmt = map[string]string{
    "insertIgnore": "INSERT IGNORE INTO `%s` (`%s`) VALUES (%s)",
}
