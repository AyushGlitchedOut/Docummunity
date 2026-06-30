# Backend
1. Migrate to MySQL
2. Change all errors to custom error using Errors.New() rather than using string comparisons. 
3. Use better data types for tables once migrate to MySQL
4. FullText Search rather than simple SQL comparison upon switching for search Functions
5. [COMPLETED] Prevent users from accessing '000' deleted User
6. [COMPLETED] Add Rate-limiting to Routes
7. [COMPLETED] Check all routes to make sure formats and messages for errors/succession are consistent
8. [COMPLETED] Handle error from preview and profile Picture saving everywhere
9. [COMPLETED] Make all DB functions return cosistent errors
10. [COMPLETED] After the backend is properly complete, make it so the delete function delets the user from firebase too
11. [COMPLETED] Make it so old profile Pictures and Previews are deleted upon updating the record/user 
12. [COMPLETED] Add comments just above the declaration of every function describing what they do so vs code shows the purpose of that function upon hovering when using it.
13. [COMPLETED] Convert all the true-false arguments in updatehandlers, delete handler etc. to query paramters (?query=true) instead of parts of the request body
14. [COMPLETED] Add file type limitations (only document files in records, only image files for preview and pfp) 
- 14) [COMPLETED] Make it so file types are scanned from a dynamic list of allowed filetypes in consts declaration, and that they are checked using MIME types/file headers rather than file extensions
15. [COMPLETED] Create a safe Quitting Configuration for the server upon pressing Ctrl+C
- 15) While Using Go Air or any another fast-reload tool, behavior of graceful shutdown seems unpredictable. But, it works perfectly in normal Binary. So, ignore it if the message for shutting down dosnt appear while using Air or any other tool. 
16. [COMPLETED]  In handlers for Creating Records/ Users, make sure files are deleted if the request fails
17. [FIXED] CRITICAL BUG: Sometimes the first user created doesnt have their files saved
- 17) Apparently what happened was that I had made a block of code that deleted the saved Profile Picture when the user creation in DB failed. Now, I had mistakenly placed that block in the error checking sequence before checking if the error was of UNIQUE CONSTRAINT (basically when the User with same ID tries to create a new user repeatedly). So, basically, whenever I clicked "create User" twice in Bruno (my api tester), even though the API returned that the user already existed and that it had handled the repeated creation, it actually deleted the stored file for the previosuly created user, making it so the user still existed but the profile Picture was deleted. I fixed that by moving the check for user's existence above the profile Picture deletion code.
18. [FIXED] BUG: When using deleteUser with preserveRecords as false, the command fails if the user doesnt have any records, even though the user exists
- 18) What was happening was that in the DB function for deleting records as well, another DB function was being used to fetch all the records for the user. When that DB function returned error since the user didnt have any records, the error was returned  by the function to the handler which was also, due to a logical error, was catching the "records not found" error and showing 404 to the user, despite the user still existing, jst without any records. To fix it, I just removed the error catching of "records not found" in the handler and in the DB function, made it so the "records not found" error case for getUserRecords() is handled properly without returning, so that the user is deleted even if they dont have any records. 
# Frontend
1. (Note to self): Try to avoid the n+1 query problem
