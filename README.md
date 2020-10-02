# trace
Trace is a framework for tracking contacts at locations. It will host a frontend
where administrators may setup a kiosk that allows students to sign in and
sign out of a location. This data can be used to generate a contact report,
showing who has been in contact with who and for how long.

## Setting Up
### With [Docker](https://www.docker.com/)
If you have docker installed, you can compile and run the server
and database with `docker-compose`

To run, cd into the repo directory and type in 

```docker-compose up -d```

The server should be running on port 80. To enable HTTP authentication, use
the `USERNAME` and `PASSWORD` environment variables:

```bash
USERNAME=admin PASSWORD=password docker-compose up -d
```

The database data is stored in a [Docker volume](https://docs.docker.com/storage/volumes/).

## Screenshots
[Scan](/.screenshots/scan.png?raw=true)

[Submitted](/.screenshots/submitted.png?raw=true)

[Manage students](/.screenshots/manageStudents.png?raw=true)

[Students at location](/.screenshots/studentsAtLocation.png?raw=true)
