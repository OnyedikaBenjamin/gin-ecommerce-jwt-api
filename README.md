Ecommerce API

This is an Ecommerce API built with Go, GORM, and MySQL. It provides features for managing users, products, orders, and order items.

Prerequisites

Before you begin, ensure you have the following installed on your system:

Go (version 1.18 or later)

MySQL

Git

Setup Instructions

Step 1: Clone the Repository

git clone <repository-url>
cd ecommerce-api

Step 2: Configure Environment Variables

Create a .env file in the root of the project and add the following content:

DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ecommerce
DB_CHARSET=utf8mb4
DB_PARSE_TIME=True
DB_LOC=Local
SERVER_PORT=8080

Replace yourpassword with your MySQL root password or another user's password.

Step 3: Install Dependencies

Use go mod to download and install project dependencies:

go mod tidy

Step 4: Start the MySQL Database

Ensure your MySQL server is running and create the database:

CREATE DATABASE ecommerce;

Step 5: Run the Application

Start the application using the following command:

go run main.go

The server will start on the port specified in the .env file (default: 8080).

Step 6: Test the API

You can test the API using tools like Postman or cURL. Example base URL:

http://localhost:8080

Project Structure

models/: Contains the database models (User, Product, Order, OrderItem).

routes/: Contains the API routes.

main.go: The main entry point of the application.

Notes

Ensure the .env file is not committed to version control by adding it to .gitignore.

The database tables will be auto-migrated when the application starts.

License

This project is built by Benjamin Onyedika Udegbunam 25th of December 2024 and licensed under the MIT License.