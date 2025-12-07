# FT Core
FIRETool Core is the main service for the FIRETool application. It has all of the financial business logic for the application, including bank financial data upload jobs.

## Setup / Installation
### Prerequisistes
* Docker
* Goose (for database migrations)
* PSQL (for inspecting the database)

### Environment variables
1. Create a .env file
2. Add the GOOSE environment variables (GOOSE_DRIVER, GOOSE_DBSTRING, GOOSE_MIGRATION_DIR, GOOSE_TABLE)
3. GOOSE_MIGRATION_DIR should be set to db/migrations
4. GOOSE_TABLE should be set to `goose_migrations`

### Database Setup
* Create a new database (ftcore_dev, ftcore_prod, etc)
* Run `goose up` to create all the tables