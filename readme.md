# Track-space

### Introduction

Track-space is a monolithic web application that provides a workspace for professionals from different fields. It includes features for collaboration, communication, and task management, enabling teams to stay organized and on track with their projects.

### Prerequisites
Before running the Track-space application, ensure that your machine has Go Compiler version 1.18.3 or higher installed. This ensures that the application runs smoothly without any errors.

### Installation
* To Clone the Track-space repository from GitHub onto your local machine.
* Open a terminal and navigate to the directory where the application is stored.
* Run the command go mod tidy in the terminal to update all the third-party packages and libraries required to run the application. Make sure that you are connected to the internet.
* Run the command go mod tidy again in the terminal to synchronize all the dependencies.
* Google mail server is integrated using goroutines
* Run the script run.sh in the terminal by typing ./run.sh. This script builds and runs the application simultaneously on your machine (PC).
* Open your favorite web browser and visit the URL http://localhost:8080 to access the Track-space application on a local server.

### Conclusion

By following the above steps, you should be able to successfully run Track-space on your local machine. The application offers a wide range of features for professionals, such as collaboration, communication, and task management, that can help streamline and optimize project workflows. Check out the application and see how it can benefit your team!

### Notice

To prevent future errors when running the application, it's necessary to set up a local mail server. If you're using Google mail server, you don't need to use mailhog. Instead, you should configure the application to use your Google mail server settings
