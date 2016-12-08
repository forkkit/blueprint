#!/bin/bash
set -e

MANAGED=managed
BP_BRANCH=blueprint
TEMP=$RANDOM

SERVICE=$1
[[ -z $SERVICE ]] && echo "Need a service" && exit 1

(
cd "$MANAGED/$SERVICE" || exit 1

# Clean up untracked things
git reset --hard HEAD
git clean -fd

# Make sure remote branches are up to date for service
git fetch -a

# Switch to blueprint branch
git checkout "$BP_BRANCH"

# Base a new branch off that
git checkout -b "${BP_BRANCH}_${TEMP}"
)

./blueprint apply service "$SERVICE"

(
# Show our diff before commit
cd "$MANAGED/$SERVICE" || exit 1
git add .
git diff "$BP_BRANCH"

# Wait on user input to continue
read -r -p "Continue and commit this change to temp branch? [y]/n " first_commit
if [[ "$first_commit" == "n" ]]; then
  echo "Aborted by user feedback"
  exit 1
fi

# Add anything new and commit
git commit -a -m "blueprint update"

# Swap to master branch and make a rebase working branch
git checkout master
git checkout -b "master_${TEMP}"

git rebase -i ${BP_BRANCH}_${TEMP}
)
