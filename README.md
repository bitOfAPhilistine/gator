# Gator
## Blog aggregator for Boot.dev
---
This is a blog aggregator made for boot.dev ([Build a Blog Aggregator](https://www.boot.dev/courses/build-blog-aggregator-golang)). It features a login system, blog following, and periodic blog aggregating! Also includes a help command that can be run with 'gator help'

## Installation
---
For wsl or linux, a pre-made bash script has been made which can automatically install Gator and set up its database. To run it, simply enter the following code in your terminal:
'sudo bash (WHEREVER YOU PUT GATOR)/gator/gator_wsl_linux_installer.sh'

For other operating systems, the installation has to be manual:
1. Install Postgres (v15 or later) [Download page](https://www.postgresql.org/download/)
2. Start the postgres server if not started by the installer
3. Create a database in postgres called gator and set its password, make sure to note down the connection url
4. Install Golang [Download page](https://go.dev/doc/install)
5. Install Goose with Go 'go install github.com/pressly/goose/v3/cmd/goose@latest'
6. Cd to the sql/schema directory in the gator directory
7. Run the following command: 'goose postgres "(THE CONNECTION STRING FROM EARLIER)" up'
8. Move the gator (no extension) file to your path or run 'go install' from the gator directory
9. Run 'gator init "(THE CONNECTION STRING FROM EARLIER)"'
10. Should be done!