# Track-space

Track-space is a project management monolithic web application which serve as a workspace for professionals ranging from technical to non-technical field such as software developers,data analyst,writers as well as set their future schedule plans and a secured interactive space to communicate with others.

### **Some of it built-in core services and activities features include :**

    • All in Editor type for users to carry out a specific type of project.
    • A simple to-do schedule to set date and time for a future plans.
    • Real time communication with interact with other available user.
    • Quick notification on every secured activities of the users.
    • Data visualization to show  statistic report based on usage.

### **Technology used are :**

#### Server-Side of the application

    • Go Complier V.1.18.3 +
    • Gin-Gonic v1.8.1 ( Golang web framework)
    • MongoDB Go Driver v1.9.1
    • MongoDB (MongoDB Atlas Cloud Database storage).
    • jwt-go for authentication.
    • Gin-contrib sessions for cookies.
    • go-simple-mail package.
    • Mailhog as local mailing server.
    • Gorilla mux web socket.

#### Client-Side of the application

    • javaScript.
    • Tiny Mce API editor
    • HTML5 
    • CSS
    • Bootstrap
    • Reconnecting web socket module
    • D3.js
    • Simple-dataTable
    • Notie.js

### **Steps to start Track-space application**

    i. Make sure you have the Go Complier 1.18.3+ installed on your machine

    ii. Go to the terminal make sure you are in the right directory of the application.

    iii. Type go mod tidy and press Enter to update all the third party packages and libraries (Be connected to the internet)

    iv. Type go mod tidy and press Enter to sync all the dependencies

    v. Type mailhog in the terminal and press enter to start the local mail server in the terminal to avoid future error which running the application.

    vi. Type ./run.sh and press Enter to build and run the application at the same time on you machine(PC)

    vii. Lastly, go to you favourite web browser and search for localhost:8080 then press Enter. You have successfully start the application on a local server on your machine. Goodluck & Enjoy.

