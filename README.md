# ArchaicReverie

Currently **"ArchaicReverie"** is a web application with the concept of game mechanics. The main idea is to implement a system where each user can create their own character and then choose one of them to solve tasks generated by the conditions of the game event (which include the characteristics of the location, time of day, weather). 
The app is designed in such a way that different types of tasks can be added.
For example, there is a task to jump an obstacle: the player sets the position of the jump itself (amplitude of arm movement, run-up, etc.), then the code calculates by a special formula all the parameters followed by a success check of the action (which is the result of calculating by an own formula the mentioned characteristics of the game event) and displays the result. 

___
### Technical Information

The API is developed using the Gin framework.

Authorization is done via a token and a cookie. 

Repository: SQL PostgreSQL DBMS

Among the libraries involved there are: 
* jmoiron/sqlx (for PostgreSQL SQL);
* gin-gonic/gin (API);
* golang-jwt/jwt (token);
* /mroth/weightedrand (as part of a random generation system) 
* 
and so on. 
