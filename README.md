# contact-list-app
Contact list application that includes a list of contacts and displays the details of the contacts. Make a search bar that helps to search for a contact from the list. The user can add, edit and delete a contact.

Project Demo link : https://drive.google.com/file/d/1qtMC7Qow5NkZFLZCCiRhAow8InDXCzJQ/view?usp=sharing

Tech Stack used :
Front-end : HTML, CSS, Javascript, Tailwindcss
backend : golang
database : postgresql
devops : Docker

Packages used :
"database/sql"
"fmt"
"html/template"
"log"
"net/http"
"os"
"github.com/gorilla/sessions"
"github.com/joho/godotenv"
_ "github.com/lib/pq"
"golang.org/x/crypto/bcrypt"


Database structure :
contactdb:
-----------------------------------------------------------------------------------------------------
    Column    |          Type          | Collation | Nullable |               Default
--------------+------------------------+-----------+----------+--------------------------------------
 id           | integer                |           | not null | nextval('contacts_id_seq'::regclass)
 username     | character varying(50)  |           | not null |
 name         | character varying(50)  |           |          |
 email        | character varying(50)  |           |          |
 phone_number | character varying(50)  |           |          |
 address      | character varying(100) |           |          |
 company_name | character varying(50)  |           |          |
Indexes:
    "contacts_pkey" PRIMARY KEY, btree (id)

userdb:
--------------------------------------------------------------------------------------------------
    Column    |          Type          | Collation | Nullable |              Default
--------------+------------------------+-----------+----------+-----------------------------------
 id           | integer                |           | not null | nextval('users_id_seq'::regclass)
 username     | character varying(50)  |           | not null | 
 phone_number | character varying(50)  |           | not null | 
 password     | character varying(255) |           | not null | 
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "users_phone_number_key" UNIQUE CONSTRAINT, btree (phone_number)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)

