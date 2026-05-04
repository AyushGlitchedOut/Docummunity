
# Backend
1. Migrate to MySQL
2. Change all errors to custom error using Errors.New() rather than using string comparisons. 
3. Use better data types for tables once migrate to MySQL
4. FullText Search rather than simple SQL comparison upon switching for search Functions
5. Prevent users from accessing '000' deleted User
6. Add Rate-limiting to Routes
7. Check all routes to make sure formats and messages for erros/succession are consistent
8. Handle error from preview and profile Picture saving everywhere
9. Make all DB functions return cosistent errors
10. After the backend is properly complete, make it so the delete function delets the user from firebase too
11. Make it so old profile Pictures and Previews are deleted upon updating the record/user