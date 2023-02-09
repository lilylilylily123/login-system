# login-system
super WIP login system demo, just to prove that i can do it

## Build Instructions
Clone this repo, and initiliaze Go:
`git clone https://github.com/lilylilylily123/login-system`
<br>
<br>
Create a new database named <code>db.users</code>, and create a table with the following fields: <br>
| userdetails  | Values |
| ------------- | ------------- |
| username  | *string*  |
| email  | *string*  |
| password | *string* |

Save your table, and then bam! You're all set.

Do <code>go build main.go</code>, and then run <code>./main</code>

Visit <a href="localhost:1000"><code href="localhost:1000">localhost:1000</code></a> in your browser!

## Route Index 


`/` -- Signups Page

`/login/` -- Login Page (will redirect to signup page if user/pass doesn't exist

`/homepage/` -- Home Page (only accessible if logged in)

There are a couple of other routes, but those are all used for internal things such as setting cookies/redirection.

## Extra

**Note:**

Since this is written with bcrypt, all passwords are automatically salted when you convert them to hash. Obviously it's still not extremely secure, but it's defenitely better than nothing.
