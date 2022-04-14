# switch no to new local branch (combine the steps above)
git checkout -b <new_branch_name>

git add . 
git commit -m "First commit"

# push to a new remote branch not yet created
git push --set-upstream <remote> <branch_name>
git push -u <remote> <branch_name>
