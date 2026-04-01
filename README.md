DOCUMMUNITY- THE DOCUMENT SHARING APP

Stack:
Backend: Go + Gin (Air for Hot Reload)
Frontend: Typescript + Vite + React

NOTE: To run the backend, do `go tool air`, as air is used in the project for hot reload
NOTE: To run the Frontned, do `npm run dev`, and terminate it using `q`, not Ctrl+C, as it runs prettier upon closing
NOTE: Before running the project, setup these:
    1. In Frontend, make a .env file and fill out all fields mentioned in the src/auth/firebaseConfig.ts file (You will get all of them from fireBase console)
    2. For backend, create a file named firebase_key.json and paste the private key for service account (You will get it from ,again, firebase developer console)