# Flow

*Development.*  Tool intended for data analysis across distributed systems.

### Settings

Program requires a `config.json` in root directory:

```javascript
{
	"password" : "MyPassword!123",
	"port"     : 1433,
	"server"   : "myserver.database.windows.net",
	"user"     : "UserMe",
	"database" : "demoDB"
}

```

### Running

Build via `go build` and run the program.

When you see terminal prompt `~>`, proceed using SQL.

Special commands:

+ `run query.sql` executes code in `./sql/query.sql`
+ `exit` or `quit` will close the terminal, exit the program, and return you to the command line

### Priorities

1) Support multiple languages 

	+ SQL
	+ Go
	+ Python
	+ Scala

2) Integrate NSQ for more flexible / faster architecture

3) Improve UI / UX


<br>
<br>

<hr>
<small>
<strong>&copy; 2015 MIT License</strong>
</small>
# Info
