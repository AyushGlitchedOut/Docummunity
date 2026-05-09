
# Backend
1. Migrate to MySQL
2. Change all errors to custom error using Errors.New() rather than using string comparisons. 
3. Use better data types for tables once migrate to MySQL
4. FullText Search rather than simple SQL comparison upon switching for search Functions
5. [COMPLETED] Prevent users from accessing '000' deleted User
6. Add Rate-limiting to Routes
7. Check all routes to make sure formats and messages for erros/succession are consistent
8. Handle error from preview and profile Picture saving everywhere
9. Make all DB functions return cosistent errors
10. After the backend is properly complete, make it so the delete function delets the user from firebase too
11. [COMPLETED] Make it so old profile Pictures and Previews are deleted upon updating the record/user 
12. Add comments just above the declaration of evwery function describing what they do so vs code shows the purpose of that function upon hovering when using it.
13. Convert all the true-false arguments in updatehandlers, delete handler etc. to query paramters (?query=true) instead of parts of the request body

# Frontend
1. (Note to self): Try to avoid the n+1 query problem