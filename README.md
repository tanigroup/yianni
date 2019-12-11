# yianni
SQL to CSV dumper utility

#### Setup:

- Ensure you are connected to the database.
- Create `config.json` with your database details.
- Save your SQL query to any text file.

#### Usage:
```bash
# ./yianni <param1:required> <param2:optional>
$ ./yianni ~/path/to/sql/file.txt
$ ./yianni ~/path/to/sql/file.txt filename.csv
```
