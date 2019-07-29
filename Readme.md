Mission Control:
overall design:

1. pkgs for individual AWS components.

2. internal for cli commands for each component.

3. configuration database for state.

4. Build flags for adding just the code you need.

IAM:
Design:

1. Add new user:

Check for config user.
if user exists:
Get permission of user.
compare user with config.
i. return user exists.
ii. Update user:
else:
create user.

# Applications:
## KMS-PASS:


