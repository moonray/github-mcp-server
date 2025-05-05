# CONTINUITY.md

## GitHub Projects V2 Tools Continuity

### Objective
Refine and robustly test the owner resolution logic for GitHub Projects V2 tools, ensuring correct handling of organizations, users, and ambiguous cases, and maintain strict TDD and modularity.

### Key Accomplishments
- Refactored owner resolution logic into a reusable `resolveOwnerID` function.
- All tools (ListOrganizationProjectsTool, ListUserProjectsTool, GetProjectTool, CreateProjectTool) now use this logic.
- Updated and extended tests for all relevant scenarios, including ambiguous and error cases.
- Fixed test panics due to mock exhaustion by queuing multiple mock responses where needed.
- Ensured all GraphQL queries and mutations use the correct owner ID (user or org) as required by GitHub's API.
- All tests for repository resource content and owner resolution now pass except for a single edge case (`owner_is_user`), which is under investigation.

### Outstanding Issues
- `TestOwnerResolutionInCreateProject/owner_is_user` fails, likely due to a mismatch between the expected and actual mutation input or mock handler logic.
- Need to confirm mutation input for `createProjectV2` uses the resolved user ID and matches the expected GraphQL structure.
- All other tests, including ambiguous and negative owner cases, pass.

### Next Steps
1. Add debug output to the test or handler to reveal the actual error message and request body for the failing case.
2. Inspect the mutation input struct and marshaling to ensure the correct field (`ownerId`) is sent.
3. Once all tests pass, refactor Projects V2 logic and tool factories into a single `projects.go` file for convention compliance and PR acceptance.

### Design and Implementation Notes
- Owner resolution is always attempted for both org and user; org is preferred if both exist.
- Mock handlers in tests must return the correct GraphQL response structure to avoid unmarshalling errors.
- All code changes follow TDD, modular, and DRY principles as per user rules.
- Security: GitHub tokens are handled via env var and not hardcoded.

### User Preferences
- Strict TDD and validation using specs/references.
- Modular, minimal, and DRY Go code.
- All logic and factories to be colocated for each feature.

---
_Last updated: 2025-04-21 03:48:58-04:00_
