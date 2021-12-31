# My solution
This is a simple solution written in GO for deleting users from the database and writing a copy to a given output file. The output file is in the csv format.

The program is seperated into `Logic` in logic all functionality not directly on a model is stored. In `Models` structures and related functionality is stored. 

## Usage
Running `go run main.go` will run the program with default settings deleting the first user in the database.  

To get a list of program arguments `-h` or `--help` can be used.  
To set the output file path use `-path=<filename>`  
To set the database connection string use `-db=<connectionString>`  
To increase the limit of query use `-limit=<count>` zero is equivalent to no limit.
To create a query use `-query=<string>`.
Query strings can contain arguments, to add an argument to a query use a `?` instead of the value.  
To parse arguments to query use `-args=<arguments>`. Arguments are seperated by `;` and if an argument contains multiple values separate them using `,`

## Security
Arbitrary queries can be passed allowing a user to modify the database in an unintended way. 