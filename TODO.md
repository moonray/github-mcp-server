# TODO.md

## GitHub Projects V2 Tools: Next Steps

### Immediate
- [ ] Debug and fix `TestOwnerResolutionInCreateProject/owner_is_user` failure:
    - [ ] Add debug output in the test handler to print the actual error and request body for this case.
    - [ ] Confirm that `createProjectV2` mutation receives the resolved user ID as `ownerId`.
    - [ ] Inspect struct tags and marshaling for mutation input.
    - [ ] Adjust either test or code until all cases pass.
- [ ] Run full test suite and validate all tests pass with no regressions.

### After All Tests Pass
- [ ] Refactor Projects V2 business logic and MCP tool factories into a single `projects.go` file, matching codebase conventions.
- [ ] Validate that all tool registrations and integrations remain functional after refactor.
- [ ] Update/validate documentation and CONTINUITY.md as needed.

### Ongoing
- [ ] Maintain strict TDD and minimal/DRY Go code.
- [ ] Ensure all GraphQL interactions match GitHub's documented API structure.
- [ ] Keep security best practices for token management.

---
_Last updated: 2025-04-21 03:49:21-04:00_
