#!/usr/bin/zsh

(
  source .env
  trap 'cd $HOME/dev/sl/server' ERR
  cd $HOME/dev/sl/server/database/sql/schemas && goose sqlite3 $DB_URL_DEV down-to 0 && goose sqlite3 $DB_URL_DEV up && cd $HOME/dev/sl/server
)
